package commits

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghz/cmd/encode"
	"github.com/urfave/cli/v2"
)

func List(c *cli.Context) error {
	path := fmt.Sprintf("%v/%v/%v/%v", c.String("dir"), c.String("owner"), c.String("repository"), Filename)
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

func Deserialize(path string) ([]github.RepositoryCommit, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %v", path)
	}

	read, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %v", path, err)
	}

	out := make([]github.RepositoryCommit, 0)
	for _, r := range strings.Split(string(read), "\n") {
		if len(r) < 1 {
			// skip empty line
			continue
		}

		var commit github.RepositoryCommit
		if err := json.Unmarshal([]byte(r), &commit); err != nil {
			return nil, fmt.Errorf("unmarshal: %v", err)
		}

		out = append(out, commit)
	}

	return out, nil
}

func CSV(c github.RepositoryCommit) string {
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

func print(format string, list []github.RepositoryCommit) error {
	sort.Slice(list, func(i, j int) bool { return list[i].Commit.Author.Date.After(*list[j].Commit.Author.Date) })

	if format == "json" {
		for _, r := range list {
			json, err := encode.JSON(r)
			if err != nil {
				return fmt.Errorf("encode: %v", err)
			}

			fmt.Println(json)
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
