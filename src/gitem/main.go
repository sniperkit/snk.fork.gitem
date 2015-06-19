package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type ReposResponse []map[string]interface{}

func execGitCommand(index, total int, dirPath, gitUrl string, args ...string) error {
	log.Printf("%d/%d - %s %s from %s\n", index, total, args[0], dirPath, gitUrl)
	cwd, err := syscall.Getwd()
	if err != nil {
		return err
	}
	err = syscall.Chdir(dirPath)
	if err != nil {
		return err
	}
	gitCmd := exec.Command("git", args...)
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

func main() {
	var userName, password, githubOrgName, rootPath string
	flag.StringVar(&rootPath, "rootPath", "", "Root path containing checked out subdirectories")
	flag.StringVar(&userName, "username", "", "Username")
	flag.StringVar(&password, "password", "", "Password")
	flag.StringVar(&githubOrgName, "githubOrg", "", "Git organization name")
	flag.Parse()

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

	url := fmt.Sprintf("https://api.github.com/orgs/%s/repos?per_page=200", githubOrgName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth(userName, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal(resp.Status)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response ReposResponse
	if err = json.Unmarshal(respBody, &response); err != nil {
		log.Fatal(err)
	}

	if err = processResponse(response, githubOrgName, rootPath); err != nil {
		log.Fatal(err)
	}
}

func processResponse(response ReposResponse, githubOrgName, rootPath string) error {
	log.Printf("Processing %d repos for org %s\n", len(response), githubOrgName)
	total := len(response)
	for index, e := range response {
		gitUrl := e["clone_url"].(string)
		size := e["size"]
		if size.(float64) == 0 {
			log.Printf("Skipping %s as it is empty.\n", gitUrl)
			continue
		}

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
				return err
			}
		} else if fi.IsDir() {
			err := execGitCommand(index, total, dirPath, gitUrl, "pull", "origin")
			if err != nil {
				return err
			}
		} else {
			log.Fatal(errors.New("Expected directory"))
		}
	}

	return nil
}
