package state

import (
	"context"
	commonmodel "card-game-server-prototype/pkg/common/model"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/common/type/gamemode"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/constant"
	model2 "card-game-server-prototype/pkg/game/zoomtxpoker/pool/model"
	service2 "card-game-server-prototype/pkg/game/zoomtxpoker/pool/service"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"time"
)

var GoInitState = &core.StateTrigger{
	Name:      "GoInitState",
	ArgsTypes: []reflect.Type{},
}

type InitState struct {
	core.State

	cfg                *config.Config
	testCFG            *config.TestConfig
	roomInfo           *commonmodel.RoomInfo
	userGroup          *commonmodel.UserGroup
	gameSetting        *model2.GameSetting
	playSettingGroup   *model2.PlaySettingGroup
	statsGroup         *model2.StatsGroup
	participantGroup   *model2.ParticipantGroup
	playedHistoryGroup *model2.PlayedHistoryGroup
	forceBuyInGroup    *model2.ForceBuyInGroup
	tableProfitsGroup  *model2.TableProfitsGroup

	gameRepoService    service2.GameRepoService
	userService        commonservice.UserService
	roomService        commonservice.RoomService
	participantService *service2.ParticipantService
	checkerService     *service2.LeaveService
}

func ProvideInitState(
	stateFactory *core.StateFactory,
	cfg *config.Config,
	testCFG *config.TestConfig,
	roomInfo *commonmodel.RoomInfo,
	userGroup *commonmodel.UserGroup,
	gameSetting *model2.GameSetting,
	playSettingGroup *model2.PlaySettingGroup,
	statsGroup *model2.StatsGroup,
	participantGroup *model2.ParticipantGroup,
	playedHistoryGroup *model2.PlayedHistoryGroup,
	forceBuyInGroup *model2.ForceBuyInGroup,
	tableProfitsGroup *model2.TableProfitsGroup,

	gameRepoService service2.GameRepoService,
	userService commonservice.UserService,
	roomService commonservice.RoomService,
	participantService *service2.ParticipantService,
	checkerService *service2.LeaveService,
) *InitState {
	return &InitState{
		State: stateFactory.Create("InitState"),

		cfg:                cfg,
		testCFG:            testCFG,
		roomInfo:           roomInfo,
		userGroup:          userGroup,
		gameSetting:        gameSetting,
		playSettingGroup:   playSettingGroup,
		statsGroup:         statsGroup,
		participantGroup:   participantGroup,
		playedHistoryGroup: playedHistoryGroup,
		forceBuyInGroup:    forceBuyInGroup,
		tableProfitsGroup:  tableProfitsGroup,

		gameRepoService:    gameRepoService,
		userService:        userService,
		roomService:        roomService,
		participantService: participantService,
		checkerService:     checkerService,
	}
}

func (state *InitState) Run(ctx context.Context, args ...any) error {
	state.roomInfo.RoomId = *state.cfg.RoomId
	state.roomInfo.ShortRoomId = *state.cfg.ShortRoomId
	state.roomInfo.GameType = *state.cfg.GameType
	state.roomInfo.GameMode = *state.cfg.GameMode
	state.roomInfo.GameMetaUid = *state.cfg.GameMetaUid

	state.Logger().Info("roomInfo updated", zap.Object("roomInfo", state.roomInfo))

	if err := state.gameRepoService.FetchGameInfo(); err != nil {
		state.Logger().Error("failed to fetch game info", zap.Error(err))
		state.GameController().GoErrorState()
		return nil
	}
	state.Logger().Info("fetched game", zap.Object("gameSetting", state.gameSetting))

	state.checkerService.RunIdleChecker()

	state.GameController().RunTicker(constant.ClearCacheInterval, func() {
		state.forceBuyInGroup.CleanExpired()
		state.tableProfitsGroup.EvictExpired()
	})

	// Handle game water
	if state.roomInfo.GameMode == gamemode.Common || state.roomInfo.GameMode == gamemode.Buddy {
		state.GameController().RunTicker(time.Minute, func() {
			if err := state.gameRepoService.FetchGameWater(); err != nil {
				state.Logger().Error("failed to fetch game water", zap.Error(err))
			}
		})
	}

	go func() {
		state.Logger().Info("starting room ping loop",
			zap.String("roomId", state.roomInfo.RoomId),
		)
		if err := state.roomService.RunPingLoop(); err != nil {
			state.Logger().Error("room ping loop failed", zap.Error(err))
			state.GameController().RunTask(state.GameController().GoErrorState)
		}
		state.Logger().Info("room ping loop ended normally",
			zap.String("roomId", state.roomInfo.RoomId),
		)
	}()

	state.GameController().GoNextState(GoMatchingState)
	return nil
}

func (state *InitState) Publish(ctx context.Context, args ...any) error {
	return nil
}

func (state *InitState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context:   &txpokergrpc.GameState_InitStateContext{InitStateContext: &txpokergrpc.InitStateContext{}},
	}
}
