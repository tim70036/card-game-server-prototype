package common

import (
	"card-game-server-prototype/pkg/common/api"
	"card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/common/server"
	"card-game-server-prototype/pkg/common/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	api.ProviderSet,
	model.ProviderSet,
	server.ProviderSet,
	service.ProviderSet,
)
