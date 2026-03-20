package middlewares

import (
	"fmt"

	"github.com/basktre/api-gateway/pkg/logger"
	"github.com/basktre/api-gateway/pkg/newrelic"
	"github.com/basktre/api-gateway/pkg/notifications/slack"
	"github.com/basktre/api-gateway/pkg/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := requestid.RequestIDFromHTTP(c.Request)
		c.Set("request_id", id)
		c.Request = c.Request.WithContext(requestid.SetRequestID(c.Request.Context(), id))
		c.Writer.Header().Set(requestid.HeaderXRequestID, id)
		c.Next()
	}
}

func ErrorReporterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		status := c.Writer.Status()
		if status < 500 {
			return
		}
		requestID := c.GetString("request_id")
		path := c.FullPath()
		method := c.Request.Method
		errMsg := fmt.Sprintf("*Status:* %d\n*Method:* %s\n*Path:* %s\n*RequestID:* %s", status, method, path, requestID)
		logger.Error("5XX error",
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("request_id", requestID),
		)
		newrelic.NoticeErrorString(c.Request.Context(), fmt.Sprintf("5XX [%d] %s %s", status, method, path))
		go func() {
			_ = slack.Get().SendAlert(c.Request.Context(), fmt.Sprintf("5XX Error on api-gateway [%d]", status), errMsg)
		}()
	}
}
