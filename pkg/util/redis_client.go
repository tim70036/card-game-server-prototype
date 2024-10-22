package util

import (
	"card-game-server-prototype/pkg/config"
	"github.com/redis/go-redis/v9"
)

func ProvideRedisClient(
	redisClientCFG *config.RedisClientConfig,
) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     *redisClientCFG.Address,
		Password: *redisClientCFG.Password,
		DB:       *redisClientCFG.Database,
	})
}
