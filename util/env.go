package util

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// LoadEnv apply .env to ENVIRONMENT VARIABLE.
func LoadEnv() {
	if _, found := os.LookupEnv("ENV_FILE"); !found {
		os.Setenv("ENV_FILE", ".env")
	}

	envfile := os.Getenv("ENV_FILE")

	if err := godotenv.Load(envfile); err != nil {
		log.Debug().Err(err).Msg("failed to load env file")
	}

	log.Debug().Interface("environmentVariables", getEnvMap()).Send()
}

// getEnvMap は環境変数をmapで返す
func getEnvMap() map[string]string {
	environ := os.Environ()
	envMap := make(map[string]string)
	for _, v := range environ {
		kv := strings.Split(v, "=")
		if kv[0] == "" {
			continue
		}
		envMap[kv[0]] = kv[1]
	}
	return envMap
}
