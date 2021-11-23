package cmd

import (
	"fmt"
	"strings"

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

	stats, err := prstats.GetStats(&in, c.Int("days"))
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
