package prstats

import (
	"context"
	"fmt"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

type ListPullRequestsInput struct {
	Owner   string
	Repo    string
	PAT     string
	Page    int
	PerPage int
	State   string
}

func ListPullRequests(ctx context.Context, in *ListPullRequestsInput) ([]*github.PullRequest, error) {
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	opts := github.PullRequestListOptions{
		State: in.State,
		ListOptions: github.ListOptions{
			Page:    in.Page,
			PerPage: in.PerPage,
		},
	}

	out := make([]*github.PullRequest, 0)
	for {
		pr, resp, err := client.PullRequests.List(ctx, in.Owner, in.Repo, &opts)
		if err != nil {
			return nil, fmt.Errorf("list PullRequests: %v", err)
		}

		out = append(out, pr...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return out, nil
}
