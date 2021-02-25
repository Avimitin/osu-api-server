package database

import "github.com/go-redis/redis/v8"

type RedisDataStore struct {
	db *redis.Client
}

func NewRedisDataStore() *RedisDataStore {
	rds := new(RedisDataStore)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "",
		Password: "",
		DB:       0,
	})
	rds.db = rdb
	return rds
}
