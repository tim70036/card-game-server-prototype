package actor

import (
	"card-game-server-prototype/pkg/core"

	"go.uber.org/zap/zapcore"
)

type ActorGroup struct {
	Data map[core.Uid]TXPokerActor
}

func ProvideActorGroup() *ActorGroup {
	return &ActorGroup{
		Data: make(map[core.Uid]TXPokerActor),
	}
}

func (g *ActorGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for uid, actor := range g.Data {
		enc.AddObject(uid.String(), actor)
	}
	return nil
}
