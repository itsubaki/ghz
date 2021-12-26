package release

import (
	"context"
	"fmt"

	"github.com/itsubaki/ghstats/pkg/release"
	"github.com/urfave/cli/v2"
)

func Fetch(c *cli.Context) error {
	in := release.FetchInput{
		Owner:      c.String("owner"),
		Repository: c.String("repository"),
		PAT:        c.String("pat"),
		Page:       c.Int("page"),
		PerPage:    c.Int("perpage"),
	}

	rel, err := release.Fetch(context.Background(), &in)
	if err != nil {
		return fmt.Errorf("fetch: %v", err)
	}

	for _, r := range rel {
		fmt.Printf("%v\n", r)
	}

	return nil
}
