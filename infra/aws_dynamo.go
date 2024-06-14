package infra

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/guregu/dynamo"

	"github.com/Yoshioka9709/yy-go-backend-template/util/env"
	"github.com/Yoshioka9709/yy-go-backend-template/util/str"
)

// NewDynamoDBClient DynamoDBの初期化
func NewDynamoDBClient() *dynamo.DB {
	// ローカルでは DynamoDB-Localを使用しているため設定を行う。
	// DYNAMODB_ENDPOINT が設定されていればローカル環境
	config := &aws.Config{}
	endpoint := env.AWSDynamoDBEndpoint()

	if endpoint != "" {
		config.Endpoint = aws.String(endpoint)
		config.Credentials = credentials.NewStaticCredentials("fakeMyKeyId", "fakeSecretAccessKey", "")
		config.Region = aws.String(endpoints.UsEast1RegionID) // us-east-1 でないと dynamodb-local は使えないためここで指定する
	} else {
		config.Region = aws.String(endpoints.ApNortheast1RegionID)
	}

	db := dynamo.New(getAWSSession(), config)
	return db
}

// TableName テーブル名を取得
func TableName(name string) string {
	return str.UpperFirst(name) + "-" + str.UpperFirst(env.Envcode())
}
