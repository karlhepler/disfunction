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
				Name:  "random",
				Usage: "Display a random function chosen from all affiliated repositories within a date range.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "owner",
						Usage: "Filter repositories by owner.",
					},
					&cli.TimestampFlag{
						Name:  "since",
						Usage: "Only commits after this date will be returned.",
						Config: cli.TimestampConfig{
							Layouts: []string{
								time.DateOnly, // 2006-01-02
								time.DateTime, // 2006-01-02 15:04:05
								time.Kitchen,  // 3:04PM
								time.RFC3339,  // 2006-01-02T15:04:05Z07:00
								time.TimeOnly, // 15:04:05
							},
						},
						Value:       time.StartOfDay(time.Now()),
						DefaultText: "beginning of today",
					},
					&cli.TimestampFlag{
						Name:  "until",
						Usage: "Only commits before this date will be returned.",
						Config: cli.TimestampConfig{
							Layouts: []string{
								time.DateOnly, // 2006-01-02
								time.DateTime, // 2006-01-02 15:04:05
								time.Kitchen,  // 3:04PM
								time.RFC3339,  // 2006-01-02T15:04:05Z07:00
								time.TimeOnly, // 15:04:05
							},
						},
						Value:       time.Now(),
						DefaultText: "now",
					},
					&cli.BoolFlag{
						Name:  "debug",
						Usage: "Output debug logging.",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					var output Output
					output.EnableDebug = cmd.Bool("debug")

					token := must.Env("GITHUB_TOKEN")
					hdl, err := disfunction.NewRandomHandler(token, output)
					if err != nil {
						hdl.HandleErr(err)
						os.Exit(1)
					}

					req := disfunction.RandomReq{
						Context: ctx,
						Owner:   cmd.String("owner"),
						Since:   cmd.Timestamp("since"),
						Until:   cmd.Timestamp("until"),
					}

					hdl.Handle(req, output)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
