package action

import (
	"encoding/json"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/type/card"
	"card-game-server-prototype/pkg/game/txpoker/type/hand"
	"card-game-server-prototype/pkg/game/txpoker/type/role"
	"go.uber.org/zap/zapcore"
)

type WinPotRecord struct {
	baseActionRecord
	Uid  core.Uid
	Role role.Role
	Chip int

	// Frontend will need to know which cards to reveal in the replay.
	// This mask help to determine which cards to reveal.
	// 0b00 -> no cards
	// 0b01 -> reveal 1st card
	// 0b10 -> reveal 2nd card
	// 0b11 -> reveal both cards
	PocketCardsMask int

	PocketCards card.CardList

	Hand hand.Hand
}

func (r *WinPotRecord) GetUid() core.Uid    { return r.Uid }
func (r *WinPotRecord) GetType() ActionType { return WinPot }

func (r *WinPotRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":            r.GetType(),
		"uid":             r.Uid,
		"role":            r.Role,
		"chip":            r.Chip,
		"pocketCardsMask": r.PocketCardsMask,
		"pocketCards":     r.PocketCards.ToHexStr(),
	})
}

func (r *WinPotRecord) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Type", r.GetType().String())
	enc.AddString("Uid", r.Uid.String())
	enc.AddString("Role", r.Role.String())
	enc.AddInt("Chip", r.Chip)
	enc.AddInt("PocketCardsMask", r.PocketCardsMask)
	enc.AddString("PocketCards", r.PocketCards.ToString())

	if r.Hand != nil {
		enc.AddObject("Hand", r.Hand)
	} else {
		enc.AddString("Hand", "nil")
	}

	return nil
}
