package action

import (
	"encoding/json"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/type/role"
	"go.uber.org/zap/zapcore"
)

type SBRecord struct {
	baseActionRecord
	Uid  core.Uid
	Role role.Role
	Chip int
}

func (r *SBRecord) GetUid() core.Uid    { return r.Uid }
func (r *SBRecord) GetType() ActionType { return SB }

func (r *SBRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type": r.GetType(),
		"uid":  r.Uid,
		"role": r.Role,
		"chip": r.Chip,
	})
}

func (r *SBRecord) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Type", r.GetType().String())
	enc.AddString("Uid", r.Uid.String())
	enc.AddString("Role", r.Role.String())
	enc.AddInt("Chip", r.Chip)
	return nil
}
