package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "random",
		Usage: "Display a random function and its author from all repositories.",
		Action: func(context.Context, *cli.Command) error {
			fmt.Println("Hello friend!")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
