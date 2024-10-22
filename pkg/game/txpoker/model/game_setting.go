package model

import (
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"time"

	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/types/known/durationpb"
)

type GameSetting struct {
	GameMetaUid string
	SmallBlind  int
	BigBlind    int

	TurnDuration                 time.Duration
	InitialExtraTurnDuration     time.Duration
	ExtraTurnRefillIntervalRound int
	RefillExtraTurnDuration      time.Duration
	MaxExtraTurnDuration         time.Duration

	InitialSitOutDuration        time.Duration
	SitOutRefillIntervalDuration time.Duration
	RefillSitOutDuration         time.Duration
	MaxSitOutDuration            time.Duration

	MinEnterLimitBB int
	MaxEnterLimitBB int
	WaterPct        int
	TableSize       int

	LeastGamePlayerAmount int
	MaxWaterLimitBB       int

	CloseAt      time.Time // See initState.go comment:CloseAt
	ForceStartAt time.Time
}

func ProvideGameSetting() *GameSetting {
	return &GameSetting{
		GameMetaUid: "Undefined",
	}
}

func (s *GameSetting) ToProto() *txpokergrpc.GameSetting {
	return &txpokergrpc.GameSetting{
		GameMetaUid: s.GameMetaUid,
		SmallBlind:  int32(s.SmallBlind),
		BigBlind:    int32(s.BigBlind),

		TurnDuration:                 durationpb.New(s.TurnDuration),
		InitialExtraTurnDuration:     durationpb.New(s.InitialExtraTurnDuration),
		ExtraTurnRefillIntervalRound: int32(s.ExtraTurnRefillIntervalRound),
		RefillExtraTurnDuration:      durationpb.New(s.RefillExtraTurnDuration),
		MaxExtraTurnDuration:         durationpb.New(s.MaxExtraTurnDuration),

		InitialSitOutDuration:        durationpb.New(s.InitialSitOutDuration),
		SitOutRefillIntervalDuration: durationpb.New(s.SitOutRefillIntervalDuration),
		RefillSitOutDuration:         durationpb.New(s.RefillSitOutDuration),
		MaxSitOutDuration:            durationpb.New(s.MaxSitOutDuration),

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
	enc.AddDuration("InitialExtraTurnDuration", s.InitialExtraTurnDuration)
	enc.AddInt("ExtraTurnRefillIntervalRound", s.ExtraTurnRefillIntervalRound)
	enc.AddDuration("RefillExtraTurnDuration", s.RefillExtraTurnDuration)
	enc.AddDuration("MaxExtraTurnDuration", s.MaxExtraTurnDuration)

	enc.AddDuration("InitialSitOutDuration", s.InitialSitOutDuration)
	enc.AddDuration("SitOutRefillIntervalDuration", s.SitOutRefillIntervalDuration)
	enc.AddDuration("RefillSitOutDuration", s.RefillSitOutDuration)
	enc.AddDuration("MaxSitOutDuration", s.MaxSitOutDuration)

	enc.AddInt("MinEnterLimitBB", s.MinEnterLimitBB)
	enc.AddInt("MaxEnterLimitBB", s.MaxEnterLimitBB)
	enc.AddInt("WaterPct", s.WaterPct)
	enc.AddInt("TableSize", s.TableSize)
	enc.AddInt("LeastGamePlayerAmount", s.LeastGamePlayerAmount)
	enc.AddInt("MaxWaterLimitBB", s.MaxWaterLimitBB)
	return nil
}
