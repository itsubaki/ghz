package releases

import (
	"context"
	"fmt"

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

	rels, err := releases.Fetch(context.Background(), &in)
	if err != nil {
		return fmt.Errorf("fetch: %v", err)
	}

	for _, r := range rels {
		fmt.Printf("%v\n", r)
	}

	return nil
}
