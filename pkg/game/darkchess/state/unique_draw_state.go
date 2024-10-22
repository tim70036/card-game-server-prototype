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
)

var GoDrawState = core.NewStateTrigger("GoDrawState")

type DrawState struct {
	core.State
}

func ProvideDrawState(
	stateFactory *core.StateFactory,
) *DrawState {
	return &DrawState{
		State: stateFactory.Create("DrawState"),
	}
}

func (state *DrawState) Run(context.Context, ...any) error {
	var duration = constant.DrawPeriod

	state.GameController().RunTimer(duration, func() {
		state.GameController().GoNextState(GoShowRoundResultState, true, core.Uid(""))
	})
	return nil
}

func (state *DrawState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	return nil
}

func (state *DrawState) ToProto(_ core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_DrawStateContext{
			DrawStateContext: &gamegrpc.DrawStateContext{},
		},
	}
}

func (state *DrawState) BeforeLeave(uid core.Uid) error {
	state.Logger().Warn("forbid to leave during game", zap.String("uid", uid.String()))
	return status.Errorf(codes.FailedPrecondition, "forbid to leave during game")
}
