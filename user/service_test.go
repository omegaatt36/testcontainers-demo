package user_test

import (
	"context"
	"testing"

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

func (s *LimiterTestSuite) testLimiter(limiter *user.Limiter) {
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		s.True(limiter.RequestToken(ctx))
	}
}

func (s *LimiterTestSuite) TestLimiterWithRealConn() {
	client := cache.NewRedisClient()

	limiter := user.NewLimiter(client, "globalTokenBucket", 10, 10)

	s.testLimiter(limiter)
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

	limiter := user.NewLimiter(client, "globalTokenBucket2", 10, 10)

	s.testLimiter(limiter)
}

func TestLimiter(t *testing.T) {
	suite.Run(t, new(LimiterTestSuite))
}
