package tags

import (
	"context"
	"fmt"

	"github.com/itsubaki/ghz/tags"
	"github.com/urfave/cli/v2"
)

func Fetch(c *cli.Context) error {
	in := tags.FetchInput{
		Owner:      c.String("owner"),
		Repository: c.String("repository"),
		PAT:        c.String("pat"),
		Page:       c.Int("page"),
		PerPage:    c.Int("perpage"),
	}

	tags, err := tags.Fetch(context.Background(), &in)
	if err != nil {
		return fmt.Errorf("fetch: %v", err)
	}

	for _, t := range tags {
		fmt.Printf("%v, %v\n", t.GetName(), t.GetCommit().GetSHA())
	}

	return nil
}
