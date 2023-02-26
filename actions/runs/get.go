package runs

import (
	"context"
	"fmt"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type GetInput struct {
	Owner      string
	Repository string
	PAT        string
	RunID      int64
}

func Get(ctx context.Context, in *GetInput) (*github.WorkflowRun, error) {
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	out, _, err := client.Actions.GetWorkflowRunByID(ctx, in.Owner, in.Repository, in.RunID)
	if err != nil {
		return nil, fmt.Errorf("get pullreq: %v", err)
	}

	return out, nil
}
