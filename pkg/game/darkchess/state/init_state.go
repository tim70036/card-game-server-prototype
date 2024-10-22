package state

import (
	"context"
	commonmodel "card-game-server-prototype/pkg/common/model"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/common/type/gamemode"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/game/darkchess/service"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoInitState = core.NewStateTrigger("GoInitState")

type InitState struct {
	core.State

	cfg             *config.Config
	testCFG         *config.TestConfig
	roomInfo        *commonmodel.RoomInfo
	gameInfo        *model.GameInfo
	userGroup       *commonmodel.UserGroup
	userService     commonservice.UserService
	gameRepoService service.GameRepoService
	roomService     commonservice.RoomService
}

func ProvideInitState(
	stateFactory *core.StateFactory,

	cfg *config.Config,
	testCFG *config.TestConfig,
	roomInfo *commonmodel.RoomInfo,
	gameInfo *model.GameInfo,
	userGroup *commonmodel.UserGroup,
	userService commonservice.UserService,
	gameRepoService service.GameRepoService,
	roomService commonservice.RoomService,
) *InitState {
	return &InitState{
		State: stateFactory.Create("InitState"),

		cfg:             cfg,
		testCFG:         testCFG,
		roomInfo:        roomInfo,
		gameInfo:        gameInfo,
		userGroup:       userGroup,
		userService:     userService,
		gameRepoService: gameRepoService,
		roomService:     roomService,
	}
}

func (state *InitState) Run(context.Context, ...any) error {
	state.roomInfo.RoomId = *state.cfg.RoomId
	state.roomInfo.ShortRoomId = *state.cfg.ShortRoomId
	state.roomInfo.GameType = *state.cfg.GameType
	state.roomInfo.GameMode = *state.cfg.GameMode
	state.roomInfo.GameMetaUid = *state.cfg.GameMetaUid

	if state.roomInfo.GameMode == gamemode.Buddy {
		if err := state.roomService.FetchRoomInfo(); err != nil {
			state.Logger().Error("failed to fetch room info", zap.Error(err))
			state.GameController().GoErrorState()
			return nil
		}
	}

	var rawValidUsers []string
	if err := jsoniter.Unmarshal([]byte(*state.cfg.ValidUsers), &rawValidUsers); err != nil {
		state.Logger().Error("failed to unmarshal valid users from cfg", zap.Error(err), zap.Object("cfg", state.testCFG))
		state.GameController().GoErrorState()
		return nil
	}

	for _, uidStr := range rawValidUsers {
		state.roomInfo.ValidUsers = append(state.roomInfo.ValidUsers, core.Uid(uidStr))
	}

	state.Logger().Debug("initialized room info from config",
		zap.Object("roomInfo", state.roomInfo),
		zap.Object("cfg", state.cfg),
	)

	if err := state.gameRepoService.FetchGameSetting(); err != nil {
		state.Logger().Error("failed to update game setting", zap.Error(err))
		state.GameController().GoErrorState()
		return nil
	}

	state.Logger().Debug("updated game setting",
		zap.Object("gameSetting", state.gameInfo.Setting),
	)

	isInitUser := !(state.roomInfo.GameMode == gamemode.Buddy || state.roomInfo.GameMode == gamemode.Elimination)

	// If buddy mode, valid users will be empty. User is added when he enters the room.
	if isInitUser {
		if err := state.userService.Init(state.roomInfo.ValidUsers...); err != nil {
			state.Logger().Error("failed to init user data", zap.Error(err), zap.Array("uids", state.roomInfo.ValidUsers))
			state.GameController().GoErrorState()
			return nil
		}
	}

	if state.testCFG.EnableSingle(string(state.roomInfo.GameType)) {
		fakeUids := core.UidList{
			core.Uid("5e415691-f78a-4f6a-a35d-100ba7f8b3fb"),
			// core.Uid("aaaaa"),
			// core.Uid("bbbbb"),
		}

		if err := state.userService.Init(fakeUids...); err != nil {
			state.Logger().Error("failed to init fake user data", zap.Error(err), zap.Array("uids", state.roomInfo.ValidUsers))
			state.GameController().GoErrorState()
			return nil
		}

		for _, uid := range fakeUids {
			state.userGroup.Data[uid].IsAI = true
		}
	}

	if isInitUser {
		if err := state.userService.FetchFromRepo(state.roomInfo.ValidUsers...); err != nil {
			state.Logger().Error("failed to refresh user data", zap.Error(err), zap.Array("uids", state.roomInfo.ValidUsers))
			state.GameController().GoErrorState()
			return nil
		}
	}

	// 非人 user 手動設定 connected and entered。
	for _, user := range state.userGroup.Data {
		if user.IsAI {
			user.IsConnected = true
			user.HasEntered = true
		}
	}

	state.Logger().Debug("user group initialized",
		zap.Object("users", state.userGroup),
	)

	go func() {
		if err := state.roomService.RunPingLoop(); err != nil {
			state.Logger().Error("room ping loop failed", zap.Error(err))
			state.GameController().RunTask(state.GameController().GoErrorState)
		}
		state.Logger().Debug("room ping loop ended normally")
	}()

	state.GameController().GoNextState(GoResetGameState)
	return nil
}

func (state *InitState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	return nil
}

func (state *InitState) ToProto(core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_InitStateContext{
			InitStateContext: &gamegrpc.InitStateContext{},
		},
	}
}
