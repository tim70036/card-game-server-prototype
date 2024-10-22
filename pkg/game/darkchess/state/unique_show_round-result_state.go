package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
)

var GoShowRoundResultState = core.NewStateTrigger("GoShowRoundResultState",
	reflect.TypeOf(false),
	core.UidType,
)

type ShowRoundResultState struct {
	core.State

	isDraw bool
	winner core.Uid
}

func ProvideShowRoundResultState(
	stateFactory *core.StateFactory,
) *ShowRoundResultState {
	return &ShowRoundResultState{
		State: stateFactory.Create("ShowRoundResultState"),
	}
}

func (state *ShowRoundResultState) Run(_ context.Context, args ...any) error {
	var duration = constant.ShowRoundResultPeriod

	state.isDraw = args[0].(bool)
	state.winner = args[1].(core.Uid)

	state.GameController().RunTimer(duration, func() {
		state.GameController().GoNextState(GoRoundScoreboardState)
	})
	return nil
}

func (state *ShowRoundResultState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	return nil
}

func (state *ShowRoundResultState) ToProto(_ core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_ShowRoundResultStateContext{
			ShowRoundResultStateContext: &gamegrpc.ShowRoundResultStateContext{
				WinnerUid: state.winner.String(),
				IsDraw:    state.isDraw,
			},
		},
	}
}

func (state *ShowRoundResultState) BeforeLeave(uid core.Uid) error {
	state.Logger().Warn("forbid to leave during game", zap.String("uid", uid.String()))
	return status.Errorf(codes.FailedPrecondition, "forbid to leave during game")
}
