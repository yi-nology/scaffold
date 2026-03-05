package resp

import (
	"scaffold/internal/pkg/errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *app.RequestContext, data interface{}) {
	c.JSON(consts.StatusOK, Response{
		Code:    errors.CodeSuccess,
		Message: errors.GetMessage(errors.CodeSuccess),
		Data:    data,
	})
}

func Error(c *app.RequestContext, code int, message string) {
	httpStatus := consts.StatusOK
	if code == errors.CodeUnauthorized {
		httpStatus = consts.StatusUnauthorized
	}
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

func BadRequest(c *app.RequestContext, message string) {
	if message == "" {
		message = errors.GetMessage(errors.CodeBadRequest)
	}
	Error(c, errors.CodeBadRequest, message)
}

func NotFound(c *app.RequestContext, message string) {
	if message == "" {
		message = errors.GetMessage(errors.CodeNotFound)
	}
	Error(c, errors.CodeNotFound, message)
}

func InternalError(c *app.RequestContext, message string) {
	if message == "" {
		message = errors.GetMessage(errors.CodeInternalError)
	}
	Error(c, errors.CodeInternalError, message)
}

func Unauthorized(c *app.RequestContext, message string) {
	if message == "" {
		message = errors.GetMessage(errors.CodeUnauthorized)
	}
	Error(c, errors.CodeUnauthorized, message)
}

func DBError(c *app.RequestContext, message string) {
	if message == "" {
		message = errors.GetMessage(errors.CodeDBError)
	}
	Error(c, errors.CodeDBError, message)
}
