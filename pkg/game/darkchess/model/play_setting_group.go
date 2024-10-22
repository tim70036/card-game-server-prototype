package model

import (
	"card-game-server-prototype/pkg/core"
	"go.uber.org/zap/zapcore"
)

type PlaySettingGroup struct {
	Data map[core.Uid]*PlaySetting
}

func ProvidePlaySettingGroup() *PlaySettingGroup {
	return &PlaySettingGroup{
		Data: make(map[core.Uid]*PlaySetting),
	}
}

func (g *PlaySettingGroup) ToProto() map[string]*PlaySetting {
	Data := make(map[string]*PlaySetting)
	for uid, playSetting := range g.Data {
		Data[uid.String()] = playSetting
	}
	return Data
}

func (g *PlaySettingGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for uid, playSetting := range g.Data {
		_ = enc.AddObject(uid.String(), playSetting)
	}
	return nil
}
