package env

import (
	"os"
)

const (
	// DynamoDB
	AWSDynamoDBEndpointKey = "AWS_DYNAMODB_ENDPOINT"
)

// AWSDynamoDBEndpoint はDynamoDBのホストを返します。
func AWSDynamoDBEndpoint() string {
	return os.Getenv(AWSDynamoDBEndpointKey)
}
