package middlewares

import (
	"net/http"
	"strings"

	"github.com/basktre/api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := config.GetStringSlice("server.security.cors.allowedOrigins")
		if len(allowedOrigins) == 0 {
			allowedOrigins = []string{"*"}
		}
		origin := c.Request.Header.Get("Origin")
		allowOrigin := "*"
		if origin != "" {
			for _, o := range allowedOrigins {
				if o == "*" || o == origin {
					allowOrigin = origin
					break
				}
			}
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func MethodRestrictionMiddleware() gin.HandlerFunc {
	allowedMethods := map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true, "HEAD": true}
	return func(c *gin.Context) {
		method := strings.ToUpper(c.Request.Method)
		if !allowedMethods[method] {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "HTTP method not allowed"})
			return
		}
		c.Next()
	}
}
