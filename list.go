package gitem

import (
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func ListRepositoriesForOrg(org, accessToken string) error {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	client := github.NewClient(tc)
	for {
		repos, resp, err := client.Repositories.ListByOrg(org, opt)
		if err != nil {
			return err
		}

		for _, repo := range repos {
			fmt.Printf("%v\n", *repo.Name)
		}

		opt.ListOptions.Page = resp.NextPage
		if opt.ListOptions.Page >= resp.LastPage {
			break
		}
	}

	return nil
}

func ListRepositoriesForUser(user string) error {
	client := github.NewClient(nil)
	opt := github.RepositoryListOptions{
		Sort: "full_name",
	}

	for {
		repos, resp, err := client.Repositories.List(user, &opt)
		if err != nil {
			return err
		}

		for _, repo := range repos {
			fmt.Printf("%v\n", *repo.Name)
		}

		opt.ListOptions.Page = resp.NextPage
		if opt.ListOptions.Page >= resp.LastPage {
			break
		}
	}

	return nil
}
