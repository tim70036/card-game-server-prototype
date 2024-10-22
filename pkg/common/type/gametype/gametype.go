package gametype

import "card-game-server-prototype/pkg/grpc/commongrpc"

type GameType string

const (
	CMJ         GameType = "0"
	DMJ         GameType = "1"
	Yablon      GameType = "2"
	CNPoker     GameType = "3"
	MJ          GameType = "4"
	AMJ         GameType = "5"
	TXPoker     GameType = "6"
	DarkChess   GameType = "7"
	HorseRaceMJ GameType = "8"
	ZoomTXPoker GameType = "9"
)

func (t GameType) String() string {
	name, ok := names[t]
	if !ok {
		return "Undefined"
	}
	return name
}

func (t GameType) ToProto() commongrpc.RoomInfo_GameType {
	return protos[t]
}

var names = map[GameType]string{
	CMJ:         "CMJ",
	DMJ:         "DMJ",
	Yablon:      "Yablon",
	CNPoker:     "CNPoker",
	MJ:          "MJ",
	AMJ:         "AMJ",
	TXPoker:     "TXPoker",
	DarkChess:   "DarkChess",
	HorseRaceMJ: "HorseRaceMJ",
	ZoomTXPoker: "ZoomTXPoker",
}

var protos = map[GameType]commongrpc.RoomInfo_GameType{
	CMJ:         commongrpc.RoomInfo_CMJ,
	DMJ:         commongrpc.RoomInfo_DMJ,
	Yablon:      commongrpc.RoomInfo_YABLON,
	CNPoker:     commongrpc.RoomInfo_CN_POKER,
	MJ:          commongrpc.RoomInfo_MJ,
	AMJ:         commongrpc.RoomInfo_AMJ,
	TXPoker:     commongrpc.RoomInfo_TX_POKER,
	DarkChess:   commongrpc.RoomInfo_DARK_CHESS,
	HorseRaceMJ: commongrpc.RoomInfo_HORSE_RACE_MJ,
	ZoomTXPoker: commongrpc.RoomInfo_TX_POKER_ZOOM,
}
