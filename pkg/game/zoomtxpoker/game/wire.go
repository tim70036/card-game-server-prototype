package game

import (
	"card-game-server-prototype/pkg/game/zoomtxpoker/game/api"
	"card-game-server-prototype/pkg/game/zoomtxpoker/game/handler"
	"card-game-server-prototype/pkg/game/zoomtxpoker/game/service"
	"card-game-server-prototype/pkg/game/zoomtxpoker/game/state"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProvideZoomTXPokerGame,

	state.ProviderSet,
	service.ProviderSet,
	handler.ProviderSet,
	api.ProviderSet,
)
