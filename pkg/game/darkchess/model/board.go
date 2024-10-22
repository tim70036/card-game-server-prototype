package model

import (
	"fmt"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"go.uber.org/zap/zapcore"
)

type Board struct {
	Cells     [][]*Cell
	TurnCount int
	IsDraw    bool
}

func ProvideBoard() *Board {
	return &Board{
		Cells: [][]*Cell{},
	}
}

func (m *Board) ToProto() *gamegrpc.Board {
	if len(m.Cells) < constant.BoardWidth {
		return &gamegrpc.Board{}
	}

	for _, row := range m.Cells {
		if len(row) < constant.BoardHeight {
			return &gamegrpc.Board{}
		}
	}

	var cells []*gamegrpc.Cell

	// 2D array
	// [
	//   [ [0,0], [0,1], [0,2], [0,3] ],
	//   [ [1,0], [1,1], [1,2], [1,3] ],
	//   [ [2,0], [2,1], [2,2], [2,3] ],
	//   [ [3,0], [3,1], [3,2], [3,3] ],
	//   [ [4,0], [4,1], [4,2], [4,3] ],
	//   [ [5,0], [5,1], [5,2], [5,3] ],
	//   [ [6,0], [6,1], [6,2], [6,3] ],
	//   [ [7,0], [7,1], [7,2], [7,3] ]
	// ]

	// 1D array (to frontend)
	// 03 13 23 33 43 53 63 73
	// -> 02 12 22 32 42 52 62 72
	// -> 01 11 21 31 41 51 61 71
	// -> 00 10 20 30 40 50 60 70

	for x := 0; x < constant.BoardWidth; x++ {
		for y := 0; y < constant.BoardHeight; y++ {
			cells = append(cells, &gamegrpc.Cell{
				GridPosition: &gamegrpc.GridPosition{
					X: int32(x),
					Y: int32(y),
				},
				Piece:      m.Cells[x][y].Piece.ToProto(),
				IsRevealed: m.Cells[x][y].IsPieceRevealed,
				IsEmpty:    m.Cells[x][y].IsEmptyCell,
			})
		}
	}

	return &gamegrpc.Board{
		Cells: cells,
	}
}

func (m *Board) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("turnCount", m.TurnCount)
	enc.AddBool("isDraw", m.IsDraw)

	_ = enc.AddArray("Cells", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		for _, row := range m.Cells {
			_ = enc.AppendArray(zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
				for _, cell := range row {

					if cell.IsPieceRevealed {
						enc.AppendString(fmt.Sprintf("R%s", cell.Piece.GetName()))
					} else {
						enc.AppendString(fmt.Sprintf("_%s", cell.Piece.GetName()))
					}
				}
				return nil
			}))
		}
		return nil
	}))
	return nil
}
