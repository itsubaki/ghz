package jobslist

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/prstats/pkg/prstats"
	"github.com/urfave/cli/v2"
)

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

	in := prstats.ListJobsInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		Page:    c.Int("page"),
		PerPage: c.Int("perpage"),
	}
	wid := c.Int64("workflow_id")

	ctx := context.Background()
	list := make([]WorkflowJob, 0)
	for _, runs := range idmap {
		for i := range runs {
			if wid > 0 && *runs[i].WorkflowID != wid {
				continue
			}

			jobs, err := prstats.ListWorkflowJobs(ctx, &in, *runs[i].ID)
			if err != nil {
				return fmt.Errorf("get WorkflowJobs List: %v", err)
			}

			list = append(list, WorkflowJob{
				WorkflowRun: runs[i],
				WorkflowJob: jobs,
			})
		}
	}

	format := strings.ToLower(c.String("format"))
	if err := print(format, list); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
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

type WorkflowJob struct {
	WorkflowRun github.WorkflowRun
	WorkflowJob []*github.WorkflowJob
}

func print(format string, list []WorkflowJob) error {
	if format == "json" {
		for _, r := range list {
			for _, j := range r.WorkflowJob {
				fmt.Println(JSON(j))
			}
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("workflow_id, workflow_name, run_id, run_number, job_id, job_name, conclusion, status, started_at, completed_at, duration(minutes)")
		for _, r := range list {
			for _, j := range r.WorkflowJob {
				fmt.Println(CSV(r.WorkflowRun, j))
			}
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}

func CSV(r github.WorkflowRun, j *github.WorkflowJob) string {
	return fmt.Sprintf(
		"%v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v",
		*r.WorkflowID,
		*r.Name,
		*r.ID,
		*r.RunNumber,
		*j.ID,
		*j.Name,
		*j.Conclusion,
		*j.Status,
		j.StartedAt.Format("2006-01-02 15:04:05"),
		j.CompletedAt.Format("2006-01-02 15:04:05"),
		j.CompletedAt.Sub(j.StartedAt.Time).Minutes(),
	)
}

func JSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(b)
}
