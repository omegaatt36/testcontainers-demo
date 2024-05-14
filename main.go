package main

import (
	"context"
	"fmt"
	"time"

	"github.com/omegaatt36/limiter/cache"
	"github.com/omegaatt36/limiter/user"
)

func main() {
	client := cache.NewRedisClient()
	maxTokens := 10

	limiter := user.NewLimiter(client, maxTokens, time.Second*5, time.Second)

	ctx := context.Background()

	for i := 0; i < 20; i++ {
		if err := limiter.AllowRequest(ctx, "55688", 1); err == nil {
			fmt.Printf("Request %02d allowed.\n", i+1)
		} else {
			fmt.Printf("Request %02d denied.\n", i+1)
		}
		time.Sleep(time.Millisecond)
	}
}
