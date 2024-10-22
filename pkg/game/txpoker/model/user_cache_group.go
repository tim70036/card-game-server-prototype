package model

import (
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"go.uber.org/zap/zapcore"
)

// UserCacheGroup
// Record userData for user left to prevent others see A default name user.
type UserCacheGroup struct {
	Data map[core.Uid]*commonmodel.User
}

func ProvideUserCacheGroup() *UserCacheGroup {
	return &UserCacheGroup{
		Data: make(map[core.Uid]*commonmodel.User),
	}
}

func (g *UserCacheGroup) ToProto() *txpokergrpc.UserCacheGroup {
	msg := &txpokergrpc.UserCacheGroup{Users: make(map[string]*commongrpc.User)}
	for uid, user := range g.Data {
		msg.Users[uid.String()] = user.ToProto()
	}
	return msg
}

func (g *UserCacheGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for uid, user := range g.Data {
		enc.AddObject(uid.String(), user)
	}
	return nil
}
