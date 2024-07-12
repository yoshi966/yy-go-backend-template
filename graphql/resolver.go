package graphql

import (
	"github.com/Yoshioka9709/yy-go-backend-template/infra"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is GraphQL resolver
type Resolver struct {
	redisClient *infra.Redis
}
