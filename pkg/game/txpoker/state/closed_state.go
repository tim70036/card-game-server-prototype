package state

import (
	"context"
	commonconstant "card-game-server-prototype/pkg/common/constant"
	commonmodel "card-game-server-prototype/pkg/common/model"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/service"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoClosedState = &core.StateTrigger{
	Name:      "GoClosedState",
	ArgsTypes: []reflect.Type{},
}

type ClosedState struct {
	core.State

	userGroup         *commonmodel.UserGroup
	seatStatusGroup   *model2.SeatStatusGroup
	userService       commonservice.UserService
	seatStatusService service.SeatStatusService
	roomService       commonservice.RoomService
	gameSetting       *model2.GameSetting
}

func ProvideClosedState(
	stateFactory *core.StateFactory,
	userGroup *commonmodel.UserGroup,
	seatStatusGroup *model2.SeatStatusGroup,
	userService commonservice.UserService,
	seatStatusService service.SeatStatusService,
	roomService commonservice.RoomService,
	gameSetting *model2.GameSetting,
) *ClosedState {
	return &ClosedState{
		State: stateFactory.Create("ClosedState"),

		userGroup:         userGroup,
		seatStatusGroup:   seatStatusGroup,
		userService:       userService,
		seatStatusService: seatStatusService,
		roomService:       roomService,
		gameSetting:       gameSetting,
	}
}

func (state *ClosedState) Run(context.Context, ...any) error {
	kickoutMsg := &txpokergrpc.Message{
		Event: &txpokergrpc.Event{
			Kickout: &commongrpc.Kickout{
				Reason: commongrpc.KickoutReason_GAME_CLOSED,
			},
		},
	}

	if now := time.Now(); !state.gameSetting.CloseAt.IsZero() &&
		(now.Equal(state.gameSetting.CloseAt) || now.After(state.gameSetting.CloseAt)) {
		kickoutMsg.Event.Kickout.Reason = commongrpc.KickoutReason_MAINTENANCE
	}

	// Need to cash out users, otherwise they will lose money.
	state.Logger().Debug("cashing out all users")
	for uid := range state.seatStatusGroup.Status {
		if _, err := state.seatStatusService.StandUp(uid); err != nil {
			state.Logger().Error("failed to stand up user", zap.Error(err), zap.String("uid", uid.String()), zap.Object("seatStatus", state.seatStatusGroup))
		}
	}

	state.Logger().Debug("destroying all users")
	if err := state.userService.Destroy(core.MessageTopic, kickoutMsg, lo.Keys(state.userGroup.Data)...); err != nil {
		state.Logger().Error("failed to destroy users", zap.Error(err))
	}

	state.Logger().Debug("closing room")
	if err := state.roomService.Close(); err != nil {
		state.Logger().Error("failed to close room", zap.Error(err))
	}

	state.GameController().RunTimer(commonconstant.ShutdownGracefulPeriod, func() {
		state.Logger().Debug("graceful period end, shutting down game", zap.Duration("ShutdownGracefulPeriod", commonconstant.ShutdownGracefulPeriod))
		state.GameController().Shutdown()
	})
	return nil
}

func (state *ClosedState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
	})
	return nil
}

func (state *ClosedState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context:   &txpokergrpc.GameState_ClosedStateContext{ClosedStateContext: &txpokergrpc.ClosedStateContext{}},
	}
}

func (state *ClosedState) BeforeConnect(uid core.Uid) error {
	state.Logger().Warn("forbid to connect due to game closed", zap.String("uid", uid.String()))
	return status.Errorf(codes.Aborted, "forbid to connect due to game closed")
}
