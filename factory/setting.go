package factory

import (
	"github.com/guregu/dynamo"
)

// サービスファクトリのキー
const (
	ServiceFactoryKey = "service_factory"
)

// ServiceFactorySettings はサービスファクトリの設定
type ServiceFactorySettings struct {
	DynamoDB *dynamo.DB
}
