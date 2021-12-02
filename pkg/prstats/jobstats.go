package prstats

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/prstats/pkg/calendar"
)

type JobStats struct {
	Name        string    `json:"name"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	RunsPerDay  float64   `json:"runs_per_day"`
	FailureRate float64   `json:"failure_rate"`
	DurationAvg float64   `json:"duration_avg"`
	DurationVar float64   `json:"duration_var"`
}

func (s JobStats) CSV() string {
	return fmt.Sprintf(
		"%v, %v, %v, %v, %v, %v, %v",
		s.Name,
		s.Start.Format("2006-01-02"),
		s.End.Format("2006-01-02"),
		s.RunsPerDay,
		s.FailureRate,
		s.DurationAvg,
		s.DurationVar,
	)
}

func (s JobStats) JSON() string {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func GetJobStats(jobs []github.WorkflowJob, weeks int, excludingWeekends bool) ([]JobStats, error) {
	out := make([]JobStats, 0)
	for _, d := range calendar.LastNWeeks(weeks) {
		start, err := calendar.Parse(d.Start)
		if err != nil {
			return nil, fmt.Errorf("parse %v: %v", d.Start, err)
		}

		end, err := calendar.Parse(d.End)
		if err != nil {
			return nil, fmt.Errorf("parse %v: %v", d.End, err)
		}

		stats, err := GetJobStatsWith(jobs, end, start, &GetJobStatsWithOptions{
			ExcludingWeekends: excludingWeekends,
		})
		if err != nil {
			return nil, fmt.Errorf("get RunStatsWith(%v~%v): %v", d.End, d.Start, err)
		}

		out = append(out, stats)
	}

	return out, nil
}

type GetJobStatsWithOptions struct {
	ExcludingWeekends bool
}

func GetJobStatsWith(jobs []github.WorkflowJob, end, start time.Time, opts *GetJobStatsWithOptions) (JobStats, error) {
	var count, failure float64
	var duration time.Duration
	for _, j := range jobs {
		if !j.CompletedAt.Before(end) || !j.CompletedAt.After(start) {
			continue
		}

		count++

		if *j.Conclusion == "failure" {
			failure++
		}

		if *j.Conclusion == "success" {
			duration += j.CompletedAt.Sub(j.StartedAt.Time)
		}
	}

	var rate, avg, variant float64
	if count > 0 {
		rate = failure / count
		avg = duration.Minutes() / count

		var sum float64
		for _, j := range jobs {
			if !j.CompletedAt.Before(end) || !j.CompletedAt.After(start) {
				continue
			}

			if *j.Conclusion != "success" {
				continue
			}

			sum = sum + math.Pow((j.CompletedAt.Sub(j.StartedAt.Time).Minutes()-avg), 2.0)
		}

		variant = sum / count
	}

	runperday := count / (end.Sub(start).Hours() / 24)
	if opts != nil && opts.ExcludingWeekends {
		runperday = count / (end.Sub(start).Hours()/24 - 2) // saturday, sunday
	}

	return JobStats{
		Name:        *jobs[0].Name,
		Start:       start,
		End:         end,
		RunsPerDay:  runperday,
		FailureRate: rate,
		DurationAvg: avg,
		DurationVar: variant,
	}, nil
}
