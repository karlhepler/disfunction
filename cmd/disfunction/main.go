package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/karlhepler/disfunction/internal/time"
	"github.com/karlhepler/disfunction/pkg/handler"
	"github.com/lithammer/dedent"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "disfunction",
		Usage: "What dis function?",
		Commands: []*cli.Command{
			{
				Name:  "random",
				Usage: "Display a random function chosen from all affiliated repos within a date range.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "owner",
						Usage: dedent.Dedent(`Filter repos by owner.
							If not defined, this will default to all repos GITHUB_TOKEN has access to.
						`),
					},
					&cli.StringSliceFlag{
						Name: "repos",
						Usage: dedent.Dedent(`Filter repos by this allow list.
							This must be used in conjunction with '--owner <value>'.
							If not defined, this will allow all repos through.
						`),
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
					token, ok := os.LookupEnv("GITHUB_TOKEN")
					if !ok {
						return fmt.Errorf("missing required env var GITHUB_TOKEN")
					}

					enableDebugMode := cmd.Bool("debug")
					console := NewConsoleLogger(enableDebugMode)
					hdl, err := handler.NewDisfunction(token, console)
					if err != nil {
						return fmt.Errorf("disfunction init failed: %w", err)
					}

					owner := cmd.String("owner")
					repos := cmd.StringSlice("repos")
					if len(repos) > 0 && owner == "" {
						return fmt.Errorf("owner is required when setting repos")
					}

					req := handler.DisfunctionReq{
						Ctx:       ctx,
						Owner:     owner,
						RepoNames: repos,
						Since:     cmd.Timestamp("since"),
						Until:     cmd.Timestamp("until"),
					}

					sender := NewConsoleSender()
					handler.Handle(hdl, req, sender) // enforce types before invoking the handler
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(fmt.Errorf("[FATAL] %w", err))
	}
}

func run()
