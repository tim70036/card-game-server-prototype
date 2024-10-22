package core

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	wire.Bind(new(GameNotifier), new(Game)),
	wire.Bind(new(GameController), new(Game)),
	ProvideStateFactory,
	ProvideMsgBus,
	ProvideGame,
	ProvideRouter,
)
