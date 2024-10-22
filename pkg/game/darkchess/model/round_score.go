package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"go.uber.org/zap/zapcore"
)

type RoundScore struct {
	Uid core.Uid

	// The original profit from the round result.
	RawProfit int

	// The final profit calculated by other factor (EX: water or  no
	// water). It will be rendered to the player.
	Profit int

	ScoreModifier int
	CapturePieces []piece.Piece
	Color         commongrpc.CnChessColorType
}

func (s *RoundScore) ToProto() *gamegrpc.RoundScore {
	var capturedPieces []commongrpc.CnChessPiece

	for _, v := range s.CapturePieces {
		capturedPieces = append(capturedPieces, v.ToProto())
	}

	return &gamegrpc.RoundScore{
		Uid:        s.Uid.String(),
		Points:     int32(len(s.CapturePieces)),
		RawProfits: int32(s.RawProfit),
		CapturedPieces: &gamegrpc.CapturedPieces{
			Pieces: capturedPieces,
		},
		ScoreModifier: gamegrpc.ScoreModifierType(s.ScoreModifier),
	}
}

func (s *RoundScore) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("uid", s.Uid.String())
	enc.AddInt("rawProfit", s.RawProfit)
	enc.AddInt("profit", s.Profit)
	enc.AddInt("ScoreModifier", s.ScoreModifier)
	enc.AddString("color", s.Color.String())
	_ = enc.AddArray("capturePieces", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		for _, v := range s.CapturePieces {
			enc.AppendString(v.GetName())
		}
		return nil
	}))
	return nil
}
