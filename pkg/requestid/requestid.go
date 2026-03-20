package requestid

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type contextKey string

const requestIDKey contextKey = "request_id"
const HeaderXRequestID = "X-Request-ID"

func GenerateRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "fallback-request-id"
	}
	return hex.EncodeToString(b)
}

func SetRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

func RequestIDFromHTTP(r *http.Request) string {
	if id := r.Header.Get(HeaderXRequestID); id != "" {
		return id
	}
	return GenerateRequestID()
}
