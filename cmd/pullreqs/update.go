package pullreqs

import (
	"context"
	"fmt"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghstats/pkg/pullreqs"
	"github.com/urfave/cli/v2"
)

func Update(c *cli.Context) error {
	path := fmt.Sprintf("%v/%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repo"), Filename)
	list, err := Deserialize(path)
	if err != nil {
		return fmt.Errorf("deserialize: %v", err)
	}

	open := make([]github.PullRequest, 0)
	for i := range list {
		if *list[i].State != "open" {
			continue
		}

		open = append(open, list[i])
	}

	fmt.Println("id, number, title, state, created_at, updated_at, merged_at, closed_at, merge_commit_sha, ")
	ctx := context.Background()
	for i := range open {
		in := pullreqs.GetInput{
			Owner:      c.String("owner"),
			Repository: c.String("repo"),
			PAT:        c.String("pat"),
			Number:     *open[i].Number,
		}

		pr, err := pullreqs.Get(ctx, &in)
		if err != nil {
			return fmt.Errorf("get pullreq: %v", err)
		}

		fmt.Println(CSV(*pr))
	}

	return nil
}
