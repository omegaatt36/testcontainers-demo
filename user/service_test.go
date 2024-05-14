package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/omegaatt36/limiter/cache"
	"github.com/omegaatt36/limiter/user"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type LimiterTestSuite struct {
	suite.Suite
}

func (s *LimiterTestSuite) TestLimiterWithRealConn() {
	client := cache.NewRedisClient()

	limiter := user.NewLimiter(client, 10, time.Second*5, time.Second)

	ctx := context.Background()

	for i := range 9 {
		s.NoError(limiter.AllowRequest(ctx, "55688", 1), "request %d should be allowed", i+1)
		time.Sleep(time.Millisecond)
	}

	s.Error(limiter.AllowRequest(ctx, "55688", 1), "request  should be denied")
}

func (s *LimiterTestSuite) TestLimiterWithTestContainers() {
	ctx := context.Background()

	request := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	s.NoError(err)

	endpoint, err := container.Endpoint(ctx, "")
	s.NoError(err)

	client := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})

	limiter := user.NewLimiter(client, 10, time.Second*5, time.Second)

	for i := range 9 {
		s.NoError(limiter.AllowRequest(ctx, "55688", 1), "request %d should be allowed", i+1)
		time.Sleep(time.Millisecond)
	}

	s.Error(limiter.AllowRequest(ctx, "55688", 1), "request should be denied")
}

func TestLimiter(t *testing.T) {
	suite.Run(t, new(LimiterTestSuite))
}
