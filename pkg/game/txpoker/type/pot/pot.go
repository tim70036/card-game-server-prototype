package pot

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"

	"github.com/samber/lo"
	"go.uber.org/zap/zapcore"
)

type Pot struct {
	HasAllIn bool // true: means this pot is all set. DO NOT change any value.
	Chips    map[core.Uid]int
	Winners  map[core.Uid]*Winner
}

func (p *Pot) ToProto() *txpokergrpc.Pot {
	return &txpokergrpc.Pot{
		Chips:       lo.MapEntries(p.Chips, func(k core.Uid, v int) (string, int32) { return k.String(), int32(v) }),
		WinnerChips: lo.MapEntries(p.Winners, func(k core.Uid, v *Winner) (string, int32) { return k.String(), int32(v.Chip) }),
	}
}

func (p *Pot) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if len(p.Chips) > 0 {
		_ = enc.AddObject("Chips", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			for uid, chips := range p.Chips {
				enc.AddInt(uid.String(), chips)
			}
			return nil
		}))
	}

	if len(p.Winners) > 0 {
		_ = enc.AddObject("Winners", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			for uid, winner := range p.Winners {
				_ = enc.AddObject(uid.String(), winner)
			}
			return nil
		}))
	}
	return nil
}
