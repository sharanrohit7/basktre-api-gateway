package middlewares

import (
	"net/http"

	"github.com/basktre/api-gateway/pkg/config"
	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware() gin.HandlerFunc {
	qps := float64(config.GetIntDefault("server.security.rateLimit.qps", 100))
	lmt := tollbooth.NewLimiter(qps, &limiter.ExpirableOptions{DefaultExpirationTTL: 0})
	lmt.SetMessage("Too many requests from this IP. Please try again later.")
	lmt.SetStatusCode(http.StatusTooManyRequests)
	return tollbooth_gin.LimitHandler(lmt)
}
