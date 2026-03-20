package routes

import (
	"github.com/basktre/api-gateway/internal/middlewares"
	"github.com/basktre/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
)

func PublicRoutes(engine *gin.Engine) {
	fwd := proxy.NewForwarder()
	api := engine.Group("/api/v1")
	{
		models := api.Group("/models", middlewares.APIKeyRequiredMiddleware())
		{
			models.POST("/call", fwd.Handler())
			models.POST("/stream", fwd.Handler())
			models.GET("", fwd.Handler())
			models.POST("/tokens/estimate", fwd.Handler())
		}

		providers := api.Group("/providers")
		{
			providers.GET("", fwd.Handler())
		}

		auth := api.Group("/auth")
		{
			auth.POST("/google-login", fwd.Handler())
		}

		workspaces := api.Group("/workspace", middlewares.GoogleAuthMiddleware())
		{
			workspaces.POST("", fwd.WorkspaceCreateHandler())
			workspaces.POST("/:id/api-key", fwd.Handler())
		}

		waitlist := api.Group("/waitlist")
		{
			waitlist.POST("", fwd.Handler())
		}
	}
}
