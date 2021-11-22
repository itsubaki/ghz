package prstats

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

type ListPRInput struct {
	Owner   string
	Repo    string
	PAT     string
	State   string
	Week    int
	PerPage int
}

type PR struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	MergedAt  *time.Time `json:"merged_at"`
	ClosedAt  *time.Time `json:"closed_at"`
}

func (r PR) String() string {
	return r.JSON()
}

func (r PR) JSON() string {
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func ListPR(in *ListPRInput) ([]PR, error) {
	ctx := context.Background()
	client := github.NewClient(nil)

	if in.PAT != "" {
		client = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: in.PAT},
		)))
	}

	opt := github.PullRequestListOptions{
		State: in.State,
		ListOptions: github.ListOptions{
			PerPage: in.PerPage,
		},
	}

	out := make([]PR, 0)
	for {
		pr, resp, err := client.PullRequests.List(ctx, in.Owner, in.Repo, &opt)
		if err != nil {
			return out, fmt.Errorf("list PR: %v", err)
		}

		for _, r := range pr {
			out = append(out, PR{
				ID:        *r.ID,
				Title:     strings.ReplaceAll(*r.Title, ",", ""),
				CreatedAt: r.CreatedAt,
				UpdatedAt: r.UpdatedAt,
				MergedAt:  r.MergedAt,
				ClosedAt:  r.ClosedAt,
			})
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return out, nil
}
