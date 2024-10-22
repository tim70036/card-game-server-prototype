package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"card-game-server-prototype/pkg/util"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"go.uber.org/zap/zapcore"
)

// Delta? https://levelup.gitconnected.com/building-immutable-data-structures-in-go-56a1068c76b2

type Buddy struct {
	Uid       core.Uid
	IsReady   bool
	IsOwner   bool
	EnterTime time.Time
}

func (b *Buddy) ToProto() *commongrpc.Buddy {
	return &commongrpc.Buddy{
		Uid:       b.Uid.String(),
		IsReady:   b.IsReady,
		IsOwner:   b.IsOwner,
		EnterTime: timestamppb.New(b.EnterTime),
	}
}

func (b *Buddy) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Uid", b.Uid.String())
	enc.AddBool("IsReady", b.IsReady)
	enc.AddBool("IsOwner", b.IsOwner)
	enc.AddTime("EnterTime", b.EnterTime)
	return nil
}

type BuddyGroup struct {
	Data map[core.Uid]*Buddy
}

func ProvideBuddyGroup(loggerFactory *util.LoggerFactory) *BuddyGroup {
	return &BuddyGroup{
		Data: make(map[core.Uid]*Buddy),
	}
}

func (g *BuddyGroup) ToProto() *commongrpc.BuddyGroup {
	msg := &commongrpc.BuddyGroup{Buddies: make(map[string]*commongrpc.Buddy)}
	for uid, buddy := range g.Data {
		msg.Buddies[uid.String()] = buddy.ToProto()
	}
	return msg
}

func (g *BuddyGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for uid, buddy := range g.Data {
		enc.AddObject(uid.String(), buddy)
	}
	return nil
}

func (g *BuddyGroup) HasAnyOwner() bool {
	for _, v := range g.Data {
		if v.IsOwner {
			return true
		}
	}
	return false
}
