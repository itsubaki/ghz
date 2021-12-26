package releases

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

func Fetch(ctx context.Context, in *FetchInput, fn ...func(list []*github.RepositoryRelease) error) ([]*github.RepositoryRelease, error) {
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

	out := make([]*github.RepositoryRelease, 0)
	for {
		rels, resp, err := client.Repositories.ListReleases(ctx, in.Owner, in.Repository, &opts)
		if err != nil {
			return nil, fmt.Errorf("list pullreqs: %v", err)
		}

		buf := make([]*github.RepositoryRelease, 0)
		var last bool
		for i := range rels {
			if rels[i].GetID() <= in.LastID {
				last = true
				break
			}

			buf = append(buf, rels[i])
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
