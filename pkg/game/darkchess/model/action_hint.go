package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"go.uber.org/zap/zapcore"
	"time"
)

type FreezeCell struct {
	Pos             *Pos
	Piece           piece.Piece
	IsPieceRevealed bool
	IsEmptyCell     bool
}

func (d *FreezeCell) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	_ = enc.AddObject("Pos", d.Pos)
	enc.AddString("Piece", d.Piece.GetName())
	enc.AddBool("IsPieceRevealed", d.IsPieceRevealed)
	enc.AddBool("IsEmptyCell", d.IsEmptyCell)
	return nil
}

func (d *FreezeCell) ToProto() *gamegrpc.Cell {
	if d == nil {
		return nil
	}

	return &gamegrpc.Cell{
		GridPosition: &gamegrpc.GridPosition{
			X: int32(d.Pos.X),
			Y: int32(d.Pos.Y),
		},
		Piece:      d.Piece.ToProto(),
		IsRevealed: d.IsPieceRevealed,
		IsEmpty:    d.IsEmptyCell,
	}
}

type RaiseData struct {
	Turn     int
	RaisedAt time.Time
}

func (d RaiseData) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("Turn", d.Turn)
	enc.AddTime("RaisedAt", d.RaisedAt)
	return nil
}

type ActionHint struct {
	Uid          core.Uid
	TurnCount    int
	TurnDuration time.Duration

	TimeExtendTurns []RaiseData
	SurrenderTurns  []RaiseData
	FreezeCell      *FreezeCell
}

func (d *ActionHint) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Uid", d.Uid.String())
	enc.AddInt("TurnCount", d.TurnCount)
	enc.AddDuration("TurnDuration", d.TurnDuration)

	if d.FreezeCell != nil {
		_ = enc.AddObject("FreezeCell", d.FreezeCell)
	}

	if len(d.TimeExtendTurns) > 0 {
		_ = enc.AddArray("TimeExtendTurns", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
			for _, v := range d.TimeExtendTurns {
				_ = enc.AppendObject(v)
			}
			return nil
		}))
	}

	if len(d.SurrenderTurns) > 0 {
		_ = enc.AddArray("SurrenderTurns", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
			for _, v := range d.SurrenderTurns {
				_ = enc.AppendObject(v)
			}
			return nil
		}))
	}

	return nil
}
