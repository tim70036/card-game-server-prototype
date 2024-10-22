package piece_test

import (
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPiece_IsGeneral(t *testing.T) {
	type args struct {
		a piece.Piece
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "GeneralRed",
			args: args{a: piece.GeneralRed},
			want: true,
		},
		{
			name: "AdvisorRed0",
			args: args{a: piece.AdvisorRed0},
			want: false,
		},
		{
			name: "AdvisorRed1",
			args: args{a: piece.AdvisorRed1},
			want: false,
		},
		{
			name: "ElephantRed0",
			args: args{a: piece.ElephantRed0},
			want: false,
		},
		{
			name: "ElephantRed1",
			args: args{a: piece.ElephantRed1},
			want: false,
		},
		{
			name: "ChariotRed0",
			args: args{a: piece.ChariotRed0},
			want: false,
		},
		{
			name: "ChariotRed1",
			args: args{a: piece.ChariotRed1},
			want: false,
		},
		{
			name: "HorseRed0",
			args: args{a: piece.HorseRed0},
			want: false,
		},
		{
			name: "HorseRed1",
			args: args{a: piece.HorseRed1},
			want: false,
		},
		{
			name: "CannonRed0",
			args: args{a: piece.CannonRed0},
			want: false,
		},
		{
			name: "CannonRed1",
			args: args{a: piece.CannonRed1},
			want: false,
		},
		{
			name: "SoldierRed0",
			args: args{a: piece.SoldierRed0},
			want: false,
		},
		{
			name: "SoldierRed1",
			args: args{a: piece.SoldierRed1},
			want: false,
		},
		{
			name: "SoldierRed2",
			args: args{a: piece.SoldierRed2},
			want: false,
		},
		{
			name: "SoldierRed3",
			args: args{a: piece.SoldierRed3},
			want: false,
		},
		{
			name: "SoldierRed4",
			args: args{a: piece.SoldierRed4},
			want: false,
		},
		{
			name: "GeneralBlack",
			args: args{a: piece.GeneralBlack},
			want: true,
		},
		{
			name: "AdvisorBlack0",
			args: args{a: piece.AdvisorBlack0},
			want: false,
		},
		{
			name: "AdvisorBlack1",
			args: args{a: piece.AdvisorBlack1},
			want: false,
		},
		{
			name: "ElephantBlack0",
			args: args{a: piece.ElephantBlack0},
			want: false,
		},
		{
			name: "ElephantBlack1",
			args: args{a: piece.ElephantBlack1},
			want: false,
		},
		{
			name: "ChariotBlack0",
			args: args{a: piece.ChariotBlack0},
			want: false,
		},
		{
			name: "ChariotBlack1",
			args: args{a: piece.ChariotBlack1},
			want: false,
		},
		{
			name: "HorseBlack0",
			args: args{a: piece.HorseBlack0},
			want: false,
		},
		{
			name: "HorseBlack1",
			args: args{a: piece.HorseBlack1},
			want: false,
		},
		{
			name: "CannonBlack0",
			args: args{a: piece.CannonBlack0},
			want: false,
		},
		{
			name: "CannonBlack1",
			args: args{a: piece.CannonBlack1},
			want: false,
		},
		{
			name: "SoldierBlack0",
			args: args{a: piece.SoldierBlack0},
			want: false,
		},
		{
			name: "SoldierBlack1",
			args: args{a: piece.SoldierBlack1},
			want: false,
		},
		{
			name: "SoldierBlack2",
			args: args{a: piece.SoldierBlack2},
			want: false,
		},
		{
			name: "SoldierBlack3",
			args: args{a: piece.SoldierBlack3},
			want: false,
		},
		{
			name: "SoldierBlack4",
			args: args{a: piece.SoldierBlack4},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.a.IsGeneral())
		})
	}
}

func TestPiece_IsAdvisor(t *testing.T) {
	type args struct {
		a piece.Piece
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "GeneralRed",
			args: args{a: piece.GeneralRed},
			want: false,
		},
		{
			name: "AdvisorRed0",
			args: args{a: piece.AdvisorRed0},
			want: true,
		},
		{
			name: "AdvisorRed1",
			args: args{a: piece.AdvisorRed1},
			want: true,
		},
		{
			name: "ElephantRed0",
			args: args{a: piece.ElephantRed0},
			want: false,
		},
		{
			name: "ElephantRed1",
			args: args{a: piece.ElephantRed1},
			want: false,
		},
		{
			name: "ChariotRed0",
			args: args{a: piece.ChariotRed0},
			want: false,
		},
		{
			name: "ChariotRed1",
			args: args{a: piece.ChariotRed1},
			want: false,
		},
		{
			name: "HorseRed0",
			args: args{a: piece.HorseRed0},
			want: false,
		},
		{
			name: "HorseRed1",
			args: args{a: piece.HorseRed1},
			want: false,
		},
		{
			name: "CannonRed0",
			args: args{a: piece.CannonRed0},
			want: false,
		},
		{
			name: "CannonRed1",
			args: args{a: piece.CannonRed1},
			want: false,
		},
		{
			name: "SoldierRed0",
			args: args{a: piece.SoldierRed0},
			want: false,
		},
		{
			name: "SoldierRed1",
			args: args{a: piece.SoldierRed1},
			want: false,
		},
		{
			name: "SoldierRed2",
			args: args{a: piece.SoldierRed2},
			want: false,
		},
		{
			name: "SoldierRed3",
			args: args{a: piece.SoldierRed3},
			want: false,
		},
		{
			name: "SoldierRed4",
			args: args{a: piece.SoldierRed4},
			want: false,
		},
		{
			name: "GeneralBlack",
			args: args{a: piece.GeneralBlack},
			want: false,
		},
		{
			name: "AdvisorBlack0",
			args: args{a: piece.AdvisorBlack0},
			want: true,
		},
		{
			name: "AdvisorBlack1",
			args: args{a: piece.AdvisorBlack1},
			want: true,
		},
		{
			name: "ElephantBlack0",
			args: args{a: piece.ElephantBlack0},
			want: false,
		},
		{
			name: "ElephantBlack1",
			args: args{a: piece.ElephantBlack1},
			want: false,
		},
		{
			name: "ChariotBlack0",
			args: args{a: piece.ChariotBlack0},
			want: false,
		},
		{
			name: "ChariotBlack1",
			args: args{a: piece.ChariotBlack1},
			want: false,
		},
		{
			name: "HorseBlack0",
			args: args{a: piece.HorseBlack0},
			want: false,
		},
		{
			name: "HorseBlack1",
			args: args{a: piece.HorseBlack1},
			want: false,
		},
		{
			name: "CannonBlack0",
			args: args{a: piece.CannonBlack0},
			want: false,
		},
		{
			name: "CannonBlack1",
			args: args{a: piece.CannonBlack1},
			want: false,
		},
		{
			name: "SoldierBlack0",
			args: args{a: piece.SoldierBlack0},
			want: false,
		},
		{
			name: "SoldierBlack1",
			args: args{a: piece.SoldierBlack1},
			want: false,
		},
		{
			name: "SoldierBlack2",
			args: args{a: piece.SoldierBlack2},
			want: false,
		},
		{
			name: "SoldierBlack3",
			args: args{a: piece.SoldierBlack3},
			want: false,
		},
		{
			name: "SoldierBlack4",
			args: args{a: piece.SoldierBlack4},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.a.IsGeneral())
		})
	}
}
