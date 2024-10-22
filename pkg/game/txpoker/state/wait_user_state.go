package state

import (
	"context"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/service"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoWaitUserState = &core.StateTrigger{
	Name:      "GoWaitUserState",
	ArgsTypes: []reflect.Type{},
}

type WaitUserState struct {
	core.State

	seatStatusService service.SeatStatusService
	gameInfo          *model2.GameInfo
	gameSetting       *model2.GameSetting
	gameRepoService   service.GameRepoService
	roomInfo          *commonmodel.RoomInfo

	cancelTryStartGameTicker context.CancelFunc
	cancelUpdateRoomTicker   context.CancelFunc
	endOnce                  *sync.Once
}

func ProvideWaitUserState(
	stateFactory *core.StateFactory,
	seatStatusService service.SeatStatusService,
	gameInfo *model2.GameInfo,
	gameSetting *model2.GameSetting,
	gameRepoService service.GameRepoService,
	roomInfo *commonmodel.RoomInfo,
) *WaitUserState {
	return &WaitUserState{
		State: stateFactory.Create("WaitUserState"),

		seatStatusService: seatStatusService,
		gameInfo:          gameInfo,
		gameSetting:       gameSetting,
		gameRepoService:   gameRepoService,
		roomInfo:          roomInfo,
	}
}

func (state *WaitUserState) Run(ctx context.Context, args ...any) error {
	state.cancelTryStartGameTicker = func() {}
	state.cancelUpdateRoomTicker = func() {}
	state.endOnce = &sync.Once{}

	if state.isExpired(time.Now()) {

		state.Logger().Warn("close room has triggered, do not start a game.",
			zap.String("room_id", state.roomInfo.RoomId),
		)
		state.endWaitUser(true)
		return nil
	}

	if state.seatStatusService.IsReadyToStartRound() {
		state.endWaitUser(false)
		return nil
	}

	state.cancelTryStartGameTicker = state.GameController().RunTicker(1*time.Second, func() {
		if state.isExpired(time.Now()) {

			if !state.seatStatusService.IsReadyToStartRound() && state == state.GameController().CurrentState() {
				state.Logger().Warn("no game, close room.",
					zap.String("room_id", state.roomInfo.RoomId),
				)
				state.endWaitUser(true)
				return
			} else {
				state.Logger().Warn("delay to close room.",
					zap.String("room_id", state.roomInfo.RoomId),
				)
			}
		}

		if state.seatStatusService.IsReadyToStartRound() &&
			state == state.GameController().CurrentState() { // In case that ticker triggered after game start.
			state.gameInfo.StreakRoundCount = 0
			state.endWaitUser(false)
		}

	})

	state.cancelUpdateRoomTicker = state.GameController().RunTicker(30*time.Second, func() {
		if err := state.gameRepoService.UpdateRoomInfo(); err != nil {
			state.Logger().Error("failed to update room info", zap.Error(err))
		}
	})

	return nil
}

func (state *WaitUserState) isExpired(now time.Time) bool {
	return !state.gameSetting.CloseAt.IsZero() &&
		(now.Equal(state.gameSetting.CloseAt) || now.After(state.gameSetting.CloseAt))
}

func (state *WaitUserState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
	})
	return nil
}

func (state *WaitUserState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context:   &txpokergrpc.GameState_WaitUserStateContext{WaitUserStateContext: &txpokergrpc.WaitUserStateContext{}},
	}
}

func (state *WaitUserState) endWaitUser(isGoCloseState bool) {
	state.endOnce.Do(func() {
		state.cancelTryStartGameTicker()
		state.cancelUpdateRoomTicker()
		if isGoCloseState {
			state.GameController().GoNextState(GoClosedState)
			return
		}
		state.GameController().GoNextState(GoStartRoundState)
	})
}
