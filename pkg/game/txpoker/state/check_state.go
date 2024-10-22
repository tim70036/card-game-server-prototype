package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoCheckState = &core.StateTrigger{
	Name:      "GoCheckState",
	ArgsTypes: []reflect.Type{core.UidType},
}

type CheckState struct {
	core.State

	actionUid core.Uid
}

func ProvideCheckState(
	stateFactory *core.StateFactory,
) *CheckState {
	return &CheckState{
		State: stateFactory.Create("CheckState"),
	}
}

func (state *CheckState) Run(ctx context.Context, args ...any) error {
	state.actionUid = args[0].(core.Uid)
	state.GameController().GoNextState(GoEvaluateActionState)
	return nil
}

func (state *CheckState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
	})
	return nil
}

func (state *CheckState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &txpokergrpc.GameState_CheckStateContext{CheckStateContext: &txpokergrpc.CheckStateContext{
			ActorUid: state.actionUid.String(),
		}},
	}
}
