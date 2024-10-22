package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoBetState = &core.StateTrigger{
	Name:      "GoBetState",
	ArgsTypes: []reflect.Type{core.UidType},
}

type BetState struct {
	core.State

	actionUid core.Uid
}

func ProvideBetState(
	stateFactory *core.StateFactory,
) *BetState {
	return &BetState{
		State: stateFactory.Create("BetState"),
	}
}

func (state *BetState) Run(ctx context.Context, args ...any) error {
	state.actionUid = args[0].(core.Uid)
	state.GameController().GoNextState(GoEvaluateActionState)
	return nil
}

func (state *BetState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
	})
	return nil
}

func (state *BetState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &txpokergrpc.GameState_BetStateContext{BetStateContext: &txpokergrpc.BetStateContext{
			ActorUid: state.actionUid.String(),
		}},
	}
}
