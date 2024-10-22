package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/type/role"
	"go.uber.org/zap/zapcore"
)

type PlayedHistory struct {
	Uid            core.Uid
	CountRoles     map[role.Role]int
	LastPlayedRole role.Role
}

func (s *PlayedHistory) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Uid", s.Uid.String())

	if len(s.CountRoles) > 0 {
		_ = enc.AddObject("PlayedHistory", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			for playedRole, count := range s.CountRoles {
				enc.AddInt(playedRole.String(), count)
			}
			return nil
		}))
	}

	return nil
}

type PlayedHistoryGroup struct {
	Data                   map[core.Uid]*PlayedHistory
	TotalBetPlayerCount    int
	TotalPlayedPlayerCount int
}

func ProvidePlayedHistoryGroup() *PlayedHistoryGroup {
	return &PlayedHistoryGroup{
		Data: make(map[core.Uid]*PlayedHistory),
	}
}

func (g *PlayedHistoryGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	_ = enc.AddObject("PlayedHistory", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for uid, history := range g.Data {
			_ = enc.AddObject(uid.String(), history)
		}
		return nil
	}))
	enc.AddInt("TotalBetPlayerCount", g.TotalBetPlayerCount)
	enc.AddInt("TotalPlayedPlayerCount", g.TotalPlayedPlayerCount)
	return nil
}
