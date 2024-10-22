package state

import (
	"context"
	"fmt"
	commonapi "card-game-server-prototype/pkg/common/api"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/common/type/gametype"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/type/event"
	"card-game-server-prototype/pkg/game/zoomtxpoker/builder"
	session2 "card-game-server-prototype/pkg/game/zoomtxpoker/game/session"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/constant"
	model3 "card-game-server-prototype/pkg/game/zoomtxpoker/pool/model"
	service2 "card-game-server-prototype/pkg/game/zoomtxpoker/pool/service"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/type/participant"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"reflect"
	"time"
)

var GoMatchingState = &core.StateTrigger{
	Name:      "GoMatchingState",
	ArgsTypes: []reflect.Type{},
}

type MatchingState struct {
	core.State

	testCFG            *config.TestConfig
	router             *core.Router
	userGroup          *commonmodel.UserGroup
	gameSetting        *model3.GameSetting
	participantGroup   *model3.ParticipantGroup
	playedHistoryGroup *model3.PlayedHistoryGroup
	statsGroup         *model3.StatsGroup
	playSettingGroup   *model3.PlaySettingGroup
	forceBuyInGroup    *model3.ForceBuyInGroup
	tableProfitsGroup  *model3.TableProfitsGroup

	participantService *service2.ParticipantService
	resyncService      *service2.ResyncService
	roleService        *service2.RoleService
	leaveService       *service2.LeaveService
	userAPI            commonapi.UserAPI
}

func ProvideMatchingState(
	stateFactory *core.StateFactory,
	testCFG *config.TestConfig,
	router *core.Router,
	userGroup *commonmodel.UserGroup,
	gameSetting *model3.GameSetting,
	participantGroup *model3.ParticipantGroup,
	playedHistoryGroup *model3.PlayedHistoryGroup,
	statsGroup *model3.StatsGroup,
	playSettingGroup *model3.PlaySettingGroup,
	forceBuyInGroup *model3.ForceBuyInGroup,
	tableProfitsGroup *model3.TableProfitsGroup,

	participantService *service2.ParticipantService,
	resyncService *service2.ResyncService,
	roleService *service2.RoleService,
	leaveService *service2.LeaveService,
	userAPI commonapi.UserAPI,
) *MatchingState {
	return &MatchingState{
		State:              stateFactory.Create("MatchingState"),
		testCFG:            testCFG,
		router:             router,
		userGroup:          userGroup,
		gameSetting:        gameSetting,
		participantGroup:   participantGroup,
		playedHistoryGroup: playedHistoryGroup,
		statsGroup:         statsGroup,
		playSettingGroup:   playSettingGroup,
		forceBuyInGroup:    forceBuyInGroup,
		tableProfitsGroup:  tableProfitsGroup,

		participantService: participantService,
		resyncService:      resyncService,
		roleService:        roleService,
		leaveService:       leaveService,
		userAPI:            userAPI,
	}
}

func (state *MatchingState) Run(ctx context.Context, args ...any) error {
	go state.GameController().RunTicker(time.Second, func() {
		state.Logger().Debug(
			"start matchmaking",
			zap.Object("participantGroup", state.participantGroup),
			zap.Object("userGroup", state.userGroup),
		)

		matchingParticipants := lo.Filter(
			lo.Values(state.participantGroup.Data),
			func(p *model3.Participant, _ int) bool {
				return p.FSM.MustState().(participant.State) == participant.MatchingState
			},
		)

		var matchedUidGroups [][]core.Uid

		matchingUids := lo.Map(matchingParticipants, func(p *model3.Participant, _ int) core.Uid { return p.Uid })
		lo.Shuffle(matchingUids)
		matchedUidGroups = lo.Filter(
			lo.Chunk(matchingUids, constant.PlayerPerMatch),
			func(uids []core.Uid, _ int) bool {
				return len(uids) == constant.PlayerPerMatch
			},
		)

		for _, uidGroup := range matchedUidGroups {
			if err := state.sendGroupToGame(uidGroup); err != nil {
				state.Logger().Error(
					"Failed to send group to game",
					zap.Error(err),
					zap.Array("uidGroup", core.UidList(uidGroup)),
				)
			}
		}
	})

	return nil
}

func (state *MatchingState) Publish(ctx context.Context, args ...any) error {
	return nil
}

