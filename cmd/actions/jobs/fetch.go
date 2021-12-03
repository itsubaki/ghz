package jobs

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/prstats/cmd/actions/runs"
	"github.com/itsubaki/prstats/pkg/actions/jobs"
	"github.com/urfave/cli/v2"
)

const Filename = "jobs.json"

func Fetch(c *cli.Context) error {
	runspath := fmt.Sprintf("%v/%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repo"), runs.Filename)
	runslist, err := runs.Deserialize(runspath)
	if err != nil {
		return fmt.Errorf("deserialize: %v", err)
	}

	idmap := make(map[int64][]github.WorkflowRun)
	for _, r := range runslist {
		runs, ok := idmap[*r.WorkflowID]
		if !ok {
			idmap[*r.WorkflowID] = make([]github.WorkflowRun, 0)
		}

		idmap[*r.WorkflowID] = append(runs, r)
	}

	dir := fmt.Sprintf("%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repo"))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	var lastID int64
	path := fmt.Sprintf("%s/%s", dir, Filename)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open %v: %v", path, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			var run github.WorkflowJob
			if err := json.Unmarshal([]byte(scanner.Text()), &run); err != nil {
				return fmt.Errorf("unmarshal: %v", err)
			}

			lastID = *run.ID
		}
	}

	in := jobs.FetchInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		Page:    c.Int("page"),
		PerPage: c.Int("perpage"),
		LastID:  lastID,
	}
	wid := c.Int64("workflow_id")

	ctx := context.Background()
	list := make([]WorkflowJob, 0)
	for _, runs := range idmap {
		for i := range runs {
			if wid > 0 && *runs[i].WorkflowID != wid {
				continue
			}

			jobs, err := jobs.Fetch(ctx, &in, *runs[i].ID)
			if err != nil {
				return fmt.Errorf("fetch: %v", err)
			}

			list = append(list, WorkflowJob{
				WorkflowRun: runs[i],
				WorkflowJob: jobs,
			})
		}
	}

	if err := serialize(path, list); err != nil {
		return fmt.Errorf("serialize: %v", err)
	}

	return nil
}

type WorkflowJob struct {
	WorkflowRun github.WorkflowRun
	WorkflowJob []*github.WorkflowJob
}

func JSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func serialize(path string, list []WorkflowJob) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("open file: %v", err)
	}
	defer file.Close()

	for _, r := range list {
		for _, j := range r.WorkflowJob {
			fmt.Fprintln(file, JSON(j))
		}
	}

	return nil
}
