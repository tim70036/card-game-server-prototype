package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/type/event"
	"go.uber.org/zap/zapcore"
)

type EventGroup struct {
	Data map[core.Uid]event.List
}

func ProvideEventGroup() *EventGroup {
	return &EventGroup{
		Data: make(map[core.Uid]event.List),
	}
}

func (g *EventGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for uid, events := range g.Data {
		_ = enc.AddArray(uid.String(), events)
	}
	return nil
}
