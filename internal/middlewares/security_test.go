package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddlewareAllowsConfiguredPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://app.example.com, http://localhost:3000")
	engine := gin.New()
	engine.Use(CORSMiddleware())
	engine.POST("/resource", func(c *gin.Context) { c.Status(http.StatusOK) })

	request := httptest.NewRequest(http.MethodOptions, "/resource", nil)
	request.Header.Set("Origin", "https://app.example.com")
	request.Header.Set("Access-Control-Request-Method", http.MethodPost)
	request.Header.Set("Access-Control-Request-Headers", "authorization,content-type")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
		t.Fatalf("allow origin = %q", got)
	}
	if got := response.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("allow credentials = %q", got)
	}
}

func TestCORSMiddlewareRejectsUnconfiguredOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://app.example.com")
	engine := gin.New()
	engine.Use(CORSMiddleware())
	engine.GET("/resource", func(c *gin.Context) { c.Status(http.StatusOK) })

	request := httptest.NewRequest(http.MethodGet, "/resource", nil)
	request.Header.Set("Origin", "https://attacker.example")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusForbidden)
	}
	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("unexpected allow origin: %q", got)
	}
}

func TestCORSMiddlewareDoesNotAddBrowserHeadersWithoutOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://app.example.com")
	engine := gin.New()
	engine.Use(CORSMiddleware())
	engine.GET("/resource", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	response := httptest.NewRecorder()
	engine.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/resource", nil))

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("unexpected allow origin: %q", got)
	}
}
