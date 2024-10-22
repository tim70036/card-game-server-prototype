package core

import (
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Request struct {
	Uid Uid
	Msg proto.Message
}

func (r *Request) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Uid", r.Uid.String())
	enc.AddString("MsgType", string(r.Msg.ProtoReflect().Descriptor().FullName()))
	enc.AddString("Msg", protojson.Format(r.Msg))
	return nil
}
