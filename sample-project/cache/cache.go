package cache

import (
	"github.com/go-redis/redis/v8"
	"os"
)

type Cache struct {
	redisClient *redis.Client
}

func NewCache(redisClient *redis.Client) *Cache {
	return &Cache{redisClient: redisClient}
}

func RedisInit() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // No password
		DB:       0,
	})
	return redisClient
}

func (c *Cache) GetCache(taskId string) (string, error) {
	return c.redisClient.Get(c.redisClient.Context(), taskId).Result()
}

func (c *Cache) SetCache(taskId string, task string) error {
	return c.redisClient.Set(c.redisClient.Context(), taskId, task, 0).Err()
}
