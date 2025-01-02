package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/karlhepler/disfunction/internal/github"
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

					ownerRepos := cmd.StringSlice("owner-repos")
					repos := make([]*github.Repository, len(ownerRepos))
					for _, ownerRepo := range ownerRepos {
						repo := &github.Repository{}

						parts := strings.Split(ownerRepo, "/")
						if ownerLogin := parts[0]; ownerLogin != "" {
							repo.Owner = &github.User{
								Login: &ownerLogin,
							}
						}
						if len(parts) > 1 {
							if repoName := parts[1]; repoName != "" {
								repo.Name = &repoName
							}
						}
					}

					req := handler.DisfunctionReq{
						Ctx:   ctx,
						Repos: repos,
						Since: cmd.Timestamp("since"),
						Until: cmd.Timestamp("until"),
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
