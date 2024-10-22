package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoCallState = &core.StateTrigger{
	Name:      "GoCallState",
	ArgsTypes: []reflect.Type{core.UidType},
}

type CallState struct {
	core.State

	actionUid core.Uid
}

func ProvideCallState(
	stateFactory *core.StateFactory,
) *CallState {
	return &CallState{
		State: stateFactory.Create("CallState"),
	}
}

func (state *CallState) Run(ctx context.Context, args ...any) error {
	state.actionUid = args[0].(core.Uid)
	state.GameController().GoNextState(GoEvaluateActionState)
	return nil
}

func (state *CallState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
	})
	return nil
}

func (state *CallState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &txpokergrpc.GameState_CallStateContext{CallStateContext: &txpokergrpc.CallStateContext{
			ActorUid: state.actionUid.String(),
		}},
	}
}
