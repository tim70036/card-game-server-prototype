package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/game/darkchess/service"
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"card-game-server-prototype/pkg/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
)

var GoMoveState = core.NewStateTrigger("GoMoveState",
	core.UidType,
	reflect.TypeOf(model2.Pos{}),
	reflect.TypeOf(model2.Pos{}),
)

type MoveState struct {
	core.State

	gameInfo        *model2.GameInfo
	boardService    *service.BoardService
	actionHintGroup *model2.ActionHintGroup
	board           *model2.Board
	replayGroup     *model2.ReplayGroup
	playerGroup     *model2.PlayerGroup

	oriMoveCell *gamegrpc.Cell
	toPos       model2.Pos
}

func ProvideMoveState(
	stateFactory *core.StateFactory,

	gameInfo *model2.GameInfo,
	boardService *service.BoardService,
	actionHintGroup *model2.ActionHintGroup,
	board *model2.Board,
	replayGroup *model2.ReplayGroup,
	playerGroup *model2.PlayerGroup,
) *MoveState {
	return &MoveState{
		State: stateFactory.Create("MoveState"),

		gameInfo:        gameInfo,
		boardService:    boardService,
		actionHintGroup: actionHintGroup,
		board:           board,
		replayGroup:     replayGroup,
		playerGroup:     playerGroup,
	}
}

func (state *MoveState) Run(_ context.Context, args ...any) error {
	var duration = constant.MovePeriod

	actUid := args[0].(core.Uid)
	fromPos := args[1].(model2.Pos)
	state.toPos = args[2].(model2.Pos)

	oriMoveCell := state.board.Cells[fromPos.X][fromPos.Y]
	state.oriMoveCell = &gamegrpc.Cell{
		GridPosition: &gamegrpc.GridPosition{
			X: int32(fromPos.X),
			Y: int32(fromPos.Y),
		},
		Piece:      oriMoveCell.Piece.ToProto(),
		IsRevealed: oriMoveCell.IsPieceRevealed,
		IsEmpty:    oriMoveCell.IsEmptyCell,
	}

	state.boardService.MovePiece(fromPos.X, fromPos.Y, state.toPos.X, state.toPos.Y)

	// 記錄炮與上一步的直線距離中的其他棋子數量
	var countPieces int
	if oriMoveCell.Piece.IsCannon() {
		if lastAction, err := lo.Last(state.replayGroup.Data); err == nil && lastAction.Move != nil {
			x := state.toPos.X
			y := state.toPos.Y
			x2 := lastAction.Move.ToX
			y2 := lastAction.Move.ToY

			if err := state.boardService.ValidCanonAndTargetPos(x, y, x2, y2); err == nil {
				countPieces = state.boardService.CountPiecesBetween(x, y, x2, y2)
			}
		}
	}

	state.replayGroup.Data = append(state.replayGroup.Data, model2.Replay{
		Uid:  actUid,
		Turn: state.board.TurnCount,
		Move: &model2.ActMove{
			X:           fromPos.X,
			Y:           fromPos.Y,
			ToX:         state.toPos.X,
			ToY:         state.toPos.Y,
			Piece:       oriMoveCell.Piece,
			CountPieces: countPieces,
		},
	})

	state.actionHintGroup.LastAction = &model2.LastActionCell{
		Pos:          &model2.Pos{X: state.toPos.X, Y: state.toPos.Y},
		Piece:        oriMoveCell.Piece,
		IsMoveAction: true,
	}

	// Reset chase same piece count
	// 和 reveal、 capture 的差別是，要多判斷移動的棋不是 freeze piece 才可以重置。
	if v, ok := state.actionHintGroup.Data[actUid]; ok && v.FreezeCell != nil &&
		!oriMoveCell.Piece.IsSame(v.FreezeCell.Piece) {
		state.actionHintGroup.Data[actUid].FreezeCell = nil
	}

	// 統計連續空步, 並且在這裡判斷是否強制和棋。
	state.evalRepeatMoves()

	// 長捉禁手
	state.evalChaseSamePiece()

	state.playerGroup.Data[actUid].HasActed = false

	state.Logger().Debug("moved",
		zap.Int("turn", state.board.TurnCount),
		zap.String("uid", actUid.String()),
		zap.String("piece", oriMoveCell.Piece.GetName()),
		zap.Object("from", fromPos),
		zap.Object("to", state.toPos),
		util.DebugField(zap.Object("board", state.board)),
		util.DebugField(zap.Array("replay", state.replayGroup)),
	)

	state.GameController().RunTimer(duration, func() {
		state.GameController().GoNextState(GoEndTurnState)
	})
	return nil
}

func (state *MoveState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	state.MsgBus().Broadcast(core.ModelTopic, &gamegrpc.Model{
		Board: state.board.ToProto(),
	})
	return nil
}

