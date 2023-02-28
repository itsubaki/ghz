package releases

import (
	"context"
	"fmt"
	"time"

	"github.com/itsubaki/ghz/releases"
	"github.com/urfave/cli/v2"
)

func Fetch(c *cli.Context) error {
	in := releases.FetchInput{
		Owner:      c.String("owner"),
		Repository: c.String("repository"),
		PAT:        c.String("pat"),
		Page:       c.Int("page"),
		PerPage:    c.Int("perpage"),
	}

	days := c.Int("days")
	if days > 0 {
		lastDay := time.Now().AddDate(0, 0, -1*days)
		in.LastDay = &lastDay
	}

	rels, err := releases.Fetch(context.Background(), &in)
	if err != nil {
		return fmt.Errorf("fetch: %v", err)
	}

	for _, r := range rels {
		fmt.Printf("%v\n", r)
	}

	return nil
}
