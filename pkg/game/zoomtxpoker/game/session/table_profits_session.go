package session

import (
	"card-game-server-prototype/pkg/core"
	"go.uber.org/zap/zapcore"
)

type TableProfitsSession struct {
	Uid             core.Uid
	Name            string
	CountGames      int
	SumBuyInChips   int
	SumWinLoseChips int
}

func (s *TableProfitsSession) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Uid", s.Uid.String())
	enc.AddString("Name", s.Name)
	enc.AddInt("CountGames", s.CountGames)
	enc.AddInt("SumBuyInChips", s.SumBuyInChips)
	enc.AddInt("SumWinLoseChips", s.SumWinLoseChips)
	return nil
}

type TableProfitsSessionGroup struct {
	Data map[core.Uid]*TableProfitsSession
}

func (g *TableProfitsSessionGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for uid, tableProfitsSession := range g.Data {
		_ = enc.AddObject(uid.String(), tableProfitsSession)
	}
	return nil
}
