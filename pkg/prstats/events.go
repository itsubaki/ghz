package prstats

import (
	"context"
	"fmt"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

type ListEventsInput struct {
	Owner   string
	Repo    string
	PAT     string
	Page    int
	PerPage int
}

func ListEvents(ctx context.Context, in *ListEventsInput) ([]*github.Event, error) {
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	opts := github.ListOptions{
		Page:    in.Page,
		PerPage: in.PerPage,
	}

	list := make([]*github.Event, 0)
	for {
		events, resp, err := client.Activity.ListRepositoryEvents(ctx, in.Owner, in.Repo, &opts)
		if err != nil {
			return nil, fmt.Errorf("list WorkflowRuns: %v", err)
		}

		list = append(list, events...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return list, nil
}
