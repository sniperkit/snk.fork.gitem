package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/stxmendez/gitem"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

func authArguments() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{Name: "user"},
		cli.StringFlag{Name: "password"},
		cli.StringFlag{Name: "token"},
	}
}

func clone(c *cli.Context) error {
	client := NewHttpClient(c)
	user := c.String("user")
	if user == "" {
		log.Fatal("user must be specified")
	}

	rootPath := c.String("rootPath")
	if rootPath == "" {
		log.Fatal("rootPath must be specified")
	}

	if user != "" {
		repos, err := gitem.ListRepositoriesForUser(client, c.String("user"))
		if err != nil {
			log.Fatal(err)
		}

		for _, repo := range repos {
			fmt.Printf("Cloning %s\n", *repo.CloneURL)
			err = gitem.Clone(repo, rootPath)
			if err != nil {
				if err == transport.ErrEmptyRemoteRepository {
					fmt.Printf("\trepository is empty\n")
				} else {
					log.Fatal(err)
				}
			}
		}
	}

	return nil
}

func listContributors(c *cli.Context) error {
	httpClient := NewHttpClient(c)
	client := github.NewClient(httpClient)
	contribs, _, err := client.Repositories.ListContributors(c.String("owner"), c.String("repo"), nil)
	if err != nil {
		return err
	}

	for _, contrib := range contribs {
		fmt.Printf("%s\n", *contrib.Login)
	}

	return nil
}

func listLanguages(c *cli.Context) error {
	httpClient := NewHttpClient(c)
	client := github.NewClient(httpClient)
	langs, _, err := client.Repositories.ListLanguages(c.String("owner"), c.String("repo"))
	if err != nil {
		return err
	}

	for k, _ := range langs {
		fmt.Printf("%s\n", k)
	}

	return nil
}

func listRepos(c *cli.Context) error {
	client := NewHttpClient(c)
	org := c.String("org")
	if org != "" {
		repos, err := gitem.ListRepositoriesForOrg(client, org)
		if err != nil {
			log.Fatal(err)
		}
		printRepos(repos)
	} else {
		repos, err := gitem.ListRepositoriesForUser(client, c.String("owner"))
		if err != nil {
			log.Fatal(err)
		}
		printRepos(repos)
	}

	return nil
}

func listTags(c *cli.Context) error {
	httpClient := NewHttpClient(c)
	client := github.NewClient(httpClient)
	opt := &github.ListOptions{}

	for {
		tags, resp, err := client.Repositories.ListTags(c.String("owner"), c.String("repo"), opt)
		if err != nil {
			return err
		}

		for _, tag := range tags {
			fmt.Printf("%s\n", *tag.Name)
		}

		opt.Page = resp.NextPage
		if opt.Page >= resp.LastPage {
			break
		}
	}

	return nil
}

func listTeams(c *cli.Context) error {
	httpClient := NewHttpClient(c)
	client := github.NewClient(httpClient)
	teams, _, err := client.Repositories.ListTeams(c.String("owner"), c.String("repo"), nil)
	if err != nil {
		return err
	}

	for _, team := range teams {
		fmt.Printf("%s\n", *team.Name)
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
			Name: "list",
			Subcommands: []cli.Command{
				{
					Action: listContributors,
					Flags:  append(authArguments(), orgArgument(), ownerArgument(), repoArgument()),
					Name:   "contributors",
				},
				{
					Action: listLanguages,
					Flags:  append(authArguments(), orgArgument(), ownerArgument(), repoArgument()),
					Name:   "languages",
				},
				{
					Action: listRepos,
					Flags:  append(authArguments(), orgArgument(), ownerArgument()),
					Name:   "repos",
					Usage:  "list repositories",
				},
				{
					Action: listTags,
					Flags:  append(authArguments(), orgArgument(), ownerArgument(), repoArgument()),
					Name:   "tags",
				},
				{
					Action: listTeams,
					Flags:  append(authArguments(), orgArgument(), ownerArgument(), repoArgument()),
					Name:   "teams",
				},
			},
		},
	}

	app.Run(os.Args)
}

func orgArgument() cli.Flag {
	return cli.StringFlag{Name: "org"}
}

func ownerArgument() cli.Flag {
	return cli.StringFlag{Name: "owner"}
}

func printRepos(repos []*github.Repository) {
	for _, repo := range repos {
		fmt.Printf("%s\n", *repo.Name)
	}
}

func repoArgument() cli.Flag {
	return cli.StringFlag{Name: "repo"}
}

// NewHttpClient returns a new http.Client which is configured for Basic or Oauth2 Authentication.
func NewHttpClient(c *cli.Context) *http.Client {
	user := c.String("user")
	pass := c.String("password")
	token := c.String("token")
	otp := c.String("otp")

	// Basic authentication
	if user != "" && pass != "" && token == "" {
		transport := github.BasicAuthTransport{Username: user, Password: pass, OTP: otp}
		return transport.Client()
	}

	// Oauth
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		return oauth2.NewClient(oauth2.NoContext, ts)
	}

	// Default case, let the github API apply its default behavior.
	return nil
}