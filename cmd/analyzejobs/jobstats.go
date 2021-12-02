package analyzejobs

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/prstats/pkg/prstats"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	jobs, err := deserialize(c.String("path"))
	if err != nil {
		return fmt.Errorf("deserialize: %v", err)
	}

	nmap := make(map[string][]github.WorkflowJob)
	for _, j := range jobs {
		list, ok := nmap[*j.Name]
		if !ok {
			nmap[*j.Name] = make([]github.WorkflowJob, 0)
		}

		nmap[*j.Name] = append(list, j)
	}

	jobstats := make(map[string][]prstats.JobStats)
	for k, v := range nmap {
		run, err := prstats.GetJobStats(v, c.Int("weeks"), c.Bool("excluding_weekends"))
		if err != nil {
			return fmt.Errorf("get JobStats: %v", err)
		}

		jobstats[k] = run
	}

	format := strings.ToLower(c.String("format"))
	if err := print(format, jobstats); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}

func deserialize(path string) ([]github.WorkflowJob, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %v", path)
	}

	read, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %v", path, err)
	}

	jobs := make([]github.WorkflowJob, 0)
	for _, r := range strings.Split(string(read), "\n") {
		if len(r) < 1 {
			// skip empty line
			continue
		}

		var job github.WorkflowJob
		if err := json.Unmarshal([]byte(r), &job); err != nil {
			return nil, fmt.Errorf("unmarshal: %v", err)
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func print(format string, list map[string][]prstats.JobStats) error {
	if format == "json" {
		for _, s := range list {
			for _, v := range s {
				fmt.Println(v.JSON())
			}
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("name, start, end, runs_per_day, failure_rate, duration_avg(m), duration_var(m)")
		for _, s := range list {
			for _, v := range s {
				fmt.Println(v.CSV())
			}
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}
