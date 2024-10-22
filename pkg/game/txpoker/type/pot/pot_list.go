package pot

import (
	"card-game-server-prototype/pkg/grpc/txpokergrpc"

	"github.com/samber/lo"
	"go.uber.org/zap/zapcore"
)

type PotList []*Pot

func (l PotList) ToProto() []*txpokergrpc.Pot {
	return lo.Map(l, func(p *Pot, _ int) *txpokergrpc.Pot { return p.ToProto() })
}

func (l PotList) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, pot := range l {
		enc.AppendObject(pot)
	}
	return nil
}
