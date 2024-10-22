package core

import (
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Response struct {
	Err error
	Msg proto.Message
}

func (r *Response) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if r == nil {
		return nil
	}

	enc.AddString("MsgType", string(r.Msg.ProtoReflect().Descriptor().FullName()))
	enc.AddString("Msg", protojson.Format(r.Msg))
	return nil
}
