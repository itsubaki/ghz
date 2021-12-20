package issues

import (
	"context"
	"fmt"

	"github.com/itsubaki/ghstats/pkg/issues"
	"github.com/urfave/cli/v2"
)

func Fetch(c *cli.Context) error {
	in := issues.FetchInput{
		Owner:      c.String("owner"),
		Repository: c.String("repo"),
		PAT:        c.String("pat"),
		Page:       c.Int("page"),
		PerPage:    c.Int("perpage"),
	}

	issues, err := issues.Fetch(context.Background(), &in)
	if err != nil {
		return fmt.Errorf("fetch: %v", err)
	}

	fmt.Println("id, number, state, title, created_at, closed_at, labels, ")
	for _, i := range issues {
		labels := make([]string, 0)
		for j := range i.Labels {
			labels = append(labels, i.Labels[j].GetName())
		}

		fmt.Printf(
			"%v, %v, %v, %v, %v, %v, %v, \n",
			i.GetID(),
			i.GetNumber(),
			i.GetState(),
			i.GetTitle(),
			i.GetCreatedAt().Format("2006-01-02 15:04:05"),
			i.GetClosedAt().Format("2006-01-02 15:04:05"),
			labels,
		)
	}

	return nil
}
