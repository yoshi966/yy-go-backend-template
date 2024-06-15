package model

// DynamoDB のテーブルプレフィックス
const (
	UsersTablePrefix = "users" // ユーザテーブル prefix
	TodosTablePrefix = "todos" // マスターテーブル prefix
)

// PK はDynamoDBのパーティションキー
type PK string

// マスターのパーティションキー
const (
	PKUser PK = "user"
	PKTodo PK = "todo"
)
