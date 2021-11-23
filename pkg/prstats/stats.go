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
		Beg  time.Time `json:"beg"`
		End  time.Time `json:"end"`
		Days int       `json:"days"`
	} `json:"range"`

	PerDay struct {
		CountPerDay float64 `json:"count_per_day"`
		Count       int     `json:"count"`
	} `json:"pr"`

	Merged struct {
		CountPerDay   float64 `json:"count_per_day"`
		Count         int     `json:"count"`
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
	Week    int
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
	out.Owner = in.Owner
	out.Repo = in.Repo
	out.Range.End = time.Now()
	out.Range.Beg = out.Range.End.AddDate(0, 0, -7*in.Week)
	out.Range.Days = 7 * in.Week

	{
		opt := github.PullRequestListOptions{
			State: in.State,
			ListOptions: github.ListOptions{
				PerPage: in.PerPage,
			},
		}

		skip := false
		list := make([]*github.PullRequest, 0)
		for {
			pr, resp, err := client.PullRequests.List(ctx, in.Owner, in.Repo, &opt)
			if err != nil {
				return nil, fmt.Errorf("list PR: %v", err)
			}

			for i := range pr {
				if pr[i].CreatedAt.Unix() < out.Range.Beg.Unix() {
					skip = true
					break
				}

				list = append(list, pr[i])
			}

			if resp.NextPage == 0 || skip {
				break
			}

			opt.Page = resp.NextPage
		}

		out.PerDay.Count = len(list)
		out.PerDay.CountPerDay = float64(out.PerDay.Count) / float64(out.Range.Days)

		var count int
		var total float64
		for _, r := range list {
			if r.MergedAt == nil {
				continue
			}

			if r.MergedAt.Unix() < out.Range.Beg.Unix() {
				continue
			}

			count++
			total += r.MergedAt.Sub(*r.CreatedAt).Hours()
		}

		out.Merged.Count = count
		out.Merged.CountPerDay = float64(out.Merged.Count) / float64(out.Range.Days)
		out.Merged.TotalHours = int(total)
		if count > 0 {
			out.Merged.HoursPerCount = int(total / float64(count))
		}
	}

	{
		opt := github.ListWorkflowRunsOptions{
			ListOptions: github.ListOptions{
				PerPage: in.PerPage,
			},
		}

		skip := false
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

				if r.UpdatedAt.Unix() < out.Range.Beg.Unix() {
					skip = true
					break
				}

				runs, ok := wmap[*r.WorkflowID]
				if !ok {
					wmap[*r.WorkflowID] = make([]*github.WorkflowRun, 0)
				}

				wmap[*r.WorkflowID] = append(runs, r)
			}

			if resp.NextPage == 0 || skip {
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

				return nil, fmt.Errorf("invalid conclusion=%v", *r.Conclusion)
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