func (state *MatchingState) sendGroupToGame(matchedUidGroup core.UidList) error {
	userSessionGroup := &session2.UserSessionGroup{
		Data: make(map[core.Uid]*session2.UserSession),
	}

	assignedRoles, err := state.roleService.Assign(matchedUidGroup)
	if err != nil {
		return err
	}

	forceBuyInSessionGroup := session2.NewForceBuyInSessionGroup()

	tableProfitsSessionGroup := &session2.TableProfitsSessionGroup{
		Data: make(map[core.Uid]*session2.TableProfitsSession),
	}

	for _, uid := range matchedUidGroup {

		if v, ok := state.forceBuyInGroup.Get(uid); ok {
			forceBuyInSessionGroup.Set(uid, v.GetBuyInChip(), v.GetExpireTime())
		}

		if v, ok := state.tableProfitsGroup.Get(uid); ok {
			tableProfitsSessionGroup.Data[uid] = &session2.TableProfitsSession{
				Uid:             v.Uid,
				Name:            v.Name,
				CountGames:      v.CountGames,
				SumBuyInChips:   v.SumBuyInChips,
				SumWinLoseChips: v.SumWinLoseChips,
			}
		}

		assignedRole := assignedRoles[uid]

		userSessionGroup.Data[uid] = &session2.UserSession{
			User: &commonmodel.User{
				Uid:         state.userGroup.Data[uid].Uid,
				ShortUid:    state.userGroup.Data[uid].ShortUid,
				Name:        state.userGroup.Data[uid].Name,
				IsAI:        state.userGroup.Data[uid].IsAI,
				IsConnected: state.userGroup.Data[uid].IsConnected,
				HasEntered:  state.userGroup.Data[uid].HasEntered,
				Cash:        state.userGroup.Data[uid].Cash,
				Level:       state.userGroup.Data[uid].Level,
				RoomCards:   state.userGroup.Data[uid].RoomCards,
			},
			Chip:            state.participantGroup.Data[uid].Chip,
			Role:            assignedRole,
			CountIdleRounds: state.participantGroup.Data[uid].CountIdleRounds,
		}
	}

	roundId := uuid.New().String()

	gameInfoSession := &session2.GameInfoSession{
		Setting: &model2.GameSetting{
			GameMetaUid:           state.gameSetting.GameMetaUid,
			SmallBlind:            state.gameSetting.SmallBlind,
			BigBlind:              state.gameSetting.BigBlind,
			TurnDuration:          state.gameSetting.TurnDuration,
			MinEnterLimitBB:       state.gameSetting.MinEnterLimitBB,
			MaxEnterLimitBB:       state.gameSetting.MaxEnterLimitBB,
			WaterPct:              state.gameSetting.WaterPct,
			TableSize:             state.gameSetting.TableSize,
			LeastGamePlayerAmount: state.gameSetting.LeastGamePlayerAmount,
			MaxWaterLimitBB:       state.gameSetting.MaxWaterLimitBB,
		},
		RoundId:                roundId,
		TotalBetPlayerCount:    state.playedHistoryGroup.TotalBetPlayerCount,
		TotalPlayedPlayerCount: state.playedHistoryGroup.TotalPlayedPlayerCount,
	}

	statsSessionGroup := &session2.StatsSessionGroup{
		Data: make(map[core.Uid]*session2.StatsSession),
	}

	for uid, st := range state.statsGroup.Data {
		statsSessionGroup.Data[uid] = &session2.StatsSession{
			Uid:                  st.Uid,
			HighestGameWinAmount: st.HighestGameWinAmount,
			EventAmountSum:       map[event.EventType]int{},
		}

		for e, count := range st.EventAmountSum {
			statsSessionGroup.Data[uid].EventAmountSum[e] = count
		}
	}

	state.Logger().Info(
		"sending to game",
		zap.Object("userSessionGroup", userSessionGroup),
		zap.Object("gameInfoSession", gameInfoSession),
		util.DebugField(zap.Object("tableProfitsSessionGroup", tableProfitsSessionGroup)),
	)

	zoomTXPokerGame, err := builder.BuildZoomTXPokerGame(
		userSessionGroup,
		gameInfoSession,
		statsSessionGroup,
		forceBuyInSessionGroup,
		tableProfitsSessionGroup,
		[]zap.Field{
			zap.String("zoom-module", "zoom-game"),
			zap.String("roundId", roundId),
		},
	)

	if err != nil {
		state.Logger().Error(
			"Failed to build game",
			zap.Error(err),
			zap.Object("playerSessionGroup", userSessionGroup),
			zap.Object("gameInfoSession", gameInfoSession),
		)
		return err
	}

	unsubFuncs := make([]core.UnsubscribeFunc, 0)

	unsub, err := zoomTXPokerGame.MsgBus().Subscribe("pool", session2.PoolCloseResultTopic, state.handlePoolCloseResult)
	if err != nil {
		state.Logger().Error("Failed to subscribe", zap.Error(err), zap.String("roundId", roundId))
		state.rollbackBind(unsubFuncs, matchedUidGroup)
		return err
	}
	unsubFuncs = append(unsubFuncs, unsub)

	for _, uid := range matchedUidGroup {
		unsub, err := zoomTXPokerGame.MsgBus().Subscribe(uid, session2.GameResultTopic, state.handleGameResult)
		if err != nil {
			state.Logger().Error("Failed to subscribe", zap.Error(err), zap.String("roundId", roundId))
			state.rollbackBind(unsubFuncs, matchedUidGroup)
			return err
		}
		unsubFuncs = append(unsubFuncs, unsub)

		unsub, err = zoomTXPokerGame.MsgBus().Subscribe(uid, session2.CashOutResultTopic, state.handleCashOutResult)
		if err != nil {
			state.Logger().Error("Failed to subscribe", zap.Error(err), zap.String("roundId", roundId))
			state.rollbackBind(unsubFuncs, matchedUidGroup)
			return err
		}
		unsubFuncs = append(unsubFuncs, unsub)

		unsub, err = zoomTXPokerGame.MsgBus().Subscribe(uid, session2.CloseResultTopic, state.handleCloseResult)
		if err != nil {
			state.Logger().Error("Failed to subscribe", zap.Error(err), zap.String("roundId", roundId))
			state.rollbackBind(unsubFuncs, matchedUidGroup)
			return err
		}
		unsubFuncs = append(unsubFuncs, unsub)
	}

	for _, uid := range matchedUidGroup {
		if err := state.router.Bind(uid, zoomTXPokerGame.MsgBus(), zoomTXPokerGame.GameNotifier()); err != nil {
			state.Logger().Error(
				"Failed to bind player to game",
				zap.Error(err),
				zap.String("uid", uid.String()),
				zap.Object("playerSessionGroup", userSessionGroup),
				zap.String("roundId", roundId),
			)

			state.rollbackBind(unsubFuncs, matchedUidGroup)
			return err
		}
	}

	chipIntoGame := lo.MapValues(
		userSessionGroup.Data,
		func(s *session2.UserSession, uid core.Uid) int { return s.Chip },
	)

	if err := state.participantService.EnterGame(chipIntoGame); err != nil {
		state.rollbackBind(unsubFuncs, matchedUidGroup)
		return err
	}

	for _, uid := range matchedUidGroup {
		state.MsgBus().Unicast(uid, core.MessageTopic, &txpokergrpc.Message{
			Model: &txpokergrpc.Model{
				Participant: state.participantGroup.Data[uid].ToProto(),
			},
		})
	}

	go func(unsubFuncs []core.UnsubscribeFunc) {
		zoomTXPokerGame.Run()

		// Unsubscribe after game end.
		for _, unsub := range unsubFuncs {
			if err := unsub(); err != nil {
				state.Logger().Error("Failed to unsubscribe", zap.Error(err))
			}
		}
	}(unsubFuncs)

	return nil
}

