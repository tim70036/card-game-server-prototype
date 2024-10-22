package config

import (
	"flag"

	"go.uber.org/zap/zapcore"
)

type ClientConfig struct {
	Address *string
	Uid     *string
	WithTLS *bool
}

func (c *ClientConfig) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("ServerAddress", *c.Address)
	enc.AddString("Uid", *c.Uid)
	enc.AddBool("WithTLS", *c.WithTLS)
	return nil
}

var ClientCFG = &ClientConfig{
	Address: flag.String("address", "localhost:8787", "server address to connect"),
	Uid:     flag.String("uid", "1", "uid of client"),
	WithTLS: tls,
}
