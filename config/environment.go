package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ApplicationConfiguration struct {
	APP_PORT      string
	ENVIRONMENT   string
	LOG_FILE_PATH string
	LOG_FILE_NAME string
}

type MySQLConfiguration struct {
	MYSQL_USER    string
	MYSQL_PASS    string
	MYSQL_HOST    string
	MYSQL_PORT    string
	MYSQL_DB_NAME string
}

type JWTConfiguration struct {
	JWT_ACCESS_TOKEN_EXPIRATION_IN_MINS int64
	JWT_ACCESS_TOKEN_SECRET             string
	JWT_REFRESH_TOKEN_SECRET            string
}

var AppConfig ApplicationConfiguration = newAppConfig()
var MySQLConfig MySQLConfiguration = newMySQLConfig()
var JWTConfig JWTConfiguration = newJWTConfig()

func newAppConfig() ApplicationConfiguration {
	godotenv.Load()
	return ApplicationConfiguration{
		APP_PORT:      getEnvString("APP_PORT"),
		ENVIRONMENT:   getEnvString("ENVIRONMENT"),
		LOG_FILE_PATH: getEnvString("LOG_FILE_PATH"),
		LOG_FILE_NAME: getEnvString("LOG_FILE_NAME"),
	}
}

func newMySQLConfig() MySQLConfiguration {
	return MySQLConfiguration{
		MYSQL_USER:    getEnvString("MYSQL_USER"),
		MYSQL_PASS:    getEnvString("MYSQL_PASS"),
		MYSQL_HOST:    getEnvString("MYSQL_HOST"),
		MYSQL_PORT:    getEnvString("MYSQL_PORT"),
		MYSQL_DB_NAME: getEnvString("MYSQL_DB_NAME"),
	}
}

func newJWTConfig() JWTConfiguration {
	return JWTConfiguration{
		JWT_ACCESS_TOKEN_EXPIRATION_IN_MINS: getEnvInt("JWT_ACCESS_TOKEN_EXPIRATION_IN_MINS"),
		JWT_ACCESS_TOKEN_SECRET:             getEnvString("JWT_ACCESS_TOKEN_SECRET"),
		JWT_REFRESH_TOKEN_SECRET:            getEnvString("JWT_REFRESH_TOKEN_SECRET"),
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

func getEnvInt(key string) int64 {
	value, exist := os.LookupEnv(key)
	if !exist {
		fmt.Printf("Environment variable [%s] does not exist.\n", key)
		os.Exit(1)
	}
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		fmt.Printf("Error: %v.\n", err)
		os.Exit(1)
	}
	return intValue
}
