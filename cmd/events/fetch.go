package events

import (
	"context"
	"fmt"

	"github.com/itsubaki/ghstats/pkg/events"
	"github.com/urfave/cli/v2"
)

func Fetch(c *cli.Context) error {
	in := events.FetchInput{
		Owner:      c.String("owner"),
		Repository: c.String("repo"),
		PAT:        c.String("pat"),
		Page:       c.Int("page"),
		PerPage:    c.Int("perpage"),
	}

	events, err := events.Fetch(context.Background(), &in)
	if err != nil {
		return fmt.Errorf("fetch: %v", err)
	}

	fmt.Println("id, login, name, created_at, type, ")
	for _, e := range events {
		fmt.Printf("%v, %v, %v, %v, %v\n", *e.ID, *e.Actor.Login, *e.Repo.Name, e.CreatedAt.Format("2006-01-02 15:04:05"), *e.Type)
	}

	return nil
}
