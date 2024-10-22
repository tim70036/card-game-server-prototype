package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/service"
	"card-game-server-prototype/pkg/game/txpoker/type/event"
	"card-game-server-prototype/pkg/game/txpoker/type/seatstatus"
	session2 "card-game-server-prototype/pkg/game/zoomtxpoker/game/session"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type ClosedState struct {
	core.State

	tableProfitsSessionGroup *session2.TableProfitsSessionGroup
	statsSessionGroup        *session2.StatsSessionGroup
	gameInfoSession          *session2.GameInfoSession

	playerGroup       *model2.PlayerGroup
	seatStatusGroup   *model2.SeatStatusGroup
	chipCacheGroup    *model2.ChipCacheGroup
	statsGroup        *model2.StatsGroup
	seatStatusService service.SeatStatusService
	gameInfo          *model2.GameInfo
	tableProfitsGroup *model2.TableProfitsGroup
}

func ProvideClosedState(
	stateFactory *core.StateFactory,

	tableProfitsSessionGroup *session2.TableProfitsSessionGroup,
	statsSessionGroup *session2.StatsSessionGroup,
	gameInfoSession *session2.GameInfoSession,

	playerGroup *model2.PlayerGroup,
	seatStatusGroup *model2.SeatStatusGroup,
	chipCacheGroup *model2.ChipCacheGroup,
	statsGroup *model2.StatsGroup,
	seatStatusService service.SeatStatusService,
	gameInfo *model2.GameInfo,
	tableProfitsGroup *model2.TableProfitsGroup,
) *ClosedState {
	return &ClosedState{
		State: stateFactory.Create("ClosedState"),

		tableProfitsSessionGroup: tableProfitsSessionGroup,
		statsSessionGroup:        statsSessionGroup,
		gameInfoSession:          gameInfoSession,

		playerGroup:       playerGroup,
		seatStatusGroup:   seatStatusGroup,
		chipCacheGroup:    chipCacheGroup,
		statsGroup:        statsGroup,
		seatStatusService: seatStatusService,
		gameInfo:          gameInfo,
		tableProfitsGroup: tableProfitsGroup,
	}
}

func (state *ClosedState) Run(ctx context.Context, args ...any) error {

	state.Logger().Debug(
		"sending results",
		zap.Object("seatStatusGroup", state.seatStatusGroup),
		zap.Object("chipCacheGroup", state.chipCacheGroup),
		zap.Object("stateGroup", state.statsGroup),
	)

	state.sendingCloseResults()

	playingSeatStatus := lo.Filter(
		lo.Values(state.seatStatusGroup.Status),
		func(s *model2.SeatStatus, _ int) bool {
			return s.FSM.MustState().(seatstatus.SeatStatusState) == seatstatus.PlayingState
		},
	)

	for _, seatStatus := range playingSeatStatus {
		if _, err := state.seatStatusService.StandUp(seatStatus.Uid); err != nil {
			state.Logger().Error("failed to stand up", zap.Error(err))
		}
	}

	for uid, cashOutChip := range state.chipCacheGroup.CashOutChips {
		cashOutResult := &session2.CashOutResultSession{
			GameId: state.gameInfo.RoundId,
			Uid:    uid,
			Chip:   cashOutChip,
		}

		state.MsgBus().Unicast(uid, session2.CashOutResultTopic, cashOutResult)
	}

	state.GameController().RunTimer(5*time.Second, state.GameController().Shutdown)

	return nil
}

func (state *ClosedState) sendingCloseResults() {
	result := &session2.CloseResultSession{
		GameId:                state.gameInfo.RoundId,
		IncrBetPlayerCount:    state.gameInfo.TotalBetPlayerCount - state.gameInfoSession.TotalBetPlayerCount,
		IncrPlayedPlayerCount: state.gameInfo.TotalPlayedPlayerCount - state.gameInfoSession.TotalPlayedPlayerCount,
	}

	state.MsgBus().Broadcast(session2.PoolCloseResultTopic, result)

	for uid := range state.seatStatusGroup.Status {
		uResult := &session2.UserCloseResultSession{
			GameId:          state.gameInfo.RoundId,
			Uid:             uid,
			IncrEventAmount: map[event.EventType]int{},
		}

		// stats
		if stats, ok := state.statsGroup.Data[uid]; ok {
			uResult.HighestGameWinAmount = stats.HighestGameWinAmount

			for e, amount := range state.calculateIncrEventAmount(uid, stats.EventAmountSum) {
				uResult.IncrEventAmount[e] = amount
			}
		}

		// table profits
		if tp, ok := state.tableProfitsGroup.Get(uid); ok {
			uResult.Name = tp.Name

			uResult.IncrCountGames = tp.CountGames
			uResult.IncrBuyInChips = tp.SumBuyInChips
			uResult.IncrWinLoseChips = tp.SumWinLoseChips

			if tpSession, ok := state.tableProfitsSessionGroup.Data[uid]; ok {
				uResult.IncrCountGames -= tpSession.CountGames
				uResult.IncrBuyInChips -= tpSession.SumBuyInChips
				uResult.IncrWinLoseChips -= tpSession.SumWinLoseChips
			}
		}
		state.MsgBus().Unicast(uid, session2.CloseResultTopic, uResult)
	}
}

func (state *ClosedState) calculateIncrEventAmount(uid core.Uid, newEventAmount map[event.EventType]int) map[event.EventType]int {
	sumEventAmount := map[event.EventType]int{}
	for e, amount := range newEventAmount {
		sumEventAmount[e] = amount
	}

	oriStats, hasOriStats := state.statsSessionGroup.Data[uid]
	if !hasOriStats {
		return sumEventAmount
	}

	for e, cur := range sumEventAmount {
		if ori, ok := oriStats.EventAmountSum[e]; ok {
			sumEventAmount[e] = cur - ori
		}
	}

	return sumEventAmount
}

func (state *ClosedState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
	})
	return nil
}

func (state *ClosedState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context:   &txpokergrpc.GameState_ClosedStateContext{ClosedStateContext: &txpokergrpc.ClosedStateContext{}},
	}
}

func (state *ClosedState) BeforeConnect(uid core.Uid) error {
	state.Logger().Warn("forbid to connect due to game closed", zap.String("uid", uid.String()))
	return status.Errorf(codes.Aborted, "forbid to connect due to game closed")
}
