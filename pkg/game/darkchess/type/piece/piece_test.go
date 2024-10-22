package piece_test

import (
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPiece(t *testing.T) {
	type args struct {
		Piece piece.Piece
	}
	type want struct {
		pieceType commongrpc.CnChessPieceType
		colorType commongrpc.CnChessColorType
		index     int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "GeneralRed",
			args: args{Piece: piece.GeneralRed},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_GENERAL,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
				index:     0,
			},
		},
		{
			name: "AdvisorRed0",
			args: args{Piece: piece.AdvisorRed0},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_ADVISOR,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
				index:     0,
			},
		},
		{
			name: "AdvisorRed1",
			args: args{Piece: piece.AdvisorRed1},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_ADVISOR,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
				index:     1,
			},
		},
		{
			name: "ElephantRed0",
			args: args{Piece: piece.ElephantRed0},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_ELEPHANT,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
				index:     0,
			},
		},
		{
			name: "ElephantRed1",
			args: args{Piece: piece.ElephantRed1},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_ELEPHANT,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
				index:     1,
			},
		},
		{
			name: "ChariotRed0",
			args: args{Piece: piece.ChariotRed0},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_CHARIOT,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
				index:     0,
			},
		},
		{
			name: "ChariotRed1",
			args: args{Piece: piece.ChariotRed1},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_CHARIOT,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
				index:     1,
			},
		},
		{
			name: "HorseRed0",
			args: args{Piece: piece.HorseRed0},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_HORSE,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
				index:     0,
			},
		},
		{
			name: "HorseRed1",
			args: args{Piece: piece.HorseRed1},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_HORSE,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
				index:     1,
			},
		},
		{
			name: "CannonRed0",
			args: args{Piece: piece.CannonRed0},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_CANNON,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
				index:     0,
			},
		},
		{
			name: "CannonRed1",
			args: args{Piece: piece.CannonRed1},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_CANNON,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
				index:     1,
			},
		},
		{
			name: "SoldierRed0",
			args: args{Piece: piece.SoldierRed0},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_SOLDIER,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
				index:     0,
			},
		},
		{
			name: "SoldierRed1",
			args: args{Piece: piece.SoldierRed1},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_SOLDIER,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
				index:     1,
			},
		},
		{
			name: "SoldierRed2",
			args: args{Piece: piece.SoldierRed2},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_SOLDIER,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
				index:     2,
			},
		},
		{
			name: "SoldierRed3",
			args: args{Piece: piece.SoldierRed3},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_SOLDIER,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
				index:     3,
			},
		},
		{
			name: "SoldierRed4",
			args: args{Piece: piece.SoldierRed4},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_SOLDIER,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
				index:     4,
			},
		},
		{
			name: "GeneralBlack",
			args: args{Piece: piece.GeneralBlack},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_GENERAL,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
				index:     0,
			},
		},
		{
			name: "AdvisorBlack0",
			args: args{Piece: piece.AdvisorBlack0},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_ADVISOR,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
				index:     0,
			},
		},
		{
			name: "AdvisorBlack1",
			args: args{Piece: piece.AdvisorBlack1},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_ADVISOR,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
				index:     1,
			},
		},
		{
			name: "ElephantBlack0",
			args: args{Piece: piece.ElephantBlack0},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_ELEPHANT,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
				index:     0,
			},
		},
		{
			name: "ElephantBlack1",
			args: args{Piece: piece.ElephantBlack1},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_ELEPHANT,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
				index:     1,
			},
		},
		{
			name: "ChariotBlack0",
			args: args{Piece: piece.ChariotBlack0},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_CHARIOT,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
				index:     0,
			},
		},
		{
			name: "ChariotBlack1",
			args: args{Piece: piece.ChariotBlack1},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_CHARIOT,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
				index:     1,
			},
		},
		{
			name: "HorseBlack0",
			args: args{Piece: piece.HorseBlack0},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_HORSE,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
				index:     0,
			},
		},
		{
			name: "HorseBlack1",
			args: args{Piece: piece.HorseBlack1},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_HORSE,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
				index:     1,
			},
		},
		{
			name: "CannonBlack0",
			args: args{Piece: piece.CannonBlack0},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_CANNON,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
				index:     0,
			},
		},
		{
			name: "CannonBlack1",
			args: args{Piece: piece.CannonBlack1},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_CANNON,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
				index:     1,
			},
		},
		{
			name: "SoldierBlack0",
			args: args{Piece: piece.SoldierBlack0},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_SOLDIER,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
				index:     0,
			},
		},
		{
			name: "SoldierBlack1",
			args: args{Piece: piece.SoldierBlack1},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_SOLDIER,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
				index:     1,
			},
		},
		{
			name: "SoldierBlack2",
			args: args{Piece: piece.SoldierBlack2},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_SOLDIER,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
				index:     2,
			},
		},
		{
			name: "SoldierBlack3",
			args: args{Piece: piece.SoldierBlack3},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_SOLDIER,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
				index:     3,
			},
		},
		{
			name: "SoldierBlack4",
			args: args{Piece: piece.SoldierBlack4},
			want: want{
				pieceType: commongrpc.CnChessPieceType_CN_CHESS_PIECE_TYPE_SOLDIER,
				colorType: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
				index:     4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pieceType := tt.args.Piece.GetType()
			if pieceType != tt.want.pieceType {
				t.Errorf("Piece.GetType() = %v, want %v", pieceType, tt.want.pieceType)
			}

			colorType := tt.args.Piece.GetColor()
			if colorType != tt.want.colorType {
				t.Errorf("Piece.GetColor() = %v, want %v", colorType, tt.want.colorType)
			}

			index := tt.args.Piece.GetIndex()
			if index != tt.want.index {
				t.Errorf("Piece.GetIndex() = %v, want %v", index, tt.want.index)
			}
		})
	}
}

