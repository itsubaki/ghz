package runs

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghstats/pkg/actions/runs"
	"github.com/urfave/cli/v2"
)

func Analyze(c *cli.Context) error {
	path := fmt.Sprintf("%v/%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repo"), Filename)
	list, err := Deserialize(path)
	if err != nil {
		return fmt.Errorf("deserialize: %v", err)
	}

	idmap := make(map[int64][]github.WorkflowRun)
	for _, r := range list {
		runs, ok := idmap[*r.WorkflowID]
		if !ok {
			idmap[*r.WorkflowID] = make([]github.WorkflowRun, 0)
		}

		idmap[*r.WorkflowID] = append(runs, r)
	}

	runstats := make(map[int64][]runs.Stats)
	for k, v := range idmap {
		run, err := runs.GetStats(v, c.Int("weeks"), c.Bool("excluding_weekends"))
		if err != nil {
			return fmt.Errorf("runstats: %v", err)
		}

		runstats[k] = run
	}

	format := strings.ToLower(c.String("format"))
	if err := printstats(format, runstats); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}

func printstats(format string, list map[int64][]runs.Stats) error {
	if format == "json" {
		for _, s := range list {
			for _, v := range s {
				fmt.Println(v.JSON())
			}
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("workflow_id, name, start, end, runs_per_day, failure_rate, duration_avg(minutes), duration_var(minutes)")
		for _, s := range list {
			for _, v := range s {
				fmt.Println(v.CSV())
			}
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}
