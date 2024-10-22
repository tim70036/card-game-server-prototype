package model

import (
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"
)

type GameSetting struct {
	GameMetaUid string
	SmallBlind  int
	BigBlind    int

	TurnDuration    time.Duration
	MinEnterLimitBB int
	MaxEnterLimitBB int
	MaxWaterLimitBB int
	WaterPct        int
	TableSize       int
	MaxUserAmount   int

	LeastGamePlayerAmount int // club 可能會用，main server 不好改，先留著。

}

func (s *GameSetting) ToProto() *txpokergrpc.GameSetting {
	return &txpokergrpc.GameSetting{
		GameMetaUid:     s.GameMetaUid,
		SmallBlind:      int32(s.SmallBlind),
		BigBlind:        int32(s.BigBlind),
		TurnDuration:    durationpb.New(s.TurnDuration),
		MinEnterLimitBb: int32(s.MinEnterLimitBB),
		MaxEnterLimitBb: int32(s.MaxEnterLimitBB),
		TableSize:       int32(s.TableSize),
	}
}

func (s *GameSetting) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("GameMetaUid", s.GameMetaUid)
	enc.AddInt("SmallBlind", s.SmallBlind)
	enc.AddInt("BigBlind", s.BigBlind)

	enc.AddDuration("TurnDuration", s.TurnDuration)
	enc.AddInt("MinEnterLimitBB", s.MinEnterLimitBB)
	enc.AddInt("MaxEnterLimitBB", s.MaxEnterLimitBB)
	enc.AddInt("WaterPct", s.WaterPct)
	enc.AddInt("TableSize", s.TableSize)
	enc.AddInt("LeastGamePlayerAmount", s.LeastGamePlayerAmount)
	enc.AddInt("MaxWaterLimitBB", s.MaxWaterLimitBB)

	enc.AddInt("MaxUserAmount", s.MaxUserAmount)
	return nil
}

func ProvideGameSetting() *GameSetting {
	return &GameSetting{
		GameMetaUid: "Undefined",
	}
}
