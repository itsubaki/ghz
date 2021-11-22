package prrate

import (
	"context"
	"encoding/json"
	"fmt"
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

type Output struct {
	Count int     `json:"count"`
	Days  int     `json:"days"`
	Rate  float64 `json:"rate"`
}

func (o Output) String() string {
	return o.JSON()
}

func (o Output) JSON() string {
	b, err := json.Marshal(o)
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

	out := Output{
		Count: len(list),
		Days:  int(list[0].CreatedAt.Sub(*list[len(list)-1].CreatedAt).Hours() / 24),
	}
	out.Rate = float64(out.Count) / float64(out.Days)

	fmt.Println(out)
	return nil
}
