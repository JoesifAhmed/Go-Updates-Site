package models

import (
	"github.com/go-redis/redis"
)

var Client *redis.Client

func InitClient() {
	Client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
