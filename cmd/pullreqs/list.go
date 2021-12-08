package pullreqs

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/google/go-github/v40/github"
	"github.com/urfave/cli/v2"
)

func List(c *cli.Context) error {
	path := fmt.Sprintf("%v/%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repo"), Filename)
	list, err := Deserialize(path)
	if err != nil {
		return fmt.Errorf("deserialize: %v", err)
	}

	format := strings.ToLower(c.String("format"))
	if err := print(format, list); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}

func Deserialize(path string) ([]github.PullRequest, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %v", path)
	}

	read, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %v", path, err)
	}

	out := make([]github.PullRequest, 0)
	for _, r := range strings.Split(string(read), "\n") {
		if len(r) < 1 {
			// skip empty line
			continue
		}

		var pr github.PullRequest
		if err := json.Unmarshal([]byte(r), &pr); err != nil {
			return nil, fmt.Errorf("unmarshal: %v", err)
		}

		out = append(out, pr)
	}

	return out, nil
}

func print(format string, list []github.PullRequest) error {
	sort.Slice(list, func(i, j int) bool { return *list[i].ID > *list[j].ID })

	if format == "json" {
		for _, r := range list {
			fmt.Println(JSON(r))
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("id, number, title, state, created_at, updated_at, merged_at, closed_at, merge_commit_sha, ")
		for _, r := range list {
			fmt.Println(CSV(r))
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}

func CSV(r github.PullRequest) string {
	out := fmt.Sprintf(
		"%v, %v, %v, %v, %v, ",
		*r.ID,
		*r.Number,
		strings.ReplaceAll(*r.Title, ",", ""),
		*r.State,
		r.CreatedAt.Format("2006-01-02 15:04:05"),
	)

	if r.UpdatedAt == nil {
		out = out + fmt.Sprintf("null, ")
	} else {
		out = out + fmt.Sprintf("%v, ", r.UpdatedAt.Format("2006-01-02 15:04:05"))
	}

	if r.MergedAt == nil {
		out = out + fmt.Sprintf("null, ")
	} else {
		out = out + fmt.Sprintf("%v, ", r.MergedAt.Format("2006-01-02 15:04:05"))
	}

	if r.ClosedAt == nil {
		out = out + fmt.Sprintf("null, ")
	} else {
		out = out + fmt.Sprintf("%v, ", r.ClosedAt.Format("2006-01-02 15:04:05"))
	}

	if r.MergeCommitSHA == nil {
		out = out + fmt.Sprintf("null, ")
	} else {
		out = out + fmt.Sprintf("%v, ", r.MergeCommitSHA)
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
