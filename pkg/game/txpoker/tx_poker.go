package txpoker

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/handler"
	"card-game-server-prototype/pkg/game/txpoker/server"
	"card-game-server-prototype/pkg/game/txpoker/state"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
)

type TXPoker struct {
	game   core.Game
	logger *zap.Logger
}

func ProvideTXPoker(
	_ state.Initiator,
	_ server.Initiator,
	_ handler.Initiator,

	game core.Game,

	loggerFactory *util.LoggerFactory,
) *TXPoker {
	txPoker := &TXPoker{
		game:   game,
		logger: loggerFactory.Create("TXPoker"),
	}

	return txPoker
}

func (txPoker *TXPoker) Run() {
	txPoker.logger.Debug("start")

	go txPoker.game.Run()

	<-txPoker.game.WaitShutdown()
	txPoker.logger.Info("game has shutdown, exiting")
}
