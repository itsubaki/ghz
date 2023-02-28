package pullreqs

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
	State      string
	LastID     int64
	LastDay    *time.Time
}

func Fetch(ctx context.Context, in *FetchInput, fn ...func(list []*github.PullRequest) error) ([]*github.PullRequest, error) {
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
		pr, resp, err := client.PullRequests.List(ctx, in.Owner, in.Repository, &opts)
		if err != nil {
			return nil, fmt.Errorf("list pullreqs: %v", err)
		}

		buf := make([]*github.PullRequest, 0)
		var last bool
		for i := range pr {
			if pr[i].GetID() <= in.LastID {
				last = true
				break
			}

			if in.LastDay != nil && pr[i].CreatedAt.Time.Before(*in.LastDay) {
				last = true
				break
			}

			buf = append(buf, pr[i])
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
