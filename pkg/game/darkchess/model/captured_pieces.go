package model

import (
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"go.uber.org/zap/zapcore"
)

type CapturedPieces struct {
	Pieces []piece.Piece // dead orders
}

func ProvideCapturedPieces() *CapturedPieces {
	return &CapturedPieces{}
}

func (m *CapturedPieces) ToProto() *gamegrpc.CapturedPieces {
	var pieces []commongrpc.CnChessPiece

	for _, v := range m.Pieces {
		pieces = append(pieces, v.ToProto())
	}

	return &gamegrpc.CapturedPieces{
		Pieces: pieces,
	}
}

func (m *CapturedPieces) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for _, v := range m.Pieces {
		arr.AppendString(v.GetName())
	}

	return nil
}
