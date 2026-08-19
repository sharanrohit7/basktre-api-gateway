package middlewares

import (
	"strings"

	brconfig "github.com/basktre/api-gateway/pkg/config"
	brerr "github.com/basktre/api-gateway/pkg/errors"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/idtoken"
)

const (
	HeaderUserEmail = "X-User-Email"
	HeaderUserSub   = "X-User-Sub"
	HeaderUserName  = "X-User-Name"
)

type googleIdentity struct {
	email string
	sub   string
	name  string
}

func GoogleAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		if token == "" || token == auth {
			brerr.RespondError(c, brerr.New(brerr.ErrUnauthorized, "missing or invalid bearer token"))
			c.Abort()
			return
		}

		clientID := brconfig.GetStringDefault("google.clientID", brconfig.GetString("GOOGLE_CLIENT_ID"))
		if clientID == "" {
			brerr.RespondError(c, brerr.New(brerr.ErrInternalServer, "gateway is not configured with GOOGLE_CLIENT_ID"))
			c.Abort()
			return
		}

		payload, err := idtoken.Validate(c.Request.Context(), token, clientID)
		if err != nil {
			brerr.RespondError(c, brerr.New(brerr.ErrUnauthorized, "invalid or expired google token"))
			c.Abort()
			return
		}

		identity, appErr := identityFromGooglePayload(payload)
		if appErr.Code != "" {
			brerr.RespondError(c, appErr)
			c.Abort()
			return
		}

		c.Set("forward_user_email", identity.email)
		c.Set("forward_user_sub", identity.sub)
		c.Set("forward_user_name", identity.name)
		c.Next()
	}
}

func identityFromGooglePayload(payload *idtoken.Payload) (googleIdentity, brerr.AppError) {
	if payload == nil {
		return googleIdentity{}, brerr.New(brerr.ErrUnauthorized, "invalid Google token payload")
	}
	if payload.Issuer != "accounts.google.com" && payload.Issuer != "https://accounts.google.com" {
		return googleIdentity{}, brerr.New(brerr.ErrUnauthorized, "invalid token issuer")
	}

	email, _ := payload.Claims["email"].(string)
	emailVerified, _ := payload.Claims["email_verified"].(bool)
	sub := payload.Subject
	if sub == "" {
		sub, _ = payload.Claims["sub"].(string)
	}
	name, _ := payload.Claims["name"].(string)
	if sub == "" {
		return googleIdentity{}, brerr.New(brerr.ErrUnauthorized, "invalid token subject")
	}
	if email == "" || !emailVerified {
		return googleIdentity{}, brerr.New(brerr.ErrUnauthorized, "google account email is missing or unverified")
	}

	return googleIdentity{email: email, sub: sub, name: name}, brerr.AppError{}
}
