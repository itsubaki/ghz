package commits

import (
	"context"
	"fmt"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
)

type GetInput struct {
	Owner      string
	Repository string
	PAT        string
	Page       int
	PerPage    int
	SHA        string
}

func Get(ctx context.Context, in GetInput) (*github.RepositoryCommit, error) {
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	commit, _, err := client.Repositories.GetCommit(ctx, in.Owner, in.Repository, in.SHA, nil)
	if err != nil {
		return nil, fmt.Errorf("get commits: %v", err)
	}

	return commit, nil
}
