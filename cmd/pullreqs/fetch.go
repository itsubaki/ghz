package pullreqs

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/google/go-github/v56/github"
	"github.com/itsubaki/ghz/cmd/encode"
	"github.com/itsubaki/ghz/pullreqs"
	"github.com/urfave/cli/v2"
)

const Filename = "pullreqs.json"

func Fetch(c *cli.Context) error {
	dir := fmt.Sprintf("%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repository"))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	path := fmt.Sprintf("%s/%s", dir, Filename)
	id, number, err := GetLastID(path)
	if err != nil {
		return fmt.Errorf("last id: %v", err)
	}

	in := pullreqs.FetchInput{
		Owner:      c.String("owner"),
		Repository: c.String("repository"),
		PAT:        c.String("pat"),
		Page:       c.Int("page"),
		PerPage:    c.Int("perpage"),
		State:      c.String("state"),
		LastID:     id,
	}

	days := c.Int("days")
	if days > 0 {
		lastDay := time.Now().AddDate(0, 0, -1*days)
		in.LastDay = &lastDay
	}

	fmt.Printf("target: %v/%v\n", in.Owner, in.Repository)
	fmt.Printf("last_id: %v(%v)\n", id, number)

	ctx := context.Background()
	list, err := pullreqs.Fetch(ctx, &in)
	if err != nil {
		return fmt.Errorf("fetch: %v", err)
	}

	if err := Serialize(path, list); err != nil {
		return fmt.Errorf("serialize: %v", err)
	}

	sort.Slice(list, func(i, j int) bool { return *list[i].ID < *list[j].ID })
	for _, r := range list {
		fmt.Printf("%v(%v)\n", *r.ID, *r.Number)
	}

	return nil
}

func Serialize(path string, list []*github.PullRequest) error {
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

func GetLastID(path string) (int64, int, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return -1, -1, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return -1, -1, fmt.Errorf("open %v: %v", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var id int64
	var number int
	for scanner.Scan() {
		var pr github.PullRequest
		if err := json.Unmarshal([]byte(scanner.Text()), &pr); err != nil {
			return -1, -1, fmt.Errorf("unmarshal: %v", err)
		}

		if *pr.ID > id {
			id = *pr.ID
			number = *pr.Number
		}
	}

	return id, number, nil
}
