package events

import (
	"context"
	"fmt"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

type FetchInput struct {
	Owner      string
	Repository string
	PAT        string
	Page       int
	PerPage    int
}

func Fetch(ctx context.Context, in *FetchInput, fn ...func(list []*github.Event) error) ([]*github.Event, error) {
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

	out := make([]*github.Event, 0)
	for {
		events, resp, err := client.Activity.ListRepositoryEvents(ctx, in.Owner, in.Repository, &opts)
		if err != nil {
			return nil, fmt.Errorf("list workflow runs: %v", err)
		}

		for i, f := range fn {
			if err := f(events); err != nil {
				return nil, fmt.Errorf("func[%v]: %v", i, err)
			}
		}

		out = append(out, events...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return out, nil
}
