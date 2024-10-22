package state

import (
	"context"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/common/type/gamemode"
	"card-game-server-prototype/pkg/core"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoEndGameState = core.NewStateTrigger("GoEndGameState")

type EndGameState struct {
	core.State

	roomInfo *commonmodel.RoomInfo
}

func ProvideEndGameState(
	stateFactory *core.StateFactory,

	roomInfo *commonmodel.RoomInfo,
) *EndGameState {
	return &EndGameState{
		State: stateFactory.Create("EndGameState"),

		roomInfo: roomInfo,
	}
}

func (state *EndGameState) Run(context.Context, ...any) error {
	state.GameController().RunTimer(1*time.Second, func() {
		if state.roomInfo.GameMode == gamemode.Buddy {
			state.GameController().GoNextState(GoResetGameState)
		} else {
			state.GameController().GoNextState(GoClosedState)
		}
	})
	return nil
}

func (state *EndGameState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	return nil
}

func (state *EndGameState) ToProto(core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_EndGameStateContext{
			EndGameStateContext: &gamegrpc.EndGameStateContext{},
		},
	}
}
