package handler

import "card-game-server-prototype/pkg/core"

type Initiator interface{}

func Init(
	game core.Game,
	connectionHandler *ConnectionHandler,
	requestHandler *RequestHandler,
) Initiator {
	game.ConfigHandler(connectionHandler)
	game.ConfigHandler(requestHandler)
	return nil
}
