package events

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type FetchInput struct {
	Owner      string
	Repository string
	PAT        string
	Page       int
	PerPage    int
	LastID     string
	LastDay    *time.Time
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

		buf := make([]*github.Event, 0)
		var last bool
		for i := range events {
			if in.LastID != "" && events[i].GetID() == in.LastID {
				last = true
				break
			}

			if in.LastDay != nil && events[i].CreatedAt.Time.Before(*in.LastDay) {
				last = true
				break
			}

			buf = append(buf, events[i])
		}

		for i, f := range fn {
			if err := f(buf); err != nil {
				return nil, fmt.Errorf("func[%v]: %v", i, err)
			}
		}

		out = append(out, buf...)
		if last || resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return out, nil
}
