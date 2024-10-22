package session

import (
	"card-game-server-prototype/pkg/game/txpoker/model"
	"go.uber.org/zap/zapcore"
)

type GameInfoSession struct {
	Setting *model.GameSetting
	RoundId string

	// pool vpip = TotalBetPlayerCount / TotalPlayedPlayerCount
	// user vpip = stats.BetGame Amount / Game Amount
	TotalBetPlayerCount    int
	TotalPlayedPlayerCount int
}

func (s *GameInfoSession) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("RoundId", s.RoundId)
	_ = enc.AddObject("Setting", s.Setting)
	enc.AddInt("TotalBetPlayerCount", s.TotalBetPlayerCount)
	enc.AddInt("TotalPlayedPlayerCount", s.TotalPlayedPlayerCount)
	return nil
}
