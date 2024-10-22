package model

import (
	"card-game-server-prototype/pkg/common/type/gamemode"
	"card-game-server-prototype/pkg/core"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"github.com/samber/lo"
	"go.uber.org/zap/zapcore"
)

type GameScore struct {
	Uid                   core.Uid
	IsBankrupt            bool
	IsDisconnected        bool
	DisconnectedRawProfit int

	// The original profit from the round result.
	RawProfit int

	// The final profit calculated by other factor (EX: water or no
	// water). It will be rendered to the player.
	Profit int

	ExpInfo *Exp // 原 EndGameInfo
}

func (s *GameScore) ToProto(mode gamemode.GameMode) *gamegrpc.GameScore {
	return &gamegrpc.GameScore{
		Uid:                s.Uid.String(),
		Profit:             lo.Ternary(mode == gamemode.Club, int32(s.Profit), int32(s.RawProfit)),
		ExpInfo:            s.ExpInfo.ToProto(),
		IsDisconnected:     s.IsDisconnected,
		DisconnectedProfit: int32(s.DisconnectedRawProfit),
	}
}

func (s *GameScore) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("uid", s.Uid.String())
	enc.AddBool("isBankrupt", s.IsBankrupt)
	enc.AddInt("rawProfit", s.RawProfit)
	enc.AddInt("DisconnectedProfit", s.DisconnectedRawProfit)
	enc.AddInt("profit", s.Profit)
	enc.AddBool("isDisconnected", s.IsDisconnected)
	_ = enc.AddObject("expInfo", s.ExpInfo)
	return nil
}
