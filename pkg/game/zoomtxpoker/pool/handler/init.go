package handler

import "card-game-server-prototype/pkg/core"

type Initiator interface{}

func Init(
	game core.Game,
	connectionHandler *ConnectionHandler,
	requestHandler *RequestHandler,
	router *core.Router,
) Initiator {
	game.ConfigHandler(connectionHandler)
	game.ConfigHandler(requestHandler)
	game.ConfigHandler(router)
	return nil
}
