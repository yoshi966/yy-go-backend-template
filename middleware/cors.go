package middleware

import (
	"time"

	"yy-go-backend-template/util/env"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// .envから値を取得出来なかった場合に利用する
const defaultMaxAge = 12 * time.Hour

// CORSMiddleware はCORSを有効にします
func CORSMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = env.CorsAllowOrigin()
	config.AllowMethods = env.CorsAllowMethods()
	config.AllowHeaders = env.CorsAllowHeaders()
	config.ExposeHeaders = env.CorsExposeHeaders()
	config.AllowCredentials = env.CorsCredentials()
	config.AllowWildcard = env.CorsAllowWildcard()
	maxAge, err := env.CorsMaxAge()
	if err != nil {
		log.Error().Err(err).Send()
		maxAge = defaultMaxAge
	}
	config.MaxAge = maxAge

	return cors.New(config)
}
