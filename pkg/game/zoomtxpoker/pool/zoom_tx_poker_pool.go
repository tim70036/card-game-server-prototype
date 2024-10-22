package pool

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/server"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/handler"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/state"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
)

type ZoomTXPokerPool struct {
	game   core.Game
	logger *zap.Logger
}

func ProvideZoomTXPokerPool(
	_ state.Initiator,
	_ server.Initiator,
	_ handler.Initiator,
	game core.Game,
	loggerFactory *util.LoggerFactory,
) *ZoomTXPokerPool {
	zoomTXPokerPool := &ZoomTXPokerPool{
		game:   game,
		logger: loggerFactory.Create("ZoomTXPokerPool"),
	}

	return zoomTXPokerPool
}

func (zoomTXPokerPool *ZoomTXPokerPool) Run() {
	zoomTXPokerPool.logger.Debug("start")
	go zoomTXPokerPool.game.Run()

	<-zoomTXPokerPool.game.WaitShutdown()
	zoomTXPokerPool.logger.Info("game has shutdown, exiting")
}
