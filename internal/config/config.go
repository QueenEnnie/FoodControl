package config

import "os"

type Config struct {
	Addr        string
	DatabaseURL string
}

func Load() Config {
	return Config{
		Addr:        getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://fcuser:fcpass@localhost:5433/fooddb"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
