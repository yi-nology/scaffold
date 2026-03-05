package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Auth(accessKey string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if accessKey == "" {
			c.Next(ctx)
			return
		}

		// Check header first
		authHeader := string(c.GetHeader("Authorization"))
		if authHeader != "" {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == accessKey {
				c.Next(ctx)
				return
			}
		}

		// Check query parameter
		queryKey := c.Query("access_key")
		if queryKey == accessKey {
			c.Next(ctx)
			return
		}

		// Check X-Access-Key header
		xKey := string(c.GetHeader("X-Access-Key"))
		if xKey == accessKey {
			c.Next(ctx)
			return
		}

		c.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    1002,
			"message": "unauthorized: invalid or missing access key",
		})
		c.Abort()
	}
}
