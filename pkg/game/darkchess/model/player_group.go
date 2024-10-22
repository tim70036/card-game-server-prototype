package model

import (
	"card-game-server-prototype/pkg/core"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"go.uber.org/zap/zapcore"
)

type PlayerGroup struct {
	Data map[core.Uid]*Player
}

func ProvidePlayerGroup() *PlayerGroup {
	return &PlayerGroup{
		Data: make(map[core.Uid]*Player),
	}
}

func (g *PlayerGroup) ToProto() *gamegrpc.PlayerGroup {
	data := make(map[string]*gamegrpc.Player)

	for uid, p := range g.Data {
		data[uid.String()] = &gamegrpc.Player{
			Uid:        p.Uid.String(),
			ChessColor: p.Color,
		}
	}

	return &gamegrpc.PlayerGroup{
		Players: data,
	}
}

func (g *PlayerGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for uid, p := range g.Data {
		_ = enc.AddObject(uid.String(), p)
	}
	return nil
}
