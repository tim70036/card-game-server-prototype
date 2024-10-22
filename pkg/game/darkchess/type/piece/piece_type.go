package piece

import (
	"card-game-server-prototype/pkg/grpc/commongrpc"
)

var (
	General = typeData{
		name:   "帥/將",
		proto:  commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_GENERAL,
		weight: 7,
	}
	Advisor = typeData{
		name:   "仕/士",
		proto:  commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_ADVISOR,
		weight: 6,
	}
	Elephant = typeData{
		name:   "相/象",
		proto:  commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_ELEPHANT,
		weight: 5,
	}
	Chariot = typeData{
		name:   "俥/車",
		proto:  commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_CHARIOT,
		weight: 4,
	}
	Horse = typeData{
		name:   "傌/馬",
		proto:  commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_HORSE,
		weight: 3,
	}
	Cannon = typeData{
		name:   "炮/包",
		proto:  commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_CANNON,
		weight: 2,
	}
	Soldier = typeData{
		name:   "兵/卒",
		proto:  commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_SOLDIER,
		weight: 1,
	}
)

var pieceTypeProtoToData = map[commongrpc.CnChessPieceType]typeData{
	commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_GENERAL:  General,
	commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_ADVISOR:  Advisor,
	commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_ELEPHANT: Elephant,
	commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_CHARIOT:  Chariot,
	commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_HORSE:    Horse,
	commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_CANNON:   Cannon,
	commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_SOLDIER:  Soldier,
}

type typeData struct {
	proto  commongrpc.CnChessPieceType
	name   string
	weight int
}

func getTypeData(pieceType commongrpc.CnChessPieceType) typeData {
	if v, ok := pieceTypeProtoToData[pieceType]; ok {
		return v
	}
	return typeData{
		proto: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_INVALID,
	}
}

func (t typeData) getProto() commongrpc.CnChessPieceType {
	return t.proto
}

func (t typeData) getName() string {
	return t.name
}

func (t typeData) isGeneral() bool {
	return t.proto == commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_GENERAL
}

func (t typeData) isAdvisor() bool {
	return t.proto == commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_ADVISOR
}

func (t typeData) isElephant() bool {
	return t.proto == commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_ELEPHANT
}

func (t typeData) isChariot() bool {
	return t.proto == commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_CHARIOT
}

func (t typeData) isHorse() bool {
	return t.proto == commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_HORSE
}

func (t typeData) isCannon() bool {
	return t.proto == commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_CANNON
}

func (t typeData) isSoldier() bool {
	return t.proto == commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_SOLDIER
}
