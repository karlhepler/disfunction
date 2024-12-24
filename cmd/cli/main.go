package main

import (
	"context"
	"log"
	"os"

	"github.com/karlhepler/disfunction/internal/must"
	"github.com/karlhepler/disfunction/internal/time"
	"github.com/karlhepler/disfunction/pkg/disfunction"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "disfunction",
		Usage: "What dis function?",
		Commands: []*cli.Command{
			{
				Name: "random",
				Usage: `Display a random function
chosen from all repositories within a date range.`,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "owner",
						Usage:    "The owner of the repositories.",
						Required: true,
					},
					&cli.TimestampFlag{
						Name: "since",
						Usage: `Only commits after this date will be returned.
	This is a timestamp in ISO 8601 format: YYYY-MM-DDTHH:MM:SSZ.
	`,
						Value:       time.StartOfDay(time.Now()),
						DefaultText: "beginning of today",
					},
					&cli.TimestampFlag{
						Name: "until",
						Usage: `Only commits before this date will be returned.
	This is a timestamp in ISO 8601 format: YYYY-MM-DDTHH:MM:SSZ.
	`,
						Value:       time.Now(),
						DefaultText: "now",
					},
					&cli.StringSliceFlag{
						Name: "kind",
						Usage: `Filter kinds of functions by this list.
	supported kinds: new`,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					var output Output

					token := must.Env("GITHUB_TOKEN")
					hdl, err := disfunction.NewRandomHandler(token)
					if err != nil {
						hdl.HandleErr(err, output)
						os.Exit(1)
					}

					req := disfunction.RandomReq{
						Context: ctx,
						Owner:   cmd.String("owner"),
						Since:   cmd.Timestamp("since"),
						Until:   cmd.Timestamp("until"),
					}

					return hdl.Handle(req, output)
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
