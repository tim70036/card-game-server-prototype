package config

import (
	"flag"
	"card-game-server-prototype/pkg/common/type/gamemode"
	"card-game-server-prototype/pkg/common/type/gametype"

	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
)

var tls = flag.Bool("tls", false, "run with tls")

type Config struct {
	RoomId      *string
	ShortRoomId *string
	GameType    *gametype.GameType
	GameMode    *gamemode.GameMode
	GameMetaUid *string
	ValidUsers  *string
	JWTKey      *string
	ClubId      *string
}

func (c *Config) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("RoomId", *c.RoomId)
	enc.AddString("ShortRoomId", *c.ShortRoomId)
	enc.AddString("GameType", c.GameType.String())
	enc.AddString("GameMode", c.GameMode.String())
	enc.AddString("GameMetaUid", *c.GameMetaUid)
	enc.AddString("ValidUsers", *c.ValidUsers)
	enc.AddString("JWTKey", *c.JWTKey)
	return nil
}

var CFG = &Config{
	RoomId:      flag.String("room-id", uuid.NewString(), "the id of the room to create"),
	ShortRoomId: flag.String("short-room-id", "undefined", "the short id of the room to create"),
	GameType:    (*gametype.GameType)(flag.String("game-type", string(gametype.TXPoker), "the type of the room to create")),
	GameMode:    (*gamemode.GameMode)(flag.String("game-mode", string(gamemode.Common), "the game mode of the room to create")),
	GameMetaUid: flag.String("game-meta-uid", "e76a7891-0db0-4d0e-8746-1771c6bd4706", "the game meta of the room to create"),
	ValidUsers:  flag.String("valid-users", "[]", "the valid users that can enter room (in uid array format)"),
	JWTKey:      flag.String("jwt-key", "youdasdasdasdasdasd1123J@@IhaveBBIGGGjjJJj", "the key to use for verifying jwt"),
}