func TestPiece_GetOppositeColor(t *testing.T) {
	type args struct {
		Piece piece.Piece
	}
	tests := []struct {
		name string
		args args
		want commongrpc.CnChessColorType
	}{
		{
			name: "GeneralRed",
			args: args{Piece: piece.GeneralRed},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
		},
		{
			name: "AdvisorRed0",
			args: args{Piece: piece.AdvisorRed0},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
		},
		{
			name: "AdvisorRed1",
			args: args{Piece: piece.AdvisorRed1},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
		},
		{
			name: "ElephantRed0",
			args: args{Piece: piece.ElephantRed0},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
		},
		{
			name: "ElephantRed1",
			args: args{Piece: piece.ElephantRed1},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
		},
		{
			name: "ChariotRed0",
			args: args{Piece: piece.ChariotRed0},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
		},
		{
			name: "ChariotRed1",
			args: args{Piece: piece.ChariotRed1},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
		},
		{
			name: "HorseRed0",
			args: args{Piece: piece.HorseRed0},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
		},
		{
			name: "HorseRed1",
			args: args{Piece: piece.HorseRed1},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK,
		},

		{
			name: "GeneralBlack",
			args: args{Piece: piece.GeneralBlack},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
		},
		{
			name: "AdvisorBlack0",
			args: args{Piece: piece.AdvisorBlack0},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
		},
		{
			name: "AdvisorBlack1",
			args: args{Piece: piece.AdvisorBlack1},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
		},
		{
			name: "ElephantBlack0",
			args: args{Piece: piece.ElephantBlack0},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
		},
		{
			name: "ElephantBlack1",
			args: args{Piece: piece.ElephantBlack1},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
		},
		{
			name: "ChariotBlack0",
			args: args{Piece: piece.ChariotBlack0},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
		},
		{
			name: "ChariotBlack1",
			args: args{Piece: piece.ChariotBlack1},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
		},
		{
			name: "HorseBlack0",
			args: args{Piece: piece.HorseBlack0},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
		},
		{
			name: "HorseBlack1",
			args: args{Piece: piece.HorseBlack1},
			want: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Piece.GetOppositeColor())
		})
	}
}