func (state *MatchingState) handleGameResult(gameResult *session2.GameResultSession) {
	state.GameController().RunTask(func() {
		state.Logger().Info("receive game result", zap.Object("gameResult", gameResult))

		// Must do it anyway.
		if err := state.router.UnBind(gameResult.Uid); err != nil {
			state.Logger().Error("Failed to unbind", zap.Error(err), zap.Object("gameResult", gameResult))
		}

		// TODO: put into service, idle round count
		if _, partExists := state.participantGroup.Data[gameResult.Uid]; partExists {

			state.participantGroup.Data[gameResult.Uid].CountIdleRounds = gameResult.CountIdleRounds

			if state.participantGroup.Data[gameResult.Uid].CountIdleRounds == 0 {
				state.Logger().Info("reset idle rounds",
					zap.Object("gameResult", gameResult),
					zap.Object("participant", state.participantGroup.Data[gameResult.Uid]),
				)
			}

		}

		// cash out chip >= min enter limit, allow force buy in
		if gameResult.Chip > 0 && gameResult.Chip >= state.gameSetting.MinEnterLimitBB*state.gameSetting.BigBlind {
			state.forceBuyInGroup.Set(gameResult.Uid, gameResult.Chip, time.Now())
		}

		if state.needCashOut(gameResult.Uid, gameResult.Chip) {
			state.Logger().Info("participant need cashing out", zap.Object("gameResult", gameResult))
			go state.cashOut(gameResult.Uid, gameResult.Chip)
			return
		}

		topUpDone, err := state.autoTopUp(gameResult)

		go func() {
			if err != nil {
				state.Logger().Error("Failed to auto top up", zap.Error(err), zap.Object("gameResult", gameResult))
			} else {
				state.Logger().Debug("waiting top up done", zap.Object("gameResult", gameResult), zap.Any("topUpDone", topUpDone))
				<-topUpDone
			}

			state.Logger().Debug("before top up done", zap.Object("gameResult", gameResult))
			state.GameController().RunTask(func() {
				state.Logger().Info("top up done", zap.Object("gameResult", gameResult))
				state.exitGame(gameResult)
			})
		}()

	})
}

