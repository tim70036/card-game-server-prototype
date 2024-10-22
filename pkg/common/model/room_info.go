package model

import (
	"card-game-server-prototype/pkg/common/type/gamemode"
	"card-game-server-prototype/pkg/common/type/gametype"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"go.uber.org/zap/zapcore"
	"time"
)

type RoomInfo struct {
	RoomId      string
	ShortRoomId string
	GameType    gametype.GameType
	GameMode    gamemode.GameMode
	GameMetaUid string
	ValidUsers  core.UidList
	ClubId      string

	CloseAt      time.Time // See initState.go comment:CloseAt
	ForceStartAt time.Time

	IsPremium           bool
	PremiumUid          string
	PremiumEndTimestamp int64
}

func ProvideRoomInfo() *RoomInfo {
	return &RoomInfo{
		RoomId:      "Undefined",
		ShortRoomId: "Undefined",
		GameType:    gametype.MJ,
		GameMode:    gamemode.Buddy,
		GameMetaUid: "Undefined",
		ValidUsers:  core.UidList{},
	}
}

func (r *RoomInfo) ToProto() *commongrpc.RoomInfo {
	return &commongrpc.RoomInfo{
		RoomId:       r.RoomId,
		ShortRoomId:  r.ShortRoomId,
		RoomGameType: r.GameType.ToProto(),
		RoomGameMode: r.GameMode.ToProto(),
		GameMetaUid:  r.GameMetaUid,
		IsPremium:    r.IsPremium,
	}
}

func (r *RoomInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("RoomId", r.RoomId)
	enc.AddString("ShortRoomId", r.ShortRoomId)
	enc.AddString("GameType", r.GameType.String())
	enc.AddString("GameMode", r.GameMode.String())
	enc.AddString("GameMetaUid", r.GameMetaUid)
	_ = enc.AddArray("ValidUsers", r.ValidUsers)
	if r.ClubId != "" {
		enc.AddString("ClubId", r.ClubId)
	}
	if !r.CloseAt.IsZero() {
		enc.AddTime("CloseAt", r.CloseAt)
	}
	if !r.ForceStartAt.IsZero() {
		enc.AddTime("ForceStartAt", r.ForceStartAt)
	}

	enc.AddBool("IsPremium", r.IsPremium)
	enc.AddString("PremiumUid", r.PremiumUid)
	enc.AddTime("PremiumEndTime", time.Unix(r.PremiumEndTimestamp, 0))
	return nil
}
