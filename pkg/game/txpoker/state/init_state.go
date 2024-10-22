package state

import (
	"context"
	"encoding/json"
	commonmodel "card-game-server-prototype/pkg/common/model"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/common/type/gamemode"
	"card-game-server-prototype/pkg/common/type/gametype"
	"card-game-server-prototype/pkg/common/type/roomstate"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/constant"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/service"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"
	"strconv"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoInitState = &core.StateTrigger{
	Name:      "GoInitState",
	ArgsTypes: []reflect.Type{},
}

type InitState struct {
	core.State

	cfg               *config.Config
	roomInfo          *commonmodel.RoomInfo
	userGroup         *commonmodel.UserGroup
	userCacheGroup    *model2.UserCacheGroup
	gameInfo          *model2.GameInfo
	gameSetting       *model2.GameSetting
	seatStatusGroup   *model2.SeatStatusGroup
	eventGroup        *model2.EventGroup
	playSettingGroup  *model2.PlaySettingGroup
	statsGroup        *model2.StatsGroup
	forceBuyInGroup   *model2.ForceBuyInGroup
	tableProfitsGroup *model2.TableProfitsGroup

	userService       commonservice.UserService
	roomService       commonservice.RoomService
	seatStatusService service.SeatStatusService
	gameRepoService   service.GameRepoService
}

func ProvideInitState(
	stateFactory *core.StateFactory,
	cfg *config.Config,
	roomInfo *commonmodel.RoomInfo,
	userGroup *commonmodel.UserGroup,
	userCacheGroup *model2.UserCacheGroup,
	gameInfo *model2.GameInfo,
	gameSetting *model2.GameSetting,
	seatStatusGroup *model2.SeatStatusGroup,
	eventGroup *model2.EventGroup,
	playSettingGroup *model2.PlaySettingGroup,
	statsGroup *model2.StatsGroup,
	forceBuyInGroup *model2.ForceBuyInGroup,
	tableProfitsGroup *model2.TableProfitsGroup,

	userService commonservice.UserService,
	roomService commonservice.RoomService,
	seatStatusService service.SeatStatusService,
	gameRepoService service.GameRepoService,
) *InitState {
	return &InitState{
		State: stateFactory.Create("InitState"),

		cfg:               cfg,
		roomInfo:          roomInfo,
		userGroup:         userGroup,
		userCacheGroup:    userCacheGroup,
		gameInfo:          gameInfo,
		gameSetting:       gameSetting,
		seatStatusGroup:   seatStatusGroup,
		eventGroup:        eventGroup,
		playSettingGroup:  playSettingGroup,
		statsGroup:        statsGroup,
		forceBuyInGroup:   forceBuyInGroup,
		tableProfitsGroup: tableProfitsGroup,

		userService:       userService,
		roomService:       roomService,
		seatStatusService: seatStatusService,
		gameRepoService:   gameRepoService,
	}
}

func (state *InitState) Run(ctx context.Context, args ...any) error {
	state.roomInfo.RoomId = *state.cfg.RoomId
	state.roomInfo.ShortRoomId = *state.cfg.ShortRoomId
	state.roomInfo.GameType = *state.cfg.GameType
	state.roomInfo.GameMode = *state.cfg.GameMode
	state.roomInfo.GameMetaUid = *state.cfg.GameMetaUid

	rawValidUsers := []string{}
	if err := json.Unmarshal([]byte(*state.cfg.ValidUsers), &rawValidUsers); err != nil {
		state.Logger().Error("failed to unmarshal valid users from cfg", zap.Error(err), zap.Object("cfg", state.cfg))
		state.GameController().GoErrorState()
		return nil
	}
	state.roomInfo.ValidUsers = lo.Map(rawValidUsers, func(uid string, _ int) core.Uid { return core.Uid(uid) })
	state.Logger().Info("initialized roomInfo", zap.Object("roomInfo", state.roomInfo), zap.Object("cfg", state.cfg))

	if err := state.gameRepoService.FetchGameInfo(); err != nil {
		state.Logger().Error("failed to fetch game info", zap.Error(err))
		state.GameController().GoErrorState()
		return nil
	}
	state.Logger().Debug("fetched gameInfo", zap.Object("gameInfo", state.gameInfo))

	state.seatStatusGroup.TableSize = state.gameSetting.TableSize

	if err := state.userService.Init(state.roomInfo.ValidUsers...); err != nil {
		state.Logger().Error("failed to init user data", zap.Error(err), zap.Array("uids", state.roomInfo.ValidUsers))
		state.GameController().GoErrorState()
		return nil
	}

	if err := state.userService.FetchFromRepo(state.roomInfo.ValidUsers...); err != nil {
		state.Logger().Error("failed to refresh user data", zap.Error(err), zap.Array("uids", state.roomInfo.ValidUsers))
		state.GameController().GoErrorState()
		return nil
	}

	state.Logger().Info(
		"user group initialized",
		zap.Object("users", state.userGroup),
		zap.Object("usersCache", state.userCacheGroup),
	)

	state.GameController().RunTicker(constant.ClearCacheInterval, func() {
		state.GameController().RunTask(func() {
			state.forceBuyInGroup.CleanExpired()
			state.tableProfitsGroup.EvictExpired()
		})
	})

	// CloseAt: CLUB MODE, BUDDY MODE, COMMON MODE
	// 因為 runTicker 有消耗，還是針對要用的 mode 才開放使用。
	// txpoker 房間一般來說不會關閉，但有些模式會由後臺控制關閉。
	// 這裡定期監聽後臺關閉訊號。
	// 收到 close 時設定關閉時間點，之後在 WaitUserState 定期檢查是否到達關閉時間點進行關閉流程。
	if state.roomInfo.GameMode == gamemode.Club ||
		state.roomInfo.GameMode == gamemode.Buddy ||
		state.roomInfo.GameMode == gamemode.Common {
		state.GameController().RunTicker(time.Minute, func() {

			// Already set closeAt, no need to check again.
			if !state.gameSetting.CloseAt.IsZero() {
				return
			}

			now := time.Now()

			d, err := state.roomService.GetDetail()
			if err != nil {
				state.Logger().Error("room get_detail failed", zap.Error(err))
				state.GameController().RunTask(state.GameController().GoErrorState)
				return
			}

			if d.State == int(roomstate.Closing) && strconv.Itoa(d.Game) == string(gametype.TXPoker) {
				switch state.roomInfo.GameMode {
				case gamemode.Club, gamemode.Buddy:
					state.gameSetting.CloseAt = now
				case gamemode.Common:
					if state.gameInfo.CreationId == 0 {
						state.Logger().Warn("creationId is ZERO!")
					}

					// 使用 CreationId 計算 Close 的時段，避免同時太多 room 被關閉。
					// 公式：每 CloseRoomPeriod 關閉 CloseRoomAmount 個房間。
					// 這邊用 CreationId 乘以 10 是為了避免 CreationId 過小，造成 CloseAt 沒有間隔。
					bufferCreationId := lo.Ternary(state.gameInfo.CreationId < 10,
						state.gameInfo.CreationId*10,
						state.gameInfo.CreationId)
					state.gameSetting.CloseAt = now.Add(time.Minute * time.Duration(
						bufferCreationId/constant.CloseRoomAmount*constant.CloseRoomPeriod,
					))
				}

				state.Logger().Warn("set closeAt",
					zap.String("mode", state.cfg.GameMode.String()),
					zap.String("roomId", state.roomInfo.RoomId),
					zap.Int("CreationId", state.gameInfo.CreationId),
					zap.Time("CloseAt", state.gameSetting.CloseAt),
				)
			}

			// Handle game water
			if state.roomInfo.GameMode == gamemode.Common || state.roomInfo.GameMode == gamemode.Buddy {
				if err := state.gameRepoService.FetchGameWater(); err != nil {
					state.Logger().Error("failed to fetch game water", zap.Error(err))
				}
			}
		})
	}

	go func() {
		if err := state.roomService.RunPingLoop(); err != nil {
			state.Logger().Error("room ping loop failed", zap.Error(err))
			state.GameController().RunTask(state.GameController().GoErrorState)
		}
		state.Logger().Info("room ping loop ended normally")
	}()

	state.seatStatusService.StartRefillSitOutDurationLoop()

	state.GameController().GoNextState(GoResetState)
	return nil
}

func (state *InitState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
	})
	return nil
}

func (state *InitState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context:   &txpokergrpc.GameState_InitStateContext{InitStateContext: &txpokergrpc.InitStateContext{}},
	}
}
