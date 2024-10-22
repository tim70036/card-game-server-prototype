package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap/zapcore"
)

// Delta? https://levelup.gitconnected.com/building-immutable-data-structures-in-go-56a1068c76b2

type User struct {
	Uid         core.Uid
	ShortUid    string
	Name        string
	IsAI        bool
	IsConnected bool
	HasEntered  bool
	Cash        int
	Level       int
	RoomCards   int
}

func (u *User) ToProto() *commongrpc.User {
	return &commongrpc.User{
		Uid:         u.Uid.String(),
		ShortUid:    u.ShortUid,
		Username:    u.Name,
		IsConnected: u.IsConnected,
		HasEntered:  u.HasEntered,
		Cash:        int32(u.Cash),
		Level:       int32(u.Level),
		RoomCards:   int32(u.RoomCards),
	}
}

func (u *User) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Uid", u.Uid.String())
	enc.AddString("ShortUid", u.ShortUid)
	enc.AddString("Name", u.Name)
	enc.AddBool("IsAI", u.IsAI)
	enc.AddBool("IsConnected", u.IsConnected)
	enc.AddBool("HasEntered", u.HasEntered)
	enc.AddInt("Cash", u.Cash)
	enc.AddInt("Level", u.Level)
	enc.AddInt("RoomCards", u.RoomCards)
	return nil
}

type UserGroup struct {
	Data map[core.Uid]*User
}

func ProvideUserGroup(loggerFactory *util.LoggerFactory) *UserGroup {
	return &UserGroup{
		Data: make(map[core.Uid]*User),
	}
}

func (g *UserGroup) ToProto() *commongrpc.UserGroup {
	msg := &commongrpc.UserGroup{Users: make(map[string]*commongrpc.User)}
	for uid, user := range g.Data {
		msg.Users[uid.String()] = user.ToProto()
	}
	return msg
}

func (g *UserGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for uid, user := range g.Data {
		enc.AddObject(uid.String(), user)
	}
	return nil
}

// HasAnyConnected 用來判斷是不是有任何一個玩家是連線中（不含 AI）
func (g *UserGroup) HasAnyConnected() bool {
	for _, user := range g.Data {
		if user.IsAI {
			continue
		}

		if user.IsConnected {
			return true
		}
	}
	return false
}
