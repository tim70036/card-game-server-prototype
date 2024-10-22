package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoDeclareShowdownState = &core.StateTrigger{
	Name:      "GoDeclareShowdownState",
	ArgsTypes: []reflect.Type{},
}

type DeclareShowdownState struct {
	core.State
}

func ProvideDeclareShowdownState(
	stateFactory *core.StateFactory,
) *DeclareShowdownState {
	return &DeclareShowdownState{
		State: stateFactory.Create("DeclareShowdownState"),
	}
}

func (state *DeclareShowdownState) Run(ctx context.Context, args ...any) error {
	state.GameController().RunTimer(1*time.Second, func() {
		state.GameController().GoNextState(GoShowdownState)
	})
	return nil
}

func (state *DeclareShowdownState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
	})
	return nil
}

func (state *DeclareShowdownState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context:   &txpokergrpc.GameState_DeclareShowdownStateContext{DeclareShowdownStateContext: &txpokergrpc.DeclareShowdownStateContext{}},
	}
}
