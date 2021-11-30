package jobslist

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v40/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

type ListJobsInput struct {
	Owner   string
	Repo    string
	PAT     string
	PerPage int
}

func ListWorkflowJobs(ctx context.Context, in *ListJobsInput, runID int64) ([]*github.WorkflowJob, error) {
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

func Action(c *cli.Context) error {
	runs, err := deserialize(c.String("path"))
	if err != nil {
		return fmt.Errorf("deserialize: %v", err)
	}

	idmap := make(map[int64][]github.WorkflowRun)
	for _, r := range runs {
		runs, ok := idmap[*r.WorkflowID]
		if !ok {
			idmap[*r.WorkflowID] = make([]github.WorkflowRun, 0)
		}

		idmap[*r.WorkflowID] = append(runs, r)
	}

	in := ListJobsInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		PerPage: c.Int("perpage"),
	}

	fmt.Println("workflow_name, run_id, run_number, job_id, job_name, conclusion, status, started_at, completed_at, duration(minutes)")
	ctx := context.Background()
	for _, runs := range idmap {
		for _, r := range runs {
			jobs, err := ListWorkflowJobs(ctx, &in, *r.ID)
			if err != nil {
				return fmt.Errorf("get WorkflowJobs List: %v", err)
			}

			for _, j := range jobs {
				fmt.Println(CSV(r, *j))
			}
		}
	}

	return nil
}

func CSV(r github.WorkflowRun, j github.WorkflowJob) string {
	return fmt.Sprintf(
		"%v, %v, %v, %v, %v, %v, %v, %v, %v, %v",
		*r.Name,
		*r.ID,
		*r.RunNumber,
		*j.ID,
		*j.Name,
		*j.Conclusion,
		*j.Status,
		*j.StartedAt,
		*j.CompletedAt,
		j.CompletedAt.Sub(j.StartedAt.Time).Minutes(),
	)
}

func deserialize(path string) ([]github.WorkflowRun, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %v", path)
	}

	read, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %v", path, err)
	}

	runs := make([]github.WorkflowRun, 0)
	for _, r := range strings.Split(string(read), "\n") {
		if len(r) < 1 {
			// skip empty line
			continue
		}

		var run github.WorkflowRun
		if err := json.Unmarshal([]byte(r), &run); err != nil {
			return nil, fmt.Errorf("unmarshal: %v", err)
		}

		runs = append(runs, run)
	}

	return runs, nil
}
