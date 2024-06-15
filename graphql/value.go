package graphql

import (
	"context"

	"github.com/Yoshioka9709/yy-go-backend-template/factory"
)

// getService はサービスファクトリを返す
func getService(ctx context.Context) factory.Service {
	c := getContext(ctx)
	return c.Service()
}
