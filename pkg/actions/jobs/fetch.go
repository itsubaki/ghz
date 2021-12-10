package jobs

import (
	"context"
	"fmt"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

type FetchInput struct {
	Owner   string
	Repo    string
	PAT     string
	Page    int
	PerPage int
}

func Fetch(ctx context.Context, in *FetchInput, runID int64) ([]*github.WorkflowJob, error) {
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	opts := github.ListWorkflowJobsOptions{
		ListOptions: github.ListOptions{
			Page:    in.Page,
			PerPage: in.PerPage,
		},
	}

	list := make([]*github.WorkflowJob, 0)
	for {
		jobs, resp, err := client.Actions.ListWorkflowJobs(ctx, in.Owner, in.Repo, runID, &opts)
		if err != nil {
			return nil, fmt.Errorf("list workflow jobs: %v", err)
		}

		list = append(list, jobs.Jobs...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return list, nil
}
