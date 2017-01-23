package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"github.com/stxmendez/gitem"
	"github.com/urfave/cli"
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
}

func printRepos(repos []*github.Repository) {
	for _, repo := range repos {
		fmt.Printf("%s\n", *repo.Name)
	}
}
