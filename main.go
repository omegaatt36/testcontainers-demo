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
	bucketKey := "globalTokenBucket"
	maxTokens := 10
	refillRate := 10

	limiter := user.NewLimiter(client, bucketKey, maxTokens, refillRate)

	ctx := context.Background()

	for i := 0; i < 20; i++ {
		if limiter.RequestToken(ctx) {
			fmt.Printf("Request %d allowed.\n", i+1)
		} else {
			fmt.Printf("Request %d denied.\n", i+1)
		}
		time.Sleep(100 * time.Microsecond)
	}

	time.Sleep(time.Second)
	fmt.Println("after 1 second")
	for i := 0; i < 20; i++ {
		if limiter.RequestToken(ctx) {
			fmt.Printf("Request %d allowed.\n", i+1)
		} else {
			fmt.Printf("Request %d denied.\n", i+1)
		}
		time.Sleep(100 * time.Microsecond)
	}
}
