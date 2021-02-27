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

// AddPlayer add given user into redis database.
// With given user's username as key and it's json
// bytes as value.
func (rds *RedisDataStore) AddPlayer(u User) error {
	jsonByte, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("marshal %v:%s", u.Username, err)
	}
	ctx := context.Background()
	sc := rds.db.Set(ctx, u.Username, jsonByte, 0)
	if sc.Err() != nil {
		return fmt.Errorf("set %s into redis: %v", u.Username, sc.Err())
	}
	return nil
}
