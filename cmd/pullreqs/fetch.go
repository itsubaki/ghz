package pullreqs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghstats/pkg/pullreqs"
	"github.com/urfave/cli/v2"
)

func Fetch(c *cli.Context) error {
	in := pullreqs.ListInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		Page:    c.Int("page"),
		PerPage: c.Int("perpage"),
		State:   c.String("state"),
	}

	list, err := pullreqs.Fetch(context.Background(), &in)
	if err != nil {
		return fmt.Errorf("fetch: %v", err)
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
			fmt.Println(JSON(r))
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("id, number, title, created_at, merged_at, duration(minutes), ")

		for _, r := range list {
			fmt.Println(CSV(r))
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}

func CSV(r *github.PullRequest) string {
	out := fmt.Sprintf(
		"%v, %v, %v, %v, ",
		*r.ID,
		*r.Number,
		strings.ReplaceAll(*r.Title, ",", ""),
		r.CreatedAt.Format("2006-01-02 15:04:05"),
	)

	if r.MergedAt != nil {
		out = out + fmt.Sprintf("%v, ", r.MergedAt.Format("2006-01-02 15:04:05"))
		out = out + fmt.Sprintf("%.4f, ", r.MergedAt.Sub(*r.CreatedAt).Minutes())
	}

	return out
}

func JSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(b)
}