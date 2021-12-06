package commits

import (
	"context"
	"fmt"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

type ListInput struct {
	Owner   string
	Repo    string
	Number  int
	PAT     string
	Page    int
	PerPage int
}

func Fetch(ctx context.Context, in *ListInput) ([]*github.RepositoryCommit, error) {
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	opts := github.ListOptions{
		Page:    in.Page,
		PerPage: in.PerPage,
	}

	out := make([]*github.RepositoryCommit, 0)
	for {
		c, resp, err := client.PullRequests.ListCommits(ctx, in.Owner, in.Repo, in.Number, &opts)
		if err != nil {
			return nil, fmt.Errorf("list PullRequests: %v", err)
		}

		out = append(out, c...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return out, nil
}
