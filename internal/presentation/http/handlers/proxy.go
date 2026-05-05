package handlers

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func NewProxyHandler(targetURL string) (*ProxyHandler, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = nil

	// Use the modern Rewrite hook instead of the deprecated Director
	proxy.Rewrite = func(pr *httputil.ProxyRequest) {
		pr.SetURL(target)
		pr.Out.Header.Set("X-Forwarded-Host", pr.In.Host)

		// Inject user identity from context
		if userID, ok := pr.In.Context().Value("public_user_id").(string); ok {
			pr.Out.Header.Set("X-User-Id", userID)
		}
		if sessionID, ok := pr.In.Context().Value("session_id").(string); ok {
			pr.Out.Header.Set("X-Session-Id", sessionID)
		}
	}

	return &ProxyHandler{
		target: target,
		proxy:  proxy,
	}, nil
}

func (h *ProxyHandler) Proxy(c *gin.Context) {
	h.proxy.ServeHTTP(c.Writer, c.Request)
}
