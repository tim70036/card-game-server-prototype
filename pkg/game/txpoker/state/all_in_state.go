package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoAllInState = &core.StateTrigger{
	Name:      "GoAllInState",
	ArgsTypes: []reflect.Type{core.UidType},
}

type AllInState struct {
	core.State

	actionUid core.Uid
}

func ProvideAllInState(
	stateFactory *core.StateFactory,
) *AllInState {
	return &AllInState{
		State: stateFactory.Create("AllInState"),
	}
}

func (state *AllInState) Run(ctx context.Context, args ...any) error {
	state.actionUid = args[0].(core.Uid)
	state.GameController().GoNextState(GoEvaluateActionState)
	return nil
}

func (state *AllInState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
	})
	return nil
}

func (state *AllInState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &txpokergrpc.GameState_AllInStateContext{AllInStateContext: &txpokergrpc.AllInStateContext{
			ActorUid: state.actionUid.String(),
		}},
	}
}
