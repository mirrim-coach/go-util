package goutils

import (
	"github.com/go-redis/redis"
	"time"
)

var _client *RClient

type RClient struct {
	client *redis.Client
}

func Redis() *RClient {
	if _client == nil {
		client := redis.NewClient(&redis.Options{
			Addr:     GetEnvVariable("REDIS_ADDRESS", "localhost:6379"),
			Password: GetEnvVariable("REDIS_PASSWORD", ""),
			DB:       0,
		})

		_client = &RClient{
			client: client,
		}

		_, err := client.Ping().Result()

		if err != nil {
			Logger().Error(err)
		}
	}
	return _client
}

func (r *RClient) Set(key string, value string, time time.Duration) error {
	return r.client.Set(key, value, time).Err()
}

func (r *RClient) Get(key string) (string, error) {
	return r.client.Get(key).Result()
}

func (r *RClient) Del(key string) error {
	return r.client.Del(key).Err()
}

func (r *RClient) Publish(channel string, message string) error {
	return r.client.Publish(channel, message).Err()
}

func (r *RClient) Expire(key string, duration time.Duration) error {
	return r.client.Expire(key, duration).Err()
}

func (r *RClient) ExpireAt(key string, timeout time.Time) error {
	return r.client.ExpireAt(key, timeout).Err()
}
