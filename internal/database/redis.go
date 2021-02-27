package database

import (
	"os"

	"github.com/go-redis/redis/v8"
)

type RedisDataStore struct {
	db *redis.Client
}

func getEnvWithFallBack(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// NewRedisDataStore return a initialized PlayerDataStore
// with redis implement
func NewRedisDataStore() *RedisDataStore {
	rds := new(RedisDataStore)

	rdb := redis.NewClient(&redis.Options{
		Addr:     getEnvWithFallBack("redis_host", "localhost:6479"),
		Password: "",
		DB:       0,
	})
	rds.db = rdb
	return rds
}
