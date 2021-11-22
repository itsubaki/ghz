package prlist

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v40/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

type PR struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	MergedAt  *time.Time `json:"merged_at"`
	ClosedAt  *time.Time `json:"closed_at"`
}

func (r PR) String() string {
	return r.JSON()
}

func (r PR) JSON() string {
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(b)
}

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

	list := make([]PR, 0)
	for {
		pr, resp, err := client.PullRequests.List(ctx, c.String("org"), c.String("repo"), &opt)
		if err != nil {
			return fmt.Errorf("list PR: %v", err)
		}

		for _, r := range pr {
			list = append(list, PR{
				ID:        *r.ID,
				Title:     *r.Title,
				CreatedAt: r.CreatedAt,
				UpdatedAt: r.UpdatedAt,
				MergedAt:  r.MergedAt,
				ClosedAt:  r.ClosedAt,
			})
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	format := c.String("format")

	if strings.ToLower(format) == "json" {
		for _, r := range list {
			fmt.Println(r)
		}

		return nil
	}

	if strings.ToLower(format) == "csv" {
		fmt.Println("id, title, created_at, merged_at, updated_at, closed_at, lead_time, ")

		for _, r := range list {
			fmt.Printf("%v, %v, %v, %v, %v, %v, ", r.ID, strings.ReplaceAll(r.Title, ",", ""), r.CreatedAt, r.MergedAt, r.UpdatedAt, r.ClosedAt)
			if r.MergedAt != nil {
				fmt.Printf("%v, ", r.MergedAt.Sub(*r.CreatedAt))
			}

			fmt.Println()
		}
	}

	return nil
}
