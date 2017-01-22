package main

import (
	//"errors"
	//"flag"
	//"log"
	//"sort"

	"github.com/stxmendez/gitem"
	//"github.com/google/go-github/github"
	"github.com/urfave/cli"
	//"golang.org/x/oauth2"
	//"gopkg.in/src-d/go-git.v4"
	//"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
)

//type ByRepoURL []*github.Repository
//
//func (a ByRepoURL) Len() int           { return len(a) }
//func (a ByRepoURL) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
//func (a ByRepoURL) Less(i, j int) bool { return *a[i].URL < *a[j].URL }

func list(c *cli.Context) error {
	org := c.String("org")
	if org != "" {
		return gitem.ListRepositoriesForOrg(org, c.String("password"))
	} else {
		return gitem.ListRepositoriesForUser(c.String("user"))
	}
}

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:  "checkout",
			Usage: "checkout repositories",
		},
		{
			Action: list,
			Flags:  []cli.Flag{
				cli.StringFlag{Name: "org"},
				cli.StringFlag{Name: "user"},
				cli.StringFlag{Name: "password"},
			},
			Name:   "list",
			Usage:  "list repositories",
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