func (state *MatchingState) handlePoolCloseResult(result *session2.CloseResultSession) {
	state.GameController().RunTask(func() {
		if result == nil {
			return
		}

		state.updatePlayedHistory(result)

		state.Logger().Info("received pool close result updated",
			zap.Object("result", result),
			util.DebugField(zap.Object("playedHistoryGroup", state.playedHistoryGroup)),
		)
	})
}

func (state *MatchingState) handleCloseResult(result *session2.UserCloseResultSession) {
	state.GameController().RunTask(func() {
		if result == nil {
			return
		}

		state.updateStats(result)
		state.updateTableProfits(result)

		state.Logger().Info("received close result updated",
			zap.Object("result", result),
			util.DebugField(zap.Object("stats", state.statsGroup)),
			util.DebugField(zap.Object("tableProfitsGroup", state.tableProfitsGroup)),
		)
	})
}

func (state *MatchingState) updateTableProfits(result *session2.UserCloseResultSession) {
	tableProfits, ok := state.tableProfitsGroup.Get(result.Uid)
	if !ok {
		tableProfits = &model2.TableProfits{
			Uid:             result.Uid,
			Name:            result.Name,
			CountGames:      0,
			SumBuyInChips:   0,
			SumWinLoseChips: 0,
		}
	}

	tableProfits.CountGames += result.IncrCountGames
	tableProfits.SumBuyInChips += result.IncrBuyInChips
	tableProfits.SumWinLoseChips += result.IncrWinLoseChips
	state.tableProfitsGroup.Save(tableProfits)
}

func (state *MatchingState) updatePlayedHistory(results *session2.CloseResultSession) {
	// 這裡會重複計算，但因為是用來算 bet / played 比例的，所以不影響。
	// 如果要個別使用，就要想新的更新機制（for pool, not for user）
	state.playedHistoryGroup.TotalBetPlayerCount += results.IncrBetPlayerCount
	state.playedHistoryGroup.TotalPlayedPlayerCount += results.IncrPlayedPlayerCount
}

func (state *MatchingState) updateStats(result *session2.UserCloseResultSession) {
	// It could be deleted in destroying user
	if _, ok := state.statsGroup.Data[result.Uid]; !ok {
		state.Logger().Warn("stats not found", zap.String("uid", result.Uid.String()))
		return
	}

	if result.HighestGameWinAmount > state.statsGroup.Data[result.Uid].HighestGameWinAmount {
		state.statsGroup.Data[result.Uid].HighestGameWinAmount = result.HighestGameWinAmount
	}

	for e, amount := range result.IncrEventAmount {
		state.statsGroup.Data[result.Uid].EventAmountSum[e] += amount
	}
}

func (state *MatchingState) handleCashOutResult(cashOutResult *session2.CashOutResultSession) {
	state.GameController().RunTask(func() {
		state.Logger().Info("receive cash out result", zap.Object("cashOutResult", cashOutResult))

		if state.needCashOut(cashOutResult.Uid, cashOutResult.Chip) {
			state.Logger().Info("participant need cashing out", zap.Object("cashOutResult", cashOutResult))
			go state.cashOut(cashOutResult.Uid, cashOutResult.Chip)
			return
		}

		part := state.participantGroup.Data[cashOutResult.Uid]

		if err := state.participantService.AddChip(part.Uid, cashOutResult.Chip); err != nil {
			state.Logger().Error("failed to add chip from cash out result", zap.Error(err), zap.Object("cashOutResult", cashOutResult))
		}

		state.MsgBus().Unicast(part.Uid, core.MessageTopic, &txpokergrpc.Message{
			Model: &txpokergrpc.Model{
				Participant: part.ToProto(),
			},
		})
	})
}

func (state *MatchingState) cashOut(uid core.Uid, chip int) {
	if err := state.userAPI.ExchangeChip(uid, gametype.ZoomTXPoker, chip); err != nil {
		state.Logger().Error("Failed to cash out", zap.Error(err), zap.String("uid", uid.String()), zap.Int("chip", chip))
	}
}

