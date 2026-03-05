package middleware

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// SPAFallback serves static files from the given directory and falls back
// to index.html for SPA routing (Vue3/React).
func SPAFallback(staticDir string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Request.URI().Path())

		// Skip API and health routes
		if strings.HasPrefix(path, "/api/") || path == "/health" {
			c.Next(ctx)
			return
		}

		// Try to serve static file
		filePath := filepath.Join(staticDir, path)
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			c.File(filePath)
			c.Abort()
			return
		}

		// SPA fallback: serve index.html for non-file requests
		indexPath := filepath.Join(staticDir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.File(indexPath)
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

// Static serves static files from a directory for a given URL prefix.
func Static(urlPrefix, dir string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Request.URI().Path())
		if !strings.HasPrefix(path, urlPrefix) {
			c.Next(ctx)
			return
		}

		filePath := filepath.Join(dir, strings.TrimPrefix(path, urlPrefix))
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			c.File(filePath)
			c.Abort()
			return
		}

		c.SetStatusCode(consts.StatusNotFound)
		c.Abort()
	}
}
