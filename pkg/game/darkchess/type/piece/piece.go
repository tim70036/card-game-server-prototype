package piece

import (
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"strconv"
	"strings"
)

var (
	InvalidPiece   = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_INVALID)
	GeneralRed     = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_GENERAL_RED)
	AdvisorRed0    = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_ADVISOR_RED_0)
	AdvisorRed1    = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_ADVISOR_RED_1)
	ElephantRed0   = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_ELEPHANT_RED_0)
	ElephantRed1   = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_ELEPHANT_RED_1)
	ChariotRed0    = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_CHARIOT_RED_0)
	ChariotRed1    = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_CHARIOT_RED_1)
	HorseRed0      = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_HORSE_RED_0)
	HorseRed1      = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_HORSE_RED_1)
	CannonRed0     = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_CANNON_RED_0)
	CannonRed1     = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_CANNON_RED_1)
	SoldierRed0    = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_SOLDIER_RED_0)
	SoldierRed1    = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_SOLDIER_RED_1)
	SoldierRed2    = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_SOLDIER_RED_2)
	SoldierRed3    = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_SOLDIER_RED_3)
	SoldierRed4    = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_SOLDIER_RED_4)
	GeneralBlack   = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_GENERAL_BLACK)
	AdvisorBlack0  = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_ADVISOR_BLACK_0)
	AdvisorBlack1  = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_ADVISOR_BLACK_1)
	ElephantBlack0 = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_ELEPHANT_BLACK_0)
	ElephantBlack1 = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_ELEPHANT_BLACK_1)
	ChariotBlack0  = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_CHARIOT_BLACK_0)
	ChariotBlack1  = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_CHARIOT_BLACK_1)
	HorseBlack0    = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_HORSE_BLACK_0)
	HorseBlack1    = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_HORSE_BLACK_1)
	CannonBlack0   = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_CANNON_BLACK_0)
	CannonBlack1   = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_CANNON_BLACK_1)
	SoldierBlack0  = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_SOLDIER_BLACK_0)
	SoldierBlack1  = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_SOLDIER_BLACK_1)
	SoldierBlack2  = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_SOLDIER_BLACK_2)
	SoldierBlack3  = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_SOLDIER_BLACK_3)
	SoldierBlack4  = Piece(commongrpc.CnChessPiece_CN_CHESS_PIECE_SOLDIER_BLACK_4)
)

var names = map[Piece]string{
	InvalidPiece: "無X",
	GeneralRed:   "帥0",
	AdvisorRed0:  "仕0",
	AdvisorRed1:  "仕1",
	ElephantRed0: "相0",
	ElephantRed1: "相1",
	ChariotRed0:  "俥0",
	ChariotRed1:  "俥1",
	HorseRed0:    "傌0",
	HorseRed1:    "傌1",
	CannonRed0:   "炮0",
	CannonRed1:   "炮1",
	SoldierRed0:  "兵0",
	SoldierRed1:  "兵1",
	SoldierRed2:  "兵2",
	SoldierRed3:  "兵3",
	SoldierRed4:  "兵4",

	GeneralBlack:   "將0",
	AdvisorBlack0:  "士0",
	AdvisorBlack1:  "士1",
	ElephantBlack0: "象0",
	ElephantBlack1: "象1",
	ChariotBlack0:  "車0",
	ChariotBlack1:  "車1",
	HorseBlack0:    "馬0",
	HorseBlack1:    "馬1",
	CannonBlack0:   "包0",
	CannonBlack1:   "包1",
	SoldierBlack0:  "卒0",
	SoldierBlack1:  "卒1",
	SoldierBlack2:  "卒2",
	SoldierBlack3:  "卒3",
	SoldierBlack4:  "卒4",
}

type Piece commongrpc.CnChessPiece

func Create(color commongrpc.CnChessColorType, pieceType commongrpc.CnChessPieceType, index int) Piece {
	return Piece(int(color) | int(pieceType) | index)
}

func New(p commongrpc.CnChessPiece) Piece {
	return Piece(p)
}

func (p Piece) GetIndex() int {
	return int(p & 0x00f)
}

func (p Piece) GetType() commongrpc.CnChessPieceType {
	return commongrpc.CnChessPieceType(p & 0x0f0)
}

func (p Piece) ToProto() commongrpc.CnChessPiece {
	return commongrpc.CnChessPiece(p)
}

func (p Piece) GetName() string {
	return names[p]
}

func (p Piece) GetColor() commongrpc.CnChessColorType {
	return commongrpc.CnChessColorType(p & 0xf00)
}

func (p Piece) GetOppositeColor() commongrpc.CnChessColorType {
	if p.GetColor() == commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED {
		return commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK
	}

	return commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED
}

func (p Piece) GetHex() string {
	return strings.ToUpper(strconv.FormatInt(int64(p.ToProto()), 16))
}

func (p Piece) GetWeight() int {
	return getTypeData(p.GetType()).weight
}

func (p Piece) IsSame(b Piece) bool {
	return p.GetColor() == b.GetColor() &&
		p.GetType() == b.GetType() &&
		p.GetIndex() == b.GetIndex()
}

func (p Piece) IsInvalid() bool {
	return p.ToProto() == commongrpc.CnChessPiece_CN_CHESS_PIECE_INVALID
}

func (p Piece) IsRed() bool {
	return p.GetColor() == commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED
}

func (p Piece) IsBlack() bool {
	return p.GetColor() == commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK
}

func (p Piece) IsGeneral() bool {
	return getTypeData(p.GetType()).isGeneral()
}

func (p Piece) IsAdvisor() bool {
	return getTypeData(p.GetType()).isAdvisor()
}

func (p Piece) IsElephant() bool {
	return getTypeData(p.GetType()).isElephant()
}

func (p Piece) IsChariot() bool {
	return getTypeData(p.GetType()).isChariot()
}

func (p Piece) IsHorse() bool {
	return getTypeData(p.GetType()).isHorse()
}

func (p Piece) IsCannon() bool {
	return getTypeData(p.GetType()).isCannon()
}

func (p Piece) IsSoldier() bool {
	return getTypeData(p.GetType()).isSoldier()
}
