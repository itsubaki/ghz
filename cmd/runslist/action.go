package runslist

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/prstats/pkg/prstats"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	in := prstats.ListWorkflowRunsInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		Page:    c.Int("page"),
		PerPage: c.Int("perpage"),
	}

	ctx := context.Background()
	runs, err := prstats.ListWorkflowRuns(ctx, &in)
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
		fmt.Println("workflow_id, workflow_name, run_id, run_number, status, conclusion, created_at, updated_at, duration(m)")

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
		*r.ID,
		*r.RunNumber,
		*r.Status,
		*r.Conclusion,
		r.CreatedAt.Format("2006-01-02 15:04:05"),
		r.UpdatedAt.Format("2006-01-02 15:04:05"),
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
