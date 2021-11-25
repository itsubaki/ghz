package main

import (
	"fmt"
	"os"

	"github.com/itsubaki/prstats/cmd/analyze"
	"github.com/itsubaki/prstats/cmd/prlist"
	"github.com/itsubaki/prstats/cmd/runslist"
	"github.com/urfave/cli/v2"
)

var date, hash, goversion string

func New(version string) *cli.App {
	app := cli.NewApp()

	app.Name = "prstats"
	app.Usage = "Github Productivity Stats"
	app.Version = version

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "pat",
			EnvVars: []string{"PAT"},
			Usage:   "Personal Access Token",
		},
		&cli.StringFlag{
			Name:    "owner",
			Aliases: []string{"o"},
			Value:   "google",
		},
		&cli.StringFlag{
			Name:    "repo",
			Aliases: []string{"r"},
			Value:   "go-github",
		},
		&cli.IntFlag{
			Name:    "perpage",
			Aliases: []string{"p"},
			Value:   30,
		},
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Value:   "json",
			Usage:   "json, csv",
		},
	}

	prlist := cli.Command{
		Name:   "prlist",
		Action: prlist.Action,
		Usage:  "List PullRequests",
		Flags:  flags,
	}

	runslist := cli.Command{
		Name:   "runslist",
		Action: runslist.Action,
		Usage:  "List WorkflowRuns",
		Flags:  flags,
	}

	analyze := cli.Command{
		Name:   "analyze",
		Action: analyze.Action,
		Usage:  "Analyze Productivity",
		Flags:  flags,
	}

	app.Commands = []*cli.Command{
		&prlist,
		&runslist,
		&analyze,
	}

	return app
}

func main() {
	v := fmt.Sprintf("%s %s %s", date, hash, goversion)
	if err := New(v).Run(os.Args); err != nil {
		panic(err)
	}
}