func (state *MoveState) ToProto(_ core.Uid) proto.Message {
	actorUid, _ := lo.FindKeyBy(state.actionHintGroup.Data, func(_ core.Uid, p *model2.ActionHint) bool {
		return p.TurnCount == state.board.TurnCount
	})

	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_MoveStateContext{
			MoveStateContext: &gamegrpc.MoveStateContext{
				ActorUid: actorUid.String(),
				Cell:     state.oriMoveCell,
				ToGridPosition: &gamegrpc.GridPosition{
					X: int32(state.toPos.X),
					Y: int32(state.toPos.Y),
				},
			},
		},
	}
}

func (state *MoveState) BeforeLeave(uid core.Uid) error {
	state.Logger().Warn("forbid to leave during game", zap.String("uid", uid.String()))
	return status.Errorf(codes.FailedPrecondition, "forbid to leave during game")
}

func (state *MoveState) evalRepeatMoves() {
	var countRepeatMove int

	for i := len(state.replayGroup.Data) - 1; i >= 0; i-- {
		if state.replayGroup.Data[i].Move == nil {
			break
		}
		countRepeatMove++
	}

	// reset in reveal state and capture state
	state.actionHintGroup.RepeatMovesCount = countRepeatMove

	// 連續空步超過 50 步，視為和棋
	if countRepeatMove >= state.gameInfo.Setting.MaxRepeatMoves {
		state.Logger().Debug("repeatMoves Max")
		state.board.IsDraw = true
	}
}

// 追擊對方同一顆棋子連續6次
// 1. 12 個連續 move
// 2. 2 個 piece
// 3. Piece 一大一小
// 4. 判斷追擊者看後手是否有機會吃掉前手
func (state *MoveState) evalChaseSamePiece() {
	var chaseActs []*model2.ActMove

	// replayGroup: old step -> latest step

	for i := len(state.replayGroup.Data) - 1; i >= 0; i-- {
		// Last one
		if i == 0 {
			break
		}

		cur := state.replayGroup.Data[i]
		prev := state.replayGroup.Data[i-1]

		if cur.Move == nil || prev.Move == nil {
			break
		}

		if state.replayGroup.Data[i].Move == nil {
			break
		}

		chaseActs = append(chaseActs, state.replayGroup.Data[i].Move)

		if len(chaseActs) == state.gameInfo.Setting.MaxChaseSamePiece*2 {
			break
		}
	}

	uniqMovePiece := lo.UniqBy(chaseActs, func(item *model2.ActMove) commongrpc.CnChessPiece {
		return item.Piece.ToProto()
	})

	var countChaseSamePiece int
	if len(chaseActs) == state.gameInfo.Setting.MaxChaseSamePiece*2 && len(uniqMovePiece) == 2 {

		// 給 canon 用的
		var prevDistance int

		// chaseActs: latest step -> old step
		for i := range chaseActs {
			// 後手/前手: 0/1, 2/3, 4/5, 6/7, 8/9, 10/11 共 6 組，檢查後手是否可吃前手。
			if i%2 == 1 {
				// 只需要處理後手
				continue
			}

			// 後手
			curAct := chaseActs[i]
			curP := piece.New(curAct.Piece.ToProto())

			// 前手
			prevAct := chaseActs[i+1]
			prevP := piece.New(prevAct.Piece.ToProto())

			if !state.boardService.IsChaseSamePiece(curP, prevP) {
				break
			}

			if curP.IsCannon() && curAct.CountPieces == 1 {

				curDistance := state.boardService.DistanceBetween(curAct.ToX, curAct.ToY, prevAct.ToX, prevAct.ToY)
				countChaseSamePiece++

				if i > 0 && prevDistance != curDistance {
					countChaseSamePiece = 0
				}

				prevDistance = curDistance

			} else {

				// curAct final position must nearby prev final position
				if state.boardService.IsNextInAxis(curAct.ToX, curAct.ToY, prevAct.ToX, prevAct.ToY) {
					countChaseSamePiece++
				}

			}

			if countChaseSamePiece == state.gameInfo.Setting.MaxChaseSamePiece {
				if usePlayer, ok := lo.FindKeyBy(state.playerGroup.Data, func(_ core.Uid, p *model2.Player) bool {
					return p.Color == curAct.Piece.GetColor()
				}); ok {

					state.actionHintGroup.Data[usePlayer].FreezeCell = &model2.FreezeCell{
						Piece: curAct.Piece,
					}

					if x, y, ok := state.boardService.GetCellPos(curAct.Piece); ok {
						state.actionHintGroup.Data[usePlayer].FreezeCell.Pos = &model2.Pos{X: x, Y: y}
						state.actionHintGroup.Data[usePlayer].FreezeCell.IsPieceRevealed = state.board.Cells[x][y].IsPieceRevealed
						state.actionHintGroup.Data[usePlayer].FreezeCell.IsEmptyCell = state.board.Cells[x][y].IsEmptyCell
					}

					state.Logger().Warn("FreezePiece",
						zap.String("uid", usePlayer.String()),
						zap.Object("freezePiece", state.actionHintGroup.Data[usePlayer].FreezeCell),
					)
				}
			}
		}
	}
}
