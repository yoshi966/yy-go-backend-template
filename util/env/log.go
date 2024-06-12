package env

import (
	"os"
	"strconv"
)

// 環境変数のキー
const (
	GraphQLQueryLogKey = "GRAPHQL_QUERY_LOG"
)

// GraphQLQueryLog はGraphQLのクエリログを出力するかを返します。
func GraphQLQueryLog() bool {
	env := os.Getenv(GraphQLQueryLogKey)
	b, _ := strconv.ParseBool(env)
	return b
}
