package action

import (
	"encoding/json"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/type/role"
	"go.uber.org/zap/zapcore"
)

type RaiseRecord struct {
	baseActionRecord
	Uid  core.Uid
	Role role.Role
	Chip int
}

func (r *RaiseRecord) GetUid() core.Uid    { return r.Uid }
func (r *RaiseRecord) GetType() ActionType { return Raise }

func (r *RaiseRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type": r.GetType(),
		"uid":  r.Uid,
		"role": r.Role,
		"chip": r.Chip,
	})
}

func (r *RaiseRecord) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Type", r.GetType().String())
	enc.AddString("Uid", r.Uid.String())
	enc.AddString("Role", r.Role.String())
	enc.AddInt("Chip", r.Chip)
	return nil
}