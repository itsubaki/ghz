package commits

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghstats/pkg/commits"
	"github.com/urfave/cli/v2"
)

func Fetch(c *cli.Context) error {
	in := commits.ListInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		Page:    c.Int("page"),
		PerPage: c.Int("perpage"),
	}

	list, err := commits.Fetch(context.Background(), &in)
	if err != nil {
		return fmt.Errorf("fetch: %v", err)
	}

	format := strings.ToLower(c.String("format"))
	if err := print(format, list); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}

func print(format string, list []*github.RepositoryCommit) error {
	if format == "json" {
		for _, r := range list {
			fmt.Println(JSON(r))
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("sha, login, date, message")

		for _, c := range list {
			fmt.Println(CSV(c))
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}

func CSV(c *github.RepositoryCommit) string {
	title := strings.Split(*c.Commit.Message, "\n")[0]
	title = strings.ReplaceAll(title, ",", " ")

	return fmt.Sprintf(
		"%v, %v, %v, %v, ",
		*c.SHA,
		*c.Commit.Author.Name,
		c.Commit.Author.Date.Format("2006-01-02 15:04:05"),
		title,
	)
}

func JSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(b)
}
