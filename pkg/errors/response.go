package errors

import (
	"github.com/basktre/api-gateway/pkg/requestid"
	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func RespondError(c *gin.Context, err AppError) {
	c.JSON(err.HTTPStatus(), gin.H{
		"success": false,
		"error": ErrorBody{
			Code:    err.Code,
			Message: err.Message,
		},
		"request_id": requestid.GetRequestID(c.Request.Context()),
	})
}
