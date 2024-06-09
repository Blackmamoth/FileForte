package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type appConfig struct {
	APP_PORT      string
	ENVIRONMENT   string
	LOG_FILE_PATH string
	LOG_FILE_NAME string
}

type mySQLConfiguration struct {
	MYSQL_USER    string
	MYSQL_PASS    string
	MYSQL_HOST    string
	MYSQL_PORT    string
	MYSQL_DB_NAME string
}

type jwtConfig struct {
	JWT_ACCESS_TOKEN_EXPIRATION_IN_MINS  int64
	JWT_ACCESS_TOKEN_SECRET              string
	JWT_REFRESH_TOKEN_SECRET             string
	JWT_REFRESH_TOKEN_EXPIRATION_IN_DAYS int64
	ACCESS_TOKEN_HEADER_NAME             string
	REFRESH_TOKEN_COOKIE_NAME            string
}

type fileUploadConfig struct {
	MAX_UPLOAD_SIZE   int64
	MEDIA_UPLOAD_PATH string
}

var AppConfig appConfig = newAppConfig()
var MySQLConfig mySQLConfiguration = newMySQLConfig()
var JWTConfig jwtConfig = newJWTConfig()
var FileUploadConfig fileUploadConfig = newFileUploadConfig()

func newAppConfig() appConfig {
	godotenv.Load()
	return appConfig{
		APP_PORT:      getEnvString("APP_PORT"),
		ENVIRONMENT:   getEnvString("ENVIRONMENT"),
		LOG_FILE_PATH: getEnvString("LOG_FILE_PATH"),
		LOG_FILE_NAME: getEnvString("LOG_FILE_NAME"),
	}
}

func newMySQLConfig() mySQLConfiguration {
	return mySQLConfiguration{
		MYSQL_USER:    getEnvString("MYSQL_USER"),
		MYSQL_PASS:    getEnvString("MYSQL_PASS"),
		MYSQL_HOST:    getEnvString("MYSQL_HOST"),
		MYSQL_PORT:    getEnvString("MYSQL_PORT"),
		MYSQL_DB_NAME: getEnvString("MYSQL_DB_NAME"),
	}
}

func newJWTConfig() jwtConfig {
	return jwtConfig{
		JWT_ACCESS_TOKEN_EXPIRATION_IN_MINS:  getEnvInt("JWT_ACCESS_TOKEN_EXPIRATION_IN_MINS"),
		JWT_ACCESS_TOKEN_SECRET:              getEnvString("JWT_ACCESS_TOKEN_SECRET"),
		JWT_REFRESH_TOKEN_SECRET:             getEnvString("JWT_REFRESH_TOKEN_SECRET"),
		JWT_REFRESH_TOKEN_EXPIRATION_IN_DAYS: getEnvInt("JWT_REFRESH_TOKEN_EXPIRATION_IN_DAYS"),
		ACCESS_TOKEN_HEADER_NAME:             getEnvString("ACCESS_TOKEN_HEADER_NAME"),
		REFRESH_TOKEN_COOKIE_NAME:            getEnvString("REFRESH_TOKEN_COOKIE_NAME"),
	}
}

func newFileUploadConfig() fileUploadConfig {
	return fileUploadConfig{
		MAX_UPLOAD_SIZE:   getEnvInt("MAX_UPLOAD_SIZE"),
		MEDIA_UPLOAD_PATH: getEnvString("MEDIA_UPLOAD_PATH"),
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
