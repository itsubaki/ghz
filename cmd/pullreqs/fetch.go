package pullreqs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghstats/pkg/pullreqs"
	"github.com/itsubaki/ghstats/pkg/pullreqs/commits"
	"github.com/urfave/cli/v2"
)

type PullRequestWithCommits struct {
	PullRequest github.PullRequest         `json:"pull_request"`
	Commits     []*github.RepositoryCommit `json:"commits"`
}

func Fetch(c *cli.Context) error {
	in := pullreqs.ListInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		Page:    c.Int("page"),
		PerPage: c.Int("perpage"),
		State:   c.String("state"),
	}

	ctx := context.Background()
	prs, err := pullreqs.Fetch(ctx, &in)
	if err != nil {
		return fmt.Errorf("fetch: %v", err)
	}

	list := make([]PullRequestWithCommits, 0)
	for i, r := range prs {
		if r.MergedAt == nil {
			continue
		}

		cmts, err := commits.Fetch(ctx, &commits.ListInput{
			Owner:   c.String("owner"),
			Repo:    c.String("repo"),
			PAT:     c.String("pat"),
			Page:    c.Int("page"),
			PerPage: c.Int("perpage"),
			Number:  *r.Number,
		})
		if err != nil {
			return fmt.Errorf("fetch commits: %v", err)
		}

		list = append(list, PullRequestWithCommits{
			PullRequest: *prs[i],
			Commits:     cmts,
		})
	}

	format := strings.ToLower(c.String("format"))
	if err := print(format, list); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}

func print(format string, list []PullRequestWithCommits) error {
	if format == "json" {
		for _, r := range list {
			fmt.Println(JSON(r))
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("id, number, title, created_at, merged_at, commit.sha, commit.login, commit.date, duration(minutes), ")

		for _, r := range list {
			for _, c := range r.Commits {
				fmt.Printf(CSV(r.PullRequest))
				fmt.Printf("%v, %v, %v, ",
					*c.SHA,
					*c.Commit.Author.Name,
					c.Commit.Author.Date.Format("2006-01-02 15:04:05"),
				)
				fmt.Printf("%.4f, ", r.PullRequest.MergedAt.Sub(*c.Commit.Author.Date).Minutes())
				fmt.Println()
			}
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}

func CSV(r github.PullRequest) string {
	return fmt.Sprintf(
		"%v, %v, %v, %v, %v, ",
		*r.ID,
		*r.Number,
		strings.ReplaceAll(*r.Title, ",", ""),
		r.CreatedAt.Format("2006-01-02 15:04:05"),
		r.MergedAt.Format("2006-01-02 15:04:05"),
	)
}

func JSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(b)
}
