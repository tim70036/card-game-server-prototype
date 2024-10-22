package action

import (
	"encoding/json"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/type/card"
	"card-game-server-prototype/pkg/game/txpoker/type/role"
	"go.uber.org/zap/zapcore"
)

type ShowdownRecord struct {
	baseActionRecord
	Uid         core.Uid
	Role        role.Role
	PocketCards card.CardList
}

func (r *ShowdownRecord) GetUid() core.Uid    { return r.Uid }
func (r *ShowdownRecord) GetType() ActionType { return Showdown }

func (r *ShowdownRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":        r.GetType(),
		"uid":         r.Uid,
		"role":        r.Role,
		"pocketCards": r.PocketCards.ToHexStr(),
	})
}

func (r *ShowdownRecord) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Type", r.GetType().String())
	enc.AddString("Uid", r.Uid.String())
	enc.AddString("Role", r.Role.String())
	enc.AddString("PocketCards", r.PocketCards.ToString())
	return nil
}
