package model

import (
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"go.uber.org/zap/zapcore"
)

type GameInfo struct {
	CreationId             int
	TotalBetPlayerCount    int
	TotalPlayedPlayerCount int
	StreakRoundCount       int
	RoundIdHistory         []string
	RoundId                string
}

func ProvideGameInfo() *GameInfo {
	return &GameInfo{
		CreationId:             0,
		TotalBetPlayerCount:    0,
		TotalPlayedPlayerCount: 0,
		StreakRoundCount:       0,
		RoundIdHistory:         []string{},
		RoundId:                "",
	}
}

func (g *GameInfo) CalculateVPIP() float64 {
	if g.TotalPlayedPlayerCount <= 0 {
		return 0
	}

	// sum 每個 round 的 (入池人數 / playing人數)
	return float64(g.TotalBetPlayerCount) / float64(g.TotalPlayedPlayerCount)
}

func (g *GameInfo) ToProto() *txpokergrpc.GameInfo {
	return &txpokergrpc.GameInfo{
		CreationId:      int32(g.CreationId),
		TotalRoundCount: int32(g.TotalPlayedPlayerCount),
		RoundIdHistory:  g.RoundIdHistory,
		RoundId:         g.RoundId,
		Vpip:            float32(g.CalculateVPIP()),
	}
}

func (g *GameInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("CreationId", g.CreationId)
	enc.AddInt("TotalPlayedPlayerCount", g.TotalPlayedPlayerCount)
	enc.AddInt("TotalBetPlayerCount", g.TotalBetPlayerCount)
	enc.AddInt("StreakRoundCount", g.StreakRoundCount)
	_ = enc.AddArray("RoundIdHistory", zapcore.ArrayMarshalerFunc(func(arrEnc zapcore.ArrayEncoder) error {
		for _, roundId := range g.RoundIdHistory {
			arrEnc.AppendString(roundId)
		}
		return nil
	}))
	enc.AddString("RoundId", g.RoundId)
	return nil
}
