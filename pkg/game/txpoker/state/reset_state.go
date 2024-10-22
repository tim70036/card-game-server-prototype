package state

import (
	"context"
	"fmt"
	commonmodel "card-game-server-prototype/pkg/common/model"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/actor"
	"card-game-server-prototype/pkg/game/txpoker/constant"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/service"
	"card-game-server-prototype/pkg/game/txpoker/type/action"
	"card-game-server-prototype/pkg/game/txpoker/type/card"
	"card-game-server-prototype/pkg/game/txpoker/type/pot"
	"card-game-server-prototype/pkg/game/txpoker/type/seatstatus"
	"card-game-server-prototype/pkg/game/txpoker/type/stage"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoResetState = &core.StateTrigger{
	Name:      "GoResetState",
	ArgsTypes: []reflect.Type{},
}

type ResetState struct {
	core.State
	testCFG          *config.TestConfig
	roomInfo         *commonmodel.RoomInfo
	userGroup        *commonmodel.UserGroup
	gameSetting      *model2.GameSetting
	table            *model2.Table
	playerGroup      *model2.PlayerGroup
	actionHintGroup  *model2.ActionHintGroup
	actorGroup       *actor.ActorGroup
	seatStatusGroup  *model2.SeatStatusGroup
	playSettingGroup *model2.PlaySettingGroup
	replay           *model2.Replay
	statsCacheGroup  *model2.StatsCacheGroup
	chipCacheGroup   *model2.ChipCacheGroup
	userCacheGroup   *model2.UserCacheGroup
	tableProfitGroup *model2.TableProfitsGroup
	forceBuyInGroup  *model2.ForceBuyInGroup

	userService       commonservice.UserService
	seatStatusService service.SeatStatusService
}

func ProvideResetState(
	stateFactory *core.StateFactory,
	testCFG *config.TestConfig,
	roomInfo *commonmodel.RoomInfo,
	userGroup *commonmodel.UserGroup,
	gameSetting *model2.GameSetting,
	table *model2.Table,
	playerGroup *model2.PlayerGroup,
	actionHintGroup *model2.ActionHintGroup,
	actorGroup *actor.ActorGroup,
	seatStatusGroup *model2.SeatStatusGroup,
	playSettingGroup *model2.PlaySettingGroup,
	replay *model2.Replay,
	statsCacheGroup *model2.StatsCacheGroup,
	chipCacheGroup *model2.ChipCacheGroup,
	userCacheGroup *model2.UserCacheGroup,
	tableProfitGroup *model2.TableProfitsGroup,
	forceBuyInGroup *model2.ForceBuyInGroup,

	userService commonservice.UserService,
	seatStatusService service.SeatStatusService,
) *ResetState {
	return &ResetState{
		State:            stateFactory.Create("ResetState"),
		testCFG:          testCFG,
		roomInfo:         roomInfo,
		userGroup:        userGroup,
		gameSetting:      gameSetting,
		table:            table,
		playerGroup:      playerGroup,
		actionHintGroup:  actionHintGroup,
		actorGroup:       actorGroup,
		seatStatusGroup:  seatStatusGroup,
		playSettingGroup: playSettingGroup,
		replay:           replay,
		statsCacheGroup:  statsCacheGroup,
		chipCacheGroup:   chipCacheGroup,
		userCacheGroup:   userCacheGroup,
		tableProfitGroup: tableProfitGroup,
		forceBuyInGroup:  forceBuyInGroup,

		userService:       userService,
		seatStatusService: seatStatusService,
	}
}

func (state *ResetState) Run(ctx context.Context, args ...any) error {
	if err := state.userService.FetchFromRepo(lo.Keys(state.userGroup.Data)...); err != nil {
		state.Logger().Error("failed to fetch user from repo", zap.Error(err), zap.Object("users", state.userGroup))
		state.GameController().GoErrorState()
		return nil
	}

	state.resetModel()
	state.Logger().Debug(
		"reset model done",
		zap.Object("users", state.userGroup),
		zap.Object("seatStatus", state.seatStatusGroup),
	)

	topUpDone := state.autoTopUp()
	go func() {
		<-topUpDone
		state.GameController().RunTask(func() {
			if err := state.standUpIdleUser(); err != nil {
				state.Logger().Error("failed to stand up idle user", zap.Error(err))
				state.GameController().GoErrorState()
				return
			}

			if err := state.standUpPoorUser(); err != nil {
				state.Logger().Error("failed to stand up poor user", zap.Error(err))
				state.GameController().GoErrorState()
				return
			}

			state.GameController().GoNextState(GoWaitUserState)
		})
	}()

	return nil
}

func (state *ResetState) Publish(ctx context.Context, args ...any) error {
	playerGroupProto := &txpokergrpc.PlayerGroup{Players: make(map[string]*txpokergrpc.Player)}
	for uid, player := range state.playerGroup.Data {
		playerGroupProto.Players[uid.String()] = player.ToProto()
		_, playerGroupProto.Players[uid.String()].HasShowdown = state.table.ShowdownPocketCards[uid]
	}

	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
		Model: &txpokergrpc.Model{
			UserGroup:         state.userGroup.ToProto(),
			SeatStatusGroup:   state.seatStatusGroup.ToProto(),
			PlayerGroup:       playerGroupProto,
			ActionHintGroup:   state.actionHintGroup.ToProto(),
			StatsCacheGroup:   state.statsCacheGroup.ToProto(),
			TableProfitsGroup: state.tableProfitGroup.ToProto(lo.Keys(state.playerGroup.Data)),
			ChipCacheGroup:    state.chipCacheGroup.ToProto(),
			UserCacheGroup:    state.userCacheGroup.ToProto(),
			Table:             state.table.ToProto(),
		},
	})
	return nil
}

