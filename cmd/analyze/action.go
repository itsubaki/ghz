package analyze

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v40/github"
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	path := c.String("path")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %v", path)
	}

	read, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %v", path, err)
	}

	runs := make([]github.WorkflowRun, 0)
	for _, r := range strings.Split(string(read), "\n") {
		if len(r) < 1 {
			continue
		}

		var run github.WorkflowRun
		if err := json.Unmarshal([]byte(r), &run); err != nil {
			return fmt.Errorf("unmarshal: %v", err)
		}

		runs = append(runs, run)
	}

	for _, r := range runs {
		fmt.Printf("%v\n", r.RunNumber)
	}

	return nil
}
