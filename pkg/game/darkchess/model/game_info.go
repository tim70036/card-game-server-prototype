package model

import (
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"go.uber.org/zap/zapcore"
)

type GameInfo struct {
	GameId     string
	RoundCount int
	Setting    *GameSetting
}

func ProvideGameInfo() *GameInfo {
	return &GameInfo{
		GameId:  "Undefined",
		Setting: &GameSetting{},
	}
}

func (m *GameInfo) ToProto() *gamegrpc.GameInfo {
	d := &gamegrpc.GameInfo{
		GameId:     m.GameId,
		RoundCount: int32(m.RoundCount),
		Setting:    &gamegrpc.GameSetting{},
	}

	if m.Setting != nil {
		d.Setting = m.Setting.ToProto()
	}

	return d
}

func (m *GameInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("GameId", m.GameId)
	enc.AddInt("RoundCount", m.RoundCount)
	_ = enc.AddObject("setting", m.Setting)
	return nil
}
