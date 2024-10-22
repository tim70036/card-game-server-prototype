package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/game/darkchess/service"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoStartTurnState = core.NewStateTrigger("GoStartTurnState")

type StartTurnState struct {
	core.State

	boardService    *service.BoardService
	actionHintGroup *model2.ActionHintGroup
	board           *model2.Board

	actUid core.Uid
}

func ProvideStartTurnState(
	stateFactory *core.StateFactory,

	boardService *service.BoardService,
	actionHintGroup *model2.ActionHintGroup,
	board *model2.Board,
) *StartTurnState {
	return &StartTurnState{
		State: stateFactory.Create("StartTurnState"),

		boardService:    boardService,
		actionHintGroup: actionHintGroup,
		board:           board,
	}
}

func (state *StartTurnState) Run(_ context.Context, _ ...any) error {
	state.actUid, _ = lo.FindKeyBy(state.actionHintGroup.Data, func(_ core.Uid, p *model2.ActionHint) bool {
		return p.TurnCount == state.board.TurnCount
	})

	if state.actUid == "" {
		state.Logger().Warn("no actUid found",
			zap.Int("turn", state.board.TurnCount),
			zap.Object("actionHint", state.actionHintGroup),
		)
		return status.Errorf(codes.Internal, "no actUid found")
	}

	state.Logger().Debug("startTurn",
		zap.Int("turn", state.board.TurnCount),
		zap.String("actUid", state.actUid.String()),
	)

	if state.board.TurnCount == 1 {
		state.GameController().RunTimer(constant.StartTurnPeriod, func() {
			state.GameController().GoNextState(GoWaitActionState, state.actUid)
		})
	} else {
		state.GameController().GoNextState(GoWaitActionState, state.actUid)
	}
	return nil
}

func (state *StartTurnState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	state.MsgBus().Broadcast(core.ModelTopic, &gamegrpc.Model{
		Board: state.board.ToProto(),
	})
	return nil
}

func (state *StartTurnState) ToProto(_ core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_StartTurnStateContext{
			StartTurnStateContext: &gamegrpc.StartTurnStateContext{
				ActorUid:  state.actUid.String(),
				TurnCount: int32(state.board.TurnCount),
			},
		},
	}
}

func (state *StartTurnState) BeforeLeave(uid core.Uid) error {
	state.Logger().Warn("forbid to leave during game", zap.String("uid", uid.String()))
	return status.Errorf(codes.FailedPrecondition, "forbid to leave during game")
}
