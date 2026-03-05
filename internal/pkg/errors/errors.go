package errors

const (
	CodeSuccess            = 0
	CodeBadRequest         = 1001
	CodeUnauthorized       = 1002
	CodeForbidden          = 1003
	CodeNotFound           = 1004
	CodeConflict           = 1005
	CodeInternalError      = 1006
	CodeServiceUnavailable = 1007
	CodeDBError            = 1008

	// Scaffold business error codes
	CodeTemplateNotFound    = 2001
	CodeTemplateCloneFailed = 2002
	CodeGenerateFailed      = 2003
	CodeTaskNotFound        = 2004
)

var codeMessages = map[int]string{
	CodeSuccess:             "success",
	CodeBadRequest:          "bad request",
	CodeUnauthorized:        "unauthorized",
	CodeForbidden:           "forbidden",
	CodeNotFound:            "not found",
	CodeConflict:            "conflict",
	CodeInternalError:       "internal server error",
	CodeServiceUnavailable:  "service unavailable",
	CodeDBError:             "database error",
	CodeTemplateNotFound:    "template not found",
	CodeTemplateCloneFailed: "template clone failed",
	CodeGenerateFailed:      "generate failed",
	CodeTaskNotFound:        "task not found",
}

type AppError struct {
	Code    int
	Message string
	Cause   error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(code int, message string, cause error) *AppError {
	return &AppError{Code: code, Message: message, Cause: cause}
}

func GetMessage(code int) string {
	if msg, ok := codeMessages[code]; ok {
		return msg
	}
	return "unknown error"
}

func BadRequest(message string) *AppError {
	if message == "" {
		message = GetMessage(CodeBadRequest)
	}
	return New(CodeBadRequest, message)
}

func Unauthorized(message string) *AppError {
	if message == "" {
		message = GetMessage(CodeUnauthorized)
	}
	return New(CodeUnauthorized, message)
}

func Forbidden(message string) *AppError {
	if message == "" {
		message = GetMessage(CodeForbidden)
	}
	return New(CodeForbidden, message)
}

func NotFound(message string) *AppError {
	if message == "" {
		message = GetMessage(CodeNotFound)
	}
	return New(CodeNotFound, message)
}

func Conflict(message string) *AppError {
	if message == "" {
		message = GetMessage(CodeConflict)
	}
	return New(CodeConflict, message)
}

func InternalError(message string) *AppError {
	if message == "" {
		message = GetMessage(CodeInternalError)
	}
	return New(CodeInternalError, message)
}

func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

func GetErrorCode(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return CodeInternalError
}
