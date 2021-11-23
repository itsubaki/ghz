package prlist

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/prstats/pkg/prstats"
	"github.com/urfave/cli/v2"
)

type PR struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	MergedAt  *time.Time `json:"merged_at"`
	ClosedAt  *time.Time `json:"closed_at"`
}

func (r PR) String() string {
	return r.JSON()
}

func (r PR) JSON() string {
	b, err := json.Marshal(r)
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

	ctx := context.Background()
	list, err := prstats.GetPRList(ctx, &in, time.Unix(0, 0))
	if err != nil {
		return fmt.Errorf("list PR: %v", err)
	}

	format := strings.ToLower(c.String("format"))
	if err := print(format, list); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}

func print(format string, list []*github.PullRequest) error {
	if format == "json" {
		for _, r := range list {
			fmt.Println(r)
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("id, title, created_at, merged_at, lead_time(hours), ")

		for _, r := range list {
			fmt.Printf("%v, %v, %v, %v, ", *r.ID, strings.ReplaceAll(*r.Title, ",", ""), r.CreatedAt, r.MergedAt)
			if r.MergedAt != nil {
				fmt.Printf("%.4f, ", r.MergedAt.Sub(*r.CreatedAt).Hours())
			}

			fmt.Println()
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}
