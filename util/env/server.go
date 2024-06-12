package env

import (
	"os"
	"strconv"
)

const (
	PortKey    = "PORT"
	OnTraceKey = "ON_TRACE"
)

// Port はサーバーのポート番号を返します
func Port() string {
	return os.Getenv(PortKey)
}

// OnTrace は各種トレース情報をヘッダーに埋め込むか否かを返します。
func OnTrace() (bool, error) {
	return strconv.ParseBool(os.Getenv(OnTraceKey))
}
