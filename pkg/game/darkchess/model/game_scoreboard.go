package model

import (
	"card-game-server-prototype/pkg/common/type/gamemode"
	"card-game-server-prototype/pkg/core"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"go.uber.org/zap/zapcore"
)

type GameScoreboard struct {
	Data     map[core.Uid]*GameScore
	WinnerId core.Uid
	IsDraw   bool
}

func ProvideGameScoreboard() *GameScoreboard {
	return &GameScoreboard{
		Data: make(map[core.Uid]*GameScore),
	}
}

func (s *GameScoreboard) ToProto(mode gamemode.GameMode) *gamegrpc.GameScoreboard {
	var scores = make([]*gamegrpc.GameScore, 0)
	for uid := range s.Data {
		scores = append(scores, s.Data[uid].ToProto(mode))
	}

	return &gamegrpc.GameScoreboard{
		Scores:    scores,
		WinnerUid: s.WinnerId.String(),
		IsDraw:    s.IsDraw,
	}
}

func (s *GameScoreboard) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for uid, v := range s.Data {
		_ = enc.AddObject(uid.String(), v)
	}
	enc.AddString("WinnerId", s.WinnerId.String())
	enc.AddBool("IsDraw", s.IsDraw)
	return nil
}
