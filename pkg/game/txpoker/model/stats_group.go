package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/type/event"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"strconv"

	"github.com/samber/lo"
	"go.uber.org/zap/zapcore"
)

type Stats struct {
	Uid                  core.Uid
	HighestGameWinAmount int
	EventAmountSum       map[event.EventType]int
}

func (s *Stats) ToProto() *txpokergrpc.Stats {
	eventSums := lo.MapToSlice(s.EventAmountSum, func(k event.EventType, v int) *txpokergrpc.Stats_EventSum {
		return &txpokergrpc.Stats_EventSum{
			Type:   k.ToProto(),
			Amount: int32(v),
		}
	})

	// TODO: there's bug in stats that some event type is undefined (sticker)
	// Filter out undefined event type here so that frontend will not crash.
	eventSums = lo.Filter(eventSums, func(e *txpokergrpc.Stats_EventSum, _ int) bool {
		return e.Type != txpokergrpc.StatsEventType_UNDEFINED
	})

	return &txpokergrpc.Stats{
		Uid:                  s.Uid.String(),
		HighestGameWinAmount: int32(s.HighestGameWinAmount),
		EventSums:            eventSums,
	}
}
func (s *Stats) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("HighestGameWinAmount", s.HighestGameWinAmount)
	enc.AddObject("EventAmountSum", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for eventType, count := range s.EventAmountSum {
			enc.AddInt(strconv.Itoa(int(eventType)), count)
		}
		return nil
	}))
	return nil
}

type StatsGroup struct {
	Data map[core.Uid]*Stats
}

func ProvideStatsGroup() *StatsGroup {
	return &StatsGroup{
		Data: make(map[core.Uid]*Stats),
	}
}

func (g *StatsGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for uid, stats := range g.Data {
		enc.AddObject(uid.String(), stats)
	}
	return nil
}
