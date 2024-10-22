package session

import (
	"card-game-server-prototype/pkg/core"
	"go.uber.org/zap/zapcore"
)

type CashOutResultSession struct {
	GameId string
	Uid    core.Uid
	Chip   int
}

func (s *CashOutResultSession) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("GameId", s.GameId)
	enc.AddString("Uid", s.Uid.String())
	enc.AddInt("Chip", s.Chip)
	return nil
}
