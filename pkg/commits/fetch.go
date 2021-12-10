package commits

import (
	"context"
	"fmt"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

type ListInput struct {
	Owner      string
	Repository string
	PAT        string
	Page       int
	PerPage    int
	LastSHA    string
}

func Fetch(ctx context.Context, in *ListInput) ([]*github.RepositoryCommit, error) {
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	opts := github.CommitsListOptions{
		ListOptions: github.ListOptions{
			Page:    in.Page,
			PerPage: in.PerPage,
		},
	}

	out := make([]*github.RepositoryCommit, 0)
	for {
		commits, resp, err := client.Repositories.ListCommits(ctx, in.Owner, in.Repository, &opts)
		if err != nil {
			return nil, fmt.Errorf("list commits: %v", err)
		}

		var last bool
		for i := range commits {
			if in.LastSHA != "" && *commits[i].SHA == in.LastSHA {
				last = true
				break
			}

			out = append(out, commits[i])
		}

		if last || resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return out, nil
}
