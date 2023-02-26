package commits

import (
	"context"
	"fmt"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type FetchInput struct {
	Owner      string
	Repository string
	PAT        string
	Page       int
	PerPage    int
	LastSHA    string
}

func Fetch(ctx context.Context, in *FetchInput, fn ...func(list []*github.RepositoryCommit) error) ([]*github.RepositoryCommit, error) {
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	opts := github.CommitsListOptions{
		ListOptions: github.ListOptions{
			Page:    in.Page,
			PerPage: in.PerPage,
		},
	}

	out := make([]*github.RepositoryCommit, 0)
	for {
		commits, resp, err := client.Repositories.ListCommits(ctx, in.Owner, in.Repository, &opts)
		if err != nil {
			return nil, fmt.Errorf("list commits: %v", err)
		}

		buf := make([]*github.RepositoryCommit, 0)
		var last bool
		for i := range commits {
			if in.LastSHA != "" && commits[i].GetSHA() == in.LastSHA {
				last = true
				break
			}

			buf = append(buf, commits[i])
		}

		for i, f := range fn {
			if err := f(buf); err != nil {
				return nil, fmt.Errorf("func[%v]: %v", i, err)
			}
		}

		out = append(out, buf...)
		if last || resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return out, nil
}
