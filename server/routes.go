package server

import (
	"fmt"
	"net/http"

	"yy-go-backend-template/util/env"

	"github.com/gin-gonic/gin"
)

const (
	apiVersion = "v1"
)

func getDebugMessage(c *gin.Context) string {
	format := "Hello, %s. host:%s"
	host := c.Request.Host
	return fmt.Sprintf(format, env.Envcode(), host)
}

func defineRoutes(r gin.IRouter) {
	// routes := r.Group("/" + apiVersion)

	// =================
	// Health Check Path
	// =================
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, getDebugMessage(c))
	})

	// =================
	// GraphQL App Handler
	// =================
	// routes.Any("/graphql", graphql.Handler())
}
