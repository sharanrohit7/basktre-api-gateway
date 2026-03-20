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

		if payload.Issuer != "accounts.google.com" && payload.Issuer != "https://accounts.google.com" {
			brerr.RespondError(c, brerr.New(brerr.ErrUnauthorized, "invalid token issuer"))
			c.Abort()
			return
		}

		email, _ := payload.Claims["email"].(string)
		sub, _ := payload.Claims["sub"].(string)
		name, _ := payload.Claims["name"].(string)
		if sub == "" {
			brerr.RespondError(c, brerr.New(brerr.ErrUnauthorized, "invalid token subject"))
			c.Abort()
			return
		}

		c.Set("forward_user_email", email)
		c.Set("forward_user_sub", sub)
		c.Set("forward_user_name", name)
		c.Next()
	}
}
