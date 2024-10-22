package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoFoldState = &core.StateTrigger{
	Name:      "GoFoldState",
	ArgsTypes: []reflect.Type{core.UidType},
}

type FoldState struct {
	core.State

	actionUid core.Uid
}

func ProvideFoldState(
	stateFactory *core.StateFactory,
) *FoldState {
	return &FoldState{
		State: stateFactory.Create("FoldState"),
	}
}

func (state *FoldState) Run(ctx context.Context, args ...any) error {
	state.actionUid = args[0].(core.Uid)
	state.GameController().GoNextState(GoEvaluateActionState)
	return nil
}

func (state *FoldState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
	})
	return nil
}

func (state *FoldState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &txpokergrpc.GameState_FoldStateContext{FoldStateContext: &txpokergrpc.FoldStateContext{
			ActorUid: state.actionUid.String(),
		}},
	}
}
