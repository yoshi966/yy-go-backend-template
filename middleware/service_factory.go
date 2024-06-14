package middleware

import (
	"github.com/Yoshioka9709/yy-go-backend-template/factory"

	"github.com/gin-gonic/gin"
)

// ServiceFactoryMiddleware サービスをコンテキストに設定する
func ServiceFactoryMiddleware(s factory.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(factory.ServiceFactoryKey, s)
		c.Next()
	}
}
