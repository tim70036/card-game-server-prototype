package state

import (
	"context"
	commonconstant "card-game-server-prototype/pkg/common/constant"
	commonmodel "card-game-server-prototype/pkg/common/model"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoClosedState = core.NewStateTrigger("GoClosedState")

type ClosedState struct {
	core.State

	userGroup   *commonmodel.UserGroup
	userService commonservice.UserService
	roomService commonservice.RoomService
}

func ProvideClosedState(
	stateFactory *core.StateFactory,

	userGroup *commonmodel.UserGroup,
	userService commonservice.UserService,
	roomService commonservice.RoomService,
) *ClosedState {
	return &ClosedState{
		State: stateFactory.Create("ClosedState"),

		userGroup:   userGroup,
		userService: userService,
		roomService: roomService,
	}
}

func (state *ClosedState) Run(_ context.Context, args ...any) error {
	kickoutMsg := &gamegrpc.Event{
		Kickout: &commongrpc.Kickout{
			Reason: commongrpc.KickoutReason_GAME_CLOSED,
		},
	}

	if len(args) > 0 {
		if reason, ok := args[0].(commongrpc.KickoutReason); ok {
			kickoutMsg.Kickout.Reason = reason
		}
	}

	state.Logger().Debug("destroying all users")
	err := state.userService.Destroy(core.EventTopic, kickoutMsg, lo.Keys(state.userGroup.Data)...)
	if err != nil {
		state.Logger().Error("failed to destroy users", zap.Error(err))
	}

	state.GameController().RunTimer(commonconstant.ShutdownGracefulPeriod, func() {
		state.Logger().Debug("graceful period end, shutting down game")
		state.GameController().Shutdown()
	})

	return nil
}

func (state *ClosedState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	return nil
}

func (state *ClosedState) ToProto(core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_CloseStateContext{
			CloseStateContext: &gamegrpc.CloseStateContext{},
		},
	}
}

func (state *ClosedState) BeforeConnect(uid core.Uid) error {
	state.Logger().Warn("forbid to connect due to game closed", zap.String("uid", uid.String()))
	return status.Errorf(codes.Aborted, "forbid to connect due to game closed")
}
