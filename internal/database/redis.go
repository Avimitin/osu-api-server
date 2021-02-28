package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
// with redis implement. dsn should be format like this:
// redis://<user>:<password>@<host>:<port>/<db_number>
// else will get panic
func NewRedisDataStore(dsn string) *RedisDataStore {
	rds := new(RedisDataStore)
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		panic(fmt.Sprintf("parse %s: %v", dsn, err))
	}
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
			os.Exit(1)
		}
	}()

	rdb := redis.NewClient(opt)
	log.Printf("redis connecting to %s", opt.Addr)
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
	return rds.set(parseSchema(u.Username, "recent"), jsonByte)
}

// CheckHealth checking connection
func (rds *RedisDataStore) CheckHealth() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	sc := rds.db.Ping(ctx)
	if e := sc.Err(); e != nil {
		return fmt.Errorf("checkhealth: ping to %v: %v", "redis server", e)
	}
	return nil
}

// GetPlayer get a user recent score by specific given name
func (rds *RedisDataStore) GetPlayer(name string) (*User, error) {
	return rds.getPlayerWithDate(name, "recent")
}

// GetPlayerOld get a user old score by specific given name
func (rds *RedisDataStore) GetPlayerOld(name string) (*User, error) {
	return rds.getPlayerWithDate(name, "old")
}

func (rds *RedisDataStore) Update(u User) error {
	if !rds.isUserExist(u.Username) {
		return errors.New("user to update not found")
	}

	jsonData, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("parse user %q data: %v", u.Username, err)
	}
	return rds.set(parseSchema(u.Username, "recent"), jsonData)
}

func (rds *RedisDataStore) UpdateOld(u User) error {
	jsonData, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("parse user %q : %v", u.Username, err)
	}

	return rds.set(parseSchema(u.Username, "old"), jsonData)
}

func (rds *RedisDataStore) set(key string, val interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	cmd := rds.db.Set(ctx, key, val, 0)
	if cmd.Err() != nil {
		return fmt.Errorf("set %q: %v", key, cmd.Err())
	}
	return nil
}

func (rds *RedisDataStore) isUserExist(name string) bool {
	user, err := rds.GetPlayer(name)
	if err != nil {
		return false
	}
	if user == nil {
		return false
	}
	return true
}

func parseSchema(key string, args ...string) string {
	return key + ":" + strings.Join(args, ":")
}

func (rds *RedisDataStore) getPlayerWithDate(name string, date string) (*User, error) {
	val, err := rds.getStr(parseSchema(name, date))
	if err != nil {
		return nil, err
	}
	return rds.parseUser(strings.NewReader(val))
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
