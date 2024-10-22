package darkchess

import (
	"card-game-server-prototype/pkg/game/darkchess/actor"
	"card-game-server-prototype/pkg/game/darkchess/api"
	"card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/game/darkchess/server"
	"card-game-server-prototype/pkg/game/darkchess/service"
	"card-game-server-prototype/pkg/game/darkchess/state"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProvideGame,

	actor.ProviderSet,
	api.ProviderSet,
	state.ProviderSet,
	server.ProviderSet,
	model.ProviderSet,
	service.ProviderSet,
)
