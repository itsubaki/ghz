package prstats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

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

type GetStatsInput struct {
	Owner   string
	Repo    string
	PAT     string
	State   string
	PerPage int
}

func GetStats(in *GetStatsInput) (*PRStats, error) {
	ctx := context.Background()
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	var out PRStats
	{
		opt := github.PullRequestListOptions{
			State: in.State,
			ListOptions: github.ListOptions{
				PerPage: in.PerPage,
			},
		}

		list := make([]*github.PullRequest, 0)
		for {
			pr, resp, err := client.PullRequests.List(ctx, in.Owner, in.Repo, &opt)
			if err != nil {
				return nil, fmt.Errorf("list PR: %v", err)
			}

			list = append(list, pr...)
			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}

		out.Owner = in.Owner
		out.Repo = in.Repo
		out.Range.Beg = list[0].CreatedAt
		out.Range.End = list[len(list)-1].CreatedAt

		out.PerDay.Count = len(list)
		out.PerDay.Days = int(list[0].CreatedAt.Sub(*list[len(list)-1].CreatedAt).Hours() / 24)
		out.PerDay.CountPerDay = float64(out.PerDay.Count) / float64(out.PerDay.Days)

		var count int
		var sum float64
		for _, r := range list {
			if r.MergedAt == nil {
				continue
			}

			count++
			sum += r.MergedAt.Sub(*r.CreatedAt).Hours()
		}

		out.Merged.Count = count
		out.Merged.Days = int(list[0].MergedAt.Sub(*list[len(list)-1].MergedAt).Hours() / 24)
		out.Merged.CountPerDay = float64(out.Merged.Count) / float64(out.Merged.Days)

		out.Merged.HoursPerCount = int(sum / float64(count))
		out.Merged.TotalHours = int(sum)
	}

	{
		opt := github.ListWorkflowRunsOptions{
			ListOptions: github.ListOptions{
				PerPage: in.PerPage,
			},
		}

		wmap := make(map[int64][]*github.WorkflowRun)
		for {
			runs, resp, err := client.Actions.ListRepositoryWorkflowRuns(ctx, in.Owner, in.Repo, &opt)
			if err != nil {
				return nil, fmt.Errorf("list Workflow Runs: %v", err)
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

			out.Workflow = append(out.Workflow, Workflow{
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

	return &out, nil
}