func TestPiece_IsSame(t *testing.T) {
	type args struct {
		a, b piece.Piece
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "GeneralRed",
			args: args{a: piece.GeneralRed, b: piece.GeneralRed},
			want: true,
		},
		{
			name: "AdvisorRed0",
			args: args{a: piece.AdvisorRed0, b: piece.AdvisorRed0},
			want: true,
		},
		{
			name: "AdvisorRed1",
			args: args{a: piece.AdvisorRed1, b: piece.AdvisorRed1},
			want: true,
		},
		{
			name: "ElephantRed0",
			args: args{a: piece.ElephantRed0, b: piece.ElephantRed0},
			want: true,
		},
		{
			name: "ElephantRed1",
			args: args{a: piece.ElephantRed1, b: piece.ElephantRed1},
			want: true,
		},
		{
			name: "ChariotRed0",
			args: args{a: piece.ChariotRed0, b: piece.ChariotRed0},
			want: true,
		},
		{
			name: "ChariotRed1",
			args: args{a: piece.ChariotRed1, b: piece.ChariotRed1},
			want: true,
		},
		{
			name: "HorseRed0",
			args: args{a: piece.HorseRed0, b: piece.HorseRed0},
			want: true,
		},
		{
			name: "HorseRed1",
			args: args{a: piece.HorseRed1, b: piece.HorseRed1},
			want: true,
		},
		{
			name: "CannonRed0",
			args: args{a: piece.CannonRed0, b: piece.CannonRed0},
			want: true,
		},
		{
			name: "CannonRed1",
			args: args{a: piece.CannonRed1, b: piece.CannonRed1},
			want: true,
		},
		{
			name: "SoldierRed0",
			args: args{a: piece.SoldierRed0, b: piece.SoldierRed0},
			want: true,
		},
		{
			name: "SoldierRed1",
			args: args{a: piece.SoldierRed1, b: piece.SoldierRed1},
			want: true,
		},
		{
			name: "SoldierRed2",
			args: args{a: piece.SoldierRed2, b: piece.SoldierRed2},
			want: true,
		},
		{
			name: "SoldierRed3",
			args: args{a: piece.SoldierRed3, b: piece.SoldierRed3},
			want: true,
		},
		{
			name: "SoldierRed4",
			args: args{a: piece.SoldierRed4, b: piece.SoldierRed4},
			want: true,
		},
		{
			name: "GeneralBlack",
			args: args{a: piece.GeneralBlack, b: piece.GeneralBlack},
			want: true,
		},
		{
			name: "AdvisorBlack0",
			args: args{a: piece.AdvisorBlack0, b: piece.AdvisorBlack0},
			want: true,
		},
		{
			name: "AdvisorBlack1",
			args: args{a: piece.AdvisorBlack1, b: piece.AdvisorBlack1},
			want: true,
		},
		{
			name: "ElephantBlack0",
			args: args{a: piece.ElephantBlack0, b: piece.ElephantBlack0},
			want: true,
		},
		{
			name: "ElephantBlack1",
			args: args{a: piece.ElephantBlack1, b: piece.ElephantBlack1},
			want: true,
		},
		{
			name: "ChariotBlack0",
			args: args{a: piece.ChariotBlack0, b: piece.ChariotBlack0},
			want: true,
		},
		{
			name: "ChariotBlack1",
			args: args{a: piece.ChariotBlack1, b: piece.ChariotBlack1},
			want: true,
		},
		{
			name: "HorseBlack0",
			args: args{a: piece.HorseBlack0, b: piece.HorseBlack0},
			want: true,
		},
		{
			name: "HorseBlack1",
			args: args{a: piece.HorseBlack1, b: piece.HorseBlack1},
			want: true,
		},
		{
			name: "CannonBlack0",
			args: args{a: piece.CannonBlack0, b: piece.CannonBlack0},
			want: true,
		},
		{
			name: "CannonBlack1",
			args: args{a: piece.CannonBlack1, b: piece.CannonBlack1},
			want: true,
		},
		{
			name: "SoldierBlack0",
			args: args{a: piece.SoldierBlack0, b: piece.SoldierBlack0},
			want: true,
		},
		{
			name: "SoldierBlack1",
			args: args{a: piece.SoldierBlack1, b: piece.SoldierBlack1},
			want: true,
		},
		{
			name: "SoldierBlack2",
			args: args{a: piece.SoldierBlack2, b: piece.SoldierBlack2},
			want: true,
		},
		{
			name: "SoldierBlack3",
			args: args{a: piece.SoldierBlack3, b: piece.SoldierBlack3},
			want: true,
		},
		{
			name: "SoldierBlack4",
			args: args{a: piece.SoldierBlack4, b: piece.SoldierBlack4},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.a.IsSame(tt.args.b))
		})
	}
}
