package config

import (
	"log"
	"strconv"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	ListenAddr  	string
	DBName      	string
	DBPassword  	string
	DBUser 			string
	JWTExpiration 	int64
	JWTSecret		string
	AWSRegion		string
	S3BucketName	string
	S3BaseUrl		string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		ListenAddr: getEnv("LISTEN_ADDR", "localhost:8080"),
		DBName: getEnv("DB_NAME", "blognest"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		JWTSecret: getEnv("JWT_SECRET", ""),
		JWTExpiration: getEnvAsInt("JWT_EXPIRATION", 3600 * 24 * 7),
		AWSRegion: getEnv("AWS_REGION", ""),
		S3BucketName: getEnv("S3_BUCKET_NAME", ""),
		S3BaseUrl: getEnv("S3_BASE_URL", ""),
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

func getEnvAsInt(key string, fallback int64) int64 {
	if val, ok := syscall.Getenv(key); ok {
		val, _ := strconv.ParseInt(val, 10, 64)
		return val 
	}

	if fallback == 0 {
		log.Fatalf("environment variable %s not set and no fallback provided", key)
	}

	return fallback
}