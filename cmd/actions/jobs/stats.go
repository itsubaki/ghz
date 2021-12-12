package jobs

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghstats/pkg/actions/jobs"
	"github.com/urfave/cli/v2"
)

func Stats(c *cli.Context) error {
	path := fmt.Sprintf("%v/%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repo"), Filename)
	list, err := Deserialize(path)
	if err != nil {
		return fmt.Errorf("deserialize: %v", err)
	}

	nmap := make(map[string][]github.WorkflowJob)
	for _, j := range list {
		list, ok := nmap[*j.Name]
		if !ok {
			nmap[*j.Name] = make([]github.WorkflowJob, 0)
		}

		nmap[*j.Name] = append(list, j)
	}

	jobstats := make(map[string][]jobs.Stats)
	for k, v := range nmap {
		run, err := jobs.GetStats(v, c.Int("weeks"), c.Bool("excluding_weekends"))
		if err != nil {
			return fmt.Errorf("stats: %v", err)
		}

		jobstats[k] = run
	}

	format := strings.ToLower(c.String("format"))
	if err := printstats(format, jobstats); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}

func printstats(format string, list map[string][]jobs.Stats) error {
	if format == "json" {
		for _, s := range list {
			for _, v := range s {
				fmt.Println(v.JSON())
			}
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("name, start, end, runs_per_day, failure_rate, duration_avg(minutes), duration_var(minutes)")
		for _, s := range list {
			for _, v := range s {
				fmt.Println(v.CSV())
			}
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}
