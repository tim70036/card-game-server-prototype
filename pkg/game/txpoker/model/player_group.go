package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/type/card"
	"card-game-server-prototype/pkg/game/txpoker/type/fold"
	"card-game-server-prototype/pkg/game/txpoker/type/hand"
	"card-game-server-prototype/pkg/game/txpoker/type/role"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"go.uber.org/zap/zapcore"
)

type Player struct {
	Uid core.Uid

	Role role.Role

	SeatId int

	PocketCards card.CardList

	// Allow to set once in the round.
	// Reset player data in start round state.
	showFold fold.ShowType

	Hand hand.Hand

	// 不需在 round 結束後重置，因為 player 會消失。
	IsIdle bool
}

func (p *Player) ToProto() *txpokergrpc.Player {
	return &txpokergrpc.Player{
		Uid:          p.Uid.String(),
		Role:         p.Role.ToProto(),
		SeatId:       int32(p.SeatId),
		PocketCards:  p.PocketCards.ToProto(),
		ShowFoldType: int32(p.showFold),
	}
}

func (p *Player) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Uid", p.Uid.String())
	enc.AddString("Role", p.Role.String())
	enc.AddInt("SeatId", p.SeatId)
	enc.AddString("PocketCards", p.PocketCards.ToString())
	enc.AddBool("IsIdle", p.IsIdle)

	if p.Hand != nil {
		enc.AddObject("Hand", p.Hand)
	} else {
		enc.AddString("Hand", "nil")
	}
	return nil
}

func (p *Player) ShowLeftFold() {
	p.showFold = fold.ShowLeft
}

func (p *Player) ShowRightFold() {
	p.showFold = fold.ShowRight
}

func (p *Player) ShowBothFold() {
	p.showFold = fold.ShowBoth
}

func (p *Player) HasShowFoldSet() bool {
	return p.showFold > 0
}

func (p *Player) GetShowFold() fold.ShowType {
	return p.showFold
}

type PlayerGroup struct {
	Data         map[core.Uid]*Player
	LastBBPlayer *Player
}

func ProvidePlayerGroup() *PlayerGroup {
	return &PlayerGroup{
		Data:         make(map[core.Uid]*Player),
		LastBBPlayer: nil,
	}
}

func (g *PlayerGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddObject("Data", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for uid, player := range g.Data {
			enc.AddObject(uid.String(), player)
		}
		return nil
	}))

	if g.LastBBPlayer != nil {
		enc.AddString("LastBBPlayer", g.LastBBPlayer.Uid.String())
	}
	return nil
}
