// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package builder

import (
	"card-game-server-prototype/pkg/common"
	"card-game-server-prototype/pkg/common/api"
	model2 "card-game-server-prototype/pkg/common/model"
	service3 "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker"
	"card-game-server-prototype/pkg/game/txpoker/actor"
	api2 "card-game-server-prototype/pkg/game/txpoker/api"
	"card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/service"
	state2 "card-game-server-prototype/pkg/game/txpoker/state"
	"card-game-server-prototype/pkg/game/zoomtxpoker/game"
	api3 "card-game-server-prototype/pkg/game/zoomtxpoker/game/api"
	"card-game-server-prototype/pkg/game/zoomtxpoker/game/handler"
	service2 "card-game-server-prototype/pkg/game/zoomtxpoker/game/service"
	"card-game-server-prototype/pkg/game/zoomtxpoker/game/session"
	"card-game-server-prototype/pkg/game/zoomtxpoker/game/state"
	"card-game-server-prototype/pkg/util"
	"github.com/google/wire"
	"go.uber.org/zap/zapcore"
)

// Injectors from build_zoom_tx_poker_game.go:

func BuildZoomTXPokerGame(userSessionGroup *session.UserSessionGroup, gameInfoSession *session.GameInfoSession, statsSessionGroup *session.StatsSessionGroup, forceBuyInSessionGroup *session.ForceBuyInSessionGroup, tableProfitsSessionGroup *session.TableProfitsSessionGroup, arg []zapcore.Field) (*game.ZoomTXPokerGame, error) {
	logConfig := _wireLogConfigValue
	loggerFactory := util.ProvideLoggerFactory(logConfig, arg)
	coreGame := core.ProvideGame(loggerFactory)
	msgBus := core.ProvideMsgBus(loggerFactory)
	stateFactory := core.ProvideStateFactory(coreGame, msgBus, loggerFactory)
	playerGroup := model.ProvidePlayerGroup()
	seatStatusGroup := model.ProvideSeatStatusGroup()
	chipCacheGroup := model.ProvideChipCacheGroup()
	statsGroup := model.ProvideStatsGroup()
	userGroup := model2.ProvideUserGroup(loggerFactory)
	actionHintGroup := model.ProvideActionHintGroup()
	playSettingGroup := model.ProvidePlaySettingGroup()
	gameSetting := model.ProvideGameSetting()
	gameInfo := model.ProvideGameInfo()
	config := _wireConfigValue
	client := util.ProvideHttpClient(loggerFactory, logConfig, config)
	apiConfig := _wireAPIConfigValue
	baseUserAPI := api.ProvideBaseUserAPI(client, apiConfig, config)
	forceBuyInGroup := model.ProvideForceBuyInGroup()
	eventGroup := model.ProvideEventGroup()
	tableProfitsGroup := model.ProvideTableProfitsGroup()
	baseSeatStatusService := service.ProvideBaseSeatStatusService(coreGame, msgBus, userGroup, seatStatusGroup, actionHintGroup, playerGroup, playSettingGroup, statsGroup, gameSetting, gameInfo, baseUserAPI, loggerFactory, forceBuyInGroup, eventGroup, tableProfitsGroup)
	serviceBaseSeatStatusService := service2.ProvideBaseSeatStatusService(baseSeatStatusService, seatStatusGroup, forceBuyInGroup, gameSetting, gameInfo, msgBus, loggerFactory)
	closedState := state.ProvideClosedState(stateFactory, tableProfitsSessionGroup, statsSessionGroup, gameInfoSession, playerGroup, seatStatusGroup, chipCacheGroup, statsGroup, serviceBaseSeatStatusService, gameInfo, tableProfitsGroup)
	roomInfo := model2.ProvideRoomInfo()
	replay := model.ProvideReplay()
	table := model.ProvideTable()
	actorGroup := actor.ProvideActorGroup()
	baseActorFactory := actor.ProvideBaseActorFactory(userGroup, seatStatusGroup, loggerFactory)
	dummyActorFactory := actor.ProvideDummyActorFactory(loggerFactory, actionHintGroup, gameSetting)
	buddyGroup := model2.ProvideBuddyGroup(loggerFactory)
	statsCacheGroup := model.ProvideStatsCacheGroup()
	userCacheGroup := model.ProvideUserCacheGroup()
	resyncService := service.ProvideResyncService(coreGame, userGroup, buddyGroup, roomInfo, gameInfo, gameSetting, table, seatStatusGroup, actionHintGroup, playSettingGroup, statsCacheGroup, chipCacheGroup, userCacheGroup, playerGroup, tableProfitsGroup, msgBus)
	initState := state.ProvideInitState(stateFactory, userSessionGroup, gameInfoSession, statsSessionGroup, forceBuyInSessionGroup, tableProfitsSessionGroup, config, roomInfo, userGroup, seatStatusGroup, chipCacheGroup, playSettingGroup, eventGroup, statsGroup, gameSetting, gameInfo, forceBuyInGroup, tableProfitsGroup, playerGroup, actionHintGroup, replay, table, actorGroup, baseActorFactory, dummyActorFactory, serviceBaseSeatStatusService, resyncService)
	resetState := state.ProvideResetState(stateFactory)
	actionHintService := service.ProvideActionHintService(actionHintGroup, seatStatusGroup, playerGroup, chipCacheGroup, gameSetting, table, replay, loggerFactory)
	baseGameAPI := api2.ProvideBaseGameAPI(client, apiConfig, roomInfo)
	apiBaseGameAPI := api3.ProvideBaseGameAPI(baseGameAPI, client, apiConfig, roomInfo)
	baseGameRepoService := service.ProvideBaseGameRepoService(roomInfo, gameInfo, gameSetting, playerGroup, seatStatusGroup, tableProfitsGroup, replay, table, apiBaseGameAPI, msgBus, loggerFactory)
	startRoundState := state.ProvideStartRoundState(stateFactory, playerGroup, actionHintGroup, seatStatusGroup, gameSetting, gameInfo, table, statsGroup, statsCacheGroup, actionHintService, baseGameRepoService)
	testConfig := _wireTestConfigValue
	dealPocketState := state2.ProvideDealPocketState(stateFactory, playerGroup, table, testConfig, roomInfo)
	evaluateActionState := state2.ProvideEvaluateActionState(stateFactory, actionHintGroup, seatStatusGroup, playerGroup, table, actionHintService)
	collectChipState := state2.ProvideCollectChipState(stateFactory, seatStatusGroup, chipCacheGroup, table, actionHintGroup, replay, actionHintService)
	dealCommunityState := state2.ProvideDealCommunityState(stateFactory, table, testConfig, roomInfo)
	waitActionState := state2.ProvideWaitActionState(stateFactory, seatStatusGroup, chipCacheGroup, actionHintGroup, actorGroup, playerGroup, gameInfo, serviceBaseSeatStatusService, actionHintService)
	foldState := state2.ProvideFoldState(stateFactory)
	checkState := state2.ProvideCheckState(stateFactory)
	betState := state2.ProvideBetState(stateFactory)
	callState := state2.ProvideCallState(stateFactory)
	raiseState := state2.ProvideRaiseState(stateFactory)
	allInState := state2.ProvideAllInState(stateFactory)
	declareShowdownState := state2.ProvideDeclareShowdownState(stateFactory)
	showdownState := state2.ProvideShowdownState(stateFactory, actionHintGroup, playerGroup, table, replay)
	dealRemainCommunityState := state2.ProvideDealRemainCommunityState(stateFactory, table, testConfig, roomInfo)
	baseJackpotService := service2.ProvideBaseJackpotService(gameSetting)
	evaluateWinnerState := state2.ProvideEvaluateWinnerState(stateFactory, gameSetting, playerGroup, actionHintGroup, table, replay, roomInfo, baseJackpotService)
	declareWinnerState := state2.ProvideDeclareWinnerState(stateFactory, roomInfo, seatStatusGroup, playerGroup, chipCacheGroup, table, baseUserAPI)
	jackpotState := state2.ProvideJackpotState(stateFactory, playerGroup, baseGameRepoService)
	baseEventService := service.ProvideBaseEventService(roomInfo, gameSetting, eventGroup, playerGroup, replay, table, apiBaseGameAPI, loggerFactory)
	serviceBaseEventService := service2.ProvideBaseEventService(baseEventService, roomInfo, eventGroup, apiBaseGameAPI, loggerFactory)
	endRoundState := state2.ProvideEndRoundState(stateFactory, gameInfo, roomInfo, seatStatusGroup, actionHintGroup, eventGroup, statsGroup, playerGroup, replay, serviceBaseSeatStatusService, baseGameRepoService, serviceBaseEventService)
	initiator := state.Init(coreGame, closedState, initState, resetState, startRoundState, dealPocketState, evaluateActionState, collectChipState, dealCommunityState, waitActionState, foldState, checkState, betState, callState, raiseState, allInState, declareShowdownState, showdownState, dealRemainCommunityState, evaluateWinnerState, declareWinnerState, jackpotState, endRoundState)
	baseRoomAPI := api.ProvideBaseRoomAPI(client, apiConfig, config)
	baseUserService := service3.ProvideBaseUserService(userGroup, msgBus, roomInfo, baseUserAPI, baseRoomAPI, loggerFactory)
	serviceBaseUserService := service.ProvideBaseUserService(baseUserService, msgBus, testConfig, roomInfo, userGroup, baseRoomAPI, seatStatusGroup, eventGroup, playSettingGroup, statsGroup, userCacheGroup, baseSeatStatusService, tableProfitsGroup)
	connectionHandler := handler.ProvideConnectionHandler(coreGame, userGroup, seatStatusGroup, serviceBaseUserService, serviceBaseSeatStatusService, resyncService, msgBus, loggerFactory)
	requestHandler := handler.ProvideRequestHandler(msgBus, eventGroup, loggerFactory)
	handlerInitiator := handler.Init(coreGame, connectionHandler, requestHandler)
	zoomTXPokerGame := game.ProvideZoomTXPokerGame(coreGame, initiator, handlerInitiator, msgBus, loggerFactory)
	return zoomTXPokerGame, nil
}

var (
	_wireLogConfigValue  = config.LogCFG
	_wireConfigValue     = config.CFG
	_wireAPIConfigValue  = config.APICFG
	_wireTestConfigValue = config.TestCFG
)

// build_zoom_tx_poker_game.go:

var zoomTXPokerGameProviderSet = wire.NewSet(config.ProviderSet, core.ProviderSet, util.ProviderSet, common.ProviderSet, txpoker.ProviderSet, game.ProviderSet)
