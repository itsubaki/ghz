package prstats

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/prstats/pkg/calendar"
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
		s.Start.Format("2006-01-02"),
		s.End.Format("2006-01-02"),
		s.RunsPerDay,
		s.FailureRate,
		s.DurationAvg,
		s.DurationVar,
	)
}

func (s RunStats) JSON() string {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func GetRunStats(runs []github.WorkflowRun, weeks int, excludingWeekends bool) ([]RunStats, error) {
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

		stats, err := GetRunStatsWith(runs, end, start, &GetRunStatsWithOptions{
			ExcludingWeekends: excludingWeekends,
		})
		if err != nil {
			return nil, fmt.Errorf("get RunStatsWith(%v~%v): %v", d.End, d.Start, err)
		}

		out = append(out, stats)
	}

	return out, nil
}

type GetRunStatsWithOptions struct {
	ExcludingWeekends bool
}

func GetRunStatsWith(runs []github.WorkflowRun, end, start time.Time, opts *GetRunStatsWithOptions) (RunStats, error) {
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

	runperday := count / (end.Sub(start).Hours() / 24)
	if opts != nil && opts.ExcludingWeekends {
		runperday = count / (end.Sub(start).Hours()/24 - 2) // saturday, sunday
	}

	return RunStats{
		WorkflowID:  *runs[0].WorkflowID,
		Name:        *runs[0].Name,
		Start:       start,
		End:         end,
		RunsPerDay:  runperday,
		FailureRate: rate,
		DurationAvg: avg,
		DurationVar: variant,
	}, nil
}
