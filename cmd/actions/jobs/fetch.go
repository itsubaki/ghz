package jobs

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghstats/cmd/actions/runs"
	"github.com/itsubaki/ghstats/pkg/actions/jobs"
	"github.com/urfave/cli/v2"
)

const Filename = "jobs.json"

func Fetch(c *cli.Context) error {
	dir := fmt.Sprintf("%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repo"))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	path := fmt.Sprintf("%s/%s", dir, Filename)
	id, err := ScanLastRunID(path)
	if err != nil {
		return fmt.Errorf("last id: %v", err)
	}

	in := jobs.FetchInput{
		Owner:     c.String("owner"),
		Repo:      c.String("repo"),
		PAT:       c.String("pat"),
		Page:      c.Int("page"),
		PerPage:   c.Int("perpage"),
		LastRunID: id,
	}
	wid := c.Int64("workflow_id")

	fmt.Printf("target: %v/%v\n", in.Owner, in.Repo)
	fmt.Printf("workflow_id: %v\n", wid)
	fmt.Printf("last_run_id: %v\n", id)

	runspath := fmt.Sprintf("%v/%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repo"), runs.Filename)
	runs, err := runs.Deserialize(runspath)
	if err != nil {
		return fmt.Errorf("deserialize: %v", err)
	}

	ctx := context.Background()
	for i := range runs {
		if wid > 0 && *runs[i].WorkflowID != wid {
			continue
		}
		if *runs[i].ID <= in.LastRunID {
			continue
		}

		jobs, err := jobs.Fetch(ctx, &in, *runs[i].ID)
		if err != nil {
			return fmt.Errorf("fetch: %v", err)
		}

		if err := serialize(path, jobs); err != nil {
			return fmt.Errorf("serialize: %v", err)
		}

		if len(jobs) > 0 {
			fmt.Printf("%v\n", *jobs[0].RunID)
		}
	}

	return nil
}

func JSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func serialize(path string, list []*github.WorkflowJob) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("open file: %v", err)
	}
	defer file.Close()

	sort.Slice(list, func(i, j int) bool { return *list[i].ID < *list[j].ID }) // asc

	for _, j := range list {
		fmt.Fprintln(file, JSON(j))
	}

	return nil
}

func ScanLastRunID(path string) (int64, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return -1, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return -1, fmt.Errorf("open %v: %v", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var id int64
	for scanner.Scan() {
		var job github.WorkflowJob
		if err := json.Unmarshal([]byte(scanner.Text()), &job); err != nil {
			return -1, fmt.Errorf("unmarshal: %v", err)
		}

		id = *job.RunID
	}

	return id, nil
}
