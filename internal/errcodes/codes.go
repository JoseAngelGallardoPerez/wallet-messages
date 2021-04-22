package errcodes

import "net/http"

const (
	CodeForbidden              = "FORBIDDEN"
	CodeMessageNotFound        = "MESSAGE_NOT_FOUND"
	CodeInvalidQueryParameters = "INVALID_QUERY_PARAMS"
	CodeUserNotFound           = "USER_NOT_FOUND"
)

func HttpStatusCodeByErrCode(code string) int {
	if status, ok := statusCodes[code]; ok {
		return status
	}
	panic("code is not present")
}

var statusCodes = map[string]int{
	CodeForbidden:              http.StatusForbidden,
	CodeMessageNotFound:        http.StatusNotFound,
	CodeInvalidQueryParameters: http.StatusBadRequest,
	CodeUserNotFound:           http.StatusNotFound,
}
