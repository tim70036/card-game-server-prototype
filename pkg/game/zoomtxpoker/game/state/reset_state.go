package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	state2 "card-game-server-prototype/pkg/game/txpoker/state"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ResetState struct {
	core.State
	hasStarted bool
}

func ProvideResetState(
	stateFactory *core.StateFactory,
) *ResetState {
	return &ResetState{
		State:      stateFactory.Create("ResetState"),
		hasStarted: false,
	}
}

func (state *ResetState) Run(ctx context.Context, args ...any) error {
	if state.hasStarted {
		state.GameController().GoNextState(state2.GoClosedState)
		return nil
	}

	state.hasStarted = true
	state.GameController().GoNextState(state2.GoStartRoundState)
	return nil
}

func (state *ResetState) Publish(ctx context.Context, args ...any) error {
	return nil
}

func (state *ResetState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context:   &txpokergrpc.GameState_ResetStateContext{ResetStateContext: &txpokergrpc.ResetStateContext{}},
	}
}
