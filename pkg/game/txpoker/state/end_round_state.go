package state

import (
	"context"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/constant"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/service"
	"card-game-server-prototype/pkg/game/txpoker/type/action"
	"card-game-server-prototype/pkg/game/txpoker/type/event"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoEndRoundState = &core.StateTrigger{
	Name:      "GoEndRoundState",
	ArgsTypes: []reflect.Type{},
}

type EndRoundState struct {
	core.State

	gameInfo        *model2.GameInfo
	roomInfo        *commonmodel.RoomInfo
	seatStatusGroup *model2.SeatStatusGroup
	actionHintGroup *model2.ActionHintGroup
	eventGroup      *model2.EventGroup
	statsGroup      *model2.StatsGroup
	playerGroup     *model2.PlayerGroup
	replay          *model2.Replay

	seatStatusService service.SeatStatusService
	gameRepoService   service.GameRepoService
	eventService      service.EventService
}

func ProvideEndRoundState(
	stateFactory *core.StateFactory,
	gameInfo *model2.GameInfo,
	roomInfo *commonmodel.RoomInfo,
	seatStatusGroup *model2.SeatStatusGroup,
	actionHintGroup *model2.ActionHintGroup,
	eventGroup *model2.EventGroup,
	statsGroup *model2.StatsGroup,
	playerGroup *model2.PlayerGroup,
	replay *model2.Replay,

	seatStatusService service.SeatStatusService,
	gameRepoService service.GameRepoService,
	eventService service.EventService,
) *EndRoundState {
	return &EndRoundState{
		State: stateFactory.Create("EndRoundState"),

		gameInfo:        gameInfo,
		roomInfo:        roomInfo,
		seatStatusGroup: seatStatusGroup,
		actionHintGroup: actionHintGroup,
		eventGroup:      eventGroup,
		statsGroup:      statsGroup,
		playerGroup:     playerGroup,
		replay:          replay,

		seatStatusService: seatStatusService,
		gameRepoService:   gameRepoService,
		eventService:      eventService,
	}
}

func (state *EndRoundState) Run(ctx context.Context, args ...any) error {
	if err := state.seatStatusService.RoundEnd(); err != nil {
		state.Logger().Error(
			"failed to end round for seat status group",
			zap.Error(err),
			zap.Object("seatStatus", state.seatStatusGroup),
		)
		state.GameController().GoErrorState()
		return nil
	}

	if err := state.gameRepoService.SubmitRoundScore(); err != nil {
		state.Logger().Error("failed to submit round score", zap.Error(err))
		state.GameController().GoErrorState()
		return nil
	}

	if err := state.eventService.EvalRoundEvents(); err != nil {
		state.Logger().Error("eval round events failed", zap.Error(err))
		state.GameController().GoErrorState()
		return nil
	}

	state.updateStats()
	state.updateGameInfo()

	if err := state.eventService.Submit(); err != nil {
		state.Logger().Error("submit round events failed", zap.Error(err))
		state.GameController().GoErrorState()
		return nil
	}

	state.GameController().RunTimer(2000*time.Millisecond, func() {
		state.GameController().GoNextState(GoResetState)
	})
	return nil
}

func (state *EndRoundState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
		Model: &txpokergrpc.Model{
			// SeatStatusGroup: state.seatStatusGroup.ToProto(), Don't
			// publish seat status group since we want frontend to
			// render player group and seat status group together in
			// next state (ResetState).
			ActionHintGroup: state.actionHintGroup.ToProto(),
			GameInfo:        state.gameInfo.ToProto(),
		},
	})
	return nil
}

func (state *EndRoundState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context:   &txpokergrpc.GameState_EndRoundStateContext{EndRoundStateContext: &txpokergrpc.EndRoundStateContext{}},
	}
}

func (state *EndRoundState) updateStats() {
	for uid, events := range state.eventGroup.Data {
		stats := state.statsGroup.Data[uid]
		for _, e := range events {
			if e.Type == event.GameWinAmount {
				stats.HighestGameWinAmount = lo.Ternary(stats.HighestGameWinAmount < e.Amount, e.Amount, stats.HighestGameWinAmount)
			}

			stats.EventAmountSum[e.Type] += e.Amount
		}
	}

	state.Logger().Debug("stats updated", zap.Object("stats", state.statsGroup))
}

func (state *EndRoundState) updateGameInfo() {
	state.gameInfo.RoundIdHistory = append(state.gameInfo.RoundIdHistory, state.gameInfo.RoundId)
	if len(state.gameInfo.RoundIdHistory) > constant.MaxRoundIdHistoryCount {
		state.gameInfo.RoundIdHistory = state.gameInfo.RoundIdHistory[1:]
	}

	state.gameInfo.TotalPlayedPlayerCount += len(state.playerGroup.Data)

	betPlayers := map[core.Uid]struct{}{}
	for _, r := range lo.Flatten(lo.Values(state.replay.ActionLog)) {
		if lo.Contains([]action.ActionType{action.Bet, action.Call, action.Raise, action.AllIn}, r.GetType()) {
			betPlayers[r.GetUid()] = struct{}{}
		}
	}

	state.gameInfo.TotalBetPlayerCount += len(betPlayers)

	state.Logger().Info("end round",
		zap.Object("gameInfo", state.gameInfo),
	)
}
