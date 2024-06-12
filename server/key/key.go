package key

// ContextKey はコンテキストに設定するキー
type ContextKey string

// String 文字列で取得
func (c ContextKey) String() string {
	return string(c)
}

// ContextKey の一覧
const (
	GinContextKey ContextKey = "ginContext" // ginのコンテキストを設定するキー
)
