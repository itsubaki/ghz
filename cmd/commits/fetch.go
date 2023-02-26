package commits

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/google/go-github/v50/github"
	"github.com/itsubaki/ghz/cmd/encode"
	"github.com/itsubaki/ghz/commits"
	"github.com/urfave/cli/v2"
)

const Filename = "commits.json"

func Fetch(c *cli.Context) error {
	dir := fmt.Sprintf("%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repository"))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	path := fmt.Sprintf("%s/%s", dir, Filename)

	sha, err := GetLastSHA(path)
	if err != nil {
		return fmt.Errorf("last id: %v", err)
	}

	in := commits.FetchInput{
		Owner:      c.String("owner"),
		Repository: c.String("repository"),
		PAT:        c.String("pat"),
		Page:       c.Int("page"),
		PerPage:    c.Int("perpage"),
		LastSHA:    sha,
	}

	list, err := commits.Fetch(context.Background(), &in)
	if err != nil {
		return fmt.Errorf("fetch: %v", err)
	}

	if err := Serialize(path, list); err != nil {
		return fmt.Errorf("serialize: %v", err)
	}

	sort.Slice(list, func(i, j int) bool { return list[i].Commit.Author.Date.Before(list[j].Commit.Author.Date.Time) })
	for _, r := range list {
		fmt.Printf("%v(%v)\n", *r.SHA, r.Commit.Author.Date)
	}

	return nil
}

func Serialize(path string, list []*github.RepositoryCommit) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("open file: %v", err)
	}
	defer file.Close()

	for _, r := range list {
		json, err := encode.JSON(r)
		if err != nil {
			return fmt.Errorf("encode: %v", err)
		}

		fmt.Fprintln(file, json)
	}

	return nil
}

func GetLastSHA(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", nil
	}

	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open %v: %v", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var sha string
	var date time.Time
	for scanner.Scan() {
		var c github.RepositoryCommit
		if err := json.Unmarshal([]byte(scanner.Text()), &c); err != nil {
			return "", fmt.Errorf("unmarshal: %v", err)
		}

		if c.Commit.Author.Date.After(date) {
			sha = *c.SHA
			date = c.Commit.Author.Date.Time
		}
	}

	return sha, nil
}
