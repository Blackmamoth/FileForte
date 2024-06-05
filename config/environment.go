package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type ApplicationConfig struct {
	APP_PORT      string
	ENVIRONMENT   string
	LOG_FILE_PATH string
	LOG_FILE_NAME string
}

var AppConfig ApplicationConfig = newAppConfig()

func newAppConfig() ApplicationConfig {
	godotenv.Load()
	return ApplicationConfig{
		APP_PORT:      getEnvString("APP_PORT"),
		ENVIRONMENT:   getEnvString("ENVIRONMENT"),
		LOG_FILE_PATH: getEnvString("LOG_FILE_PATH"),
		LOG_FILE_NAME: getEnvString("LOG_FILE_NAME"),
	}
}

func getEnvString(key string) string {
	value, exist := os.LookupEnv(key)
	if !exist {
		fmt.Printf("Environment variable [%s] does not exist.\n", key)
		os.Exit(1)
	}
	return value
}
