package commits

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghstats/pkg/commits"
	"github.com/urfave/cli/v2"
)

const Filename = "commits.json"

func Fetch(c *cli.Context) error {
	dir := fmt.Sprintf("%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repo"))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	path := fmt.Sprintf("%s/%s", dir, Filename)

	sha, err := scanLastSHA(path)
	if err != nil {
		return fmt.Errorf("last id: %v", err)
	}

	in := commits.ListInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		Page:    c.Int("page"),
		PerPage: c.Int("perpage"),
		LastSHA: sha,
	}

	list, err := commits.Fetch(context.Background(), &in)
	if err != nil {
		return fmt.Errorf("fetch: %v", err)
	}

	if err := serialize(path, list); err != nil {
		return fmt.Errorf("serialize: %v", err)
	}

	sort.Slice(list, func(i, j int) bool { return list[i].Commit.Author.Date.Before(*list[j].Commit.Author.Date) })
	for _, r := range list {
		fmt.Printf("%v(%v)\n", *r.SHA, r.Commit.Author.Date)
	}

	return nil
}

func serialize(path string, list []*github.RepositoryCommit) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("open file: %v", err)
	}
	defer file.Close()

	for _, r := range list {
		fmt.Fprintln(file, JSON(r))
	}

	return nil
}

func scanLastSHA(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", nil
	}

	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open %v: %v", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var date time.Time
	var sha string
	for scanner.Scan() {
		var c github.RepositoryCommit
		if err := json.Unmarshal([]byte(scanner.Text()), &c); err != nil {
			return "", fmt.Errorf("unmarshal: %v", err)
		}

		if c.Commit.Author.Date.After(date) {
			sha = *c.SHA
			date = *c.Commit.Author.Date
		}
	}

	return sha, nil
}
