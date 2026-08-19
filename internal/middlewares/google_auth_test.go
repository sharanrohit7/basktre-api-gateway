package middlewares

import (
	"testing"

	"google.golang.org/api/idtoken"
)

func TestIdentityFromGooglePayload(t *testing.T) {
	payload := &idtoken.Payload{
		Issuer:  "https://accounts.google.com",
		Subject: "google-subject-1",
		Claims: map[string]interface{}{
			"email":          "user@example.com",
			"email_verified": true,
			"name":           "Test User",
		},
	}

	identity, appErr := identityFromGooglePayload(payload)
	if appErr.Code != "" {
		t.Fatalf("unexpected error: %v", appErr)
	}
	if identity.email != "user@example.com" || identity.sub != "google-subject-1" || identity.name != "Test User" {
		t.Fatalf("unexpected identity: %+v", identity)
	}
}

func TestIdentityFromGooglePayloadRejectsUnverifiedEmail(t *testing.T) {
	payload := &idtoken.Payload{
		Issuer:  "accounts.google.com",
		Subject: "google-subject-1",
		Claims: map[string]interface{}{
			"email":          "user@example.com",
			"email_verified": false,
		},
	}

	if _, appErr := identityFromGooglePayload(payload); appErr.Code == "" {
		t.Fatal("expected unverified email to be rejected")
	}
}

func TestIdentityFromGooglePayloadRejectsInvalidIssuer(t *testing.T) {
	payload := &idtoken.Payload{
		Issuer:  "https://attacker.example",
		Subject: "google-subject-1",
		Claims: map[string]interface{}{
			"email":          "user@example.com",
			"email_verified": true,
		},
	}

	if _, appErr := identityFromGooglePayload(payload); appErr.Code == "" {
		t.Fatal("expected invalid issuer to be rejected")
	}
}
