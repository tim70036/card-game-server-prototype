package service

import (
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/util"
	"testing"
)

func TestBoardService_ValidPos(t *testing.T) {
	inst := ProvideBoardService(util.ProvideLoggerFactory(config.LogCFG), nil, nil, nil)
	arr := [][]string{
		{"00", "01", "02", "03"},
		{"10", "11", "12", "13"},
		{"20", "21", "22", "23"},
		{"30", "31", "32", "33"},
		{"40", "41", "42", "43"},
		{"50", "51", "52", "53"},
		{"60", "61", "62", "63"},
		{"70", "71", "72", "73"},
	}

	var getPos = func(s string, arr [][]string) model.Pos {
		for x, row := range arr {
			for y, v := range row {
				if s == v {
					return model.Pos{X: x, Y: y}
				}
			}
		}
		return model.Pos{X: -1, Y: -1}
	}

	type args struct {
		a, b string
		n    int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{a: "00", b: "05", n: 1},
		},
		{
			name: "2",
			args: args{a: "00", b: "05", n: 2},
		},
		{
			name: "3",
			args: args{a: "33", b: "31", n: 2},
		},
		{
			name: "4",
			args: args{a: "33", b: "31", n: 2},
		},
		{
			name: "5",
			args: args{a: "00", b: "70", n: 1},
		},
		{
			name: "6",
			args: args{a: "53", b: "50", n: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for _, row := range arr {
				t.Logf("%+v", row)
			}

			t.Logf("a: %v, b: %v", tt.args.a, tt.args.b)

			a := getPos(tt.args.a, arr)
			b := getPos(tt.args.b, arr)
			t.Logf("a, b isInAxis: %v", inst.isInAxis(a.X, a.Y, b.X, b.Y))
			t.Logf("a, b IsNextInAxis: %v", inst.IsNextInAxis(a.X, a.Y, b.X, b.Y))
			t.Logf("a, b isSamePos: %v", inst.isSamePos(a.X, a.Y, b.X, b.Y))

			{
				v, ok := inst.moveRightNCells(a.X, tt.args.n)
				t.Logf("%v R %v cell is %v, isMoved:%v", tt.args.a, tt.args.n, arr[v][a.Y], ok)
			}

			{
				v, ok := inst.moveLeftNCells(a.X, tt.args.n)
				t.Logf("%v L %v cell is %v, isMoved:%v", tt.args.a, tt.args.n, arr[v][a.Y], ok)
			}

			{
				v, ok := inst.moveUpNCells(a.Y, tt.args.n)
				t.Logf("%v U %v cell is %v, isMoved:%v", tt.args.a, tt.args.n, arr[a.X][v], ok)
			}

			{
				v, ok := inst.moveDownNCells(a.Y, tt.args.n)
				t.Logf("%v D %v cell is %v, isMoved:%v", tt.args.a, tt.args.n, arr[a.X][v], ok)
			}

			{
				if inst.isInAxis(a.X, a.Y, b.X, b.Y) && !inst.isSamePos(a.X, a.Y, b.X, b.Y) {
					for _, v := range inst.posBetween(a.X, a.Y, b.X, b.Y) {
						t.Logf("[%v]", arr[v.X][v.Y])
					}
				}
			}
		})
	}
}
