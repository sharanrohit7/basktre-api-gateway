package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/basktre/api-gateway/pkg/config"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := config.GetStringSlice("server.security.cors.allowedOrigins")
		if configured := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS")); configured != "" {
			allowedOrigins = strings.Split(configured, ",")
		}

		origin := strings.TrimSpace(c.Request.Header.Get("Origin"))
		if origin == "" {
			c.Next()
			return
		}

		allowed := false
		for _, configuredOrigin := range allowedOrigins {
			configuredOrigin = strings.TrimSpace(configuredOrigin)
			if configuredOrigin == "*" || configuredOrigin == origin {
				allowed = true
				break
			}
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "origin is not allowed"})
			return
		}

		c.Writer.Header().Add("Vary", "Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-CSRF-Token, X-Requested-With, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Max-Age", "600")
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
