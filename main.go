package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
	"syscall"

	"github.com/golang/gddo/httputil/header"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type ReposResponse []map[string]interface{}

func execGitCommand(index, total int, dirPath, gitUrl string, args ...string) error {
	log.Printf("%d/%d - git %s %s\n", index, total, args[0], gitUrl)
	cwd, err := syscall.Getwd()
	if err != nil {
		return err
	}
	err = syscall.Chdir(dirPath)
	if err != nil {
		return err
	}
	gitCmd := exec.Command("git", args...)
	log.Printf("%v\n", gitCmd)
	gitIn, _ := gitCmd.StdinPipe()
	gitOut, _ := gitCmd.StdoutPipe()
	gitCmd.Start()
	gitIn.Close()
	buf, err := ioutil.ReadAll(gitOut)
	if err != nil {
		return err
	}
	err = gitCmd.Wait()
	if err != nil {
		fmt.Print(string(buf))
	}
	syscall.Chdir(cwd)
	return err
}

func fetchRepos(debug bool, url, userName, password string) (nextPage string, response ReposResponse, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.SetBasicAuth(userName, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("%v - url=%v", resp.Status, url)
		return
	}

	params := header.ParseList(resp.Header, "Link")
	if len(params) > 0 {
		for _, param := range params {
			if strings.Index(param, "rel=\"next\"") != -1 {
				s, e := strings.Index(param, "<"), strings.Index(param, ">")
				if s != -1 && e != -1 {
					nextPage = param[s+1 : e]
				}
			}
		}
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if debug {
		var out bytes.Buffer
		json.Indent(&out, respBody, "", "\t")
		out.WriteTo(os.Stdout)
	}

	err = json.Unmarshal(respBody, &response)

	return nextPage, response, err
}

type ByRepoURL []*github.Repository

func (a ByRepoURL) Len() int           { return len(a) }
func (a ByRepoURL) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRepoURL) Less(i, j int) bool { return *a[i].URL < *a[j].URL }

func main() {
	var debug bool
	var userName, password, githubOrgName, rootPath string
	flag.BoolVar(&debug, "debug", false, "Print debug logging information")
	flag.StringVar(&rootPath, "rootPath", "", "Root path containing checked out subdirectories")
	flag.StringVar(&userName, "username", "", "Username")
	flag.StringVar(&password, "password", "", "Password")
	flag.StringVar(&githubOrgName, "githubOrg", "", "Git organization name")
	flag.Parse()

	log.SetFlags(log.Flags() | log.Lshortfile)

	if rootPath == "" {
		log.Fatal(errors.New("rootPath must be specified"))
	}

	if userName == "" {
		log.Fatal(errors.New("username must be specified"))
	}

	if password == "" {
		log.Fatal(errors.New("password must be specified"))
	}

	if githubOrgName == "" {
		log.Fatal(errors.New("Github orgname must be specified"))
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "142d328e59a56d3cdea11b51e5a0d9cf711042e3"},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	allRepos := []*github.Repository{}
	client := github.NewClient(tc)
	for {
		repos, resp, err := client.Repositories.ListByOrg(githubOrgName, opt)
		if err != nil {
			log.Fatal(err)
		}

		allRepos = append(allRepos, repos...)

		opt.ListOptions.Page = resp.NextPage
		if opt.ListOptions.Page >= resp.LastPage {
			break
		}
	}

	log.Printf("%d total repos", len(allRepos))
	sort.Sort(ByRepoURL(allRepos))

	for i, repo := range allRepos {
		log.Printf("\trepo=%v\n", *repo.URL)
		// func execGitCommand(index, total int, dirPath, gitUrl string, args ...string) error {
		err := execGitCommand(i, len(allRepos), rootPath, *repo.CloneURL, "clone", *repo.CloneURL)
		if err != nil {
			log.Fatal(err)
		}
	}

	//// Fetch the initial set of repos and continue until no more links are returned.
	//url := fmt.Sprintf("https://api.github.com/orgs/%s/repos?per_page=200&type=all", githubOrgName)
	//var aggregatedResponse ReposResponse
	//for nextPage := url; nextPage != ""; {
	//	var response ReposResponse
	//	var err error
	//	nextPage, response, err = fetchRepos(debug, nextPage, userName, password)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	aggregatedResponse = append(aggregatedResponse, response...)
	//}
	//
	//if err := processResponse(aggregatedResponse, githubOrgName, rootPath); err != nil {
	//	log.Fatal(err)
	//}
}

func processResponse(response ReposResponse, githubOrgName, rootPath string) error {
	log.Printf("Processing %d repos for org %s\n", len(response), githubOrgName)
	total := len(response)
	for index, e := range response {
		gitUrl := e["clone_url"].(string)

		if !strings.HasSuffix(gitUrl, ".git") {
			log.Fatal(errors.New("Missing '.git' from URL.."))
		}

		i := strings.LastIndex(gitUrl, "/")
		if i == -1 {
			log.Fatal(errors.New("Invalid url format..."))
		}

		dirName := gitUrl[i+1 : len(gitUrl)-4]
		dirPath := fmt.Sprintf("%s%s%s", rootPath, string(os.PathSeparator), dirName)
		fi, err := os.Stat(dirPath)
		if err != nil {
			// Doesn't exist, checkout
			err := execGitCommand(index, total, rootPath, gitUrl, "clone", gitUrl)
			if err != nil {
				log.Println(err)
			}
		} else if fi.IsDir() {
			err := execGitCommand(index, total, dirPath, gitUrl, "pull", "origin")
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Fatal(errors.New("Expected directory"))
		}
	}

	return nil
}
