package pool

import (
	"card-game-server-prototype/pkg/game/txpoker/server"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/api"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/handler"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/model"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/service"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/state"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProvideZoomTXPokerPool,

	state.ProviderSet,
	server.ProviderSet,
	model.ProviderSet,
	service.ProviderSet,
	handler.ProviderSet,
	api.ProviderSet,
)
