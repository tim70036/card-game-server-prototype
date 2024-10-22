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

var GoRevealState = core.NewStateTrigger("GoRevealState",
	core.UidType,
	reflect.TypeOf(model2.Pos{}),
)

type RevealState struct {
	core.State

	boardService    *service.BoardService
	playerGroup     *model2.PlayerGroup
	actionHintGroup *model2.ActionHintGroup
	board           *model2.Board
	replayGroup     *model2.ReplayGroup

	revealPos model2.Pos
}

func ProvideRevealState(
	stateFactory *core.StateFactory,

	boardService *service.BoardService,
	playerGroup *model2.PlayerGroup,
	actionHintGroup *model2.ActionHintGroup,
	board *model2.Board,
	replayGroup *model2.ReplayGroup,
) *RevealState {
	return &RevealState{
		State: stateFactory.Create("RevealState"),

		boardService:    boardService,
		playerGroup:     playerGroup,
		actionHintGroup: actionHintGroup,
		board:           board,
		replayGroup:     replayGroup,
	}
}

func (state *RevealState) Run(_ context.Context, args ...any) error {
	var duration = constant.RevealPeriod

	actUid := args[0].(core.Uid)
	state.revealPos = args[1].(model2.Pos)

	oriRevealCell := state.board.Cells[state.revealPos.X][state.revealPos.Y]

	state.boardService.RevealPiece(state.revealPos.X, state.revealPos.Y)

	state.actionHintGroup.RepeatMovesCount = 0

	// Set user color if not set
	state.setPlayersUseColor(actUid, oriRevealCell.Piece)

	// Reset chase same piece count
	if v, ok := state.actionHintGroup.Data[actUid]; ok && v.FreezeCell != nil {
		state.actionHintGroup.Data[actUid].FreezeCell = nil
	}

	state.replayGroup.Data = append(state.replayGroup.Data, model2.Replay{
		Uid:  actUid,
		Turn: state.board.TurnCount,
		Reveal: &model2.ActReveal{
			X:     state.revealPos.X,
			Y:     state.revealPos.Y,
			Piece: oriRevealCell.Piece,
		},
	})

	state.actionHintGroup.LastAction = &model2.LastActionCell{
		Pos:            &model2.Pos{X: state.revealPos.X, Y: state.revealPos.Y},
		Piece:          oriRevealCell.Piece,
		IsRevealAction: true,
	}

	state.playerGroup.Data[actUid].HasActed = false

	state.Logger().Debug("revealed",
		zap.Int("turn", state.board.TurnCount),
		zap.String("uid", actUid.String()),
		zap.String("piece", oriRevealCell.Piece.GetName()),
		zap.Object("pos", state.revealPos),
		util.DebugField(zap.Object("board", state.board)),
		util.DebugField(zap.Array("replay", state.replayGroup)),
	)

	state.GameController().RunTimer(duration, func() {
		state.GameController().GoNextState(GoEndTurnState)
	})
	return nil
}

func (state *RevealState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	state.MsgBus().Broadcast(core.ModelTopic, &gamegrpc.Model{
		Board:       state.board.ToProto(),
		PlayerGroup: state.playerGroup.ToProto(),
	})
	return nil
}

func (state *RevealState) ToProto(_ core.Uid) proto.Message {
	actorUid, _ := lo.FindKeyBy(state.actionHintGroup.Data, func(_ core.Uid, p *model2.ActionHint) bool {
		return p.TurnCount == state.board.TurnCount
	})

	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_RevealStateContext{
			RevealStateContext: &gamegrpc.RevealStateContext{
				ActorUid: actorUid.String(),
				Cell: &gamegrpc.Cell{
					GridPosition: &gamegrpc.GridPosition{
						X: int32(state.revealPos.X),
						Y: int32(state.revealPos.Y),
					},
					Piece:      state.board.Cells[state.revealPos.X][state.revealPos.Y].Piece.ToProto(),
					IsRevealed: state.board.Cells[state.revealPos.X][state.revealPos.Y].IsPieceRevealed,
					IsEmpty:    state.board.Cells[state.revealPos.X][state.revealPos.Y].IsEmptyCell,
				},
			},
		},
	}
}

func (state *RevealState) BeforeLeave(uid core.Uid) error {
	state.Logger().Warn("forbid to leave during game", zap.String("uid", uid.String()))
	return status.Errorf(codes.FailedPrecondition, "forbid to leave during game")
}

func (state *RevealState) setPlayersUseColor(uid core.Uid, pickPiece piece.Piece) {
	if v, ok := state.playerGroup.Data[uid]; ok && v.Color == commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_INVALID {
		anotherPlayer := lo.Without(lo.Keys(state.playerGroup.Data), uid)[0]

		state.playerGroup.Data[uid].Color = pickPiece.GetColor()
		state.playerGroup.Data[anotherPlayer].Color = pickPiece.GetOppositeColor()
	}
}
