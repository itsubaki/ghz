package main

import (
	"fmt"
	"os"

	"github.com/itsubaki/prstats/cmd/actions"
	"github.com/itsubaki/prstats/cmd/pr"
	"github.com/urfave/cli/v2"
)

var date, hash, goversion string

func New(version string) *cli.App {
	app := cli.NewApp()

	app.Name = "prstats"
	app.Usage = "Github PR stats"
	app.Version = version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name: "pat",
			EnvVars: []string{
				"PAT",
			},
		},
	}

	pr := cli.Command{
		Name:    "pr",
		Aliases: []string{"p"},
		Action:  pr.Action,
		Usage:   "PR stats",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "org",
				Aliases: []string{"o"},
				Value:   "github",
			},
			&cli.StringFlag{
				Name:    "repo",
				Aliases: []string{"r"},
				Value:   "docs",
			},
			&cli.StringFlag{
				Name:    "state",
				Aliases: []string{"s"},
				Value:   "all",
			},
			&cli.IntFlag{
				Name:    "perpage",
				Aliases: []string{"p"},
				Value:   30,
			},
		},
	}

	actions := cli.Command{
		Name:    "actions",
		Aliases: []string{"a"},
		Action:  actions.Action,
	}

	app.Commands = []*cli.Command{
		&pr,
		&actions,
	}

	return app
}

func main() {
	v := fmt.Sprintf("%s %s %s", date, hash, goversion)
	if err := New(v).Run(os.Args); err != nil {
		panic(err)
	}
}
