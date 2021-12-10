package runs

import (
	"context"
	"fmt"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

type FetchInput struct {
	Owner      string
	Repository string
	PAT        string
	Page       int
	PerPage    int
	LastID     int64
}

func Fetch(ctx context.Context, in *FetchInput) ([]*github.WorkflowRun, error) {
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
		runs, resp, err := client.Actions.ListRepositoryWorkflowRuns(ctx, in.Owner, in.Repository, &opts)
		if err != nil {
			return nil, fmt.Errorf("list workflow runs: %v", err)
		}

		var last bool
		for i := range runs.WorkflowRuns {
			if *runs.WorkflowRuns[i].ID <= in.LastID {
				last = true
				break
			}

			list = append(list, runs.WorkflowRuns[i])
		}

		if last || resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return list, nil
}
