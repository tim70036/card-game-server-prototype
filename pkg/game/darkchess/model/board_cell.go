package model

import (
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	"go.uber.org/zap/zapcore"
)

type Pos struct {
	X, Y int
}

func (p Pos) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("X", p.X)
	enc.AddInt("Y", p.Y)
	return nil
}

type Cell struct {
	Piece           piece.Piece
	IsPieceRevealed bool
	IsPieceDead     bool
	IsEmptyCell     bool
}

func NewEmptyCell() *Cell {
	return &Cell{
		Piece:       piece.InvalidPiece,
		IsEmptyCell: true,
	}
}

func (p *Cell) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Piece", p.Piece.GetName())
	enc.AddBool("IsPieceRevealed", p.IsPieceRevealed)
	enc.AddBool("IsPieceDead", p.IsPieceDead)
	enc.AddBool("IsEmptyCell", p.IsEmptyCell)
	return nil
}
