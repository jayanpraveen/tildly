package datastore

import (
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

func DialRedisCache() *cache.Cache {

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	rch := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, time.Hour),
	})

	return rch
}