func (state *MatchingState) needCashOut(uid core.Uid, chips int) bool {
	if chips == 0 {
		return false
	}

	part, partExists := state.participantGroup.Data[uid]
	if !partExists {
		return true
	}

	if !lo.Contains(
		[]participant.State{participant.MatchingState, participant.PlayingState},
		part.FSM.MustState().(participant.State),
	) {
		return true
	}

	return false
}

func (state *MatchingState) autoTopUp(gameResult *session2.GameResultSession) (<-chan struct{}, error) {
	part, ok := state.participantGroup.Data[gameResult.Uid]
	if !ok {
		return nil, fmt.Errorf("participant not found uid %v", gameResult.Uid)
	}

	state.Logger().Info("start auto top up",
		zap.Object("gameResult", gameResult),
		zap.Int("MaxEnterLimitBB", state.gameSetting.MaxEnterLimitBB),
		zap.Int("BigBlind", state.gameSetting.BigBlind),
		zap.Object("playSetting", state.playSettingGroup.Data[gameResult.Uid]),
		zap.Object("participant", part),
	)

	maxEnterLimitChip := state.gameSetting.MaxEnterLimitBB * state.gameSetting.BigBlind

	// 之後的 exit game 會加回去 gameResult.Chip，這裡先用來計算
	newPartChip := gameResult.Chip + part.Chip

	// > max, skip
	if newPartChip > maxEnterLimitChip {
		state.Logger().Info("no need to top up",
			zap.Object("gameResult", gameResult),
			zap.Int("partChip", newPartChip),
			zap.Int("maxEnterLimitChip", maxEnterLimitChip),
		)
		topUpDone := make(chan struct{})
		close(topUpDone)
		return topUpDone, nil
	}

	// if part.QueuedTopUpChip + part.Chip > max chip, then top up to max chip
	topUpChip := part.QueuedTopUpChip
	if topUpChip+newPartChip > maxEnterLimitChip {
		topUpChip = maxEnterLimitChip - newPartChip
	}

	playSetting := state.playSettingGroup.Data[gameResult.Uid]

	if playSetting.AutoTopUp &&
		topUpChip+newPartChip < int(float64(maxEnterLimitChip)*playSetting.AutoTopUpThresholdPercent) {
		topUpChip = int(float64(maxEnterLimitChip)*playSetting.AutoTopUpChipPercent) - newPartChip
	}

	part.QueuedTopUpChip = 0

	if topUpChip <= 0 {
		state.Logger().Info("no need to top up", zap.Object("gameResult", gameResult))
		topUpDone := make(chan struct{})
		close(topUpDone)
		return topUpDone, nil
	}

	topUpDone, err := state.participantService.TopUp(gameResult.Uid, topUpChip)
	if err != nil {
		state.Logger().Error(
			"failed to auto top up user", zap.Error(err),
			zap.Object("gameResult", gameResult),
		)
		return nil, err
	}

	return topUpDone, nil
}

func (state *MatchingState) exitGame(gameResult *session2.GameResultSession) {
	part, ok := state.participantGroup.Data[gameResult.Uid]
	if !ok {
		state.Logger().Warn("exitGame but participant not found", zap.Object("gameResult", gameResult))
		return
	}

	if err := state.participantService.ExitGame(gameResult.Uid, gameResult.Chip); err != nil {
		state.Logger().Error("Failed to exit game", zap.Error(err), zap.Object("gameResult", gameResult))
	}

	// Must after ExitGame, make sure user is not in PlayingState.
	state.resyncService.Send(gameResult.Uid)

	// Idle or Not Enough Chip(無法負擔BigBlind Role 的費用)
	if part.IsIdleRoundsReachMax() || part.Chip <= state.gameSetting.BigBlind {
		// User will leave matching but return to observing
		if _, err := state.participantService.ExitMatch(part.Uid, false); err != nil {
			state.Logger().Error("Failed to exit match for poor user", zap.Error(err), zap.Object("gameResult", gameResult), zap.Object("participant", part))
		}
	}
}

func (state *MatchingState) rollbackBind(unsubFuncs []core.UnsubscribeFunc, uids core.UidList) {
	for _, unsub := range unsubFuncs {
		if err := unsub(); err != nil {
			state.Logger().Error("Failed to unsubscribe", zap.Error(err))
		}
	}

	for _, uid := range uids {
		if err := state.router.UnBind(uid); err != nil {
			state.Logger().Error("Failed to unbind", zap.Error(err),
				zap.String("uid", uid.String()),
			)
		}
	}
}
