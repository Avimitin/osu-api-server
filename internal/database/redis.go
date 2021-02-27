package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

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
		Addr:     getEnvWithFallBack("redis_host", "localhost:6379"),
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

func (rds *RedisDataStore) CheckHealth() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	sc := rds.db.Ping(ctx)
	if e := sc.Err(); e != nil {
		return fmt.Errorf("checkhealth: ping to %v: %v", "redis server", e)
	}
	return nil
}

func (rds *RedisDataStore) GetPlayer(name string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	val, err := rds.db.Get(ctx, name).Result()
	switch {
	case err == redis.Nil:
		return nil, errors.New("user not found")
	case err != nil:
		return nil, fmt.Errorf("get %q failed: %v", name, err)
	case val == "":
		return nil, errors.New("get nil value")
	}

	var u *User
	err = json.Unmarshal([]byte(val), &u)
	if err != nil {
		return nil, fmt.Errorf("parse %q: %v", val, err)
	}
	return u, nil
}
