package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"go.uber.org/zap/zapcore"
)

type PlaySetting struct {
	Uid                 core.Uid
	IsAuto              bool
	IsTriggerAutoByIdle bool
}

func (d *PlaySetting) ToProto() *darkchessgrpc.PlaySetting {
	return &darkchessgrpc.PlaySetting{
		IsAuto: d.IsAuto,
	}
}

func (d *PlaySetting) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("uid", d.Uid.String())
	enc.AddBool("isAuto", d.IsAuto)
	enc.AddBool("IsTriggerAutoByIdle", d.IsTriggerAutoByIdle)
	return nil
}
