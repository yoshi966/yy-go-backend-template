package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/Yoshioka9709/yy-go-backend-template/util/env"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

const (
	projectName = "yy-go-backend-template"
	envcode     = "test"

	// dockerを強制終了する秒数
	dockerExpireSeconds uint = 900

	// DynamoDBの1トランザクションのアイテム数上限
	dynamodbTransactionItemLimit = 25
)

var (
	dockerPoolOnce    sync.Once
	startDynamoDBOnce sync.Once
)

var (
	dockerPool       *dockertest.Pool
	dynamoDBResource *dockertest.Resource
)

func getDockerPool(t testing.TB) *dockertest.Pool {
	dockerPoolOnce.Do(func() {
		var err error
		dockerPool, err = dockertest.NewPool("")
		require.NoError(t, err)
	})
	return dockerPool
}

// getDockerAuth はDocker Hubのログイン情報を返す
func getDockerAuth() docker.AuthConfiguration {
	auths, err := docker.NewAuthConfigurationsFromDockerCfg()
	if err != nil {
		return docker.AuthConfiguration{}
	}
	return auths.Configs["https://index.docker.io/v1/"]
}

// GetProjectDir はプロジェクトルートのディレクトリを取得する
func GetProjectDir() string {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)

	// `go.mod` のあるディレクトリを探す
	for ; dir != "/" && dir != "."; dir = filepath.Dir(dir) {
		_, err := os.Stat(filepath.Join(dir, "go.mod"))
		if err == nil {
			break
		}
	}
	return dir
}

// InitDynamoDB はテスト用のDynamoDBを初期化します
func InitDynamoDB(t testing.TB) {
	startDynamoDB(t)

	// DynamoDB
	as := session.Must(session.NewSession())
	ac := &aws.Config{
		Endpoint:    aws.String(env.AWSDynamoDBEndpoint()),
		Credentials: credentials.NewStaticCredentials("fakeMyKeyId", "fakeSecretAccessKey", ""),
		Region:      aws.String(endpoints.UsEast1RegionID),
	}
	db := dynamodb.New(as, ac)

	// データファイル取得
	dir := filepath.Join(GetProjectDir(), "fixtures/dynamo")
	files, _ := filepath.Glob(filepath.Join(dir, "structure/*.json"))
	for _, file := range files {
		table := filepath.Base(file)
		envcode := env.Envcode()
		tableName := table[0:len(table)-len(filepath.Ext(table))] + "-" + strings.ToUpper(string(envcode[0])) + (envcode[1:])

		// テーブル削除（エラーは無視）
		_, _ = db.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: &tableName,
		})

		// テーブル定義ファイルのロード
		r, err := os.Open(file) //nolint:gosec
		require.NoError(t, err)
		defer func() {
			err = r.Close()
			require.NoError(t, err)
		}()

		ct := &dynamodb.CreateTableInput{}
		err = json.NewDecoder(r).Decode(ct)
		require.NoError(t, err)

		ct.TableName = &tableName
		_, err = db.CreateTable(ct)
		require.NoError(t, err)

		// TTLの設定
		_, err = db.UpdateTimeToLive(&dynamodb.UpdateTimeToLiveInput{
			TableName: &tableName,
			TimeToLiveSpecification: &dynamodb.TimeToLiveSpecification{
				AttributeName: aws.String("TTL"),
				Enabled:       aws.Bool(true),
			},
		})
		require.NoError(t, err)

		// テストデータファイルのロード
		dataFile := filepath.Join(dir, "data", "test", tableName+".json")
		if _, err := os.Stat(dataFile); err == nil {
			r, err := os.Open(dataFile) // nolint: gosec // テスト用なので無視
			require.NoError(t, err)
			defer func() {
				err = r.Close()
				require.NoError(t, err)
			}()

			items := []*dynamodb.TransactWriteItem{}
			err = json.NewDecoder(r).Decode(&items)
			require.NoError(t, err)

			for i := 0; i < len(items); i += dynamodbTransactionItemLimit {
				end := i + dynamodbTransactionItemLimit
				if end > len(items) {
					end = len(items)
				}

				_, err = db.TransactWriteItems(&dynamodb.TransactWriteItemsInput{
					TransactItems: items[i:end],
				})
				require.NoError(t, err)
			}
		}
	}
}

func startDynamoDB(t testing.TB) {
	startDynamoDBOnce.Do(func() {
		pool := getDockerPool(t)

		var err error
		dynamoDBResource, err = pool.RunWithOptions(&dockertest.RunOptions{
			Repository: "amazon/dynamodb-local",
			Auth:       getDockerAuth(),
			Env: []string{
				"ENVCODE=test",
			},
		})
		require.NoError(t, err)

		err = dynamoDBResource.Expire(dockerExpireSeconds)
		require.NoError(t, err)

		// 準備できるまで待機
		endpoint := fmt.Sprintf("http://%s", dynamoDBResource.GetHostPort("8000/tcp"))
		err = pool.Retry(func() error {
			resp, err := http.Get(endpoint) //nolint:gosec
			if err != nil {
				return err
			}
			return resp.Body.Close()
		})
		require.NoError(t, err)

		// nolint: errcheck
		os.Setenv(env.ProjectNameKey, projectName)
		// nolint: errcheck
		os.Setenv(env.EnvcodeKey, envcode)
		// nolint: errcheck
		os.Setenv(env.AWSDynamoDBEndpointKey, endpoint)
	})
}
