package service

import (
	"github.com/urfave/cli/v2"
)

var environment string

var debugAddr string
var httpAddr string
var grpcAddr string

var zipkinAddr string

var Commands = cli.Commands{
	{
		Name:        "server",
		Aliases:     []string{"serve"},
		Description: "start the API boilerplate",
		Usage:       "start the API boilerplate server",
		UsageText:   "aerobisoft server [command options]",
		HelpName:    "server",
		Action:      Run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "debug-addr",
				Usage:       "Debug and metrics listen address",
				EnvVars:     []string{"DEBUG_ADDR"},
				Required:    true,
				DefaultText: ":8080",
				Destination: &debugAddr,
			},
			&cli.StringFlag{
				Name:        "environment",
				Usage:       "Environment application is running in",
				EnvVars:     []string{"ENVIRONMENT"},
				Required:    true,
				DefaultText: "development",
				Destination: &environment,
			},
			&cli.StringFlag{
				Name:        "grpc-addr",
				Usage:       "gRPC listen address",
				EnvVars:     []string{"GRPC_ADDR"},
				DefaultText: ":8082",
				Required:    true,
				Destination: &grpcAddr,
			},
			&cli.StringFlag{
				Name:        "http-addr",
				Usage:       "HTTP listen address",
				EnvVars:     []string{"HTTP_ADDR"},
				DefaultText: ":8081",
				Required:    true,
				Destination: &httpAddr,
			},
			&cli.StringFlag{
				Name:        "zipkin-addr",
				Usage:       "Enable Zipkin tracing",
				EnvVars:     []string{"ZIPKIN_ADDR"},
				Destination: &zipkinAddr,
			},
		},
	},
}
