package cheat

import (
	"fmt"
	"card-game-server-prototype/pkg/core"
	card2 "card-game-server-prototype/pkg/game/txpoker/type/card"
	"github.com/tidwall/gjson"
	"go.uber.org/zap/zapcore"
)

type CheatData struct {
	PlayerPocketCards map[core.Uid]card2.CardList
	CommunityCards    card2.CardList
}

func FromRawCheatData(rawCheatData string) (*CheatData, error) {
	playerPocketCards := map[core.Uid]card2.CardList{}
	for rawUid, rawCards := range gjson.Get(rawCheatData, "playerPocketCards").Map() {
		uid := core.Uid(rawUid)

		cards := card2.CardList{}
		for _, rawCard := range rawCards.Array() {
			card, err := card2.FromHex(rawCard.String(), 0)
			if err != nil {
				return nil, err
			}
			cards = append(cards, card)
		}

		if len(cards) != 2 {
			return nil, fmt.Errorf("uid %v pocket cards length must be 2, but is: %v", uid, len(cards))
		}

		playerPocketCards[uid] = cards
	}

	communityCards := card2.CardList{}
	for _, rawCard := range gjson.Get(rawCheatData, "communityCards").Array() {
		card, err := card2.FromHex(rawCard.String(), 0)
		if err != nil {
			return nil, err
		}
		communityCards = append(communityCards, card)
	}

	if len(communityCards) != 5 {
		return nil, fmt.Errorf("community cards length must be 5, but is: %v", len(communityCards))
	}

	return &CheatData{
		PlayerPocketCards: playerPocketCards,
		CommunityCards:    communityCards,
	}, nil
}

func (d *CheatData) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddObject("PlayerPocketCards", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for uid, cards := range d.PlayerPocketCards {
			enc.AddString(uid.String(), cards.ToString())
		}
		return nil
	}))

	enc.AddString("CommunityCards", d.CommunityCards.ToString())
	return nil
}
