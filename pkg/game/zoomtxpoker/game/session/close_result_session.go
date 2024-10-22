package session

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/type/event"
	"go.uber.org/zap/zapcore"
)

type UserCloseResultSession struct {
	GameId string
	Uid    core.Uid
	Name   string

	// Stats
	HighestGameWinAmount int
	IncrEventAmount      map[event.EventType]int

	// TableProfits: result-amount - session-amount
	IncrCountGames   int
	IncrBuyInChips   int
	IncrWinLoseChips int
}

func (s *UserCloseResultSession) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("GameId", s.GameId)
	enc.AddString("Uid", s.Uid.String())
	enc.AddString("Name", s.Name)
	enc.AddInt("HighestGameWinAmount", s.HighestGameWinAmount)
	_ = enc.AddObject("IncrEventAmount", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for e, count := range s.IncrEventAmount {
			if count > 0 {
				enc.AddInt(e.String(), count)
			}
		}
		return nil
	}))
	if s.IncrCountGames > 0 {
		enc.AddInt("IncrCountGames", s.IncrCountGames)
	}
	if s.IncrBuyInChips > 0 {
		enc.AddInt("IncrBuyInChips", s.IncrBuyInChips)
	}
	if s.IncrWinLoseChips > 0 {
		enc.AddInt("SumWinLoseChips", s.IncrWinLoseChips)
	}
	return nil
}

type CloseResultSession struct {
	GameId string

	// VPIP
	IncrBetPlayerCount    int
	IncrPlayedPlayerCount int
}

func (s *CloseResultSession) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("GameId", s.GameId)
	enc.AddInt("TotalBetPlayerCount", s.IncrBetPlayerCount)
	enc.AddInt("TotalPlayedPlayerCount", s.IncrPlayedPlayerCount)
	return nil
}
