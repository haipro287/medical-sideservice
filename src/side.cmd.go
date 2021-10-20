package main

import (
	"github.com/sirupsen/logrus"
	"github.com/sonntuet1997/medical-chain-utils/common"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

const (
	serviceName = "gateway-service"
)

func main() {
	logger = logrus.New()
	if err := makeApp().Run(os.Args); err != nil {
		logger.WithField("err", err).Error("shutting down due to error")
		_ = os.Stderr.Sync()
		os.Exit(1)
	}
}

func makeApp() *cli.App {
	app := &cli.App{
		Name:                 serviceName,
		Version:              "v1.0.1",
		EnableBashCompletion: true,
		Compiled:             time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Son Nguyen",
				Email: "sonntuet1997@gmail.com",
			},
		},
		Copyright: "(c) 2021 SOTANEXT inc.",
		Action:    runMain,
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Usage:   "run server",
				Action:  runMain,
			},
			{
				Name:    "gen-mnemonic",
				Aliases: []string{"g"},
				Usage:   "generate new mnemonic phrase",
				Action:  genMnemonic,
			},
			//{
			//	Name:    "seed-data",
			//	Aliases: []string{"sd"},
			//	Usage:   "seed data",
			//	Action:  seedAuthData,
			//	Flags: []cli.Flag{
			//		&cli.BoolFlag{
			//			Name:    "clean",
			//			EnvVars: []string{"CLEAN_DB"},
			//			Usage:   "Clean DB before seeding",
			//		},
			//	},
			//},
			//{
			//	Name:    "clean",
			//	Aliases: []string{"c"},
			//	Usage:   "clean DB",
			//	Action:  clean,
			//},
		},
		Flags: append([]cli.Flag{
			&cli.StringFlag{
				Name:    "runtime-version",
				EnvVars: []string{"RUNTIME_VERSION"},
				Value:   "v1.0.0",
			},
			&cli.IntFlag{
				Name:    "port",
				Value:   8080,
				EnvVars: []string{"PORT"},
				Usage:   "The port for exposing the http endpoints for accessing the side service",
			},
			&cli.IntFlag{
				Name:    "grpc-port",
				Value:   7070,
				EnvVars: []string{"GRPC_PORT"},
				Usage:   "The port for exposing the gRPC endpoints for accessing the side service",
			},
			&cli.IntFlag{
				Name:    "pprof-port",
				Value:   6060,
				EnvVars: []string{"PPROF_PORT"},
				Usage:   "The port for exposing pprof endpoints",
			},
			&cli.StringFlag{
				Name:    "cosmos-endpoint",
				Value:   "localhost:9090",
				EnvVars: []string{"COSMOS_ENDPOINT"},
				Usage:   "Cosmos GRPC endpoint",
			},
			&cli.StringFlag{
				Name:    "mnemonic",
				Value:   "",
				EnvVars: []string{"MNEMONIC"},
				Usage:   "Mnemonic of 3rd service account",
			},
			&cli.StringFlag{
				Name:    "chain-id",
				Value:   "medichain",
				EnvVars: []string{"CHAIN_ID"},
				Usage:   "Cosmos blockchain id",
			},
			&cli.BoolFlag{
				Name:    "disable-tracing",
				EnvVars: []string{"DISABLE_TRACING"},
				Usage:   "disable-tracing",
			},
			&cli.BoolFlag{
				Name:    "disable-profiler",
				EnvVars: []string{"DISABLE_PROFILER"},
				Usage:   "disable-profiler",
			},
			&cli.BoolFlag{
				Name:    "disable-stats",
				EnvVars: []string{"DISABLE_STATS"},
				Usage:   "disable-stats",
			},
			&cli.BoolFlag{
				Name:    "allow-kill",
				EnvVars: []string{"ALLOW_KILL"},
				Usage:   "allow remote request to kill server",
			},
		},
			common.LoggerFlag...),
	}
	return app
}
