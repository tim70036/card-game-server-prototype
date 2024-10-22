package action

import (
	"encoding/json"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/type/role"
	"go.uber.org/zap/zapcore"
)

type AllInRecord struct {
	baseActionRecord
	Uid     core.Uid
	Role    role.Role
	BetType ActionType
	Chip    int
}

func (r *AllInRecord) GetUid() core.Uid    { return r.Uid }
func (r *AllInRecord) GetType() ActionType { return AllIn }

func (r *AllInRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    r.GetType(),
		"uid":     r.Uid,
		"role":    r.Role,
		"betType": r.BetType,
		"chip":    r.Chip,
	})
}

func (r *AllInRecord) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Type", r.GetType().String())
	enc.AddString("Uid", r.Uid.String())
	enc.AddString("Role", r.Role.String())
	enc.AddString("BetType", r.BetType.String())
	enc.AddInt("Chip", r.Chip)
	return nil
}
