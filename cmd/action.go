package cmd

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

type PRStats struct {
	Org  string `json:"org"`
	Repo string `json:"repo"`

	Range struct {
		Beg *time.Time `json:"beg"`
		End *time.Time `json:"end"`
	} `json:"range"`

	PRPerDay struct {
		Count int     `json:"count"`
		Days  int     `json:"days"`
		Rate  float64 `json:"rate"`
	} `json:"pr"`

	Lifetime struct {
		AverageHours int `json:"average_hours"`
		TotalHours   int `json:"total_hours"`
		Count        int `json:"count"`
	} `json:"lifetime"`

	Deploy struct {
		Count       int     `json:"count"`
		Succeeded   int     `json:"succeeded"`
		Failed      int     `json:"failed"`
		Cancelled   int     `json:"cancelled"`
		SuccessRate float64 `json:"rate"`
	} `json:"deploy"`
}

func (s PRStats) String() string {
	return s.JSON()
}

func (s PRStats) JSON() string {
	b, err := json.Marshal(s)
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

	list := make([]*github.PullRequest, 0)
	for {
		pr, resp, err := client.PullRequests.List(ctx, c.String("org"), c.String("repo"), &opt)
		if err != nil {
			return fmt.Errorf("list PR: %v", err)
		}

		list = append(list, pr...)
		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	var stats PRStats
	stats.Org = c.String("org")
	stats.Repo = c.String("repo")
	stats.Range.Beg = list[0].CreatedAt
	stats.Range.End = list[len(list)-1].CreatedAt

	stats.PRPerDay.Count = len(list)
	stats.PRPerDay.Days = int(list[0].CreatedAt.Sub(*list[len(list)-1].CreatedAt).Hours() / 24)
	stats.PRPerDay.Rate = float64(stats.PRPerDay.Count) / float64(stats.PRPerDay.Days)

	var count int
	var sum float64
	for _, r := range list {
		if r.MergedAt == nil {
			continue
		}

		count++
		sum += r.MergedAt.Sub(*r.CreatedAt).Hours()
	}

	stats.Lifetime.Count = count
	stats.Lifetime.TotalHours = int(sum)
	stats.Lifetime.AverageHours = int(sum / float64(count))

	fmt.Println(stats)
	return nil
}
