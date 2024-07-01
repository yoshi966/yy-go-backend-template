API_BIN_NAME:=bin/yy-go-backend-template
TRIMPATH:=-trimpath
API_ENTRY_POINT:=cmd/main.go
LDFLAGS:=-s -w -extldflags="static"
GO_LDFLAGS_VERSION:=-X '${VERSION_PACKAGE_NAME}.version=${VERSION}'
VERSION_PACKAGE_NAME:=github.com/Yoshioka9709/yy-go-backend-template/server
VERSION:=$(shell git rev-parse --short HEAD)
BUILDFLAGS:=-gcflags 'all=-N -l'

GOBIN?=$(shell go env GOPATH)/bin

# サブモジュールのスキーマからコードを自動生成
.PHONY: generate-graphql
generate-graphql:
	@go install github.com/99designs/gqlgen
	$(GOBIN)/gqlgen --config .gqlgen.yml

# 初期データを投入
.PHONY: init-db-local
init-db-local: docker-up init-dynamo init-redis
	@echo "$(BOLD)+-------------------------------------------+$(RESET)"
	@echo "$(BOLD)|                                           |$(RESET)"
	@echo "$(LBLUE)|  dynamo-db-admin: http://localhost:8002/  |$(RESET)"
	@echo "$(LRED)|  RedisInsight:    http://localhost:8001/  |$(RESET)"
	@echo "$(BOLD)|                                           |$(RESET)"
	@echo "$(BOLD)+-------------------------------------------+$(RESET)"

.PHONY: init-dynamo
init-dynamo: docker-up
	AWS_DEFAULT_REGION= AWS_ACCESS_KEY_ID=fake AWS_SECRET_ACCESS_KEY=fake \
	fixtures/dynamo/init_tables.sh

.PHONY: init-redis
init-redis: docker-up
	docker compose exec redis redis-cli flushall

# ローカル開発で作成されたバイナリを削除
.PHONY: remove-bin
remove-bin:
	rm -rf ./bin

# 起動中のdockerコンテナとボリュームを削除
.PHONY: remove-container
remove-container: docker-down
	docker ps -q | \
	xargs docker stop | \
	xargs docker rm

# バイナリ・コンテナ・環境変数をまとめて削除
.PHONY: clean
clean: remove-cache remove-bin

# ビルドコマンド
.PHONY: build-amd64
build-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BFF_BIN_NAME) -tags 'netgo' -installsuffix netgo $(TRIMPATH)\
		-ldflags '$(LDFLAGS) $(GO_LDFLAGS_VERSION)' \
		$(BUILDFLAGS) \
		$(BFF_ENTRY_POINT)

# api用ビルドコマンド
.PHONY: build-api
build-api:
	go build -o $(API_BIN_NAME) -tags 'netgo' -installsuffix netgo $(TRIMPATH)\
		-ldflags '$(LDFLAGS) $(GO_LDFLAGS_VERSION)' \
		$(BUILDFLAGS) \
		$(API_ENTRY_POINT)

# ローカルサーバーを起動
.PHONY: start
start: docker-up
	@go install github.com/air-verse/air
	ulimit -n 1024 && $(GOBIN)/air -c .air.toml

.PHONY: build-local
build-local:
	CGO_ENABLED=0 GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) $(MAKE) build-api TRIMPATH=

# ローカル開発に必要なDocker環境を実行
.PHONY: docker-up
docker-up:
	docker compose up -d

# ローカル開発で利用しているDocker環境を停止
.PHONY: docker-down
docker-down:
	docker compose down --remove-orphans
	docker volume rm yy-go-backend-template-dynamodb-local
	docker volume rm yy-go-backend-template-dynamodb-admin
	docker volume rm yy-go-backend-template-dynamodb-redis
	docker volume rm yy-go-backend-template-dynamodb-redisinsight

