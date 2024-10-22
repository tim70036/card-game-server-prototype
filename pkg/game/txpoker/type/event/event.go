package event

import "go.uber.org/zap/zapcore"

type Event struct {
	Type   EventType
	Amount int
	Extra  string
}

func (e *Event) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Type", e.Type.String())
	enc.AddInt("Amount", e.Amount)
	if e.Extra != "" {
		enc.AddString("Extra", e.Extra)
	}
	return nil
}

type EventList []*Event

func (l EventList) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, e := range l {
		enc.AppendObject(e)
	}
	return nil
}
