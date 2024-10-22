package action

import (
	"encoding/json"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/type/role"
	"go.uber.org/zap/zapcore"
)

type CheckRecord struct {
	baseActionRecord
	Uid  core.Uid
	Role role.Role
}

func (r *CheckRecord) GetUid() core.Uid    { return r.Uid }
func (r *CheckRecord) GetType() ActionType { return Check }

func (r *CheckRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type": r.GetType(),
		"uid":  r.Uid,
		"role": r.Role,
	})
}

func (r *CheckRecord) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Type", r.GetType().String())
	enc.AddString("Uid", r.Uid.String())
	enc.AddString("Role", r.Role.String())
	return nil
}
