package middleware

import (
	"context"
	"fmt"
	"runtime"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.uber.org/zap"

	"scaffold/internal/pkg/logger"
)

func Recovery() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if r := recover(); r != nil {
				stack := make([]byte, 4096)
				n := runtime.Stack(stack, false)

				logger.Error("Panic recovered",
					zap.Any("error", r),
					zap.ByteString("stack", stack[:n]),
				)

				c.JSON(consts.StatusInternalServerError, map[string]interface{}{
					"code":    1006,
					"message": fmt.Sprintf("Internal Server Error: %v", r),
				})
				c.Abort()
			}
		}()
		c.Next(ctx)
	}
}
