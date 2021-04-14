package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	gobit "github.com/pot-code/gobit/pkg"
)

// RedisClient .
type RedisClient struct {
	conn *redis.Client
}

// interface assertion
var _ CacheDB = &RedisClient{}

// NewRedisClient create an redis client
func NewRedisClient(host string, port int, password string) *RedisClient {
	conn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
	})
	return &RedisClient{
		conn: conn,
	}
}

func (rdb *RedisClient) SetExp(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rdb.conn.Set(ctx, key, value, expiration).Err()
}

func (rdb *RedisClient) Set(ctx context.Context, key string, value interface{}) error {
	return rdb.conn.Set(ctx, key, value, 0).Err()
}

func (rdb *RedisClient) Get(ctx context.Context, key string) (string, error) {
	cmd := rdb.conn.Get(ctx, key)
	v, err := cmd.Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return v, err
}

func (rdb *RedisClient) Delete(ctx context.Context, key ...string) (bool, error) {
	cmd := rdb.conn.Del(ctx, key...)
	code, err := cmd.Result()
	return code == 1, err
}

func (rdb *RedisClient) Incr(ctx context.Context, key string) (bool, error) {
	cmd := rdb.conn.Incr(ctx, key)
	code, err := cmd.Result()
	return code == 1, err
}

func (rdb *RedisClient) IncrBy(ctx context.Context, key string, val int64) (bool, error) {
	cmd := rdb.conn.IncrBy(ctx, key, val)
	code, err := cmd.Result()
	return code == 1, err
}

func (rdb *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	cmd := rdb.conn.Exists(ctx, key)
	code, err := cmd.Result()

	return code == 1, err
}

func (rdb *RedisClient) Ping(ctx context.Context) error {
	cmd := rdb.conn.Ping(ctx)
	if r, err := cmd.Result(); err != nil {
		return err
	} else if r != "PONG" {
		return gobit.ErrInternalError
	}
	return nil
}

func (rdb *RedisClient) Close(ctx context.Context) error {
	return rdb.conn.Close()
}
