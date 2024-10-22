package core

import (
	"fmt"
	"reflect"

	"go.uber.org/zap/zapcore"
)

type StateTrigger struct {
	Name      string
	ArgsTypes []reflect.Type
}

func NewStateTrigger(name string, args ...reflect.Type) *StateTrigger {
	return &StateTrigger{
		Name:      name,
		ArgsTypes: args,
	}
}

func (t *StateTrigger) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Name", t.Name)
	for i, argType := range t.ArgsTypes {
		enc.AddString(fmt.Sprintf("arg%v", i), argType.String())
	}
	return nil
}
