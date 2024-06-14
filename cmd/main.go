// BFF Server
package main

import (
	"fmt"
	"os"

	"github.com/Yoshioka9709/yy-go-backend-template/server"
	"github.com/Yoshioka9709/yy-go-backend-template/util"
)

func main() {
	util.LoadEnv()

	if err := server.Start(); err != nil {
		_ = fmt.Errorf("server error : %w", err)
		os.Exit(1)
	}
}
