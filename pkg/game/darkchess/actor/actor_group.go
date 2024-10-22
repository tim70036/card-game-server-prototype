package actor

import (
	"card-game-server-prototype/pkg/core"
	"go.uber.org/zap/zapcore"
)

type Group struct {
	Data map[core.Uid]Actor
}

func ProvideActorGroup() *Group {
	return &Group{
		Data: make(map[core.Uid]Actor),
	}
}

func (g *Group) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for uid, actor := range g.Data {
		_ = enc.AddObject(uid.String(), actor)
	}
	return nil
}
