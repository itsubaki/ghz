package pr

import (
	"context"
	"fmt"

	"github.com/google/go-github/v40/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

func Action(c *cli.Context) error {
	ctx := context.Background()
	client := github.NewClient(nil)

	pat := c.String("pat")
	if pat != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: pat},
		)))
	}

	opt := github.PullRequestListOptions{
		State: c.String("state"),
		ListOptions: github.ListOptions{
			PerPage: c.Int("perpage"),
		},
	}

	for {
		pr, resp, err := client.PullRequests.List(ctx, c.String("org"), c.String("repo"), &opt)
		if err != nil {
			return fmt.Errorf("list PR: %v", err)
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage

		for _, r := range pr {
			fmt.Printf("ID: %v, title: %v, created: %v, merged: %v, closed: %v\n", *r.ID, *r.Title, r.CreatedAt, r.MergedAt, r.ClosedAt)
		}
	}

	return nil
}
