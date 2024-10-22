package model

import (
	"card-game-server-prototype/pkg/core"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"go.uber.org/zap/zapcore"
)

type RoundScoreboard struct {
	Data     map[core.Uid]*RoundScore
	WinnerId core.Uid
	IsDraw   bool
}

func ProvideRoundScoreboard() *RoundScoreboard {
	return &RoundScoreboard{
		Data: map[core.Uid]*RoundScore{},
	}
}

func (s *RoundScoreboard) ToProto() *gamegrpc.RoundScoreboard {
	var scores []*gamegrpc.RoundScore
	for _, v := range s.Data {
		scores = append(scores, v.ToProto())
	}

	return &gamegrpc.RoundScoreboard{
		Scores:    scores,
		WinnerUid: s.WinnerId.String(),
		IsDraw:    s.IsDraw,
	}
}

func (s *RoundScoreboard) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for uid, v := range s.Data {
		_ = enc.AddObject(uid.String(), v)
	}
	enc.AddBool("IsDraw", s.IsDraw)
	enc.AddString("WinnerId", s.WinnerId.String())
	return nil
}
