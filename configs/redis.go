package configs

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic("Impossible de se connecter à Redis: " + err.Error())
	}

	log.Println("Redis connecté............✅")
}

func SetCache(key string, value any, ttl time.Duration) error {
	return RedisClient.Set(Ctx, key, value, ttl).Err()
}

func GetCache(key string) (string, error) {
	return RedisClient.Get(Ctx, key).Result()
}

func DeleteCache(key string) error {
	return RedisClient.Del(Ctx, key).Err()
}
