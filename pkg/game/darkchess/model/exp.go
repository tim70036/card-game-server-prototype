package model

import (
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"go.uber.org/zap/zapcore"
)

type Exp struct {
	BeforeLevel    int
	BeforeExp      int
	LevelUpExp     int
	NextLevel      int
	NextLevelUpExp int
	IncreaseExp    int
}

func (e *Exp) ToProto() *commongrpc.ExpInfo {
	// Prevent frontend crash.
	var beforeLevel int32 = 1
	var nextLevel int32 = 1

	if e.BeforeLevel > 0 {
		beforeLevel = int32(e.BeforeLevel)
	}

	if e.NextLevel > 0 {
		nextLevel = int32(e.NextLevel)
	}

	return &commongrpc.ExpInfo{
		BeforeLevel:    beforeLevel,
		BeforeExp:      int32(e.BeforeExp),
		LevelUpExp:     int32(e.LevelUpExp),
		NextLevel:      nextLevel,
		NextLevelUpExp: int32(e.NextLevelUpExp),
		IncreaseExp:    int32(e.IncreaseExp),
	}
}

func (e *Exp) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if e == nil {
		return nil
	}
	enc.AddInt("beforeLevel", e.BeforeLevel)
	enc.AddInt("beforeExp", e.BeforeExp)
	enc.AddInt("levelUpExp", e.LevelUpExp)
	enc.AddInt("nextLevel", e.NextLevel)
	enc.AddInt("nextLevelUpExp", e.NextLevelUpExp)
	enc.AddInt("increaseExp", e.IncreaseExp)
	return nil
}
