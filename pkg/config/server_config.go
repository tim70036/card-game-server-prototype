package config

import (
	"flag"

	"go.uber.org/zap/zapcore"
)

type ServerConfig struct {
	Port    *int
	WithTLS *bool
	TLSCert *string
	TLSKey  *string
}

func (c *ServerConfig) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("Port", *c.Port)
	enc.AddBool("WithTLS", *c.WithTLS)
	enc.AddString("TLSCert", *c.TLSCert)
	enc.AddString("TLSKey", *c.TLSKey)
	return nil
}

var ServerCFG = &ServerConfig{
	Port:    flag.Int("port", 8787, "port for server to listen"),
	WithTLS: tls,
	TLSCert: flag.String("tls-cert", "", "tls public cert"),
	TLSKey:  flag.String("tls-key", "", "tls private key"),
}
