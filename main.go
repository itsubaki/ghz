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
			Name:    "format",
			Aliases: []string{"f"},
			Value:   "json",
			Usage:   "json, csv",
		},
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
			Usage:   "all, open, closed",
		},
		&cli.StringFlag{
			Name:    "workflow",
			Aliases: []string{"w"},
			Usage:   "workflow name of deployment",
		},
		&cli.IntFlag{
			Name:    "perpage",
			Aliases: []string{"p"},
			Value:   100,
		},
	}

	prlist := cli.Command{
		Name:    "list",
		Aliases: []string{"p"},
		Action:  prlist.Action,
		Usage:   "PR list",
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
