package util

import (
	"crypto/rand"

	"github.com/oklog/ulid/v2"
)

// ULIDEntropy はULID生成のエントロピー
//
// テストで固定の値が欲しい場合は上書きしてください
var ULIDEntropy = rand.Reader

// NewULID はULIDを生成します
func NewULID() ulid.ULID {
	ms := ulid.Timestamp(GetTimeNow())

	// エラーは無視
	return ulid.MustNew(ms, ULIDEntropy)
}

// GetULIDString 文字列で ULID を出力
func GetULIDString() string {
	return NewULID().String()
}
