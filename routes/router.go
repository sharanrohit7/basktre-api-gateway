package routes

import (
	"net/http"

	"github.com/basktre/api-gateway/internal/middlewares"
	"github.com/basktre/api-gateway/pkg/newrelic"
	"github.com/gin-gonic/gin"
	nrgin "github.com/newrelic/go-agent/v3/integrations/nrgin"
)

func Initialize(engine *gin.Engine) error {
	engine.Use(gin.Recovery())
	engine.Use(middlewares.RequestIDMiddleware())
	engine.Use(middlewares.CORSMiddleware())
	engine.Use(middlewares.MethodRestrictionMiddleware())
	engine.Use(middlewares.RateLimitMiddleware())
	engine.Use(middlewares.CSRFMiddleware())
	engine.Use(middlewares.RequestLoggerMiddleware())
	engine.Use(middlewares.ErrorReporterMiddleware())

	if newrelic.IsEnabled() {
		engine.Use(nrgin.Middleware(newrelic.GetApplication()))
	}

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	PublicRoutes(engine)
	InternalRoutes(engine)
	return nil
}
