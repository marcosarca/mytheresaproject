package config

import (
	"os"
)

const (
	dbFile = "DB_FILE"
	port   = "HTTP_PORT"
)

type Config struct {
	DbFile string
	Port   string
}

func New() Config {
	return Config{
		DbFile: GetEnvString(dbFile, ""),
		Port:   GetEnvString(port, "8080"),
	}
}

func GetEnvString(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultValue
}
