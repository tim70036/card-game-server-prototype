package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"go.uber.org/zap/zapcore"
)

type PickBoard struct {
	Pieces       []piece.Piece
	PickedPieces map[core.Uid]piece.Piece
	FirstPlayer  core.Uid
}

func ProvidePickBoard() *PickBoard {
	return &PickBoard{
		PickedPieces: map[core.Uid]piece.Piece{},
	}
}

func (b *PickBoard) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	_ = enc.AddArray("Pieces", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		for _, p := range b.Pieces {
			enc.AppendString(p.GetName())
		}
		return nil
	}))

	_ = enc.AddObject("PickedPieces", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for u, p := range b.PickedPieces {
			enc.AddString(u.String(), p.GetName())
		}
		return nil
	}))

	enc.AddString("FirstPlayer", b.FirstPlayer.String())

	return nil
}

func (b *PickBoard) ToProto() *gamegrpc.PickResult {
	data := &gamegrpc.PickResult{
		Pieces: map[string]commongrpc.CnChessPiece{},
	}

	if b.FirstPlayer != "" {
		tmp := b.FirstPlayer.String()
		data.FirstUid = &tmp
	}

	for uid, picked := range b.PickedPieces {
		data.Pieces[uid.String()] = picked.ToProto()
	}

	return data
}
