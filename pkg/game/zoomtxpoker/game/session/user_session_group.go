package session

import (
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/type/role"
	"go.uber.org/zap/zapcore"
)

type UserSession struct {
	User *commonmodel.User

	Chip int

	Role role.Role

	CountIdleRounds int
}

func (s *UserSession) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddObject("User", s.User)
	enc.AddInt("Chip", s.Chip)
	enc.AddString("Role", s.Role.String())
	enc.AddInt("CountIdleRounds", s.CountIdleRounds)
	return nil
}

type UserSessionGroup struct {
	Data map[core.Uid]*UserSession
}

func (g *UserSessionGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddObject("Data", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for uid, playerSession := range g.Data {
			enc.AddObject(uid.String(), playerSession)
		}
		return nil
	}))
	return nil
}
