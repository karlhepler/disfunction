package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "devex",
		Usage: "The DevEx Toolkit",
		Commands: []*cli.Command{
			{
				Name:  "git/hooks",
				Usage: "Manage //:.git/hooks/*",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					// interface is list, get, <//:.git/hooks/[n], ...> (required)
					// 	- interface is
					// var output Output
					// output.EnableDebug = cmd.Bool("debug")

					// token := must.Env("GITHUB_TOKEN")
					// hdl, err := disfunction.NewRandomHandler(token, output)
					// if err != nil {
					// 	hdl.HandleErr(err)
					// 	os.Exit(1)
					// }

					// req := disfunction.RandomReq{
					// 	Context: ctx,
					// 	Owner:   cmd.String("owner"),
					// 	Since:   cmd.Timestamp("since"),
					// 	Until:   cmd.Timestamp("until"),
					// }

					// hdl.Handle(req, output)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
