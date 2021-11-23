package main

import (
	"fmt"
	"os"

	"github.com/itsubaki/prstats/cmd"
	"github.com/itsubaki/prstats/cmd/prlist"
	"github.com/urfave/cli/v2"
)

var date, hash, goversion string

func New(version string) *cli.App {
	app := cli.NewApp()

	app.Name = "prstats"
	app.Usage = "Github PR stats"
	app.Version = version
	app.Action = cmd.Action
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "pat",
			EnvVars: []string{"PAT"},
		},
		&cli.StringFlag{
			Name:    "owner",
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
			Usage:   "all, open, closed",
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
		&cli.IntFlag{
			Name:    "week",
			Aliases: []string{"w"},
			Value:   1,
		},
	}

	prlist := cli.Command{
		Name:    "list",
		Aliases: []string{"p"},
		Action:  prlist.Action,
		Usage:   "PR list",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "json",
				Usage:   "json, csv",
			},
		},
	}

	app.Commands = []*cli.Command{
		&prlist,
	}

	return app
}

func main() {
	v := fmt.Sprintf("%s %s %s", date, hash, goversion)
	if err := New(v).Run(os.Args); err != nil {
		panic(err)
	}
}
