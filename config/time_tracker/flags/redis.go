package flags

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisFlags struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
}

func (f RedisFlags) Init(ctx context.Context) (*redis.Client, error) {
	cfg := redis.Options{Addr: f.Addr, Password: f.Password, DB: 0}
	redisClient := redis.NewClient(&cfg)

	err := redisClient.Ping(ctx).Err()

	if err != nil {
		return nil, err
	}
	return redisClient, nil
}
