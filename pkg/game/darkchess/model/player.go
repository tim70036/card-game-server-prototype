package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"go.uber.org/zap/zapcore"
)

type Player struct {
	Uid core.Uid

	// The accumulated raw profit from each round.
	RawProfit int

	// The accumulated final profit from each round.
	Profit int

	// Player will bankrupt if cash < 0 after a round.
	IsBankrupt bool

	// Record how long the player has disconnected. If too long, the
	// game will end.
	DisconnectedRoundCount int

	// Idle 到結束遊戲算「斷線1場」
	IdleTurns []int

	Color    commongrpc.CnChessColorType
	IsWinner bool
	HasActed bool // has raised any reveal, move, capture
}

func (d *Player) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Uid", d.Uid.String())
	enc.AddBool("IsBankrupt", d.IsBankrupt)
	enc.AddInt("RawProfit", d.RawProfit)
	enc.AddInt("Profit", d.Profit)
	enc.AddInt("DisconnectedRoundCount", d.DisconnectedRoundCount)
	if len(d.IdleTurns) > 0 {
		_ = enc.AddArray("IdleTurns", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
			for _, turn := range d.IdleTurns {
				enc.AppendInt(turn)
			}
			return nil
		}))
	}
	enc.AddString("Color", d.Color.String())
	enc.AddBool("IsWinner", d.IsWinner)
	return nil
}
