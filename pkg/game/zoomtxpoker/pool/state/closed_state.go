package state

import (
	"context"
	commonconstant "card-game-server-prototype/pkg/common/constant"
	commonmodel "card-game-server-prototype/pkg/common/model"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var GoClosedState = &core.StateTrigger{
	Name:      "GoClosedState",
	ArgsTypes: []reflect.Type{},
}

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

func (state *ClosedState) Run(ctx context.Context, args ...any) error {
	kickoutMsg := &txpokergrpc.Message{
		Event: &txpokergrpc.Event{
			Kickout: &commongrpc.Kickout{
				Reason: commongrpc.KickoutReason_GAME_CLOSED,
			},
		},
	}

	if len(args) > 0 {
		if reason, ok := args[0].(commongrpc.KickoutReason); ok {
			kickoutMsg.Event.Kickout.Reason = reason
		}
	}

	// Need to cash out users, otherwise they will lose money.
	// state.Logger().Info("cashing out all users")
	// for uid := range state.seatStatusGroup.Status {
	//	if _, err := state.seatStatusService.StandUp(uid); err != nil {
	//		state.Logger().Error("failed to stand up user", zap.Error(err), zap.String("uid", uid.String()), zap.Object("seatStatus", state.seatStatusGroup))
	//	}
	// }

	state.Logger().Info("destroying all users")
	if err := state.userService.Destroy(core.ModelTopic, kickoutMsg, lo.Keys(state.userGroup.Data)...); err != nil {
		state.Logger().Error("failed to destroy users", zap.Error(err))
	}

	state.Logger().Info("closing room")
	if err := state.roomService.Close(); err != nil {
		state.Logger().Error("failed to close room", zap.Error(err))
	}

	state.GameController().RunTimer(commonconstant.ShutdownGracefulPeriod, func() {
		state.Logger().Info("graceful period end, shutting down game", zap.Duration("ShutdownGracefulPeriod", commonconstant.ShutdownGracefulPeriod))
		state.GameController().Shutdown()
	})
	return nil
}

func (state *ClosedState) Publish(ctx context.Context, args ...any) error {
	return nil
}

func (state *ClosedState) BeforeConnect(uid core.Uid) error {
	state.Logger().Warn("forbid to connect due to game closed", zap.String("uid", uid.String()))
	return status.Errorf(codes.Aborted, "forbid to connect due to game closed")
}
