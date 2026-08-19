package middlewares

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/basktre/api-gateway/pkg/logger"
	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx := c.Request.Context()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		logger.InfofWithContext(ctx, "[API Incoming] %s %s", c.Request.Method, path)
		if len(requestBody) > 0 && !strings.Contains(c.Request.URL.Path, "/byok/credentials/") {
			logger.InfofWithContext(ctx, "[API Request Payload] %s", string(requestBody))
		}

		w := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		latency := time.Since(start)
		logger.InfofWithContext(ctx, "[API Completed] Status: %d | Latency: %v", c.Writer.Status(), latency)
		if w.body.Len() > 0 {
			logger.InfofWithContext(ctx, "[API Response Payload] %s", w.body.String())
		}
	}
}
