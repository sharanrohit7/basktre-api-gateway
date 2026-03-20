package server

import (
	"context"
	"os"
	"strings"

	brconfig "github.com/basktre/api-gateway/pkg/config"
	"github.com/basktre/api-gateway/pkg/logger"
	"github.com/basktre/api-gateway/routes"
	"github.com/gin-gonic/gin"
)

func Initialize(_ context.Context) {
	env := brconfig.GetStringDefault("server.environment", "local")
	if env == "prod" || env == "stage" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	if err := routes.Initialize(engine); err != nil {
		logger.Fatalf("failed to initialize routes: %v", err)
	}

	port := brconfig.GetStringDefault("server.port", ":8081")
	if envPort := strings.TrimSpace(os.Getenv("PORT")); envPort != "" {
		if strings.HasPrefix(envPort, ":") {
			port = envPort
		} else {
			port = ":" + envPort
		}
	}
	logger.Infof("Starting api-gateway HTTP server on %s", port)

	if err := engine.Run(port); err != nil {
		logger.Fatalf("server exited with error: %v", err)
	}
}
