package jobs

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
	JobID      int64
}

func Get(ctx context.Context, in *GetInput) (*github.WorkflowJob, error) {
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	out, _, err := client.Actions.GetWorkflowJobByID(ctx, in.Owner, in.Repository, in.JobID)
	if err != nil {
		return nil, fmt.Errorf("get pullreq: %v", err)
	}

	return out, nil
}
