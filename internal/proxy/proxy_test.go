package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAttachGatewayAuthentication(t *testing.T) {
	t.Setenv(gatewaySecretEnv, "shared-secret")
	request, err := http.NewRequest(http.MethodGet, "http://router.test", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set(gatewaySecretHeader, "client-controlled-value")

	if err := attachGatewayAuthentication(request); err != nil {
		t.Fatalf("attachGatewayAuthentication returned error: %v", err)
	}
	if got := request.Header.Get(gatewaySecretHeader); got != "shared-secret" {
		t.Fatalf("gateway secret = %q, want shared-secret", got)
	}
}

func TestGoogleLoginHandlerForwardsOnlyVerifiedIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv(gatewaySecretEnv, "shared-secret")

	downstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get(gatewaySecretHeader); got != "shared-secret" {
			t.Errorf("gateway secret = %q, want shared-secret", got)
		}
		if got := r.Header.Get("X-User-Email"); got != "verified@example.com" {
			t.Errorf("forwarded email = %q, want verified@example.com", got)
		}
		if got := r.Header.Get("X-User-Sub"); got != "verified-google-subject" {
			t.Errorf("forwarded subject = %q, want verified-google-subject", got)
		}
		if got := r.Header.Get("Authorization"); got != "" {
			t.Errorf("authorization header leaked downstream: %q", got)
		}
		w.Header().Set("Access-Control-Allow-Origin", "https://downstream.invalid")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer downstream.Close()
	t.Setenv("DOWNSTREAM_BASE_URL", downstream.URL)

	engine := gin.New()
	forwarder := NewForwarder()
	engine.POST("/api/v1/auth/google-login", func(c *gin.Context) {
		c.Set("forward_user_email", "verified@example.com")
		c.Set("forward_user_sub", "verified-google-subject")
		c.Set("forward_user_name", "Verified User")
		c.Next()
	}, forwarder.GoogleLoginHandler())

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/google-login", nil)
	request.Header.Set("Authorization", "Bearer google-id-token")
	request.Header.Set(gatewaySecretHeader, "client-spoofed-secret")
	request.Header.Set("X-User-Email", "attacker@example.com")
	request.Header.Set("X-User-Sub", "attacker-subject")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d; body=%s", response.Code, http.StatusNoContent, response.Body.String())
	}
	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("downstream CORS header leaked through gateway: %q", got)
	}
}

func TestAttachGatewayAuthenticationRequiresConfiguration(t *testing.T) {
	t.Setenv(gatewaySecretEnv, "")
	request, err := http.NewRequest(http.MethodGet, "http://router.test", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := attachGatewayAuthentication(request); err == nil {
		t.Fatal("expected missing gateway secret error")
	}
}
