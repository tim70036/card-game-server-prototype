package config

import (
	"flag"

	"go.uber.org/zap/zapcore"
)

type APIConfig struct {
	MainServerHost   *string
	MainServerAPIKey *string
}

func (c *APIConfig) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return nil
}

var APICFG = &APIConfig{
	MainServerHost:   flag.String("main-server-host", "main-server-staging.game-soul-swe.com", "main server host to send request to"),
	MainServerAPIKey: flag.String("main-server-api-key", "base64base641234a1231231231123123asdz", "main server api key for sending request"),
}
