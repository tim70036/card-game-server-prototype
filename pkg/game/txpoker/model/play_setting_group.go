package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"

	"go.uber.org/zap/zapcore"
)

type PlaySetting struct {
	Uid core.Uid

	WaitBB bool

	AutoTopUp bool

	AutoTopUpThresholdPercent float64

	AutoTopUpChipPercent float64
}

func (p *PlaySetting) ToProto() *txpokergrpc.PlaySetting {
	return &txpokergrpc.PlaySetting{
		WaitBb:                    p.WaitBB,
		AutoTopUp:                 p.AutoTopUp,
		AutoTopUpThresholdPercent: p.AutoTopUpThresholdPercent,
		AutoTopUpChipPercent:      p.AutoTopUpChipPercent,
	}
}

func (p *PlaySetting) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Uid", p.Uid.String())
	enc.AddBool("WaitBB", p.WaitBB)
	enc.AddBool("AutoTopUp", p.AutoTopUp)
	enc.AddFloat64("AutoTopUpThresholdPercent", p.AutoTopUpThresholdPercent)
	enc.AddFloat64("AutoTopUpChipPercent", p.AutoTopUpChipPercent)
	return nil
}

type PlaySettingGroup struct {
	Data map[core.Uid]*PlaySetting
}

func ProvidePlaySettingGroup() *PlaySettingGroup {
	return &PlaySettingGroup{
		Data: make(map[core.Uid]*PlaySetting),
	}
}

func (g *PlaySettingGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for uid, playSetting := range g.Data {
		enc.AddObject(uid.String(), playSetting)
	}
	return nil
}
