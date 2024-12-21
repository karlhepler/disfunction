package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
					&cli.TimestampFlag{
						Name: "since",
						Usage: `Only commits after this date will be returned.
	This is a timestamp in ISO 8601 format: YYYY-MM-DDTHH:MM:SSZ.
	`,
						Value:       GetStartOfDay(time.Now()),
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
				Action: func(context.Context, *cli.Command) error {
					fmt.Println("Hello friend!")
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func GetStartOfDay(t time.Time) time.Time {
	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		0, 0, 0, 0, // hour, minute, second, nanosecond
		t.Location(),
	)
}
