package middlewares

import (
	"github.com/basktre/api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
)

func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.GetBool("server.security.csrf.enabled") {
			c.Next()
			return
		}
		if c.Request.Method == "GET" || c.Request.Method == "OPTIONS" || c.Request.Method == "HEAD" {
			c.Next()
			return
		}
		c.Next()
	}
}
