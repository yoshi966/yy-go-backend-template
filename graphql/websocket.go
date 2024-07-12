package graphql

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"
)

const (
	pingInterval     = 30 * time.Second
	handshakeTimeout = 4000 * time.Second
	readBufferSize   = 1 << 10
	writeBufferSize  = 1 << 10
)

// WebSocket の設定
func transportWebsocket() *transport.Websocket {
	return &transport.Websocket{
		KeepAlivePingInterval: pingInterval,
		Upgrader: websocket.Upgrader{
			HandshakeTimeout: handshakeTimeout,
			CheckOrigin: func(_ *http.Request) bool {
				// ひとまずCorsは全部許可
				return true
			},
			EnableCompression: true,
			ReadBufferSize:    readBufferSize,
			WriteBufferSize:   writeBufferSize,
		},
	}
}
