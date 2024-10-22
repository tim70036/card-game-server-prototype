package game

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/zoomtxpoker/game/handler"
	"card-game-server-prototype/pkg/game/zoomtxpoker/game/state"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
)

type ZoomTXPokerGame struct {
	game   core.Game
	msgBus core.MsgBus
	logger *zap.Logger
}

func ProvideZoomTXPokerGame(
	game core.Game,
	_ state.Initiator,
	_ handler.Initiator,

	msgBus core.MsgBus,
	loggerFactory *util.LoggerFactory,
) *ZoomTXPokerGame {
	return &ZoomTXPokerGame{
		game:   game,
		msgBus: msgBus,
		logger: loggerFactory.Create("ZoomTXPokerGame"),
	}
}

func (zoomTXPokerGame *ZoomTXPokerGame) MsgBus() core.MsgBus {
	return zoomTXPokerGame.msgBus
}

func (zoomTXPokerGame *ZoomTXPokerGame) GameNotifier() core.GameNotifier {
	return zoomTXPokerGame.game
}

func (zoomTXPokerGame *ZoomTXPokerGame) Run() {
	zoomTXPokerGame.logger.Debug("start")
	go zoomTXPokerGame.game.Run()
	<-zoomTXPokerGame.game.WaitShutdown()
	zoomTXPokerGame.logger.Info("game has shutdown, exiting")
}
