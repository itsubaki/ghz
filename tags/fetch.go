package tags

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
	LastName   string
}

func Fetch(ctx context.Context, in *FetchInput, fn ...func(list []*github.RepositoryTag) error) ([]*github.RepositoryTag, error) {
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

	out := make([]*github.RepositoryTag, 0)
	for {
		tags, resp, err := client.Repositories.ListTags(ctx, in.Owner, in.Repository, &opts)
		if err != nil {
			return nil, fmt.Errorf("list tags: %v", err)
		}

		buf := make([]*github.RepositoryTag, 0)
		var last bool
		for i := range tags {
			if in.LastName != "" && tags[i].GetName() == in.LastName {
				last = true
				break
			}

			buf = append(buf, tags[i])
		}

		for i, f := range fn {
			if err := f(buf); err != nil {
				return nil, fmt.Errorf("func[%v]: %v", i, err)
			}
		}

		out = append(out, tags...)
		if last || resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return out, nil
}
