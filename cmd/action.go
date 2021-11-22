package cmd

import (
	"context"
	"fmt"

	"github.com/google/go-github/v40/github"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	client := github.NewClient(nil)

	pr, res, err := client.PullRequests.List(context.Background(), "itsubaki", "q", &github.PullRequestListOptions{
		State: "all",
	})
	if err != nil {
		return fmt.Errorf("list PullRequests: %v", err)
	}

	fmt.Printf("res: %v\n", res)

	for _, r := range pr {
		fmt.Printf("ID: %v, title: %v, created: %v, merged: %v, closed: %v\n", *r.ID, *r.Title, r.CreatedAt, r.MergedAt, r.ClosedAt)
	}

	return nil
}
