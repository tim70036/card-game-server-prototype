package state

import (
	"context"
	"fmt"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	actor2 "card-game-server-prototype/pkg/game/txpoker/actor"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	service2 "card-game-server-prototype/pkg/game/txpoker/service"
	txpokerstate "card-game-server-prototype/pkg/game/txpoker/state"
	action2 "card-game-server-prototype/pkg/game/txpoker/type/action"
	"card-game-server-prototype/pkg/game/txpoker/type/card"
	event2 "card-game-server-prototype/pkg/game/txpoker/type/event"
	"card-game-server-prototype/pkg/game/txpoker/type/pot"
	role2 "card-game-server-prototype/pkg/game/txpoker/type/role"
	"card-game-server-prototype/pkg/game/txpoker/type/stage"
	session2 "card-game-server-prototype/pkg/game/zoomtxpoker/game/session"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InitState struct {
	core.State

	userSessionGroup         *session2.UserSessionGroup
	gameInfoSession          *session2.GameInfoSession
	statsSessionGroup        *session2.StatsSessionGroup
	forceBuyInSessionGroup   *session2.ForceBuyInSessionGroup
	tableProfitsSessionGroup *session2.TableProfitsSessionGroup

	cfg              *config.Config
	roomInfo         *commonmodel.RoomInfo
	userGroup        *commonmodel.UserGroup
	seatStatusGroup  *model2.SeatStatusGroup
	chipCacheGroup   *model2.ChipCacheGroup
	playSettingGroup *model2.PlaySettingGroup
	eventGroup       *model2.EventGroup
	statsGroup       *model2.StatsGroup
	gameSetting      *model2.GameSetting
	gameInfo         *model2.GameInfo
	forceBuyInGroup  *model2.ForceBuyInGroup
	tableProfitGroup *model2.TableProfitsGroup

	playerGroup     *model2.PlayerGroup
	actionHintGroup *model2.ActionHintGroup
	replay          *model2.Replay
	table           *model2.Table

	actorGroup        *actor2.ActorGroup
	baseActorFactory  *actor2.BaseActorFactory
	dummyActorFactory *actor2.DummyActorFactory

	seatStatusService service2.SeatStatusService
	resyncService     *service2.ResyncService
}

func ProvideInitState(
	stateFactory *core.StateFactory,
	userSessionGroup *session2.UserSessionGroup,
	gameInfoSession *session2.GameInfoSession,
	statsSessionGroup *session2.StatsSessionGroup,
	forceBuyInSessionGroup *session2.ForceBuyInSessionGroup,
	tableProfitsSessionGroup *session2.TableProfitsSessionGroup,

	cfg *config.Config,
	roomInfo *commonmodel.RoomInfo,
	userGroup *commonmodel.UserGroup,
	seatStatusGroup *model2.SeatStatusGroup,
	chipCacheGroup *model2.ChipCacheGroup,
	playSettingGroup *model2.PlaySettingGroup,
	eventGroup *model2.EventGroup,
	statsGroup *model2.StatsGroup,
	gameSetting *model2.GameSetting,
	gameInfo *model2.GameInfo,
	forceBuyInGroup *model2.ForceBuyInGroup,
	tableProfitGroup *model2.TableProfitsGroup,

	playerGroup *model2.PlayerGroup,
	actionHintGroup *model2.ActionHintGroup,
	replay *model2.Replay,
	table *model2.Table,

	actorGroup *actor2.ActorGroup,
	baseActorFactory *actor2.BaseActorFactory,
	dummyActorFactory *actor2.DummyActorFactory,

	seatStatusService service2.SeatStatusService,
	resyncService *service2.ResyncService,
) *InitState {
	return &InitState{
		State:                    stateFactory.Create("InitState"),
		userSessionGroup:         userSessionGroup,
		gameInfoSession:          gameInfoSession,
		statsSessionGroup:        statsSessionGroup,
		forceBuyInSessionGroup:   forceBuyInSessionGroup,
		tableProfitsSessionGroup: tableProfitsSessionGroup,

		cfg:              cfg,
		roomInfo:         roomInfo,
		userGroup:        userGroup,
		seatStatusGroup:  seatStatusGroup,
		chipCacheGroup:   chipCacheGroup,
		playSettingGroup: playSettingGroup,
		eventGroup:       eventGroup,
		statsGroup:       statsGroup,
		gameSetting:      gameSetting,
		gameInfo:         gameInfo,
		forceBuyInGroup:  forceBuyInGroup,
		tableProfitGroup: tableProfitGroup,

		playerGroup:     playerGroup,
		actionHintGroup: actionHintGroup,
		replay:          replay,
		table:           table,

		actorGroup:        actorGroup,
		baseActorFactory:  baseActorFactory,
		dummyActorFactory: dummyActorFactory,

		seatStatusService: seatStatusService,
		resyncService:     resyncService,
	}
}

func (state *InitState) Run(ctx context.Context, args ...any) error {
	if err := state.initModel(); err != nil {
		state.Logger().Error("failed to init model", zap.Error(err))
		state.GameController().GoErrorState()
		return nil
	}

	for uid := range state.userSessionGroup.Data {
		state.resyncService.Send(uid)
	}

	state.GameController().GoNextState(txpokerstate.GoResetState)
	return nil
}

func (state *InitState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic, &txpokergrpc.Message{
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

func (state *InitState) initModel() error {
	state.roomInfo.RoomId = *state.cfg.RoomId
	state.roomInfo.ShortRoomId = *state.cfg.ShortRoomId
	state.roomInfo.GameType = *state.cfg.GameType
	state.roomInfo.GameMode = *state.cfg.GameMode
	state.roomInfo.GameMetaUid = *state.cfg.GameMetaUid

	state.Logger().Info("game roomInfo updated",
		zap.Object("cfg", state.cfg),
		zap.Object("roomInfo", state.roomInfo),
	)

	state.gameSetting.GameMetaUid = state.gameInfoSession.Setting.GameMetaUid
	state.gameSetting.SmallBlind = state.gameInfoSession.Setting.SmallBlind
	state.gameSetting.BigBlind = state.gameInfoSession.Setting.BigBlind
	state.gameSetting.TurnDuration = state.gameInfoSession.Setting.TurnDuration
	state.gameSetting.InitialExtraTurnDuration = state.gameInfoSession.Setting.InitialExtraTurnDuration
	state.gameSetting.ExtraTurnRefillIntervalRound = state.gameInfoSession.Setting.ExtraTurnRefillIntervalRound
	state.gameSetting.RefillExtraTurnDuration = state.gameInfoSession.Setting.RefillExtraTurnDuration
	state.gameSetting.MaxExtraTurnDuration = state.gameInfoSession.Setting.MaxExtraTurnDuration
	state.gameSetting.InitialSitOutDuration = state.gameInfoSession.Setting.InitialSitOutDuration
	state.gameSetting.SitOutRefillIntervalDuration = state.gameInfoSession.Setting.SitOutRefillIntervalDuration
	state.gameSetting.RefillSitOutDuration = state.gameInfoSession.Setting.RefillSitOutDuration
	state.gameSetting.MaxSitOutDuration = state.gameInfoSession.Setting.MaxSitOutDuration
	state.gameSetting.MinEnterLimitBB = state.gameInfoSession.Setting.MinEnterLimitBB
	state.gameSetting.MaxEnterLimitBB = state.gameInfoSession.Setting.MaxEnterLimitBB
	state.gameSetting.WaterPct = state.gameInfoSession.Setting.WaterPct
	state.gameSetting.TableSize = state.gameInfoSession.Setting.TableSize
	state.gameSetting.LeastGamePlayerAmount = state.gameInfoSession.Setting.LeastGamePlayerAmount
	state.gameSetting.MaxWaterLimitBB = state.gameInfoSession.Setting.MaxWaterLimitBB

	state.gameInfo.RoundId = state.gameInfoSession.RoundId
	state.gameInfo.TotalBetPlayerCount = state.gameInfoSession.TotalBetPlayerCount
	state.gameInfo.TotalPlayedPlayerCount = state.gameInfoSession.TotalPlayedPlayerCount

	state.seatStatusGroup.TableSize = state.gameSetting.TableSize

	for uid, userSession := range state.userSessionGroup.Data {
		state.userGroup.Data[uid] = &commonmodel.User{
			Uid:         userSession.User.Uid,
			ShortUid:    userSession.User.ShortUid,
			Name:        userSession.User.Name,
			IsAI:        userSession.User.IsAI,
			IsConnected: userSession.User.IsConnected,
			HasEntered:  userSession.User.HasEntered,
			Cash:        userSession.User.Cash,
			Level:       userSession.User.Level,
			RoomCards:   userSession.User.RoomCards,
		}

		state.chipCacheGroup.SeatStatusChips[uid] = userSession.Chip

		state.playSettingGroup.Data[uid] = &model2.PlaySetting{
			Uid:                       userSession.User.Uid,
			WaitBB:                    true,
			AutoTopUp:                 false,
			AutoTopUpThresholdPercent: 0.0,
			AutoTopUpChipPercent:      0.0,
		}

		state.eventGroup.Data[uid] = event2.EventList{}

		state.statsGroup.Data[uid] = &model2.Stats{
			Uid:            uid,
			EventAmountSum: map[event2.EventType]int{},
		}

		if st, ok := state.statsSessionGroup.Data[uid]; ok {
			state.statsGroup.Data[uid].HighestGameWinAmount = st.HighestGameWinAmount

			for e, count := range st.EventAmountSum {
				state.statsGroup.Data[uid].EventAmountSum[e] = count
			}
		}

		if v, ok := state.forceBuyInSessionGroup.Get(uid); ok {
			state.forceBuyInGroup.Set(uid, v.GetBuyInChip(), v.GetExpireTime())
		}

		tableProfit := &model2.TableProfits{
			Uid:             uid,
			Name:            userSession.User.Name,
			CountGames:      0,
			SumBuyInChips:   0,
			SumWinLoseChips: 0,
		}

		if tableProfitsSession, ok := state.tableProfitsSessionGroup.Data[uid]; ok {
			tableProfit.CountGames = tableProfitsSession.CountGames
			tableProfit.SumBuyInChips = tableProfitsSession.SumBuyInChips
			tableProfit.SumWinLoseChips = tableProfitsSession.SumWinLoseChips
		}

		state.tableProfitGroup.Save(tableProfit)

		state.playerGroup.Data[uid] = &model2.Player{
			Uid:         userSession.User.Uid,
			Role:        userSession.Role,
			PocketCards: card.CardList{},
			Hand:        nil,
		}

		state.actionHintGroup.Hints[uid] = &model2.ActionHint{
			Uid:              userSession.User.Uid,
			BetChip:          0,
			RaiseChip:        0,
			CallingChip:      0,
			MinRaisingChip:   0,
			Action:           action2.Undefined,
			AvailableActions: []action2.ActionType{},
			Duration:         state.gameSetting.TurnDuration,
		}

		if userSession.User.IsAI {
			state.actorGroup.Data[uid] = state.dummyActorFactory.Create(uid)
		} else {
			state.actorGroup.Data[uid] = state.baseActorFactory.Create(uid)
		}
	}

	for uid, userSession := range state.userSessionGroup.Data {
		state.seatStatusGroup.Status[uid] = state.seatStatusService.NewSeatStatus(uid)

		seatId, err := state.evalSeatId(userSession.Role, len(state.userSessionGroup.Data))
		if err != nil {
			return fmt.Errorf("failed to eval seat id for uid %v: %w", uid, err)
		}

		state.seatStatusGroup.Status[uid].CountIdleRounds = state.userSessionGroup.Data[uid].CountIdleRounds

		if err := state.seatStatusService.SitDown(uid, seatId); err != nil {
			return fmt.Errorf("failed to sit down for uid %v : %w", uid, err)
		}

		if err := state.seatStatusService.BuyIn(uid, userSession.Chip); err != nil {
			return fmt.Errorf("failed to buy in for uid %v : %w", uid, err)
		}

		state.playerGroup.Data[uid].SeatId = seatId
	}

	for uid, userSession := range state.userSessionGroup.Data {
		state.replay.PlayerRecords[uid] = &model2.PlayerRecord{
			Uid:                     userSession.User.Uid,
			Role:                    userSession.Role,
			SeatId:                  state.playerGroup.Data[uid].SeatId,
			PocketCards:             card.CardList{},
			InitSeatStatusChip:      state.seatStatusGroup.Status[uid].Chip,
			BeforeWinSeatStatusChip: 0,
			IsWinner:                false,
			HasShowdown:             false,
		}
	}

	state.replay.RoundId = state.gameInfo.RoundId
	state.replay.CommunityCards = card.CardList{}
	for s := stage.AnteStage; s <= stage.ShowdownStage; s++ {
		state.replay.ActionLog[s] = []action2.ActionRecord{}
		state.replay.StagePotChip[s] = 0
	}

	state.table.Deck = card.GenerateShuffledDeck(0)
	state.table.CommunityCards = card.CardList{}
	state.table.ShowdownPocketCards = map[core.Uid]card.CardList{}
	state.table.BetStageFSM = stage.NewBetStageFSM()
	state.table.Pots = pot.PotList{}

	state.actionHintGroup.RaiserHint = nil

	state.Logger().Debug(
		"init model",
		zap.Object("userSessionGroup", state.userSessionGroup),
		zap.Object("gameInfoSession", state.gameInfoSession),
		zap.Object("gameSetting", state.gameSetting),
		zap.Object("gameInfo", state.gameInfo),
		zap.Object("userGroup", state.userGroup),
		zap.Object("seatStatusGroup", state.seatStatusGroup),
		zap.Object("chipCacheGroup", state.chipCacheGroup),
		zap.Object("playSettingGroup", state.playSettingGroup),
		zap.Object("eventGroup", state.eventGroup),
		zap.Object("statsGroup", state.statsGroup),
		zap.Object("playerGroup", state.playerGroup),
		zap.Object("actionHintGroup", state.actionHintGroup),
		zap.Object("replay", state.replay),
		zap.Object("table", state.table),
	)

	return nil
}

func (state *InitState) evalSeatId(playerRole role2.Role, playerCount int) (int, error) {
	roles, err := role2.GetRoles(playerCount)
	if err != nil {
		return 0, fmt.Errorf("failed eval seat id %d: %w", playerCount, err)
	}

	_, seatId, ok := lo.FindIndexOf(roles, func(role role2.Role) bool {
		return playerRole == role
	})

	if !ok {
		return 0, fmt.Errorf("seat id not found in roles setting [%v] for playerRole role: %s", roles, playerRole)
	}

	return seatId, nil
}
