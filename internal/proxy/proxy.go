package proxy

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	brconfig "github.com/basktre/api-gateway/pkg/config"
	brerr "github.com/basktre/api-gateway/pkg/errors"
	"github.com/basktre/api-gateway/pkg/requestid"
	"github.com/gin-gonic/gin"
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

func (f *Forwarder) handler(enrichWorkspaceCreate bool) gin.HandlerFunc {
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

		bodyReader, err := buildDownstreamBody(c, enrichWorkspaceCreate)
		if err != nil {
			brerr.RespondError(c, brerr.Wrap(brerr.ErrBadRequest, "invalid request body", err))
			return
		}

		req, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, targetURL, bodyReader)
		if err != nil {
			brerr.RespondError(c, brerr.Wrap(brerr.ErrInternalServer, "failed to build downstream request", err))
			return
		}

		req.Header = c.Request.Header.Clone()
		req.Header.Del("Authorization")
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
			for _, v := range values {
				c.Writer.Header().Add(k, v)
			}
		}
		c.Status(resp.StatusCode)
		_, _ = io.Copy(c.Writer, resp.Body)
	}
}

func buildDownstreamBody(c *gin.Context, enrichWorkspaceCreate bool) (io.Reader, error) {
	if !enrichWorkspaceCreate {
		return c.Request.Body, nil
	}

	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	defer c.Request.Body.Close()

	body := map[string]interface{}{}
	if len(bytes.TrimSpace(raw)) > 0 {
		if err := json.Unmarshal(raw, &body); err != nil {
			return nil, err
		}
	}

	if email, ok := c.Get("forward_user_email"); ok {
		if v, ok := email.(string); ok && v != "" {
			body["email"] = v
		}
	}
	if name, ok := c.Get("forward_user_name"); ok {
		if v, ok := name.(string); ok && v != "" {
			body["user_name"] = v
		}
	}
	if sub, ok := c.Get("forward_user_sub"); ok {
		if v, ok := sub.(string); ok && v != "" {
			body["google_id"] = v
		}
	}

	merged, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("Content-Length", strconv.Itoa(len(merged)))
	return bytes.NewReader(merged), nil
}
