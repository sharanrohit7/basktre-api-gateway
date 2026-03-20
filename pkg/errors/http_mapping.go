package errors

import "net/http"

var ErrorCodeToHTTPStatus = map[ErrorCode]int{
	ErrBadRequest:     http.StatusBadRequest,
	ErrUnauthorized:   http.StatusUnauthorized,
	ErrExternalAPI:    http.StatusBadGateway,
	ErrInternalServer: http.StatusInternalServerError,
}

func HTTPStatusForCode(code ErrorCode) int {
	if status, ok := ErrorCodeToHTTPStatus[code]; ok {
		return status
	}
	return http.StatusInternalServerError
}
