package commits

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghstats/cmd/pullreqs"
	"github.com/itsubaki/ghstats/pkg/pullreqs/commits"
	"github.com/urfave/cli/v2"
)

const Filename = "pullreqs_commits.json"

func Fetch(c *cli.Context) error {
	prpath := fmt.Sprintf("%v/%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repo"), pullreqs.Filename)
	prs, err := pullreqs.Deserialize(prpath)
	if err != nil {
		return fmt.Errorf("deserialize: %v", err)
	}

	lastNum, err := scanLastNumber(prpath)
	if err != nil {
		return fmt.Errorf("last id: %v", err)
	}

	in := commits.FetchInput{
		Owner:   c.String("owner"),
		Repo:    c.String("repo"),
		PAT:     c.String("pat"),
		Page:    c.Int("page"),
		PerPage: c.Int("perpage"),
	}

	fmt.Printf("target: %v/%v\n", in.Owner, in.Repo)
	fmt.Printf("last_number: %v\n", lastNum)

	dir := fmt.Sprintf("%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repo"))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	path := fmt.Sprintf("%s/%s", dir, Filename)

	ctx := context.Background()
	for i := range prs {
		if *prs[i].Number <= lastNum {
			continue
		}

		cmts, err := commits.Fetch(ctx, &in, *prs[i].Number)
		if err != nil {
			return fmt.Errorf("fetch: %v", err)
		}

		if err := serialize(path, cmts); err != nil {
			return fmt.Errorf("serialize: %v", err)
		}

		if len(cmts) > 0 {
			fmt.Printf("%v\n", *prs[i].Number)
		}
	}

	return nil
}

func serialize(path string, list []*github.RepositoryCommit) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("open file: %v", err)
	}
	defer file.Close()

	sort.Slice(list, func(i, j int) bool { return list[i].Commit.Author.Date.Unix() < list[i].Commit.Author.Date.Unix() }) // asc

	for _, j := range list {
		fmt.Fprintln(file, JSON(j))
	}

	return nil
}

func scanLastNumber(path string) (int, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return -1, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return -1, fmt.Errorf("open %v: %v", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var number int
	for scanner.Scan() {
		var pr github.PullRequest
		if err := json.Unmarshal([]byte(scanner.Text()), &pr); err != nil {
			return -1, fmt.Errorf("unmarshal: %v", err)
		}

		number = *pr.Number
	}

	return number, nil
}

func JSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(b)
}
