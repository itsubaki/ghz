package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/itsubaki/prstats/pkg/prstats"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	in := prstats.GetStatsInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		State:   c.String("state"),
		PerPage: c.Int("perpage"),
	}

	end, begin, err := timerange(c.Int("days"), c.String("end"), c.String("begin"))
	if err != nil {
		return fmt.Errorf("timerange: %v", err)
	}

	stats, err := prstats.GetStats(context.Background(), &in, end, begin)
	if err != nil {
		return fmt.Errorf("get stats: %v", err)
	}

	format := strings.ToLower(c.String("format"))
	if err := print(format, stats); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}

func timerange(days int, end, begin string) (time.Time, time.Time, error) {
	now := time.Now()
	if begin == "" || end == "" {
		return now, now.AddDate(0, 0, -1*days), nil
	}

	pbegin, err := time.Parse("2006-01-02", begin)
	if err != nil {
		return now, now, fmt.Errorf("parse time=%s: %v", begin, err)
	}

	pend, err := time.Parse("2006-01-02", end)
	if err != nil {
		return now, now, fmt.Errorf("parse time=%s: %v", end, err)
	}

	return pend, pbegin, nil
}

func print(format string, stats *prstats.Stats) error {
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
