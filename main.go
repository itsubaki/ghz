package main

import (
	"fmt"
	"os"

	"github.com/itsubaki/ghstats/cmd/actions/jobs"
	"github.com/itsubaki/ghstats/cmd/actions/runs"
	"github.com/itsubaki/ghstats/cmd/commits"
	"github.com/itsubaki/ghstats/cmd/events"
	"github.com/itsubaki/ghstats/cmd/issues"
	"github.com/itsubaki/ghstats/cmd/pullreqs"
	prcommits "github.com/itsubaki/ghstats/cmd/pullreqs/commits"
	"github.com/urfave/cli/v2"
)

var date, hash, goversion string

func New(version string) *cli.App {
	app := cli.NewApp()
	app.Name = "ghstats"
	app.Usage = "Github Stats"
	app.Version = version

	dir := cli.StringFlag{
		Name:    "dir",
		Aliases: []string{"d"},
		Value:   fmt.Sprintf("/var/tmp/%v", app.Name),
	}

	own := cli.StringFlag{
		Name:    "owner",
		Aliases: []string{"o"},
		Value:   "itsubaki",
	}

	repo := cli.StringFlag{
		Name:    "repository",
		Aliases: []string{"r"},
		Value:   "q",
	}

	format := cli.StringFlag{
		Name:    "format",
		Aliases: []string{"f"},
		Value:   "json",
		Usage:   "json, csv",
	}

	pat := cli.StringFlag{
		Name:    "pat",
		Aliases: []string{"t"},
		EnvVars: []string{"PAT"},
		Usage:   "Personal Access Token",
	}

	page := cli.IntFlag{
		Name:  "page",
		Value: 0,
	}

	perpage := cli.IntFlag{
		Name:  "perpage",
		Value: 100,
	}

	runs := cli.Command{
		Name:    "runs",
		Aliases: []string{"r"},
		Subcommands: []*cli.Command{
			{
				Name:    "fetch",
				Aliases: []string{"f"},
				Action:  runs.Fetch,
				Flags: []cli.Flag{
					&dir,
					&own,
					&repo,
					&pat,
					&page,
					&perpage,
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Action:  runs.List,
				Flags: []cli.Flag{
					&dir,
					&own,
					&repo,
					&format,
				},
			},
		},
	}

	jobs := cli.Command{
		Name:    "jobs",
		Aliases: []string{"j"},
		Subcommands: []*cli.Command{
			{
				Name:    "fetch",
				Aliases: []string{"f"},
				Action:  jobs.Fetch,
				Flags: []cli.Flag{
					&dir,
					&own,
					&repo,
					&pat,
					&page,
					&perpage,
					&cli.Int64Flag{
						Name:    "workflow_id",
						Aliases: []string{"wid"},
						Value:   -1,
					},
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Action:  jobs.List,
				Flags: []cli.Flag{
					&dir,
					&own,
					&repo,
					&format,
				},
			},
		},
	}

	actions := cli.Command{
		Name:    "actions",
		Aliases: []string{"a"},
		Subcommands: []*cli.Command{
			&jobs,
			&runs,
		},
	}

	pullreqs := cli.Command{
		Name:    "pullreqs",
		Aliases: []string{"pr"},
		Subcommands: []*cli.Command{
			{
				Name:    "fetch",
				Aliases: []string{"f"},
				Action:  pullreqs.Fetch,
				Flags: []cli.Flag{
					&dir,
					&own,
					&repo,
					&pat,
					&page,
					&perpage,
					&cli.StringFlag{
						Name:    "state",
						Aliases: []string{"s"},
						Value:   "all",
						Usage:   "all, open, closed",
					},
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Action:  pullreqs.List,
				Flags: []cli.Flag{
					&dir,
					&own,
					&repo,
					&format,
				},
			},
			{
				Name:    "commits",
				Aliases: []string{"c"},
				Subcommands: []*cli.Command{
					{
						Name:    "fetch",
						Aliases: []string{"f"},
						Action:  prcommits.Fetch,
						Flags: []cli.Flag{
							&dir,
							&own,
							&repo,
							&pat,
							&page,
							&perpage,
						},
					},
					{
						Name:    "list",
						Aliases: []string{"l"},
						Action:  prcommits.List,
						Flags: []cli.Flag{
							&dir,
							&own,
							&repo,
							&format,
						},
					},
				},
			},
			{
				Name:    "update",
				Aliases: []string{"u"},
				Action:  pullreqs.Update,
				Flags: []cli.Flag{
					&dir,
					&own,
					&repo,
					&pat,
					&format,
				},
			},
		},
	}

	commits := cli.Command{
		Name:    "commits",
		Aliases: []string{"c"},
		Subcommands: []*cli.Command{
			{
				Name:    "fetch",
				Aliases: []string{"f"},
				Action:  commits.Fetch,
				Flags: []cli.Flag{
					&dir,
					&own,
					&repo,
					&pat,
					&page,
					&perpage,
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Action:  commits.List,
				Flags: []cli.Flag{
					&dir,
					&own,
					&repo,
					&format,
				},
			},
		},
	}

	issues := cli.Command{
		Name:    "issues",
		Aliases: []string{"i"},
		Subcommands: []*cli.Command{
			{
				Name:    "fetch",
				Aliases: []string{"f"},
				Action:  issues.Fetch,
				Flags: []cli.Flag{
					&dir,
					&own,
					&repo,
					&pat,
					&page,
					&perpage,
				},
			},
		},
	}

	events := cli.Command{
		Name:    "events",
		Aliases: []string{"e"},
		Subcommands: []*cli.Command{
			{
				Name:    "fetch",
				Aliases: []string{"f"},
				Action:  events.Fetch,
				Flags: []cli.Flag{
					&dir,
					&own,
					&repo,
					&pat,
					&page,
					&perpage,
					&format,
				},
			},
		},
	}

	app.Commands = []*cli.Command{
		&actions,
		&pullreqs,
		&commits,
		&events,
		&issues,
	}

	return app
}

func main() {
	v := fmt.Sprintf("%s %s %s", date, hash, goversion)
	if err := New(v).Run(os.Args); err != nil {
		panic(err)
	}
}
