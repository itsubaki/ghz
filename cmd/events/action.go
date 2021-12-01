package events

import (
	"context"
	"fmt"

	"github.com/itsubaki/prstats/pkg/prstats"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	in := prstats.ListEventsInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		Page:    c.Int("page"),
		PerPage: c.Int("perpage"),
	}

	events, err := prstats.ListEvents(context.Background(), &in)
	if err != nil {
		return fmt.Errorf("get Events List: %v", err)
	}

	for _, e := range events {
		fmt.Printf("%v %v %v %v %v\n", *e.ID, *e.Actor, *e.Repo, *e.Type, string(*e.RawPayload))
	}

	return nil
}
