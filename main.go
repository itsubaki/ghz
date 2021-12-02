package main

import (
	"fmt"
	"os"

	"github.com/itsubaki/prstats/cmd/analyze"
	"github.com/itsubaki/prstats/cmd/analyzejobs"
	"github.com/itsubaki/prstats/cmd/events"
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
			Aliases: []string{"t"},
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
			Name:  "page",
			Value: 0,
		},
		&cli.IntFlag{
			Name:  "perpage",
			Value: 1000,
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
		Flags: append(flags, []cli.Flag{
			&cli.StringFlag{
				Name:  "state",
				Value: "all",
				Usage: "all, open, closed",
			},
		}...),
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
		Flags: append(flags, []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Value: "out.json",
			},
			&cli.Int64Flag{
				Name:    "workflow_id",
				Aliases: []string{"wid"},
				Value:   -1,
			},
		}...),
	}

	events := cli.Command{
		Name:   "events",
		Action: events.Action,
		Usage:  "List Events",
		Flags:  flags,
	}

	runstats := cli.Command{
		Name:   "analyze",
		Action: analyze.Action,
		Usage:  "Analyze Productivity with Runslist",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Value:   "out.json",
			},
			&cli.StringFlag{
				Name:    "weeks",
				Aliases: []string{"w"},
				Value:   "52",
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "json",
				Usage:   "json, csv",
			},
			&cli.BoolFlag{
				Name:    "excluding_weekends",
				Aliases: []string{"ew"},
				Value:   false,
			},
		},
	}

	jobstats := cli.Command{
		Name:   "analyze-jobs",
		Action: analyzejobs.Action,
		Usage:  "Analyze Productivity with Jobslist",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Value:   "out_jobs.json",
			},
			&cli.StringFlag{
				Name:    "weeks",
				Aliases: []string{"w"},
				Value:   "52",
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "json",
				Usage:   "json, csv",
			},
			&cli.BoolFlag{
				Name:    "excluding_weekends",
				Aliases: []string{"ew"},
				Value:   false,
			},
		},
	}

	app.Commands = []*cli.Command{
		&prlist,
		&runslist,
		&jobslist,
		&events,
		&runstats,
		&jobstats,
	}

	return app
}

func main() {
	v := fmt.Sprintf("%s %s %s", date, hash, goversion)
	if err := New(v).Run(os.Args); err != nil {
		panic(err)
	}
}
