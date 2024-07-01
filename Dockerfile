# アプリケーションビルド用
FROM golang:1.22.4-alpine AS build_stage

WORKDIR /src

# hadolint ignore=DL3018
RUN apk add --no-cache make git

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN make build-amd64

# アプリケーション 実行用

# hadolint ignore=DL3007
FROM gcr.io/distroless/base:latest AS runtime_api

WORKDIR /app

COPY --from=build_stage /src/bin/yy-go-backend-template /app/yy-go-backend-template
COPY --from=build_stage /src/.env /app/.env

EXPOSE 12340

CMD ["./yy-go-backend-template"]