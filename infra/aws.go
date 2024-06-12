package infra

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

// getAWSSession はAWSセッションを取得します
// https://docs.aws.amazon.com/sdk-for-go/api/aws/session/
func getAWSSession() *session.Session {
	return session.Must(session.NewSession())
}
