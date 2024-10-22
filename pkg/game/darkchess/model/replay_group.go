package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	"go.uber.org/zap/zapcore"
)

type ActReveal struct {
	X, Y  int
	Piece piece.Piece
}

type ActMove struct {
	X, Y, ToX, ToY int
	Piece          piece.Piece
	CountPieces    int
}

type ActCapture struct {
	X, Y, ToX, ToY int
	Piece          piece.Piece
	TargetPiece    piece.Piece
}

type Replay struct {
	Uid     core.Uid
	Turn    int
	Reveal  *ActReveal
	Move    *ActMove
	Capture *ActCapture
}

type ReplayGroup struct {
	Data []Replay
}

func ProvideReplayGroup() *ReplayGroup {
	return &ReplayGroup{}
}

func (m *ReplayGroup) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, replay := range m.Data {
		_ = enc.AppendObject(zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			enc.AddString("Uid", replay.Uid.String())

			if replay.Reveal != nil {
				_ = enc.AddObject("Reveal", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
					enc.AddString("Piece", replay.Reveal.Piece.GetName())
					enc.AddInt("X", replay.Reveal.X)
					enc.AddInt("Y", replay.Reveal.Y)
					return nil
				}))
			}

			if replay.Move != nil {
				_ = enc.AddObject("Move", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
					enc.AddString("Piece", replay.Move.Piece.GetName())
					enc.AddInt("X", replay.Move.X)
					enc.AddInt("Y", replay.Move.Y)
					enc.AddInt("ToX", replay.Move.ToX)
					enc.AddInt("ToY", replay.Move.ToY)
					enc.AddInt("CountPieces", replay.Move.CountPieces)
					return nil
				}))
			}

			if replay.Capture != nil {
				_ = enc.AddObject("Capture", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
					enc.AddString("Piece", replay.Capture.Piece.GetName())
					enc.AddString("TargetPiece", replay.Capture.TargetPiece.GetName())
					enc.AddInt("X", replay.Capture.X)
					enc.AddInt("Y", replay.Capture.Y)
					enc.AddInt("ToX", replay.Capture.ToX)
					enc.AddInt("ToY", replay.Capture.ToY)
					return nil
				}))
			}

			return nil
		}))
	}
	return nil
}
