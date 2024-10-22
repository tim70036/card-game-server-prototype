package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoRaiseState = &core.StateTrigger{
	Name:      "GoRaiseState",
	ArgsTypes: []reflect.Type{core.UidType},
}

type RaiseState struct {
	core.State

	actionUid core.Uid
}

func ProvideRaiseState(
	stateFactory *core.StateFactory,
) *RaiseState {
	return &RaiseState{
		State: stateFactory.Create("RaiseState"),
	}
}

func (state *RaiseState) Run(ctx context.Context, args ...any) error {
	state.actionUid = args[0].(core.Uid)
	state.GameController().GoNextState(GoEvaluateActionState)
	return nil
}

func (state *RaiseState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
	})
	return nil
}

func (state *RaiseState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &txpokergrpc.GameState_RaiseStateContext{RaiseStateContext: &txpokergrpc.RaiseStateContext{
			ActorUid: state.actionUid.String(),
		}},
	}
}
