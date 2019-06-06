package cache

import (
	"fmt"

	"github.com/go-redis/redis"
)

func NewCache(host string, port int) (*redis.Client, error) {
	client := redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	res := client.Ping()
	if res.Err() != nil {
		return nil, fmt.Errorf("failed to connect to redis : %v", res.Err())
	}
	return client, nil
}
