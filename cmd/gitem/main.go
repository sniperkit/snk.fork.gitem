package main

import (
	//"errors"
	//"flag"
	//"log"
	//"sort"

	"github.com/google/go-github/github"
	"github.com/stxmendez/gitem"
	//"github.com/google/go-github/github"
	"github.com/urfave/cli"
	//"golang.org/x/oauth2"
	//"gopkg.in/src-d/go-git.v4"
	//"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"fmt"
	"log"
	"os"
)

func clone(c *cli.Context) error {
	user := c.String("user")
	if user == "" {
		log.Fatal("user must be specified")
	}

	rootPath := c.String("rootPath")
	if rootPath == "" {
		log.Fatal("rootPath must be specified")
	}

	if user != "" {
		repos, err := gitem.ListRepositoriesForUser(c.String("user"))
		if err != nil {
			log.Fatal(err)
		}

		for _, repo := range repos {
			fmt.Printf("Cloning %s\n", *repo.CloneURL)
			err = gitem.Clone(repo, rootPath)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
}

func printRepos(repos []*github.Repository) {
	for _, repo := range repos {
		fmt.Printf("%s\n", *repo.Name)
	}
}

func list(c *cli.Context) error {
	org := c.String("org")
	if org != "" {
		repos, err := gitem.ListRepositoriesForOrg(org, c.String("password"))
		if err != nil {
			log.Fatal(err)
		}
		printRepos(repos)
	} else {
		repos, err := gitem.ListRepositoriesForUser(c.String("user"))
		if err != nil {
			log.Fatal(err)
		}
		printRepos(repos)
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Action: clone,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "org"},
				cli.StringFlag{Name: "user"},
				cli.StringFlag{Name: "password"},
				cli.StringFlag{Name: "rootPath"},
			},
			Name:  "clone",
			Usage: "clone repositories",
		},
		{
			Action: list,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "org"},
				cli.StringFlag{Name: "user"},
				cli.StringFlag{Name: "password"},
			},
			Name:  "list",
			Usage: "list repositories",
		},
	}

	app.Run(os.Args)

	//var debug bool
	//var accessToken, userName, password, orgName, rootPath string
	//flag.BoolVar(&debug, "debug", false, "Print debug logging information")
	//flag.StringVar(&rootPath, "rootPath", "", "Root path containing checked out subdirectories")
	//flag.StringVar(&userName, "username", "", "Username")
	//flag.StringVar(&password, "password", "", "Password")
	//flag.StringVar(&orgName, "org", "", "git organization name")
	//flag.StringVar(&accessToken, "accessToken", "", "git access token")
	//flag.Parse()
	//
	//log.SetFlags(log.Flags() | log.Lshortfile)
	//
	//if rootPath == "" {
	//	log.Fatal(errors.New("rootPath must be specified"))
	//}
	//
	//if userName == "" {
	//	log.Fatal(errors.New("username must be specified"))
	//}
	//
	//if password == "" {
	//	log.Fatal(errors.New("password must be specified"))
	//}
	//
	//if orgName == "" {
	//	log.Fatal(errors.New("Github orgname must be specified"))
	//}
	//
	//ts := oauth2.StaticTokenSource(
	//	&oauth2.Token{AccessToken: accessToken},
	//)
	//tc := oauth2.NewClient(oauth2.NoContext, ts)
	//
	//opt := &github.RepositoryListByOrgOptions{
	//	ListOptions: github.ListOptions{PerPage: 50},
	//}
	//
	//allRepos := []*github.Repository{}
	//client := github.NewClient(tc)
	//for {
	//	repos, resp, err := client.Repositories.ListByOrg(orgName, opt)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	allRepos = append(allRepos, repos...)
	//
	//	opt.ListOptions.Page = resp.NextPage
	//	if opt.ListOptions.Page >= resp.LastPage {
	//		break
	//	}
	//}
	//
	//log.Printf("%d total repos", len(allRepos))
	//sort.Sort(ByRepoURL(allRepos))
	//
	//for _, repo := range allRepos {
	//	log.Printf("\trepo=%v\n", *repo.URL)
	//
	//	r, err := git.NewFilesystemRepository(rootPath + "/" + *repo.Name)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	r.IsEmpty()
	//	//r := git.NewMemoryRepository()
	//
	//	log.Printf("%#v\n", r)
	//
	//	cloneOptions := &git.CloneOptions{
	//		Auth:       http.NewBasicAuth(userName, accessToken),
	//		RemoteName: "origin",
	//		URL:        *repo.CloneURL,
	//	}
	//
	//	cfg, err := r.Config()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	log.Printf("%#v\n", cfg)
	//	err = r.Clone(cloneOptions)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	// Debugging break early
	//	break
	//}
}
