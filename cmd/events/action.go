package events

import (
	"context"
	"fmt"

	"github.com/google/go-github/v40/github"
	"github.com/urfave/cli/v2"
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

func Action(c *cli.Context) error {
	in := ListEventsInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		Page:    c.Int("page"),
		PerPage: c.Int("perpage"),
	}

	events, err := ListEvents(context.Background(), &in)
	if err != nil {
		return fmt.Errorf("get Events List: %v", err)
	}

	for _, e := range events {
		fmt.Printf("%v %v %v %v %v\n", *e.ID, *e.Actor, *e.Repo, *e.Type, string(*e.RawPayload))
	}

	return nil
}
