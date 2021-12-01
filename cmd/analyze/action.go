package analyze

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/prstats/pkg/prstats"
	"github.com/urfave/cli/v2"
)

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

	runstats := make(map[int64][]prstats.RunStats)
	for k, v := range idmap {
		run, err := prstats.GetRunStats(v, c.Int("weeks"), c.Bool("excludingWeekends"))
		if err != nil {
			return fmt.Errorf("get RunStats: %v", err)
		}

		runstats[k] = run
	}

	format := strings.ToLower(c.String("format"))
	if err := print(format, runstats); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}

func print(format string, list map[int64][]prstats.RunStats) error {
	if format == "json" {
		for _, s := range list {
			for _, v := range s {
				fmt.Println(v.JSON())
			}
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("workflow_ID, name, start, end, runs_per_day, failure_rate, duration_avg(m), duration_var(m)")
		for _, s := range list {
			for _, v := range s {
				fmt.Println(v.CSV())
			}
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}
