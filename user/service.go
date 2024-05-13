package user

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	client *redis.Client

	key             string
	maxTokens       int
	tokensPerSecond int

	script string
}

func NewLimiter(client *redis.Client, key string, maxTokens, tokensPerSecond int) *Limiter {
	script := `
local bucket = KEYS[1]
local maxTokens = tonumber(ARGV[1])
local refillRate = tonumber(ARGV[2])
local period = 1  -- 一秒鐘
local lastRefillTimeKey = bucket .. ':lastRefillTime'
local tokensKey = bucket .. ':tokens'

local currentTime = tonumber(redis.call('time')[1])
local lastRefillTime = tonumber(redis.call('get', lastRefillTimeKey) or currentTime)
local elapsed = currentTime - lastRefillTime

local currentTokens = tonumber(redis.call('get', tokensKey) or maxTokens)

if elapsed >= period then
	local newTokens = math.min(maxTokens, refillRate * math.floor(elapsed / period))
	currentTokens = math.min(maxTokens, currentTokens + newTokens)
	redis.call('set', tokensKey, currentTokens)
	redis.call('set', lastRefillTimeKey, currentTime)
end

if currentTokens > 0 then
	redis.call('decr', tokensKey)
	return 1
else
	return 0
end
	`

	return &Limiter{
		client:          client,
		key:             key,
		maxTokens:       maxTokens,
		tokensPerSecond: tokensPerSecond,
		script:          script,
	}
}

func (limiter *Limiter) RequestToken(ctx context.Context) bool {
	if limiter == nil {
		return false
	}

	res, err := limiter.client.Eval(ctx, limiter.script, []string{limiter.key},
		limiter.maxTokens, limiter.tokensPerSecond).Result()
	if err != nil {
		fmt.Println("Error executing script:", err)
		return false
	}

	return res.(int64) == 1
}
