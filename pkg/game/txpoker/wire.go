package txpoker

import (
	"card-game-server-prototype/pkg/game/txpoker/actor"
	"card-game-server-prototype/pkg/game/txpoker/api"
	"card-game-server-prototype/pkg/game/txpoker/handler"
	"card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/server"
	"card-game-server-prototype/pkg/game/txpoker/service"
	"card-game-server-prototype/pkg/game/txpoker/state"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProvideTXPoker,

	actor.ProviderSet,
	api.ProviderSet,
	state.ProviderSet,
	server.ProviderSet,
	model.ProviderSet,
	service.ProviderSet,
	handler.ProviderSet,
)
