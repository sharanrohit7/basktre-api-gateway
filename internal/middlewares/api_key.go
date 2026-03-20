package middlewares

import (
	"strings"

	brerr "github.com/basktre/api-gateway/pkg/errors"
	"github.com/gin-gonic/gin"
)

// APIKeyRequiredMiddleware ensures api-key header is present for routes
// that depend on downstream api-key authentication.
func APIKeyRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := strings.TrimSpace(c.GetHeader("api-key"))
		if apiKey == "" {
			brerr.RespondError(c, brerr.New(brerr.ErrUnauthorized, "missing api-key header"))
			c.Abort()
			return
		}
		c.Next()
	}
}
