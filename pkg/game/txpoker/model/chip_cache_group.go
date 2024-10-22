package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"

	"go.uber.org/zap/zapcore"
)

type ChipCacheGroup struct {
	SeatStatusChips map[core.Uid]int

	// Record chip for winners that standup (could sit down in
	// another seat, and become reserving state). These will be cash
	// out when round end.
	CashOutChips map[core.Uid]int
}

func ProvideChipCacheGroup() *ChipCacheGroup {
	return &ChipCacheGroup{
		SeatStatusChips: make(map[core.Uid]int),
		CashOutChips:    make(map[core.Uid]int),
	}
}

func (g *ChipCacheGroup) ToProto() *txpokergrpc.ChipCacheGroup {
	msg := &txpokergrpc.ChipCacheGroup{SeatStatusChips: make(map[string]int32)}
	for uid, chip := range g.SeatStatusChips {
		msg.SeatStatusChips[uid.String()] = int32(chip)
	}
	return msg
}

func (g *ChipCacheGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddObject("SeatStatusChips", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for uid, chip := range g.SeatStatusChips {
			enc.AddInt(uid.String(), chip)
		}
		return nil
	}))

	enc.AddObject("CashOutChips", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for uid, chip := range g.CashOutChips {
			enc.AddInt(uid.String(), chip)
		}
		return nil
	}))
	return nil
}
