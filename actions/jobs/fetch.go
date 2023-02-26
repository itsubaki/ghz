package jobs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type FetchInput struct {
	Owner      string
	Repository string
	PAT        string
	Page       int
	PerPage    int
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

	out := make([]*github.WorkflowJob, 0)
	for {
		jobs, resp, err := client.Actions.ListWorkflowJobs(ctx, in.Owner, in.Repository, runID, &opts)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return out, nil
			}

			return nil, fmt.Errorf("list workflow jobs: %v", err)
		}

		out = append(out, jobs.Jobs...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return out, nil
}
