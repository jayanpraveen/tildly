package datastore

import (
	"log"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

func DialRedisClient() *cache.Cache {

	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	log.Print("connected to ", rdb.Options().Addr)

	rch := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, time.Hour),
	})

	return rch
}
