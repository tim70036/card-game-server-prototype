//go:build wireinject
// +build wireinject

package main

import (
	"card-game-server-prototype/pkg/common"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/util"

	"github.com/google/wire"
	"go.uber.org/zap"
)

func BuildGameMaker() (*GameMaker, error) {
	wire.Build(wire.NewSet(
		config.ProviderSet,
		core.ProviderSet,
		common.ProviderSet,
		util.ProviderSet,
		wire.Value([]zap.Field{}),
		ProvideGameMaker,
		ProvideAgones,
	))
	return nil, nil
}
