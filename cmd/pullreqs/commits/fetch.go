package commits

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghstats/cmd/pullreqs"
	"github.com/itsubaki/ghstats/pkg/pullreqs/commits"
	"github.com/urfave/cli/v2"
)

const Filename = "pullreqs_commits.json"

type CommitWithPRID struct {
	PullReqID     int64 `json:"pullreq_id,omitempty"`
	PullReqNumber int   `json:"pullreq_number,omitempty"`
	github.RepositoryCommit
}

func (c CommitWithPRID) CSV() string {
	title := strings.Split(*c.Commit.Message, "\n")[0]
	title = strings.ReplaceAll(title, ",", " ")

	return fmt.Sprintf("%v, %v, %v, %v, %v, %v, ",
		c.PullReqID,
		c.PullReqNumber,
		*c.SHA,
		*c.Commit.Author.Name,
		c.Commit.Author.Date.Format("2006-01-02 15:04:05"),
		title,
	)
}

func Fetch(c *cli.Context) error {
	dir := fmt.Sprintf("%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repository"))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	path := fmt.Sprintf("%s/%s", dir, Filename)

	lastID, lastNum, err := GetLastNumber(path)
	if err != nil {
		return fmt.Errorf("last id: %v", err)
	}

	in := commits.FetchInput{
		Owner:      c.String("owner"),
		Repository: c.String("repository"),
		PAT:        c.String("pat"),
		Page:       c.Int("page"),
		PerPage:    c.Int("perpage"),
	}

	fmt.Printf("target: %v/%v\n", in.Owner, in.Repository)
	fmt.Printf("last_id: %v(%v)\n", lastID, lastNum)

	prpath := fmt.Sprintf("%v/%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repository"), pullreqs.Filename)
	prs, err := pullreqs.Deserialize(prpath)
	if err != nil {
		return fmt.Errorf("deserialize: %v", err)
	}
	sort.Slice(prs, func(i, j int) bool { return *prs[i].ID < *prs[j].ID })

	ctx := context.Background()
	for i := range prs {
		if *prs[i].Number <= lastNum {
			continue
		}

		list, err := commits.Fetch(ctx, &in, *prs[i].Number)
		if err != nil {
			return fmt.Errorf("fetch: %v", err)
		}

		clist := make([]CommitWithPRID, 0)
		for j := range list {
			clist = append(clist, CommitWithPRID{
				PullReqID:        *prs[i].ID,
				PullReqNumber:    *prs[i].Number,
				RepositoryCommit: *list[j],
			})
		}

		if err := Serialize(path, clist); err != nil {
			return fmt.Errorf("serialize: %v", err)
		}

		if len(list) > 0 {
			fmt.Printf("%v(%v)\n", *prs[i].ID, *prs[i].Number)
		}
	}

	return nil
}

func Serialize(path string, list []CommitWithPRID) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("open file: %v", err)
	}
	defer file.Close()

	for _, j := range list {
		fmt.Fprintln(file, JSON(j))
	}

	return nil
}

func GetLastNumber(path string) (int64, int, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return -1, -1, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return -1, -1, fmt.Errorf("open %v: %v", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var number int
	var id int64
	for scanner.Scan() {
		var c CommitWithPRID
		if err := json.Unmarshal([]byte(scanner.Text()), &c); err != nil {
			return -1, -1, fmt.Errorf("unmarshal: %v", err)
		}

		if c.PullReqID > id {
			id = c.PullReqID
			number = c.PullReqNumber
		}
	}

	return id, number, nil
}

func JSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return string(b)
}
