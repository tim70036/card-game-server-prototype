package actor

import (
	"card-game-server-prototype/pkg/core"
)

type Actor interface {
	core.Actor
	DecidePick() (core.ActorRequestList, error)
	DecideAction() (core.ActorRequestList, error)
	DecideSkipScoreboard() (core.ActorRequestList, error)
}
