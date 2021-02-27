package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
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
	sc := rds.db.Set(ctx, u.Username+":recent", jsonByte, 0)
	if sc.Err() != nil {
		return fmt.Errorf("set %s into redis: %v", u.Username, sc.Err())
	}
	return nil
}

// CheckHealth checking connection
func (rds *RedisDataStore) CheckHealth() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	sc := rds.db.Ping(ctx)
	if e := sc.Err(); e != nil {
		return fmt.Errorf("checkhealth: ping to %v: %v", "redis server", e)
	}
	return nil
}

func (rds *RedisDataStore) getPlayerWithDate(name string, date string) (*User, error) {
	val, err := rds.getStr(parseSchema(name, date))
	if err != nil {
		return nil, err
	}
	return rds.parseUser(strings.NewReader(val))
}

// GetPlayer get a user recent score by specific given name
func (rds *RedisDataStore) GetPlayer(name string) (*User, error) {
	return rds.getPlayerWithDate(name, "recent")
}

// GetPlayerOld get a user old score by specific given name
func (rds *RedisDataStore) GetPlayerOld(name string) (*User, error) {
	return rds.getPlayerWithDate(name, "old")
}

func parseSchema(key string, args ...string) string {
	return key + ":" + strings.Join(args, ":")
}

func (rds *RedisDataStore) getStr(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	val, err := rds.db.Get(ctx, key).Result()
	switch {
	case err == redis.Nil:
		return "", fmt.Errorf("%q not found", key)
	case err != nil:
		return "", fmt.Errorf("get %q failed: %v", key, err)
	case val == "":
		return "", errors.New("get nil value")
	}
	return val, nil
}

func (rds *RedisDataStore) parseUser(data io.Reader) (*User, error) {
	jsonbyte, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, fmt.Errorf("parseUser: read data: %v", err)
	}
	var u *User
	err = json.Unmarshal(jsonbyte, &u)
	if err != nil {
		return nil, fmt.Errorf("parse %q: %v", jsonbyte, err)
	}
	return u, nil
}
