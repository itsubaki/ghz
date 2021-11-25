package analyze

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v40/github"
	"github.com/urfave/cli/v2"
)

type RunStats struct {
	WorkflowID  string
	Name        string
	Begin       time.Time
	End         time.Time
	RunPerDay   float64
	FailureRate float64
	DurationAvg float64
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

	fmt.Println("workflow_ID, name, number, run_ID, conclusion, status, created_at, updated_at, duration(hours)")
	for _, v := range idmap {
		for _, r := range v {
			fmt.Printf("%v, %v, %v, %v, %v, %v, %v, %v, %v\n", *r.WorkflowID, *r.Name, *r.RunNumber, *r.ID, *r.Conclusion, *r.Status, r.CreatedAt, r.UpdatedAt, r.UpdatedAt.Sub(r.CreatedAt.Time).Hours())
		}
	}

	runstats := make(map[int64][]RunStats)
	for k, v := range idmap {
		run, err := GetRunStats(v)
		if err != nil {
			return fmt.Errorf("get RunStats: %v", err)
		}

		runstats[k] = run
	}

	for _, s := range runstats {
		for _, v := range s {
			fmt.Printf("%v\n", v)
		}
	}

	return nil
}

func GetRunStats(runs []github.WorkflowRun) ([]RunStats, error) {
	out := make([]RunStats, 0)

	date := calender.Last12Months()
	for _, d := range date {
		fmt.Printf("%v %v %v\n", d.Period, d.Start, d.End)
	}

	return out, nil
}

func GetRunStatsWith(runs []github.WorkflowRun, end, begin string) ([]RunStats, error) {

	return nil, nil
}
