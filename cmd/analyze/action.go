package analyze

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/prstats/pkg/calendar"
	"github.com/urfave/cli/v2"
)

type RunStats struct {
	WorkflowID  int64     `json:"workflow_id"`
	Name        string    `json:"name"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	RunsPerDay  float64   `json:"runs_per_day"`
	FailureRate float64   `json:"failure_rate"`
	DurationAvg float64   `json:"duration_avg"`
	DurationVar float64   `json:"duration_var"`
}

func (s RunStats) CSV() string {
	return fmt.Sprintf(
		"%v, %v, %v, %v, %v, %v, %v, %v",
		s.WorkflowID,
		s.Name,
		s.Start,
		s.End,
		s.RunsPerDay,
		s.FailureRate,
		s.DurationAvg,
		s.DurationVar,
	)
}

func (s RunStats) String() string {
	return s.JSON()
}

func (s RunStats) JSON() string {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	return string(b)
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

	runstats := make(map[int64][]RunStats)
	for k, v := range idmap {
		run, err := GetRunStats(v, c.Int("weeks"))
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

func print(format string, list map[int64][]RunStats) error {
	if format == "json" {
		for _, s := range list {
			for _, v := range s {
				fmt.Println(v)
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

func GetRunStats(runs []github.WorkflowRun, weeks int) ([]RunStats, error) {
	out := make([]RunStats, 0)
	for _, d := range calendar.LastNWeeks(weeks) {
		start, err := calendar.Parse(d.Start)
		if err != nil {
			return nil, fmt.Errorf("parse %v: %v", d.Start, err)
		}

		end, err := calendar.Parse(d.End)
		if err != nil {
			return nil, fmt.Errorf("parse %v: %v", d.End, err)
		}

		stats, err := GetRunStatsWith(runs, end, start)
		if err != nil {
			return nil, fmt.Errorf("get RunStatsWith(%v~%v): %v", d.End, d.Start, err)
		}

		out = append(out, stats)
	}

	return out, nil
}

func GetRunStatsWith(runs []github.WorkflowRun, end, start time.Time) (RunStats, error) {
	var count, failure float64
	var duration time.Duration
	for _, r := range runs {
		if !r.UpdatedAt.Time.Before(end) || !r.UpdatedAt.Time.After(start) {
			continue
		}

		count++

		if *r.Conclusion == "failure" {
			failure++
		}

		if *r.Conclusion == "success" {
			duration += r.UpdatedAt.Time.Sub(r.CreatedAt.Time)
		}
	}

	var rate, avg, variant float64
	if count > 0 {
		rate = failure / count
		avg = duration.Minutes() / count

		var sum float64
		for _, r := range runs {
			if !r.UpdatedAt.Time.Before(end) || !r.UpdatedAt.Time.After(start) {
				continue
			}

			if *r.Conclusion != "success" {
				continue
			}

			sum = sum + math.Pow((r.UpdatedAt.Time.Sub(r.CreatedAt.Time).Minutes()-avg), 2.0)
		}

		variant = sum / count
	}

	return RunStats{
		WorkflowID:  *runs[0].WorkflowID,
		Name:        *runs[0].Name,
		Start:       start,
		End:         end,
		RunsPerDay:  count / (end.Sub(start).Hours()/24 - 2), // exclude (saturday, sunday)
		FailureRate: rate,
		DurationAvg: avg,
		DurationVar: variant,
	}, nil
}
