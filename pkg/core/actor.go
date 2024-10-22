package core

import (
	"time"

	"go.uber.org/zap/zapcore"
)

type Actor interface {
	Uid() Uid
	MarshalLogObject(enc zapcore.ObjectEncoder) error
}

type ActorRequest struct {
	Delay time.Duration
	Req   *Request
}

func (r *ActorRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddDuration("Delay", r.Delay)
	enc.AddObject("Req", r.Req)
	return nil
}

type ActorRequestList []*ActorRequest

func (l ActorRequestList) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, req := range l {
		enc.AppendObject(req)
	}
	return nil
}
