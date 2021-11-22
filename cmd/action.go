package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/itsubaki/prstats/pkg/prstats"
	"github.com/urfave/cli/v2"
)

type PRStats struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`

	Range struct {
		Beg *time.Time `json:"beg"`
		End *time.Time `json:"end"`
	} `json:"range"`

	PerDay struct {
		CountPerDay float64 `json:"count_per_day"`
		Count       int     `json:"count"`
		Days        int     `json:"days"`
	} `json:"pr"`

	Merged struct {
		CountPerDay   float64 `json:"count_per_day"`
		Count         int     `json:"count"`
		Days          int     `json:"days"`
		HoursPerCount int     `json:"hours_per_count"`
		TotalHours    int     `json:"total_hours"`
	} `json:"merged"`

	Workflow []Workflow `json:"workflow"`
}

type Workflow struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	FailureRate float64 `json:"failure_rate"`
	Count       int     `json:"count"`
	Success     int     `json:"success"`
	Failure     int     `json:"failure"`
	Skipped     int     `json:"skipped"`
	Cancelled   int     `json:"cancelled"`
}

func (s PRStats) String() string {
	return s.JSON()
}

func (s PRStats) JSON() string {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func Action(c *cli.Context) error {
	in := prstats.GetStatsInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		State:   c.String("state"),
		PerPage: c.Int("perpage"),
	}

	stats, err := prstats.GetStats(&in)
	if err != nil {
		return fmt.Errorf("get stats: %v", err)
	}

	format := strings.ToLower(c.String("format"))
	if err := print(format, stats); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}

func print(format string, stats *prstats.PRStats) error {
	if format == "json" {
		fmt.Println(stats)
		return nil
	}

	if format == "csv" {
		fmt.Println("not implemented.")
		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}
