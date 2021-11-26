package runslist

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v40/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

type ListWorkflowRunsInput struct {
	Owner   string
	Repo    string
	PAT     string
	PerPage int
}

func ListWorkflowRuns(ctx context.Context, in *ListWorkflowRunsInput) ([]*github.WorkflowRun, error) {
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

	list := make([]*github.WorkflowRun, 0)
	for {
		runs, resp, err := client.Actions.ListRepositoryWorkflowRuns(ctx, in.Owner, in.Repo, &opt)
		if err != nil {
			return nil, fmt.Errorf("list WorkflowRuns: %v", err)
		}

		list = append(list, runs.WorkflowRuns...)
		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return list, nil
}

func Action(c *cli.Context) error {
	in := ListWorkflowRunsInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		PerPage: c.Int("perpage"),
	}

	ctx := context.Background()
	runs, err := ListWorkflowRuns(ctx, &in)
	if err != nil {
		return fmt.Errorf("get WorkflowRuns List: %v", err)
	}

	format := strings.ToLower(c.String("format"))
	if err := print(format, runs); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}

func print(format string, list []*github.WorkflowRun) error {
	if format == "json" {
		for _, r := range list {
			fmt.Println(JSON(r))
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("workflow_ID, name, number, run_ID, conclusion, status, created_at, updated_at, duration(m)")

		for _, r := range list {
			fmt.Println(CSV(r))
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}

func CSV(r *github.WorkflowRun) string {
	return fmt.Sprintf(
		"%v, %v, %v, %v, %v, %v, %v, %v, %v",
		*r.WorkflowID,
		*r.Name,
		*r.RunNumber,
		*r.ID,
		*r.Conclusion,
		*r.Status,
		r.CreatedAt,
		r.UpdatedAt,
		r.UpdatedAt.Sub(r.CreatedAt.Time).Minutes(),
	)
}

func JSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func ListWorkflowJobs(ctx context.Context, in *ListWorkflowRunsInput, runID int64) ([]*github.WorkflowJob, error) {
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	opt := github.ListWorkflowJobsOptions{
		ListOptions: github.ListOptions{
			PerPage: in.PerPage,
		},
	}

	list := make([]*github.WorkflowJob, 0)
	for {
		jobs, resp, err := client.Actions.ListWorkflowJobs(ctx, in.Owner, in.Repo, runID, &opt)
		if err != nil {
			return nil, fmt.Errorf("list WorkflowJobs: %v", err)
		}

		list = append(list, jobs.Jobs...)
		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return list, nil
}
