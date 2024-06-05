package config

import "syscall"

type Config struct {
	ListenAddr string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		ListenAddr: getEnv("LISTEN_ADDR", "localhost:8080"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := syscall.Getenv(key); ok {
		return val
	}

	return fallback
}