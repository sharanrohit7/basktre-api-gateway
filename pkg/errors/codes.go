package errors

type ErrorCode string

const (
	ErrBadRequest     ErrorCode = "BAD_REQUEST"
	ErrUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrInternalServer ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrExternalAPI    ErrorCode = "EXTERNAL_API_ERROR"
)
