package prstats

import (
	"context"
	"fmt"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

type ListWorkflowRunsInput struct {
	Owner   string
	Repo    string
	PAT     string
	Page    int
	PerPage int
}

func ListWorkflowRuns(ctx context.Context, in *ListWorkflowRunsInput) ([]*github.WorkflowRun, error) {
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	opts := github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{
			Page:    in.Page,
			PerPage: in.PerPage,
		},
	}

	list := make([]*github.WorkflowRun, 0)
	for {
		runs, resp, err := client.Actions.ListRepositoryWorkflowRuns(ctx, in.Owner, in.Repo, &opts)
		if err != nil {
			return nil, fmt.Errorf("list WorkflowRuns: %v", err)
		}

		list = append(list, runs.WorkflowRuns...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return list, nil
}
