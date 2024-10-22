package config

import (
	"flag"

	"go.uber.org/zap/zapcore"
)

type RedisClientConfig struct {
	Address  *string
	Password *string
	Database *int
}

func (c *RedisClientConfig) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("ServerAddress", *c.Address)
	enc.AddString("Password", *c.Password)
	enc.AddInt("Database", *c.Database)
	return nil
}

var RedisClientCFG = &RedisClientConfig{
	Address:  flag.String("redis-address", "localhost:6379", "redis address to connect"),
	Password: flag.String("redis-password", "", "password for redis auth"),
	Database: flag.Int("redis-database", 0, "redis database to connect"),
}
