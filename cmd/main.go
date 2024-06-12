// BFF Server
package main

import (
	"fmt"
	"os"

	"yy-go-backend-template/server"
	"yy-go-backend-template/util"
)

func main() {
	util.LoadEnv()

	if err := server.Start(); err != nil {
		_ = fmt.Errorf("server error : %w", err)
		os.Exit(1)
	}
}
