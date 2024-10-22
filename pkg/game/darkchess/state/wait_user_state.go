package state

import (
	"context"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoWaitUserState = core.NewStateTrigger("GoWaitUserState")

type WaitUserState struct {
	core.State

	userGroup *commonmodel.UserGroup

	cancelTimer context.CancelFunc
	endOnce     *sync.Once
}

func ProvideWaitUserState(
	stateFactory *core.StateFactory,

	userGroup *commonmodel.UserGroup,
) *WaitUserState {
	return &WaitUserState{
		State: stateFactory.Create("WaitUserState"),

		userGroup: userGroup,
	}
}

func (state *WaitUserState) Run(context.Context, ...any) error {
	state.cancelTimer = func() {}
	state.endOnce = &sync.Once{}
	state.tryStartGame()

	state.cancelTimer = state.GameController().RunTimer(constant.WaitUserEnterGracefulPeriod, func() {
		state.Logger().Debug("wait user enter graceful period expired, starting game")
		state.endWaitUser()
	})
	return nil
}

func (state *WaitUserState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	return nil
}

func (state *WaitUserState) ToProto(core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_WaitUserStateContext{
			WaitUserStateContext: &gamegrpc.WaitUserStateContext{},
		},
	}
}

func (state *WaitUserState) HandleEnter(core.Uid) error {
	state.tryStartGame()
	return nil
}

func (state *WaitUserState) BeforeLeave(uid core.Uid) error {
	state.Logger().Warn("forbid to leave during game", zap.String("uid", uid.String()))
	return status.Errorf(codes.FailedPrecondition, "forbid to leave during game")
}

func (state *WaitUserState) tryStartGame() {
	countEnteredUsers := len(lo.PickBy(state.userGroup.Data, func(uid core.Uid, user *commonmodel.User) bool {
		return user.HasEntered
	}))

	if countEnteredUsers >= constant.StartGameUserCount {
		state.Logger().Debug("enough users have entered, starting game",
			zap.Int("users", len(state.userGroup.Data)),
		)
		state.endWaitUser()
	}
}

func (state *WaitUserState) endWaitUser() {
	state.endOnce.Do(func() {
		state.Logger().Debug("end waiting user")
		state.cancelTimer()
		state.GameController().GoNextState(GoStartGameState)
	})
}
