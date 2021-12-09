package pullreqs

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghstats/pkg/pullreqs"
	"github.com/urfave/cli/v2"
)

const Filename = "pullreqs.json"

func Fetch(c *cli.Context) error {
	dir := fmt.Sprintf("%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repo"))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	path := fmt.Sprintf("%s/%s", dir, Filename)
	id, number, err := GetLastID(path)
	if err != nil {
		return fmt.Errorf("last id: %v", err)
	}

	in := pullreqs.ListInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		Page:    c.Int("page"),
		PerPage: c.Int("perpage"),
		State:   c.String("state"),
		LastID:  id,
	}

	fmt.Printf("target: %v/%v\n", in.Owner, in.Repo)
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
		fmt.Fprintln(file, JSON(r))
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
