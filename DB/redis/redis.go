package redis

import (
	"log"

	"github.com/go-redis/redis/v7"
)

var client *redis.Client

func NewClient() *redis.Client {
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		pong, err := client.Ping().Result()
		if err == nil && pong == "PONG" {
			log.Println("redis连接成功")
		}
	}
	return client
	// Output: PONG <nil>
}
