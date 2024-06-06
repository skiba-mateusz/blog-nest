package config

import (
	"log"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	ListenAddr  string
	DBName      string
	DBPassword  string
	DBUser 		string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		ListenAddr: getEnv("LISTEN_ADDR", "localhost:8080"),
		DBName: getEnv("DB_NAME", "blognest"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := syscall.Getenv(key); ok {
		return val
	}

	if fallback == "" {
		log.Fatalf("environment variable %s not set and no fallback provided", key)
	}

	return fallback
}