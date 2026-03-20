package errors

import "fmt"

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e AppError) HTTPStatus() int { return HTTPStatusForCode(e.Code) }

func New(code ErrorCode, msg string) AppError {
	return AppError{Code: code, Message: msg}
}

func Wrap(code ErrorCode, msg string, err error) AppError {
	return AppError{Code: code, Message: msg, Err: err}
}
