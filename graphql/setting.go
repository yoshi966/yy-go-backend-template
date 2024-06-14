package graphql

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"sync"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/Yoshioka9709/yy-go-backend-template/factory"
	"github.com/Yoshioka9709/yy-go-backend-template/graphql/generated"
	"github.com/Yoshioka9709/yy-go-backend-template/infra"
	"github.com/Yoshioka9709/yy-go-backend-template/server/key"
	"github.com/Yoshioka9709/yy-go-backend-template/util/env"
	"github.com/Yoshioka9709/yy-go-backend-template/util/errs"
	"github.com/Yoshioka9709/yy-go-backend-template/util/errs/codes"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// クエリログ整形フォーマット
var reGraphQLRawQuery = regexp.MustCompile(`\n *| {2,}`)

const (
	// 許容するクエリの複雑さ
	complexityLimit = 196
	// クエリキャッシュサイズ
	queryCacheSize = 1024
	// APQキャッシュサイズ
	apqCacheSize = 128

	// チャット の pubsub用Redisキー
	chatPubSubKey string = "chat:pubsub"

	mb            int64 = 1 << 20
	maxMemory           = 32 * mb
	maxUploadSize       = 50 * mb
)

// AppContext 独自実装したコンテキスト
// サーバ内で持ち回したい場合はここに追加する
type AppContext struct {
	*gin.Context
}

// getContext は独自定義したコンテキストを返す
func getContext(ctx context.Context) AppContext {
	return AppContext{
		Context: ctx.Value(key.GinContextKey).(*gin.Context),
	}
}

// Service コンテキストにセットされたサービスを取得
func (a *AppContext) Service() factory.Service {
	return a.MustGet(factory.ServiceFactoryKey).(factory.Service)
}

// Handler は GraphQL の処理を行うハンドラーです
func Handler() gin.HandlerFunc {
	h := newHandler()
	return func(c *gin.Context) {
		// context に必要な値をセット
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, key.GinContextKey, c)
		c.Request = c.Request.WithContext(ctx)

		// 処理を実行
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (r *Resolver) startSubscribingRedis() {
	go func() {
		ctx := context.Background()
		subscriber := r.redisClient.NewSubscribe(ctx, chatPubSubKey)
		defer func() {
			err := subscriber.Close()
			if err != nil {
				log.Debug().Err(err).Send()
			}
		}()
		for {
			receive, err := subscriber.Receive(ctx)
			if err != nil {
				log.Warn().Err(err).Msg("pubsub redis connection refused")
			}
			log.Debug().Interface("value", receive).Msg("pubsub")
		}
	}()
}

func newResolver() *Resolver {
	r := Resolver{
		redisClient: infra.NewRedisClient(),
	}
	return &r
}

func newGraphQLConfig() generated.Config {
	resolver := newResolver()

	if !env.IsTest() {
		// サーバ起動時にPubSubを起動
		resolver.startSubscribingRedis()
	}

	return generated.Config{
		Resolvers: resolver,
	}
}

var (
	h         *handler.Server
	setupOnce sync.Once
)

func newHandler() *handler.Server {
	// NOTE newResolver() が2回呼ばれないための対策
	setupOnce.Do(setupHandler)
	return h
}

func setupHandler() {
	config := newGraphQLConfig()

	h = handler.New(generated.NewExecutableSchema(config))

	// 開発環境のみ Introspection を有効にする
	if env.IsDebugAPIEnabled() {
		h.Use(extension.Introspection{})
	}
	h.Use(extension.FixedComplexityLimit(complexityLimit))
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(apqCacheSize),
	})

	h.SetErrorPresenter(errorPresenter)
	h.SetRecoverFunc(panicRecorder)
	h.SetQueryCache(lru.New(queryCacheSize))

	// queryログの出力
	if env.GraphQLQueryLog() {
		h.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
			oc := graphql.GetOperationContext(ctx)
			// queryの整形
			q := reGraphQLRawQuery.ReplaceAllString(oc.RawQuery, " ")
			// 行末の半角スペースを削除
			q = strings.TrimRight(q, " ")
			log.Debug().Str("query", q).Send()
			return next(ctx)
		})
	}

	h.AddTransport(transport.POST{})
	h.AddTransport(transport.MultipartForm{
		MaxMemory:     maxMemory,
		MaxUploadSize: maxUploadSize,
	})

	h.AddTransport(transportWebsocket())
}

func errorPresenter(ctx context.Context, e error) *gqlerror.Error {
	err := graphql.DefaultErrorPresenter(ctx, e)

	if errcode.KindProtocol == errcode.GetErrorKind(gqlerror.List{err}) {
		// gqlgen文法エラーはそのまま返す
		log.Warn().
			Err(err).
			Msg("GraphQL protocol error")

		// ErrorMiddleware で処理しないタイプに設定
		getContext(ctx).Error(err).SetType(gin.ErrorTypePublic)
		return err
	}

	log.Info().
		Err(err).
		Msg("GraphQL error")

	var ce errs.CodeError
	if !errors.As(err, &ce) {
		ce = errs.Wrap(codes.InternalError, errors.Unwrap(err))
	}

	extensions := map[string]any{
		"code": ce.Code().String(),
	}

	// スタックトレース
	if env.IsErrorStackEnabled() {
		extensions["stack"] = ce.Stack()
	}
	err.Extensions = extensions

	// ErrorMiddleware で処理しないタイプに設定
	getContext(ctx).Error(ce).SetType(gin.ErrorTypePublic)

	return err
}

func panicRecorder(_ context.Context, e any) error {
	return errs.NewPanicError(e)
}
