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
				Usage: "Display a random function from all repos within a date range.",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name: "repos",
						Aliases: []string{
							"r",
							"repo",
							"owner/repo",
							"owner/repos",
						},
						Usage: dedent.Dedent(`Filter repos by this allow list.
							If not defined, this will allow all accessible repos through.

							Examples:
								--repos foo/,bar/						# allow all repos owned by foo and bar
								--repos fizz,buzz 					# allow all repos named fizz and buzz from all owners
								--repos foo/fizz,bar/buzz 	# only allow owner/repos foo/fizz and bar/buzz
						`),
					},
					&cli.StringSliceFlag{
						Name: "files",
						Aliases: []string{
							"f",
							"file",
							"pattern",
							"patterns",
						},
						Usage: dedent.Dedent(`Filter commits by file pattern.
							If not defined, this will allow all commits.

							Examples:
								--files foo.go,bar.js
								--files '*.go,*.js'
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

					req := handler.DisfunctionReq{
						Ctx:          ctx,
						AllowedRepos: parseRepos(cmd.StringSlice("repos")),
						AllowedFiles: cmd.StringSlice("files"),
						Since:        cmd.Timestamp("since"),
						Until:        cmd.Timestamp("until"),
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
