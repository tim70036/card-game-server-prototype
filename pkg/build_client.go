//go:build wireinject
// +build wireinject

package main

import (
	"card-game-server-prototype/pkg/client"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"

	"github.com/google/wire"
)

func BuildClient() (*client.Client, error) {
	wire.Build(wire.NewSet(
		config.ProviderSet,
		util.ProviderSet,
		wire.Value([]zap.Field{}),

		client.ProvideClient,
		client.ProvideAuthInterceptor,
		client.ProvideConnectionServiceClient,
		client.ProvideMessageServiceClient,
	))
	return nil, nil
}
