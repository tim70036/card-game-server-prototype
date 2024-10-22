package event

import (
	"go.uber.org/zap/zapcore"
)

type Event struct {
	Code   Code
	Amount int
	Extra  string
}

func (e *Event) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Code", e.Code.String())
	enc.AddInt("Amount", e.Amount)
	enc.AddString("Extra", e.Extra)
	return nil
}

type List []*Event

func (l List) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, e := range l {
		_ = enc.AppendObject(e)
	}
	return nil
}
