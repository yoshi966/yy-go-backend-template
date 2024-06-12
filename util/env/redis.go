package env

import (
	"os"
	"strconv"
)

// 環境変数のキー
const (
	RedisEndpointKey = "REDIS_ENDPOINT"
	RedisPoolSizeKey = "REDIS_POOL_SIZE"
)

// RedisEndpoint はRedisのホストを返します。
func RedisEndpoint() string {
	return os.Getenv(RedisEndpointKey)
}

// RedisPoolSize はRedisのプールサイズを返します。
func RedisPoolSize() int {
	n, _ := strconv.Atoi(os.Getenv(RedisPoolSizeKey))
	return n
}
