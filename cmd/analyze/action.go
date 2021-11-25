package analyze

import (
	"encoding/json"
	"fmt"
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
	RunPerDay   float64   `json:"run_per_day"`
	FailureRate float64   `json:"failure_rate"`
	DurationAvg float64   `json:"duration_avg"`
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

	// fmt.Println("workflow_ID, name, number, run_ID, conclusion, status, created_at, updated_at, duration(hours)")
	// for _, v := range idmap {
	// 	for _, r := range v {
	// 		fmt.Printf("%v, %v, %v, %v, %v, %v, %v, %v, %v\n", *r.WorkflowID, *r.Name, *r.RunNumber, *r.ID, *r.Conclusion, *r.Status, r.CreatedAt, r.UpdatedAt, r.UpdatedAt.Sub(r.CreatedAt.Time).Hours())
	// 	}
	// }

	runstats := make(map[int64][]RunStats)
	for k, v := range idmap {
		run, err := GetRunStats(v, c.Int("weeks"))
		if err != nil {
			return fmt.Errorf("get RunStats: %v", err)
		}

		runstats[k] = run
	}

	for _, s := range runstats {
		for _, v := range s {
			fmt.Println(v)
		}
	}

	return nil
}

func GetRunStats(runs []github.WorkflowRun, weeks int) ([]RunStats, error) {
	out := make([]RunStats, 0)
	for _, d := range calendar.LastNWeeks(weeks) {
		start, _ := calendar.Parse(d.Start)
		end, _ := calendar.Parse(d.End)

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

	var rate, avg float64
	if count > 0 {
		rate = failure / count
		avg = duration.Hours() / count
	}

	return RunStats{
		WorkflowID:  *runs[0].WorkflowID,
		Name:        *runs[0].Name,
		Start:       start,
		End:         end,
		RunPerDay:   count / (end.Sub(start).Hours() / 24),
		FailureRate: rate,
		DurationAvg: avg,
	}, nil
}
