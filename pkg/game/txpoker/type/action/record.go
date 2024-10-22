package action

import (
	"card-game-server-prototype/pkg/core"

	"go.uber.org/zap/zapcore"
)

type ActionRecord interface {
	GetUid() core.Uid
	GetType() ActionType
	MarshalJSON() ([]byte, error)
	MarshalLogObject(enc zapcore.ObjectEncoder) error
	mustEmbedBaseActionRecord()
}

type baseActionRecord struct{}

func (b *baseActionRecord) mustEmbedBaseActionRecord() {}
