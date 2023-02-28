package issues

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
	LastID     int64
	LastDay    *time.Time
}

func Fetch(ctx context.Context, in *FetchInput) ([]*github.Issue, error) {
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	opts := github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{
			Page:    in.Page,
			PerPage: in.PerPage,
		},
	}

	out := make([]*github.Issue, 0)
	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, in.Owner, in.Repository, &opts)
		if err != nil {
			return nil, fmt.Errorf("list workflow runs: %v", err)
		}

		buf := make([]*github.Issue, 0)
		var last bool
		for i := range issues {
			if issues[i].GetID() <= in.LastID {
				last = true
				break
			}

			if in.LastDay != nil && issues[i].CreatedAt.Time.Before(*in.LastDay) {
				last = true
				break
			}

			buf = append(buf, issues[i])
		}

		out = append(out, buf...)
		if last || resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return out, nil

}