func (state *ResetState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context:   &txpokergrpc.GameState_ResetStateContext{ResetStateContext: &txpokergrpc.ResetStateContext{}},
	}
}

func (state *ResetState) autoTopUp() <-chan struct{} {
	topUpDones := make([]<-chan struct{}, 0)

	sittingStatus := lo.Filter(lo.Values(state.seatStatusGroup.Status), func(status *model2.SeatStatus, _ int) bool {
		return status.FSM.MustState().(seatstatus.SeatStatusState) == seatstatus.JoiningState ||
			status.FSM.MustState().(seatstatus.SeatStatusState) == seatstatus.SittingOutState ||
			status.FSM.MustState().(seatstatus.SeatStatusState) == seatstatus.WaitingState
	})

	for _, seatStatus := range sittingStatus {
		topUpChip := state.seatStatusGroup.TopUpQueue[seatStatus.Uid]
		playSetting := state.playSettingGroup.Data[seatStatus.Uid]
		maxEnterLimitChip := state.gameSetting.MaxEnterLimitBB * state.gameSetting.BigBlind

		if playSetting.AutoTopUp &&
			seatStatus.Chip+topUpChip < int(float64(maxEnterLimitChip)*playSetting.AutoTopUpThresholdPercent) {
			topUpChip = int(float64(maxEnterLimitChip)*playSetting.AutoTopUpChipPercent) - seatStatus.Chip
		}

		if topUpChip <= 0 {
			continue
		}

		topUpDone, err := state.seatStatusService.TopUp(seatStatus.Uid, topUpChip)
		if err != nil {
			state.Logger().Error(
				"failed to top up user", zap.Error(err),
				zap.Object("seatStatus", seatStatus),
				zap.Object("playSetting", playSetting),
				zap.Int("topUpChip", topUpChip),
				zap.Int("topUpChipInQueue", state.seatStatusGroup.TopUpQueue[seatStatus.Uid]),
			)

			state.MsgBus().Unicast(seatStatus.Uid, core.MessageTopic, &txpokergrpc.Message{
				Event: &txpokergrpc.Event{
					Warning: &txpokergrpc.Warning{Reason: txpokergrpc.WarningReason_AUTO_TOP_UP_FAILED},
				},
			})
			continue
		}

		topUpDones = append(topUpDones, topUpDone)
		state.Logger().Info(
			"topping up user",
			zap.Object("seatStatus", seatStatus),
			zap.Object("playSetting", playSetting),
			zap.Int("topUpChip", topUpChip),
			zap.Int("topUpChipInQueue", state.seatStatusGroup.TopUpQueue[seatStatus.Uid]),
		)
	}

	maps.Clear(state.seatStatusGroup.TopUpQueue)
	return lo.FanIn(len(topUpDones), topUpDones...)
}

