package utils

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	_ = godotenv.Load(".env")
}

func GetEnv(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	return value
}

func GetIntEnv(key string, def int) int {
	value := os.Getenv(key)
	if value == "" {
		return def
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return def
	}
	return intValue
}
