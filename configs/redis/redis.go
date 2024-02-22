package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"task-one/helpers"
	"time"
)

type Redis interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (string, error)
}

type RedisClient struct {
	rdb *redis.Client
}

func InitRedis() *RedisClient {
	env := helpers.GetConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     env.Redis.Host,
		Password: env.Redis.Password,
		DB:       env.Redis.Db,
	})

	res := client.Ping(context.Background())
	log.Println(res)

	return &RedisClient{rdb: client}
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.rdb.Set(ctx, key, data, 10*time.Minute).Err()
	return err

}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, err

}
