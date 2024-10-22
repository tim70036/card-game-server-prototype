package rawevent

import (
	"go.uber.org/zap/zapcore"
)

type RawEvent struct {
	EventId int    `json:"eventId"`
	Amount  int    `json:"amount"`
	Extra   string `json:"extra"`
}

func (e *RawEvent) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("EventId", e.EventId)
	enc.AddInt("Amount", e.Amount)
	enc.AddString("Extra", e.Extra)
	return nil
}

type RawEventList []*RawEvent

func (l RawEventList) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, e := range l {
		_ = enc.AppendObject(e)
	}
	return nil
}
