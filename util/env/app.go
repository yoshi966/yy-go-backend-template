package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// 環境変数のキー
const (
	EnvcodeKey           = "ENVCODE"
	ProjectNameKey       = "PROJECT_NAME"
	DebugEnabledKey      = "DEBUG_ENABLED"
	ErrorStackEnabledKey = "ERROR_STACK_ENABLED"
	CorsAllowOriginKey   = "CORS_ALLOW_ORIGINS"
	CorsAllowMethodsKey  = "CORS_ALLOW_METHODS"
	CorsAllowHeadersKey  = "CORS_ALLOW_HEADERS"
	CorsAllowWildcardKey = "CORS_ALLOW_WILDCARD"
	CorsExposeHeadersKey = "CORS_EXPOSE_HEADERS"
	CorsCredentialsKey   = "CORS_CREDENTIALS"
	CorsMaxAgeKey        = "CORS_MAX_AGE"
)

// ProjectName は環境変数よりプロジェクト名を返します。
func ProjectName() string {
	return os.Getenv(ProjectNameKey)
}

// Envcode は環境変数より環境名を返します。
func Envcode() string {
	return os.Getenv(EnvcodeKey)
}

// IsTest はテスト環境で動作しているか返します。
func IsTest() bool {
	envCode := Envcode()
	return envCode == "test"
}

// IsLocal はローカル環境で動作しているか返します。
func IsLocal() bool {
	envCode := Envcode()
	return envCode == "local" || IsTest()
}

// IsDevelopment は開発環境（ローカルを含む）で動作しているか返します。
func IsDevelopment() bool {
	envCode := Envcode()
	return envCode == "dev" || IsLocal()
}

// IsProduction は本番環境で動作しているか返します。
func IsProduction() bool {
	envCode := Envcode()
	return envCode == "prd"
}

// IsDebugAPIEnabled はデバッグAPIが有効かを返します。
func IsDebugAPIEnabled() bool {
	env := os.Getenv(DebugEnabledKey)
	b, _ := strconv.ParseBool(env)
	return b
}

// IsErrorStackEnabled はAPIのエラースタックが有効かを返します。
func IsErrorStackEnabled() bool {
	env := os.Getenv(ErrorStackEnabledKey)
	b, _ := strconv.ParseBool(env)
	return b
}

// CorsAllowOrigin は環境変数よりAllowOriginを返します。
func CorsAllowOrigin() []string {
	env := os.Getenv(CorsAllowOriginKey)
	if env == "" {
		return []string{}
	}
	return strings.Split(env, ",")
}

// CorsAllowMethods は環境変数よりAllowMethodsを返します。
func CorsAllowMethods() []string {
	env := os.Getenv(CorsAllowMethodsKey)
	if env == "" {
		return []string{}
	}
	return strings.Split(env, ",")
}

// CorsAllowHeaders は環境変数よりAllowHeadersを返します。
func CorsAllowHeaders() []string {
	env := os.Getenv(CorsAllowHeadersKey)
	if env == "" {
		return []string{}
	}
	return strings.Split(env, ",")
}

// CorsAllowWildcard はOriginのワイルドカードを許容する設定を返します。
func CorsAllowWildcard() bool {
	env := os.Getenv(CorsAllowWildcardKey)
	b, _ := strconv.ParseBool(env)
	return b
}

// CorsExposeHeaders は環境変数よりExposeHeadersを返します。
func CorsExposeHeaders() []string {
	env := os.Getenv(CorsExposeHeadersKey)
	if env == "" {
		return []string{}
	}
	return strings.Split(env, ",")
}

// CorsCredentials は環境変数よりCorsCredentialsを返します。
func CorsCredentials() bool {
	env := os.Getenv(CorsCredentialsKey)
	b, _ := strconv.ParseBool(env)
	return b
}

// CorsMaxAge は環境変数よりCorsMaxAgeを返します。
func CorsMaxAge() (time.Duration, error) {
	env := os.Getenv(CorsMaxAgeKey)
	return time.ParseDuration(env)
}
