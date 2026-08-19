package proxy

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	brconfig "github.com/basktre/api-gateway/pkg/config"
	brerr "github.com/basktre/api-gateway/pkg/errors"
	"github.com/basktre/api-gateway/pkg/requestid"
	"github.com/gin-gonic/gin"
)

const (
	gatewaySecretEnv    = "GATEWAY_INTERNAL_SECRET"
	gatewaySecretHeader = "X-Basktre-Gateway-Secret"
)

type Forwarder struct {
	client *http.Client
}

func NewForwarder() *Forwarder {
	return &Forwarder{client: &http.Client{}}
}

func (f *Forwarder) Handler() gin.HandlerFunc {
	return f.handler(false)
}

func (f *Forwarder) WorkspaceCreateHandler() gin.HandlerFunc {
	return f.handler(true)
}

// GoogleLoginHandler forwards only identities verified by GoogleAuthMiddleware
// and authenticates the gateway itself to the router.
func (f *Forwarder) GoogleLoginHandler() gin.HandlerFunc {
	return f.handler(true)
}

// BYOKHandler authenticates the gateway to the router for sensitive credential operations.
func (f *Forwarder) BYOKHandler() gin.HandlerFunc {
	return f.handler(true)
}

func (f *Forwarder) handler(authenticateGateway bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		baseURL := strings.TrimSuffix(brconfig.GetString("DOWNSTREAM_BASE_URL"), "/")
		if baseURL == "" {
			baseURL = strings.TrimSuffix(brconfig.GetString("downstream.baseURL"), "/")
		}
		if baseURL == "" {
			brerr.RespondError(c, brerr.New(brerr.ErrInternalServer, "missing DOWNSTREAM_BASE_URL"))
			return
		}

		targetURL := baseURL + c.Request.URL.Path
		if raw := c.Request.URL.RawQuery; raw != "" {
			targetURL += "?" + raw
		}

		req, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, targetURL, c.Request.Body)
		if err != nil {
			brerr.RespondError(c, brerr.Wrap(brerr.ErrInternalServer, "failed to build downstream request", err))
			return
		}

		req.Header = c.Request.Header.Clone()
		req.Header.Del("Authorization")
		if authenticateGateway {
			if err := attachGatewayAuthentication(req); err != nil {
				brerr.RespondError(c, brerr.Wrap(brerr.ErrInternalServer, "gateway authentication is not configured", err))
				return
			}
		}
		req.Header.Set(requestid.HeaderXRequestID, requestid.GetRequestID(c.Request.Context()))

		if v, ok := c.Get("forward_user_email"); ok {
			req.Header.Set("X-User-Email", v.(string))
		}
		if v, ok := c.Get("forward_user_sub"); ok {
			req.Header.Set("X-User-Sub", v.(string))
		}
		if v, ok := c.Get("forward_user_name"); ok {
			req.Header.Set("X-User-Name", v.(string))
		}

		resp, err := f.client.Do(req)
		if err != nil {
			brerr.RespondError(c, brerr.Wrap(brerr.ErrExternalAPI, "downstream request failed", err))
			return
		}
		defer resp.Body.Close()

		for k, values := range resp.Header {
			if isDownstreamCORSHeader(k) {
				continue
			}
			for _, v := range values {
				c.Writer.Header().Add(k, v)
			}
		}
		c.Status(resp.StatusCode)
		_, _ = io.Copy(c.Writer, resp.Body)
	}
}

func isDownstreamCORSHeader(name string) bool {
	return strings.HasPrefix(strings.ToLower(name), "access-control-")
}

func attachGatewayAuthentication(req *http.Request) error {
	secret := strings.TrimSpace(os.Getenv(gatewaySecretEnv))
	if secret == "" {
		return errors.New("GATEWAY_INTERNAL_SECRET is empty")
	}
	req.Header.Set(gatewaySecretHeader, secret)
	return nil
}
