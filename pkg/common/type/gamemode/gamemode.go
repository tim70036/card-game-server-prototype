package gamemode

import "card-game-server-prototype/pkg/grpc/commongrpc"

type GameMode string

const (
	Buddy       GameMode = "0"
	Common      GameMode = "1"
	Club        GameMode = "2"
	Rank        GameMode = "3"
	Carnival    GameMode = "4"
	Qualifier   GameMode = "5"
	Elimination GameMode = "6"
)

func (m GameMode) String() string {
	name, ok := names[m]
	if !ok {
		return "Undefined"
	}
	return name
}

func (m GameMode) Val() string {
	return string(m)
}

func (m GameMode) ToProto() commongrpc.RoomInfo_GameMode {
	return protos[m]
}

var names = map[GameMode]string{
	Buddy:       "Buddy",
	Common:      "Common",
	Club:        "Club",
	Rank:        "Rank",
	Carnival:    "Carnival",
	Qualifier:   "Qualifier",
	Elimination: "Elimination",
}

var protos = map[GameMode]commongrpc.RoomInfo_GameMode{
	Buddy:       commongrpc.RoomInfo_BUDDY,
	Common:      commongrpc.RoomInfo_COMMON,
	Club:        commongrpc.RoomInfo_CLUB,
	Rank:        commongrpc.RoomInfo_RANK,
	Carnival:    commongrpc.RoomInfo_CARNIVAL,
	Qualifier:   commongrpc.RoomInfo_QUALIFIER,
	Elimination: commongrpc.RoomInfo_ELIMINATION,
}
