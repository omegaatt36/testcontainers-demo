package cache

import "github.com/redis/go-redis/v9"

// NewRedisClient 創建一個 redis 客戶端的實例
func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return client
}
