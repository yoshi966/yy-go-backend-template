package env

import (
	"os"
)

// 環境変数のキー
const (
	// DynamoDB
	AWSDynamoDBEndpointKey = "AWS_DYNAMODB_ENDPOINT"
)

// AWSDynamoDBEndpoint はDynamoDBのホストを返します。
func AWSDynamoDBEndpoint() string {
	return os.Getenv(AWSDynamoDBEndpointKey)
}
