package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/game/darkchess/service"
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

var GoCaptureState = core.NewStateTrigger("GoCaptureState",
	core.UidType,
	reflect.TypeOf(model2.Pos{}),
	reflect.TypeOf(model2.Pos{}),
)

type CaptureState struct {
	core.State

	boardService    *service.BoardService
	board           *model2.Board
	actionHintGroup *model2.ActionHintGroup
	capturedPieces  *model2.CapturedPieces
	replayGroup     *model2.ReplayGroup
	playerGroup     *model2.PlayerGroup

	oriHunterCell *gamegrpc.Cell
	oriTargetCell *gamegrpc.Cell
}

func ProvideCaptureState(
	stateFactory *core.StateFactory,

	boardService *service.BoardService,
	board *model2.Board,
	actionHintGroup *model2.ActionHintGroup,
	capturedPieces *model2.CapturedPieces,
	replayGroup *model2.ReplayGroup,
	playerGroup *model2.PlayerGroup,
) *CaptureState {
	return &CaptureState{
		State: stateFactory.Create("CaptureState"),

		boardService:    boardService,
		board:           board,
		actionHintGroup: actionHintGroup,
		capturedPieces:  capturedPieces,
		replayGroup:     replayGroup,
		playerGroup:     playerGroup,
	}
}

func (state *CaptureState) Run(_ context.Context, args ...any) error {
	var duration = constant.CapturePeriod

	actUid := args[0].(core.Uid)
	fromPos := args[1].(model2.Pos)
	toPos := args[2].(model2.Pos)

	oriHunterCell := state.board.Cells[fromPos.X][fromPos.Y]
	oriTargetCell := state.board.Cells[toPos.X][toPos.Y]

	state.oriHunterCell = &gamegrpc.Cell{
		GridPosition: &gamegrpc.GridPosition{
			X: int32(fromPos.X),
			Y: int32(fromPos.Y),
		},
		Piece:      oriHunterCell.Piece.ToProto(),
		IsRevealed: oriHunterCell.IsPieceRevealed,
		IsEmpty:    oriHunterCell.IsEmptyCell,
	}

	state.oriTargetCell = &gamegrpc.Cell{
		GridPosition: &gamegrpc.GridPosition{
			X: int32(toPos.X),
			Y: int32(toPos.Y),
		},
		Piece:      oriTargetCell.Piece.ToProto(),
		IsRevealed: oriTargetCell.IsPieceRevealed,
		IsEmpty:    oriTargetCell.IsEmptyCell,
	}

	state.boardService.CapturePiece(fromPos.X, fromPos.Y, toPos.X, toPos.Y)

	state.actionHintGroup.RepeatMovesCount = 0

	// Reset chase same piece count
	if v, ok := state.actionHintGroup.Data[actUid]; ok && v.FreezeCell != nil {
		state.actionHintGroup.Data[actUid].FreezeCell = nil
	}

	state.replayGroup.Data = append(state.replayGroup.Data, model2.Replay{
		Uid:  actUid,
		Turn: state.board.TurnCount,
		Capture: &model2.ActCapture{
			X:           fromPos.X,
			Y:           fromPos.Y,
			ToX:         toPos.X,
			ToY:         toPos.Y,
			Piece:       oriHunterCell.Piece,
			TargetPiece: oriTargetCell.Piece,
		},
	})

	state.actionHintGroup.LastAction = &model2.LastActionCell{
		Pos:             &model2.Pos{X: toPos.X, Y: toPos.Y},
		Piece:           oriHunterCell.Piece,
		IsCaptureAction: true,
	}

	state.playerGroup.Data[actUid].HasActed = false

	state.Logger().Debug("captured",
		zap.Int("turn", state.board.TurnCount),
		zap.String("uid", actUid.String()),
		zap.String("hunter", oriHunterCell.Piece.GetName()),
		zap.String("target", oriTargetCell.Piece.GetName()),
		zap.Object("from", fromPos),
		zap.Object("to", toPos),
		zap.Array("body", state.capturedPieces),
		util.DebugField(zap.Object("board", state.board)),
		util.DebugField(zap.Array("replay", state.replayGroup)),
	)

	state.GameController().RunTimer(duration, func() {
		state.GameController().GoNextState(GoEndTurnState)
	})
	return nil
}

func (state *CaptureState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	state.MsgBus().Broadcast(core.ModelTopic, &gamegrpc.Model{
		Board:          state.board.ToProto(),
		CapturedPieces: state.capturedPieces.ToProto(),
	})
	return nil
}

func (state *CaptureState) ToProto(_ core.Uid) proto.Message {
	actorUid, _ := lo.FindKeyBy(state.actionHintGroup.Data, func(_ core.Uid, p *model2.ActionHint) bool {
		return p.TurnCount == state.board.TurnCount
	})

	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_CaptureStateContext{
			CaptureStateContext: &gamegrpc.CaptureStateContext{
				ActorUid:     actorUid.String(),
				Cell:         state.oriHunterCell,
				CapturedCell: state.oriTargetCell,
			},
		},
	}
}

func (state *CaptureState) BeforeLeave(uid core.Uid) error {
	state.Logger().Warn("forbid to leave during game", zap.String("uid", uid.String()))
	return status.Errorf(codes.FailedPrecondition, "forbid to leave during game")
}
