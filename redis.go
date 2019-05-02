package goutils

import (
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
)

var _client *RClient
var connection_options *redis.Options

type RClient struct {
	client *redis.Client
}

func Redis() *RClient {
	if _client == nil {
		if GetEnvVariable("ENV", "development") == "test" {
			connection_options = TestRedisOptions()
		} else {
			connection_options = RedisOptions()
		}
		client := redis.NewClient(connection_options)

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

func RedisOptions() *redis.Options {
	return &redis.Options{
		Addr:     GetEnvVariable("REDIS_ADDRESS", "localhost:6379"),
		Password: GetEnvVariable("REDIS_PASSWORD", ""),
		DB:       0,
	}
}

// newTestRedis returns a redis.Cmdable.
func TestRedisOptions() *redis.Options {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	return &redis.Options{Addr: mr.Addr()}
}
