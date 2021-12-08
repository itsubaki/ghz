package pullreqs

import (
	"context"
	"fmt"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

type GetInput struct {
	Owner  string
	Repo   string
	PAT    string
	Number int
}

func Get(ctx context.Context, in *GetInput) (*github.PullRequest, error) {
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	pr, _, err := client.PullRequests.Get(ctx, in.Owner, in.Repo, in.Number)
	if err != nil {
		return nil, fmt.Errorf("get pull rquests: %v", err)
	}

	return pr, nil
}
