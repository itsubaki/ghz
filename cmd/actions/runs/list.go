package runs

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

	sort.Slice(list, func(i, j int) bool { return *list[i].ID > *list[j].ID }) // desc

	format := strings.ToLower(c.String("format"))
	if err := print(format, list); err != nil {
		return fmt.Errorf("print: %v", err)
	}

	return nil
}

func Deserialize(path string) ([]github.WorkflowRun, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %v", path)
	}

	read, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %v", path, err)
	}

	runs := make([]github.WorkflowRun, 0)
	for _, r := range strings.Split(string(read), "\n") {
		if len(r) < 1 {
			// skip empty line
			continue
		}

		var run github.WorkflowRun
		if err := json.Unmarshal([]byte(r), &run); err != nil {
			return nil, fmt.Errorf("unmarshal: %v", err)
		}

		runs = append(runs, run)
	}

	return runs, nil
}

func print(format string, list []github.WorkflowRun) error {
	if format == "json" {
		for _, r := range list {
			fmt.Println(JSON(r))
		}

		return nil
	}

	if format == "csv" {
		fmt.Println("workflow_id, workflow_name, run_id, run_number, status, conclusion, created_at, updated_at, head_commit.sha, head_commit.message")

		for _, r := range list {
			fmt.Println(CSV(r))
		}

		return nil
	}

	return fmt.Errorf("invalid format=%v", format)
}

func CSV(r github.WorkflowRun) string {
	title := strings.Split(strings.ReplaceAll(*r.HeadCommit.Message, ",", " "), "\n")[0]

	return fmt.Sprintf(
		"%v, %v, %v, %v, %v, %v, %v, %v, %v, %v, ",
		*r.WorkflowID,
		*r.Name,
		*r.ID,
		*r.RunNumber,
		*r.Status,
		*r.Conclusion,
		r.CreatedAt.Format("2006-01-02 15:04:05"),
		r.UpdatedAt.Format("2006-01-02 15:04:05"),
		*r.HeadCommit.ID,
		title,
	)
}