// Stand up those who has not enough chip
func (state *ResetState) standUpPoorUser() error {
	joiningStatus := lo.Filter(lo.Values(state.seatStatusGroup.Status), func(status *model2.SeatStatus, _ int) bool {
		return status.FSM.MustState().(seatstatus.SeatStatusState) == seatstatus.JoiningState
	})

	for _, status := range joiningStatus {
		if status.Chip < state.gameSetting.BigBlind {
			_, err := state.seatStatusService.StandUp(status.Uid)
			if err != nil {
				return fmt.Errorf("failed to stand up user with not enough chip uid %v: %w", status.Uid, err)
			}

			// Make user show buy in, not force buy in, GCS-4299
			state.GameController().RunTask(func() {
				state.forceBuyInGroup.Delete(status.Uid)
			})

			state.Logger().Info(
				"standing up, not enough chip",
				zap.String("uid", status.Uid.String()),
				zap.Int("$", status.Chip),
			)

			state.MsgBus().Unicast(status.Uid, core.MessageTopic, &txpokergrpc.Message{
				Event: &txpokergrpc.Event{
					Warning: &txpokergrpc.Warning{Reason: txpokergrpc.WarningReason_INSUFFICIENT_BALANCE},
				},
			})
		}
	}

	return nil
}

func (state *ResetState) standUpIdleUser() error {
	for uid, seatStatus := range state.seatStatusGroup.Status {

		if seatStatus.CountIdleRounds >= constant.MaxIdleRounds {
			state.Logger().Info("standing up, idle too many times",
				zap.String("uid", uid.String()),
				zap.Int("countIdleRounds", seatStatus.CountIdleRounds),
			)

			if _, err := state.seatStatusService.StandUp(uid); err != nil {
				state.Logger().Error("idle stand up failed", zap.Error(err), zap.Object("seatStatus", seatStatus))
				return err
			}

			state.MsgBus().Unicast(seatStatus.Uid, core.MessageTopic, &txpokergrpc.Message{
				Event: &txpokergrpc.Event{
					Warning: &txpokergrpc.Warning{Reason: txpokergrpc.WarningReason_IDLE},
				},
			})
		}
	}

	return nil
}

func (state *ResetState) resetModel() {
	maps.Clear(state.playerGroup.Data)
	maps.Clear(state.actorGroup.Data)

	maps.Clear(state.actionHintGroup.Hints)
	state.actionHintGroup.RaiserHint = nil

	maps.Clear(state.statsCacheGroup.Data)
	maps.Clear(state.chipCacheGroup.SeatStatusChips)
	maps.Clear(state.chipCacheGroup.CashOutChips)
	maps.Clear(state.userCacheGroup.Data)

	for uid, u := range state.userGroup.Data {
		state.userCacheGroup.Data[uid] = &commonmodel.User{
			Uid:         u.Uid,
			ShortUid:    u.ShortUid,
			Name:        u.Name,
			IsAI:        u.IsAI,
			IsConnected: u.IsConnected,
			HasEntered:  u.HasEntered,
			Cash:        u.Cash,
			Level:       u.Level,
			RoomCards:   u.RoomCards,
		}
	}

	state.table.Deck = card.GenerateShuffledDeck(0) // 暫時只有一副 deck，不特別 log
	state.table.CommunityCards = card.CardList{}
	state.table.ShowdownPocketCards = map[core.Uid]card.CardList{}
	state.table.BetStageFSM = stage.NewBetStageFSM()
	state.table.Pots = pot.PotList{}

	maps.Clear(state.replay.PlayerRecords)
	state.replay.CommunityCards = card.CardList{}
	for s := stage.AnteStage; s <= stage.ShowdownStage; s++ {
		state.replay.ActionLog[s] = []action.ActionRecord{}
		state.replay.StagePotChip[s] = 0
	}
}
