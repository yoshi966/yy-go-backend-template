package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"yy-go-backend-template/factory"
	"yy-go-backend-template/infra"
	"yy-go-backend-template/middleware"
	"yy-go-backend-template/util/env"

	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
)

const (
	shutdownTimeOut   = 10 * time.Second
	readHeaderTimeout = 3600 * time.Second
)

var (
	dynamoClient *dynamo.DB
)

func initClients() {
	var err error
	dynamoClient = infra.NewDynamoDBClient()

	if err != nil {
		panic(err)
	}
}

func getServiceSetting() *factory.ServiceFactorySettings {
	initClients()
	return &factory.ServiceFactorySettings{
		DynamoDB: dynamoClient,
	}
}

// Start starts api server
func Start() error {
	// Ginの初期化
	r := gin.New()

	gin.SetMode(gin.DebugMode)

	// Service
	setting := getServiceSetting()

	// サービスファクトリの初期化
	serviceRegistory := factory.NewService(setting)
	r.Use(middleware.ServiceFactoryMiddleware(serviceRegistory))

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorMiddleware())
	r.Use(trace())

	defineRoutes(r)

	port := env.Port()
	if port == "" {
		// 環境変数でポートの指定が無ければGinのデフォルトのポートを指定
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           r,
	}

	go func() {
		// Start server
		fmt.Printf("🐞🍣\trunning server version\t: %v\n", version)
		fmt.Printf("🚀💫\tserver listen on port\t: %v\n", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	// Wait for "interrupt" or "kill" signal to gracefully shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit

	fmt.Printf("Shutdown Server with Signal %v\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeOut)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server Shutdown: %v\n", err)
	}
	fmt.Println("Server exiting")

	return nil
}
