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

type PRStats struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`

	Range struct {
		Beg *time.Time `json:"beg"`
		End *time.Time `json:"end"`
	} `json:"range"`

	PerDay struct {
		CountPerDay float64 `json:"count_per_day"`
		Count       int     `json:"count"`
		Days        int     `json:"days"`
	} `json:"pr"`

	Merged struct {
		CountPerDay   float64 `json:"count_per_day"`
		Count         int     `json:"count"`
		Days          int     `json:"days"`
		HoursPerCount int     `json:"hours_per_count"`
		TotalHours    int     `json:"total_hours"`
	} `json:"merged"`

	Workflow []Workflow `json:"workflow"`
}

type Workflow struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	FailureRate float64 `json:"failure_rate"`
	Count       int     `json:"count"`
	Success     int     `json:"success"`
	Failure     int     `json:"failure"`
	Skipped     int     `json:"skipped"`
	Cancelled   int     `json:"cancelled"`
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
	owner := c.String("owner")
	repo := c.String("repo")

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
			pr, resp, err := client.PullRequests.List(ctx, owner, repo, &opt)
			if err != nil {
				return fmt.Errorf("list PR: %v", err)
			}

			list = append(list, pr...)
			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}

		stats.Owner = owner
		stats.Repo = repo
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

		stats.Merged.Count = count
		stats.Merged.Days = int(list[0].MergedAt.Sub(*list[len(list)-1].MergedAt).Hours() / 24)
		stats.Merged.CountPerDay = float64(stats.Merged.Count) / float64(stats.Merged.Days)

		stats.Merged.HoursPerCount = int(sum / float64(count))
		stats.Merged.TotalHours = int(sum)
	}

	{
		opt := github.ListWorkflowRunsOptions{
			ListOptions: github.ListOptions{
				PerPage: c.Int("perpage"),
			},
		}

		wmap := make(map[int64][]*github.WorkflowRun)
		for {
			runs, resp, err := client.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, &opt)
			if err != nil {
				return fmt.Errorf("list Workflow Runs: %v", err)
			}

			for _, r := range runs.WorkflowRuns {
				if r.Conclusion == nil {
					continue
				}

				runs, ok := wmap[*r.WorkflowID]
				if !ok {
					wmap[*r.WorkflowID] = make([]*github.WorkflowRun, 0)
				}

				wmap[*r.WorkflowID] = append(runs, r)
			}

			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}

		for k, v := range wmap {
			if len(v) < 1 {
				continue
			}

			var success, failure, skipped, cancelled int
			for _, r := range v {
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

			stats.Workflow = append(stats.Workflow, Workflow{
				ID:          k,
				Name:        *v[0].Name,
				FailureRate: float64(failure) / float64(len(v)),
				Count:       len(v),
				Success:     success,
				Failure:     failure,
				Skipped:     skipped,
				Cancelled:   cancelled,
			})
		}

	}

	fmt.Println(stats)
	return nil
}
