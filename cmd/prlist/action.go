package prlist

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v40/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

type ListPullRequestsInput struct {
	Owner   string
	Repo    string
	PAT     string
	Page    int
	PerPage int
	State   string
}

func ListPullRequests(ctx context.Context, in *ListPullRequestsInput) ([]*github.PullRequest, error) {
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	opts := github.PullRequestListOptions{
		State: in.State,
		ListOptions: github.ListOptions{
			Page:    in.Page,
			PerPage: in.PerPage,
		},
	}

	out := make([]*github.PullRequest, 0)
	for {
		pr, resp, err := client.PullRequests.List(ctx, in.Owner, in.Repo, &opts)
		if err != nil {
			return nil, fmt.Errorf("list PullRequests: %v", err)
		}

		out = append(out, pr...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return out, nil
}

func Action(c *cli.Context) error {
	in := ListPullRequestsInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		Page:    c.Int("page"),
		PerPage: c.Int("perpage"),
		State:   c.String("state"),
	}

	list, err := ListPullRequests(context.Background(), &in)
	if err != nil {
		return fmt.Errorf("get PullRequest List: %v", err)
	}

	format := strings.ToLower(c.String("format"))
	if err := print(format, list); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}

func print(format string, list []*github.PullRequest) error {
	if format == "json" {
		for _, r := range list {
			fmt.Println(JSON(r))
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("id, title, created_at, merged_at, duration(m), ")

		for _, r := range list {
			fmt.Println(CSV(r))
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}

func CSV(r *github.PullRequest) string {
	out := fmt.Sprintf(
		"%v, %v, %v, %v, ",
		*r.ID,
		strings.ReplaceAll(*r.Title, ",", ""),
		r.CreatedAt.Format("2006-01-02 15:04:05"),
		r.MergedAt.Format("2006-01-02 15:04:05"),
	)

	if r.MergedAt != nil {
		out = out + fmt.Sprintf("%.4f, ", r.MergedAt.Sub(*r.CreatedAt).Minutes())
	}

	return out
}

func JSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(b)
}
