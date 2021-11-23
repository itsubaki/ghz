package prstats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

type WorkflowRun struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	CountPerDay    float64 `json:"count_per_day"`
	FailureRate    float64 `json:"failure_rate"`
	Success        int     `json:"success"`
	Failure        int     `json:"failure"`
	Skipped        int     `json:"skipped"`
	Cancelled      int     `json:"cancelled"`
	ActionRequired int     `json:"action_required"`
	Count          int     `json:"count"`
}

type Range struct {
	Beg  time.Time `json:"beg"`
	End  time.Time `json:"end"`
	Days int       `json:"days"`
}

type Created struct {
	CountPerDay float64 `json:"count_per_day"`
	Count       int     `json:"count"`
}

type Merged struct {
	CountPerDay   float64 `json:"count_per_day"`
	HoursPerCount float64 `json:"hours_per_count"`
	TotalHours    float64 `json:"total_hours"`
	Count         int     `json:"count"`
}

type PRStats struct {
	Owner        string        `json:"owner"`
	Repo         string        `json:"repo"`
	Range        Range         `json:"range"`
	Created      Created       `json:"created"`
	Merged       Merged        `json:"merged"`
	WorkflowRuns []WorkflowRun `json:"workflow_runs"`
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

func GetPRList(ctx context.Context, in *GetStatsInput, begin time.Time) ([]*github.PullRequest, error) {
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	opt := github.PullRequestListOptions{
		State: in.State,
		ListOptions: github.ListOptions{
			PerPage: in.PerPage,
		},
	}

	skip := false
	out := make([]*github.PullRequest, 0)
	for {
		pr, resp, err := client.PullRequests.List(ctx, in.Owner, in.Repo, &opt)
		if err != nil {
			return nil, fmt.Errorf("list PR: %v", err)
		}

		for i := range pr {
			if pr[i].CreatedAt.Unix() < begin.Unix() {
				skip = true
				break
			}

			out = append(out, pr[i])
		}

		if resp.NextPage == 0 || skip {
			break
		}

		opt.Page = resp.NextPage
	}

	return out, nil
}

func GetMergedCount(list []*github.PullRequest) (int, float64, float64) {
	var count int
	var total float64
	for _, r := range list {
		if r.MergedAt == nil {
			continue
		}

		count++
		total += r.MergedAt.Sub(*r.CreatedAt).Hours()
	}

	percount := 0.0
	if count > 0 {
		percount = total / float64(count)
	}

	return count, total, percount
}

func GetWorflowRunsList(ctx context.Context, in *GetStatsInput, begin time.Time) ([]WorkflowRun, error) {
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	opt := github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{
			PerPage: in.PerPage,
		},
	}

	skip := false
	list := make(map[int64][]*github.WorkflowRun)
	for {
		runs, resp, err := client.Actions.ListRepositoryWorkflowRuns(ctx, in.Owner, in.Repo, &opt)
		if err != nil {
			return nil, fmt.Errorf("list Workflow Runs: %v", err)
		}

		for _, r := range runs.WorkflowRuns {
			if r.Conclusion == nil {
				continue
			}

			if r.UpdatedAt.Unix() < begin.Unix() {
				skip = true
				break
			}

			runs, ok := list[*r.WorkflowID]
			if !ok {
				list[*r.WorkflowID] = make([]*github.WorkflowRun, 0)
			}

			list[*r.WorkflowID] = append(runs, r)
		}

		if resp.NextPage == 0 || skip {
			break
		}

		opt.Page = resp.NextPage
	}

	out := make([]WorkflowRun, 0)
	for k, v := range list {
		if len(v) < 1 {
			continue
		}

		var success, failure, skipped, cancelled, required int
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

			if *r.Conclusion == "action_required" {
				required++
				continue
			}

			return nil, fmt.Errorf("invalid conclusion=%v", *r.Conclusion)
		}

		w := WorkflowRun{
			ID:             k,
			Name:           *v[0].Name,
			Success:        success,
			Failure:        failure,
			Skipped:        skipped,
			Cancelled:      cancelled,
			ActionRequired: required,
			Count:          len(v),
		}

		if w.Count > 0 {
			w.FailureRate = float64(w.Failure) / float64(w.Count)
		}

		out = append(out, w)
	}

	return out, nil
}

func GetStats(in *GetStatsInput, days int) (*PRStats, error) {
	end := time.Now()
	beg := end.AddDate(0, 0, -1*days)

	ctx := context.Background()
	created, err := GetPRList(ctx, in, beg)
	if err != nil {
		return nil, fmt.Errorf("get PR list: %v", err)
	}
	merged, total, percount := GetMergedCount(created)

	runs, err := GetWorflowRunsList(ctx, in, beg)
	if err != nil {
		return nil, fmt.Errorf("get WorkflowRuns list: %v", err)
	}

	for i := range runs {
		runs[i].CountPerDay = float64(runs[i].Count) / float64(days)
	}

	return &PRStats{
		Owner: in.Owner,
		Repo:  in.Repo,
		Range: Range{
			Beg:  beg,
			End:  end,
			Days: days,
		},
		Created: Created{
			CountPerDay: float64(len(created)) / float64(days),
			Count:       len(created),
		},
		Merged: Merged{
			CountPerDay:   float64(merged) / float64(days),
			HoursPerCount: percount,
			TotalHours:    total,
			Count:         merged,
		},
		WorkflowRuns: runs,
	}, nil
}
