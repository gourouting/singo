package serializer

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Response is the base serializer.
type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg"`
	Error string      `json:"error,omitempty"`
}

// TrackedErrorResponse is an error response with tracking information.
type TrackedErrorResponse struct {
	Response
	TrackID string `json:"track_id"`
}

// Three-digit error codes reuse the original HTTP meanings.
// Five-digit error codes are application-defined errors.
// Five-digit codes starting with 5 are server-side errors, such as database operation failures.
// Five-digit codes starting with 4 are client-side errors, sometimes caused by client code and sometimes by user actions.
const (
	// CodeCheckLogin means the user is not logged in.
	CodeCheckLogin = 401
	// CodeNoRightErr means unauthorized access.
	CodeNoRightErr = 403
	// CodeDBError means a database operation failed.
	CodeDBError = 50001
	// CodeEncryptError means encryption failed.
	CodeEncryptError = 50002
	// CodeParamErr means a parameter error.
	CodeParamErr = 40001
)

// CheckLogin checks whether the user is logged in.
func CheckLogin() Response {
	return Response{
		Code: CodeCheckLogin,
		Msg:  "Not logged in",
	}
}

// Err handles common errors.
func Err(errCode int, msg string, err error) Response {
	res := Response{
		Code: errCode,
		Msg:  msg,
	}
	// Hide underlying errors in production.
	if err != nil && gin.Mode() != gin.ReleaseMode {
		res.Error = fmt.Sprintf("%+v", err)
	}
	return res
}

// DBErr handles database operation failures.
func DBErr(msg string, err error) Response {
	if msg == "" {
		msg = "Database operation failed"
	}
	return Err(CodeDBError, msg, err)
}

// ParamErr handles parameter errors.
func ParamErr(msg string, err error) Response {
	if msg == "" {
		msg = "Parameter error"
	}
	return Err(CodeParamErr, msg, err)
}
