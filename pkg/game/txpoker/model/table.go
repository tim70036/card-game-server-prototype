package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/type/card"
	"card-game-server-prototype/pkg/game/txpoker/type/pot"
	"card-game-server-prototype/pkg/game/txpoker/type/stage"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"github.com/qmuntal/stateless"
	"github.com/samber/lo"
	"go.uber.org/zap/zapcore"
)

type Table struct {
	Deck                card.CardList
	CommunityCards      card.CardList
	ShowdownPocketCards map[core.Uid]card.CardList

	BetStageFSM *stateless.StateMachine
	Pots        pot.PotList
}

func ProvideTable() *Table {
	return &Table{
		Deck:                make(card.CardList, 0),
		CommunityCards:      make(card.CardList, 0),
		ShowdownPocketCards: make(map[core.Uid]card.CardList),
		BetStageFSM:         stage.NewBetStageFSM(),
		Pots:                make(pot.PotList, 0),
	}
}

func (t *Table) ToProto() *txpokergrpc.Table {
	showdownPocketCards := lo.MapEntries(
		t.ShowdownPocketCards,
		func(uid core.Uid, cards card.CardList) (string, *txpokergrpc.CardList) {
			return uid.String(), &txpokergrpc.CardList{Cards: cards.ToProto()}
		},
	)

	return &txpokergrpc.Table{
		CommunityCards:      t.CommunityCards.ToProto(),
		ShowdownPocketCards: showdownPocketCards,
		Pots:                t.Pots.ToProto(),
	}
}

func (t *Table) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if len(t.CommunityCards) > 0 {
		enc.AddString("CommunityCards", t.CommunityCards.ToString())
	}
	if len(t.ShowdownPocketCards) > 0 {
		_ = enc.AddObject("ShowdownPocketCards", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			for uid, cards := range t.ShowdownPocketCards {
				enc.AddString(uid.String(), cards.ToString())
			}
			return nil
		}))
	}
	enc.AddString("BetStageFSM", t.BetStageFSM.MustState().(stage.Stage).String())
	if len(t.Pots) > 0 {
		_ = enc.AddArray("Pots", t.Pots)
	}
	return nil
}
