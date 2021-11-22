package cmd

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

type PRStats struct {
	Org  string `json:"org"`
	Repo string `json:"repo"`

	Range struct {
		Beg *time.Time `json:"beg"`
		End *time.Time `json:"end"`
	} `json:"range"`

	PerDay struct {
		Count       int     `json:"count"`
		Days        int     `json:"days"`
		CountPerDay float64 `json:"count_per_day"`
	} `json:"pr"`

	Lifetime struct {
		AverageHours int `json:"average_hours"`
		TotalHours   int `json:"total_hours"`
		MergedCount  int `json:"merged_count"`
	} `json:"lifetime"`

	Deploy struct {
		WorkflowName string  `json:"workflow_name"`
		Count        int     `json:"count"`
		Success      int     `json:"success"`
		Failure      int     `json:"failure"`
		Skipped      int     `json:"skipped"`
		Cancelled    int     `json:"cancelled"`
		FailureRate  float64 `json:"failure_rate"`
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

	var stats PRStats
	{

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

		stats.Org = c.String("org")
		stats.Repo = c.String("repo")
		stats.Range.Beg = list[0].CreatedAt
		stats.Range.End = list[len(list)-1].CreatedAt

		stats.PerDay.Count = len(list)
		stats.PerDay.Days = int(list[0].CreatedAt.Sub(*list[len(list)-1].CreatedAt).Hours() / 24)
		stats.PerDay.CountPerDay = float64(stats.PerDay.Count) / float64(stats.PerDay.Days)

		var count int
		var sum float64
		for _, r := range list {
			if r.MergedAt == nil {
				continue
			}

			count++
			sum += r.MergedAt.Sub(*r.CreatedAt).Hours()
		}

		stats.Lifetime.MergedCount = count
		stats.Lifetime.TotalHours = int(sum)
		stats.Lifetime.AverageHours = int(sum / float64(count))
	}

	{
		opt := github.ListWorkflowRunsOptions{
			ListOptions: github.ListOptions{
				PerPage: c.Int("perpage"),
			},
		}

		list := make([]*github.WorkflowRun, 0)
		for {
			runs, resp, err := client.Actions.ListRepositoryWorkflowRuns(ctx, c.String("org"), c.String("repo"), &opt)
			if err != nil {
				return fmt.Errorf("list PR: %v", err)
			}

			for _, r := range runs.WorkflowRuns {
				if strings.ToLower(c.String("workflow")) != *r.Name {
					continue
				}

				if r.Conclusion == nil {
					continue
				}

				list = append(list, runs.WorkflowRuns...)
			}

			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}

		var success, failure, skipped, cancelled int
		for _, r := range list {
			if *r.Conclusion == "success" {
				success++
				continue
			}

			if *r.Conclusion == "failure" {
				failure++
				continue
			}

			if *r.Conclusion == "skipped" {
				skipped++
				continue
			}

			if *r.Conclusion == "cancelled" {
				cancelled++
				continue
			}
		}

		stats.Deploy.WorkflowName = c.String("workflow")
		if len(list) > 0 {
			stats.Deploy.Count = len(list)
			stats.Deploy.Success = success
			stats.Deploy.Failure = failure
			stats.Deploy.Skipped = skipped
			stats.Deploy.Cancelled = cancelled
			stats.Deploy.FailureRate = float64(failure) / float64(len(list))
		}
	}

	fmt.Println(stats)
	return nil
}
