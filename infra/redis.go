package infra

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/Yoshioka9709/yy-go-backend-template/util/env"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	lockPrefix  = "LOCK:"
	lockTimeout = 1 * time.Minute
)

// Redis はRedisクライアント
type Redis struct {
	client redis.UniversalClient
}

// NewRedisClient RedisClientの初期化
func NewRedisClient() *Redis {
	options := &redis.UniversalOptions{
		Addrs:    []string{env.RedisEndpoint()},
		PoolSize: env.RedisPoolSize(),
	}
	return &Redis{redis.NewClient(options.Simple())}
}

// Close はRedisクライアントを閉じる
func (r *Redis) Close() error {
	return r.client.Close()
}

// Publish はRedisにメッセージを通知する
func (r *Redis) Publish(ctx context.Context, channelKey string, obj any) error {
	b, err := json.Marshal(obj)
	if err != nil {
		log.Error().Msgf("can not marshal json %+v", obj)
		return err
	}
	return r.client.Publish(ctx, channelKey, b).Err()
}

// NewSubscribe はRedisからメッセージを受け取るサブスクライバを生成する
func (r *Redis) NewSubscribe(ctx context.Context, channelKey string) *redis.PubSub {
	return r.client.Subscribe(ctx, channelKey)
}

// Set はRedisに値を設定する
func (r *Redis) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	buf := &bytes.Buffer{}

	// nil 以外はエンコードして保存
	if !r.isNil(value) {
		err := gob.NewEncoder(buf).Encode(value)
		if err != nil {
			return err
		}
	}

	return r.client.Set(ctx, key, buf.Bytes(), expiration).Err()
}

// Get はRedisから値を取得する
func (r *Redis) Get(ctx context.Context, key string, obj any) error {
	b, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	// 空なら nil 扱い
	if len(b) == 0 {
		rv := reflect.ValueOf(obj)
		if rv.Kind() == reflect.Pointer {
			rv = rv.Elem()
			rv.Set(reflect.New(rv.Type()).Elem())
		}
		return nil
	}

	return gob.NewDecoder(bytes.NewBuffer(b)).Decode(obj)
}

// Del はRedisのキーを削除する
func (r *Redis) Del(ctx context.Context, keys ...string) error {
	_, err := r.DelCount(ctx, keys...)
	return err
}

// DelCount はRedisのキーを削除して、削除した件数を返す
func (r *Redis) DelCount(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Del(ctx, keys...).Result()
}

// Exists はRedisのキーが存在するか
//
// エラーは無視します
func (r *Redis) Exists(ctx context.Context, key string) bool {
	return r.client.Exists(ctx, key).Val() == 1
}

// Keys はRedisのキーを検索する
func (r *Redis) Keys(ctx context.Context, pattern string) ([]string, error) {
	return r.client.Keys(ctx, pattern).Result()
}

// TTL はRedisの有効時間を返す
//
// 無期限の場合は -1 を返す。存在しない場合は -2 を返す。
func (r *Redis) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// Lock は指定したキーで排他ロックを実行する
//
// Redisの `SETNX` を使用したロックを実行。ロックが取得できない場合はエラーを返す。
// ロックが取得できた場合は、アンロック関数を返す。
func (r *Redis) Lock(ctx context.Context, key string) (func(ctx context.Context), error) {
	return r.LockWithTimeout(ctx, key, lockTimeout)
}

// LockWithTimeout は指定したキーとタイムアウトで排他ロックを実行する
//
// Redisの `SETNX` を使用したロックを実行。ロックが取得できない場合はエラーを返す。
// ロックが取得できた場合は、アンロック関数を返す。
func (r *Redis) LockWithTimeout(ctx context.Context, key string, timeout time.Duration) (func(ctx context.Context), error) {
	lockKey := lockPrefix + key
	b, err := r.client.SetNX(ctx, lockKey, key, timeout).Result()
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, fmt.Errorf("%s is locked", key)
	}

	unlock := func(ctx context.Context) {
		r.client.Del(ctx, lockKey)
	}
	return unlock, nil
}

// isNil は値が nil かどうかを判定する
func (r *Redis) isNil(value any) bool {
	if value == nil {
		return true
	}

	rv := reflect.ValueOf(value)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return true
		}
		rv = rv.Elem()
	}
	return false
}

// GetLockAddLikeCountKey いいね数追加時のロックに使うキーを返します。
func GetLockAddLikeCountKey(userID string) string {
	return "LockAddLikeCount:" + userID
}

// GetIsReadMatchKey マッチの既読管理に使うキーを返します。
func GetIsReadMatchKey(userID, matchID string) string {
	return "IsReadMatch:" + userID + ":" + matchID
}
