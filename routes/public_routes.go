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

		// Publicly accessible via Google Auth for dashboard integration
		dashboardModels := api.Group("/models", middlewares.GoogleAuthMiddleware())
		{
			dashboardModels.GET("/list", fwd.Handler())
		}

		providers := api.Group("/providers")
		{
			providers.GET("", fwd.Handler())
		}

		auth := api.Group("/auth", middlewares.GoogleAuthMiddleware())
		{
			auth.POST("/google-login", fwd.GoogleLoginHandler())
		}

		workspaces := api.Group("/workspace", middlewares.GoogleAuthMiddleware())
		{
			workspaces.POST("", fwd.WorkspaceCreateHandler())
			workspaces.GET("/user/:user_id", fwd.Handler())
			workspaces.PUT("/:id", fwd.Handler())
			workspaces.DELETE("/:id", fwd.Handler())

			workspaces.POST("/:id/api-key", fwd.Handler())
			workspaces.GET("/:id/api-key", fwd.Handler())
			workspaces.PUT("/:id/api-key/:key_id", fwd.Handler())
			workspaces.DELETE("/:id/api-key/:key_id", fwd.Handler())

			workspaces.PUT("/:id/byok/credentials/:provider", fwd.BYOKHandler())
			workspaces.GET("/:id/byok/credentials", fwd.BYOKHandler())
			workspaces.DELETE("/:id/byok/credentials/:provider", fwd.BYOKHandler())
		}

		wallet := api.Group("/wallet", middlewares.GoogleAuthMiddleware())
		{
			wallet.GET("/:workspace_id", fwd.Handler())
		}

		waitlist := api.Group("/waitlist")
		{
			waitlist.POST("", fwd.Handler())
		}

		dashboard := api.Group("/dashboard", middlewares.GoogleAuthMiddleware())
		{
			dashboard.GET("/workspace/:id/usage", fwd.Handler())
			dashboard.GET("/workspace/:id/usage/daily", fwd.Handler())
			dashboard.GET("/workspace/:id/transactions", fwd.Handler())
		}
	}
}
