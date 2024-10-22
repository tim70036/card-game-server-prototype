package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"

	"go.uber.org/zap/zapcore"
)

type StatsCacheGroup struct {
	Data map[core.Uid]*Stats
}

func ProvideStatsCacheGroup() *StatsCacheGroup {
	return &StatsCacheGroup{
		Data: make(map[core.Uid]*Stats),
	}
}

func (g *StatsCacheGroup) ToProto() *txpokergrpc.StatsCacheGroup {
	msg := &txpokergrpc.StatsCacheGroup{Stats: make(map[string]*txpokergrpc.Stats)}
	for uid, stats := range g.Data {
		msg.Stats[uid.String()] = stats.ToProto()
	}
	return msg
}

func (g *StatsCacheGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for uid, stats := range g.Data {
		enc.AddObject(uid.String(), stats)
	}
	return nil
}
