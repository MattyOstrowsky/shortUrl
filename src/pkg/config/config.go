package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	defaultDbHost         = "localhost"
	defaultCollectionName = "urls"
	defaultDbName         = "urldb"
	defaultMinLength      = 6
	defaultMaxLength      = 8
	defaultLength         = 6
)

type UrlConfig struct {
	ServerPort     string
	DbName         string
	DbUser         string
	DbPass         string
	DbHost         string
	DbPort         int
	CollectionName string
	MinLength      int
	MaxLength      int
	DefaultLength  int
}

var Config UrlConfig

func init() {
	var err error
	Config, err = GetConfig()
	if err != nil {
		fmt.Println("Error during configuration initialization:", err)
		os.Exit(1)
	}
}
func GetConfig() (UrlConfig, error) {
	serverPort, err := getEnvWithError("SERVER_PORT")
	if err != nil {
		return UrlConfig{}, err
	}
	dbUser, err := getEnvWithError("DB_USER")
	if err != nil {
		return UrlConfig{}, err
	}

	dbPass, err := getEnvWithError("DB_PASS")
	if err != nil {
		return UrlConfig{}, err
	}

	dbPort, err := getIntEnvWithError("DB_PORT")
	if err != nil {
		return UrlConfig{}, err
	}
	dbHost := getEnv("DB_HOST", defaultDbHost)
	dbName := getEnv("DB_NAME", defaultDbName)
	collectionName := getEnv("COLLECTION_NAME", defaultCollectionName)
	minLength := getIntEnv("MIN_LENGTH", defaultMinLength)
	maxLength := getIntEnv("MAX_LENGTH", defaultMaxLength)
	defaultLength := getIntEnv("DEFAULT_LENGTH", defaultLength)

	config := UrlConfig{
		ServerPort:     serverPort,
		DbName:         dbName,
		DbUser:         dbUser,
		DbPass:         dbPass,
		DbHost:         dbHost,
		DbPort:         dbPort,
		CollectionName: collectionName,
		MinLength:      minLength,
		MaxLength:      maxLength,
		DefaultLength:  defaultLength,
	}
	return config, nil
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func getIntEnvWithError(key string) (int, error) {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return 0, fmt.Errorf("%s environment variable not set", key)
	}
	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("error converting %s to int: %s", key, err)
	}
	return valueInt, nil
}
func getEnvWithError(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", errors.New(key + " environment variable not set")
	}
	return value, nil
}

func getIntEnv(key string, fallback int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		return fallback
	}
	return valueInt
}
