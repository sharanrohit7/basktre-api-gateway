package main

import (
	"context"
	"flag"

	appinit "github.com/basktre/api-gateway/init"
	brconfig "github.com/basktre/api-gateway/pkg/config"
	"github.com/basktre/api-gateway/pkg/logger"
	"github.com/basktre/api-gateway/server"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./config/config.yaml", "Path to config file")
	flag.Parse()

	brconfig.MustInit(configPath)

	ctx := context.Background()
	appinit.Initialize(ctx)

	logger.Infof("Starting in HTTP mode")
	server.Initialize(ctx)
}
