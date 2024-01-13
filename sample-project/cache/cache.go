package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"os"
	"sample-project/models"
	"time"
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

func (c *Cache) GetCache(taskId string) (*models.Task, error) {
	msg := c.redisClient.Get(c.redisClient.Context(), taskId).Val()
	task := &models.Task{}
	err := json.Unmarshal([]byte(msg), task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (c *Cache) SetCache(taskId string, task *models.Task) error {
	t, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return c.redisClient.Set(context.Background(), taskId, t, 5*time.Minute).Err()
}

func (c *Cache) DelCache(taskId string) error {
	return c.redisClient.Del(context.Background(), taskId).Err()
}
