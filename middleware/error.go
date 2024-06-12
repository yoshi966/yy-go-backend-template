package middleware

import (
	"errors"

	"yy-go-backend-template/util/env"
	"yy-go-backend-template/util/errs"
	"yy-go-backend-template/util/errs/codes"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ErrorMiddleware はエラーハンドリングのミドルウェア
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		ginErr := c.Errors.ByType(gin.ErrorTypePrivate).Last()
		if ginErr != nil {
			err := ginErr.Err

			log.Debug().Err(err).Msg("Handler error")

			var ce errs.CodeError
			ok := errors.As(err, &ce)
			if !ok {
				ce = errs.Wrap(codes.InternalError, err)
			}
			c.AbortWithStatusJSON(ce.Code().HTTPStatus(), newErrorResponse(ce))
		}
	}
}

type errorResponse struct {
	Error errorResponseBody `json:"error"`
}

type errorResponseBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Stack   string `json:"stack,omitempty"`
}

func newErrorResponse(ce errs.CodeError) *errorResponse {
	res := &errorResponse{
		errorResponseBody{
			Code:    ce.Code().String(),
			Message: ce.Message(env.IsDevelopment()),
		},
	}

	// スタックトレース
	if env.IsErrorStackEnabled() {
		res.Error.Stack = ce.Stack()
	}

	return res
}
