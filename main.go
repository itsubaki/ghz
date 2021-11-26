package main

import (
	"fmt"
	"os"

	"github.com/itsubaki/prstats/cmd/analyze"
	"github.com/itsubaki/prstats/cmd/jobslist"
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
			Name:  "owner",
			Value: "google",
		},
		&cli.StringFlag{
			Name:  "repo",
			Value: "go-github",
		},
		&cli.IntFlag{
			Name:  "perpage",
			Value: 100,
		},
		&cli.StringFlag{
			Name:  "format",
			Value: "json",
			Usage: "json, csv",
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

	jobslist := cli.Command{
		Name:   "jobslist",
		Action: jobslist.Action,
		Usage:  "List WorkflowRun Jobs",
		Flags: append(flags, &cli.StringFlag{
			Name:  "path",
			Value: "out.json",
		}),
	}

	analyze := cli.Command{
		Name:   "analyze",
		Action: analyze.Action,
		Usage:  "Analyze Productivity",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Value: "out.json",
			},
			&cli.StringFlag{
				Name:  "weeks",
				Value: "52",
			},
			&cli.StringFlag{
				Name:  "format",
				Value: "json",
				Usage: "json, csv",
			},
		},
	}

	app.Commands = []*cli.Command{
		&prlist,
		&runslist,
		&jobslist,
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
