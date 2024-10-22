package session

import (
	"card-game-server-prototype/pkg/core"
	"go.uber.org/zap/zapcore"
)

type GameResultSession struct {
	GameId string

	Uid core.Uid

	Chip int

	CountIdleRounds int
}

func (s *GameResultSession) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("GameId", s.GameId)
	enc.AddString("Uid", s.Uid.String())
	enc.AddInt("Chip", s.Chip)
	enc.AddInt("CountIdleRounds", s.CountIdleRounds)
	return nil
}
