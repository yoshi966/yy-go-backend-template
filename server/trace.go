package server

import (
	"github.com/Yoshioka9709/yy-go-backend-template/util/env"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// trace 各種トレース情報をヘッダーに埋め込む
func trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("x-amzn-trace-id")
		c.Writer.Header().Set("x-request-id", traceID)
		onTrace, err := env.OnTrace()
		if err != nil {
			log.Info().
				Err(err).
				Msg("There is an error in the environment variable ON_TRACE.")
		}

		if onTrace {
			c.Writer.Header().Set("x-api-version", version)
		}
		c.Next()
	}
}
