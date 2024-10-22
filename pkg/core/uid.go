package core

import (
	"reflect"

	"go.uber.org/zap/zapcore"
)

type Uid string

func (u Uid) String() string { return string(u) }

var UidType = reflect.TypeOf(Uid(""))

type UidList []Uid

func (l UidList) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, uid := range l {
		enc.AppendString(uid.String())
	}
	return nil
}
