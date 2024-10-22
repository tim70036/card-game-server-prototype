//go:build wireinject
// +build wireinject

package commontest

import (
	"card-game-server-prototype/pkg/common"
	commonapi "card-game-server-prototype/pkg/common/api"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/util"

	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	config.ProviderSet,
	core.ProviderSet,
	common.ProviderSet,
	util.ProviderSet,
)

func BuildRegistry() (*Registry, error) {
	wire.Build(wire.NewSet(
		wire.Bind(new(commonapi.UserAPI), new(*commonapi.BaseUserAPI)),
		wire.Bind(new(commonapi.RoomAPI), new(*commonapi.BaseRoomAPI)),

		ProvideRegistry,
		ProvideServiceRegistry,
		ProvideAPIRegistry,
		providerSet,
	))
	return nil, nil
}
