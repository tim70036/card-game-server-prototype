package actor

import (
	"card-game-server-prototype/pkg/core"
)

type TXPokerActor interface {
	core.Actor
	DecideAction() (core.ActorRequestList, error)
}
