/*
Sniperkit-Bot
- Status: analyzed
*/

package gitem

import (
	"net/http"
	"sort"

	"github.com/google/go-github/github"
)

// ByRepoURL struct used to sort repositories by URL.
type ByRepoURL []*github.Repository

func (a ByRepoURL) Len() int           { return len(a) }
func (a ByRepoURL) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRepoURL) Less(i, j int) bool { return *a[i].URL < *a[j].URL }

// ListRepositoriesForOrg list all repos associated with the given organization.
func ListRepositoriesForOrg(httpClient *http.Client, org string) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	// ListByOrg does not allow you to specify a sort order, so we collect them all and then sort.  Less optimal but
	// okay for our purposes.

	allRepos := []*github.Repository{}
	client := github.NewClient(httpClient)
	for {
		repos, resp, err := client.Repositories.ListByOrg(org, opt)
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		opt.ListOptions.Page = resp.NextPage
		if opt.ListOptions.Page >= resp.LastPage {
			break
		}
	}

	sort.Sort(ByRepoURL(allRepos))

	return allRepos, nil
}

// ListRepositoriesForUser lists all repositories owned by the specified user.
func ListRepositoriesForUser(httpClient *http.Client, user string) ([]*github.Repository, error) {
	client := github.NewClient(httpClient)
	opt := github.RepositoryListOptions{
		Type: "all",
	}

	allRepos := []*github.Repository{}
	for {
		repos, resp, err := client.Repositories.List(user, &opt)
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		opt.ListOptions.Page = resp.NextPage
		if opt.ListOptions.Page >= resp.LastPage {
			break
		}
	}

	return allRepos, nil
}
