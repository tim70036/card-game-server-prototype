package session

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/type/event"
	"go.uber.org/zap/zapcore"
	"strconv"
)

type StatsSession struct {
	Uid                  core.Uid
	HighestGameWinAmount int
	EventAmountSum       map[event.EventType]int
}

func (s *StatsSession) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("HighestGameWinAmount", s.HighestGameWinAmount)
	_ = enc.AddObject("EventAmountSum", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for eventType, count := range s.EventAmountSum {
			enc.AddInt(strconv.Itoa(int(eventType)), count)
		}
		return nil
	}))
	return nil
}

type StatsSessionGroup struct {
	Data map[core.Uid]*StatsSession
}

func (g *StatsSessionGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	_ = enc.AddObject("Data", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for uid, statsSession := range g.Data {
			_ = enc.AddObject(uid.String(), statsSession)
		}
		return nil
	}))
	return nil
}
