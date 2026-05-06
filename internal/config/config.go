package config

import (
	"os"
)

func GetDBConnStr() string {
	if val := os.Getenv("DATABASE_URL"); val != "" {
		return val
	}
	return "host=postgres port=5432 user=user password=password dbname=habr_db sslmode=disable"
}

func GetKafkaBroker() string {
	return getEnv("KAFKA_BROKER", "kafka:9092")
}

func GetRedisAddr() string {
	return getEnv("REDIS_ADDR", "redis:6379")
}
func GetAuthHost() string {
	return getEnv("AUTH_HOST", "http://auth:8081")
}

func GetServerPort(defaultPort string) string {
	return getEnv("PORT", defaultPort)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
