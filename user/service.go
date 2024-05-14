package user

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	client *redis.Client

	limit         int
	limitPeriod   time.Duration // 1 hour for limitPeriod
	counterWindow time.Duration // 1 minute for example, 1/60 of the period
}

// 新建一個滑窗限流器
func NewLimiter(client *redis.Client, limit int, period, expiry time.Duration) *Limiter {
	return &Limiter{
		client: client,

		limit:         limit,
		limitPeriod:   period,
		counterWindow: expiry,
	}
}

func (r *Limiter) AllowRequest(ctx context.Context, key string, incr int) error {
	now := time.Now()
	timestamp := fmt.Sprint(now.Truncate(r.counterWindow).Unix())

	val, err := r.client.HIncrBy(ctx, key, timestamp, int64(incr)).Result()
	if err != nil {
		return err
	}

	if val >= int64(r.limit) {
		return ErrRateLimitExceeded(0, r.limit, r.limitPeriod, now.Add(r.limitPeriod))
	}

	r.client.Expire(ctx, key, r.limitPeriod)

	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return err
	}

	threshold := fmt.Sprint(now.Add(-r.limitPeriod).Unix())

	total := 0
	for k, v := range result {
		if k > threshold {
			i, _ := strconv.Atoi(v)
			total += i
		} else {
			r.client.HDel(ctx, key, k)
		}
	}

	if total >= int(r.limit) {
		return ErrRateLimitExceeded(0, r.limit, r.limitPeriod, now.Add(r.limitPeriod))
	}

	return nil
}

type RateLimitExceeded struct {
	Remaining int
	Limit     int
	Period    time.Duration
	Reset     time.Time
}

func ErrRateLimitExceeded(remaining int, limit int, period time.Duration, reset time.Time) error {
	return RateLimitExceeded{
		Remaining: remaining,
		Limit:     limit,
		Period:    period,
		Reset:     reset,
	}
}

func (e RateLimitExceeded) Error() string {
	return fmt.Sprintf(
		"rate limit of %d per %v has been exceeded and resets at %v",
		e.Limit, e.Period, e.Reset)
}
