package main

import (
	"log"
	"os"

	abservice "aerobisoft.com/platform/cmd/service"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "aerobisoft",
		Description: "API microservice boilerplate",
		Usage:       "API microservice boilerplate",
		UsageText:   "aerobisoft [global options] command [command options] [arguments...]",
		Authors: []*cli.Author{
			{
				Name:  "New Author",
				Email: "info@aerobisoft.com",
			},
		},
		Copyright: "MIT License",
		Commands:  abservice.Commands,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
