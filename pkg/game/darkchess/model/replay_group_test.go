package model

import (
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	"testing"
)

func Test_ChaseSamePiece(t *testing.T) {
	inst := ProvideReplayGroup()

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

	for _, v := range arr {
		t.Logf("%+v", v)
	}

	inst.Data = append(inst.Data,
		Replay{
			Uid:  "2",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.ElephantRed0,
				X:     1,
				Y:     1,
				ToX:   2,
				ToY:   1,
			},
		},
		Replay{
			Uid:  "1",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.AdvisorBlack1,
				X:     1,
				Y:     0,
				ToX:   2,
				ToY:   0,
			},
		},
		Replay{
			Uid:  "2",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.ElephantRed0,
				X:     2,
				Y:     1,
				ToX:   1,
				ToY:   1,
			},
		},
		Replay{
			Uid:  "1",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.AdvisorBlack1,
				X:     2,
				Y:     0,
				ToX:   1,
				ToY:   0,
			},
		},
		Replay{
			Uid:  "2",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.ElephantRed0,
				X:     1,
				Y:     1,
				ToX:   2,
				ToY:   1,
			},
		},
		Replay{
			Uid:  "1",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.AdvisorBlack1,
				X:     1,
				Y:     0,
				ToX:   2,
				ToY:   0,
			},
		},
		Replay{
			Uid:  "2",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.ElephantRed0,
				X:     2,
				Y:     1,
				ToX:   1,
				ToY:   1,
			},
		},
		Replay{
			Uid:  "1",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.AdvisorBlack1,
				X:     2,
				Y:     0,
				ToX:   1,
				ToY:   0,
			},
		},
		Replay{
			Uid:  "2",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.ElephantRed0,
				X:     1,
				Y:     1,
				ToX:   2,
				ToY:   1,
			},
		},
		Replay{
			Uid:  "1",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.AdvisorBlack1,
				X:     1,
				Y:     0,
				ToX:   2,
				ToY:   0,
			},
		},
		Replay{
			Uid:  "2",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.ElephantRed0,
				X:     2,
				Y:     1,
				ToX:   1,
				ToY:   1,
			},
		},
		Replay{
			Uid:  "1",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.AdvisorBlack1,
				X:     2,
				Y:     0,
				ToX:   1,
				ToY:   0,
			},
		},
		Replay{
			Uid:  "2",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.ElephantRed0,
				X:     1,
				Y:     1,
				ToX:   2,
				ToY:   1,
			},
		},
		Replay{
			Uid:  "1",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.AdvisorBlack1,
				X:     1,
				Y:     0,
				ToX:   2,
				ToY:   0,
			},
		},
		Replay{
			Uid:  "2",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.ElephantRed0,
				X:     2,
				Y:     1,
				ToX:   1,
				ToY:   1,
			},
		},
		Replay{
			Uid:  "1",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.AdvisorBlack1,
				X:     2,
				Y:     0,
				ToX:   1,
				ToY:   0,
			},
		},
		Replay{
			Uid:  "2",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.ElephantRed0,
				X:     1,
				Y:     1,
				ToX:   2,
				ToY:   1,
			},
		},
		Replay{
			Uid:  "1",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.AdvisorBlack1,
				X:     1,
				Y:     0,
				ToX:   2,
				ToY:   0,
			},
		},
		Replay{
			Uid:  "2",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.ElephantRed0,
				X:     2,
				Y:     1,
				ToX:   1,
				ToY:   1,
			},
		},
		Replay{
			Uid:  "1",
			Turn: 1,
			Move: &ActMove{
				Piece: piece.AdvisorBlack1,
				X:     2,
				Y:     0,
				ToX:   1,
				ToY:   0,
			},
		},
	)

	if inst.evalChaseSamePiece() {
		t.Log("pass")
	} else {
		t.Error("fail")
	}
}
