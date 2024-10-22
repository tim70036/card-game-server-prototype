package state

import (
	"context"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/common/type/gamemode"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/actor"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	service2 "card-game-server-prototype/pkg/game/darkchess/service"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/durationpb"
	"reflect"
	"sort"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoGameScoreboardState = core.NewStateTrigger("GoGameScoreboardState")

type GameScoreboardState struct {
	core.State

	roomInfo               *commonmodel.RoomInfo
	playerGroup            *model2.PlayerGroup
	gameScoreboard         *model2.GameScoreboard
	roundScoreboardRecords *model2.RoundScoreboardRecords
	gameRepoService        service2.GameRepoService
	actorGroup             *actor.Group
	eventGroup             *model2.EventGroup
	eventService           service2.EventService

	duration time.Duration

	// Record the player requests to skip waiting duration for
	// scoreboard.
	requestToSkipPlayers map[core.Uid]struct{}

	cancelTimer context.CancelFunc
	endOnce     *sync.Once
}

func ProvideGameScoreboardState(
	stateFactory *core.StateFactory,

	roomInfo *commonmodel.RoomInfo,
	gameScoreboard *model2.GameScoreboard,
	roundScoreboardRecords *model2.RoundScoreboardRecords,
	gameRepoService service2.GameRepoService,
	playerGroup *model2.PlayerGroup,
	actorGroup *actor.Group,
	eventGroup *model2.EventGroup,
	eventService service2.EventService,
) *GameScoreboardState {
	return &GameScoreboardState{
		State: stateFactory.Create("GameScoreboardState"),

		roomInfo:               roomInfo,
		gameScoreboard:         gameScoreboard,
		roundScoreboardRecords: roundScoreboardRecords,
		gameRepoService:        gameRepoService,
		playerGroup:            playerGroup,
		actorGroup:             actorGroup,
		eventGroup:             eventGroup,
		eventService:           eventService,
		duration: lo.Ternary(
			roomInfo.GameMode == gamemode.Buddy,
			10*time.Second,
			50*time.Second,
		),
	}
}

func (state *GameScoreboardState) Run(context.Context, ...any) error {
	state.cancelTimer = func() {}
	state.endOnce = &sync.Once{}
	state.requestToSkipPlayers = map[core.Uid]struct{}{}

	// sum round scores
	for _, scores := range state.roundScoreboardRecords.Data {
		for uid, score := range scores.Data {
			if _, ok := state.gameScoreboard.Data[uid]; !ok {
				continue
			}
			state.gameScoreboard.Data[uid].RawProfit += score.RawProfit
			state.gameScoreboard.Data[uid].Profit += score.Profit
		}
	}

	// get last round raw profit
	lastRoundScore, _ := lo.Last(state.roundScoreboardRecords.Data)
	lastRoundRawProfit := map[core.Uid]int{}
	for u, v := range lastRoundScore.Data {
		lastRoundRawProfit[u] = v.RawProfit
	}

	for uid, curPlayer := range state.playerGroup.Data {
		// update disconnected
		if curPlayer.DisconnectedRoundCount > 0 {
			state.gameScoreboard.Data[uid].IsDisconnected = true
			state.gameScoreboard.Data[uid].DisconnectedRawProfit = lastRoundRawProfit[uid]
		}

		// update bankrupt
		state.gameScoreboard.Data[uid].IsBankrupt = curPlayer.IsBankrupt
	}

	// Eval winner or draw
	var totalRawProfits []int
	for _, v := range state.gameScoreboard.Data {
		totalRawProfits = append(totalRawProfits, v.RawProfit)
	}

	sort.Slice(totalRawProfits, func(i, j int) bool {
		return totalRawProfits[i] > totalRawProfits[j]
	})

	if totalRawProfits[0]-totalRawProfits[1] == 0 {
		state.gameScoreboard.IsDraw = true
	} else {
		state.gameScoreboard.WinnerId = lo.Keys(lo.PickBy(state.gameScoreboard.Data, func(key core.Uid, score *model2.GameScore) bool {
			return score.RawProfit == totalRawProfits[0]
		}))[0]
	}

	state.Logger().Debug("game scoreboard evaluated",
		zap.Object("gameScoreboard", state.gameScoreboard),
	)

	if err := state.gameRepoService.SubmitGameScore(); err != nil {
		state.Logger().Error("failed to submit game score", zap.Error(err))
		state.GameController().GoErrorState()
		return nil
	}

	// Need delay... since the game result is eventual consistent.
	time.Sleep(1000 * time.Millisecond)

	if err := state.gameRepoService.FetchGameResult(); err != nil {
		state.Logger().Error("failed to update from game result", zap.Error(err))
		state.GameController().GoErrorState()
		return nil
	}

	state.Logger().Debug("game scoreboard updated from submit result",
		zap.Object("gameScoreboard", state.gameScoreboard),
	)

	if err := state.eventService.EvalGameEvents(); err != nil {
		state.Logger().Error("eval game events failed", zap.Error(err))
		state.GameController().GoErrorState()
		return nil
	}

	state.Logger().Debug("game events evaluated", zap.Object("events", state.eventGroup))

	if err := state.eventService.Submit(); err != nil {
		state.Logger().Error("submit game events failed", zap.Error(err))
		state.GameController().GoErrorState()
		return nil
	}

	state.cancelTimer = state.GameController().RunTimer(state.duration, state.goNextState)

	for _, curActor := range state.actorGroup.Data {
		reqs, err := curActor.DecideSkipScoreboard()
		if err != nil {
			state.Logger().Error("actor cannot decide call",
				zap.Error(err),
				zap.Object("actor", curActor),
			)
			state.GameController().GoErrorState()
			return nil
		}

		state.GameController().RunActorRequests(reqs)
	}

	return nil
}

func (state *GameScoreboardState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	return nil
}

func (state *GameScoreboardState) ToProto(core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_GameScoreboardStateContext{
			GameScoreboardStateContext: &gamegrpc.GameScoreboardStateContext{
				Duration:   durationpb.New(state.duration),
				Scoreboard: state.gameScoreboard.ToProto(state.roomInfo.GameMode),
			}},
	}
}

func (state *GameScoreboardState) AcceptRequestTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(&gamegrpc.SkipScoreboardRequest{}),
	}
}

func (state *GameScoreboardState) HandleRequest(req *core.Request) error {
	myPlayer, ok := state.playerGroup.Data[req.Uid]
	if !ok {
		state.Logger().Warn("HandleRequest failed, player not found",
			zap.Object("req", req),
		)
		return status.Errorf(codes.NotFound, "player not found")
	}

	if _, ok := state.requestToSkipPlayers[myPlayer.Uid]; ok {
		state.Logger().Debug("player already request to skip game scoreboard", zap.Object("req", req))
		return status.Errorf(codes.AlreadyExists, "player already request to skip game scoreboard")
	}

	switch req.Msg.(type) {
	case *gamegrpc.SkipScoreboardRequest:
		state.requestToSkipPlayers[myPlayer.Uid] = struct{}{}
	default:
		state.Logger().Warn("not supported request", zap.Object("req", req))
		return status.Errorf(codes.Unimplemented, "not supported request")
	}

	if len(state.requestToSkipPlayers) == len(state.playerGroup.Data) {
		state.Logger().Debug("all players have requested to skip game scoreboard")
		state.GameController().RunTask(state.goNextState)
	}

	return nil
}

func (state *GameScoreboardState) goNextState() {
	state.endOnce.Do(func() {
		state.Logger().Debug("ending scoreboard")

		state.cancelTimer()
		state.GameController().GoNextState(GoEndGameState)
	})
}
