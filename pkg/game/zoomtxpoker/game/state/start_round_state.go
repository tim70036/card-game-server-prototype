package state

import (
	"context"
	"fmt"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	service2 "card-game-server-prototype/pkg/game/txpoker/service"
	txpokerstate "card-game-server-prototype/pkg/game/txpoker/state"
	"card-game-server-prototype/pkg/game/txpoker/type/action"
	"card-game-server-prototype/pkg/game/txpoker/type/role"
	"card-game-server-prototype/pkg/game/txpoker/type/stage"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StartRoundState struct {
	core.State

	playerGroup     *model2.PlayerGroup
	actionHintGroup *model2.ActionHintGroup
	seatStatusGroup *model2.SeatStatusGroup
	gameSetting     *model2.GameSetting
	gameInfo        *model2.GameInfo
	table           *model2.Table
	statsGroup      *model2.StatsGroup
	statsCacheGroup *model2.StatsCacheGroup

	actionHintService *service2.ActionHintService
	gameRepoService   service2.GameRepoService
}

func ProvideStartRoundState(
	stateFactory *core.StateFactory,

	playerGroup *model2.PlayerGroup,
	actionHintGroup *model2.ActionHintGroup,
	seatStatusGroup *model2.SeatStatusGroup,
	gameSetting *model2.GameSetting,
	gameInfo *model2.GameInfo,
	table *model2.Table,
	statsGroup *model2.StatsGroup,
	statsCacheGroup *model2.StatsCacheGroup,

	actionHintService *service2.ActionHintService,
	gameRepoService service2.GameRepoService,
) *StartRoundState {
	return &StartRoundState{
		State: stateFactory.Create("StartRoundState"),

		playerGroup:     playerGroup,
		actionHintGroup: actionHintGroup,
		seatStatusGroup: seatStatusGroup,
		gameSetting:     gameSetting,
		gameInfo:        gameInfo,
		table:           table,
		statsGroup:      statsGroup,
		statsCacheGroup: statsCacheGroup,

		actionHintService: actionHintService,
		gameRepoService:   gameRepoService,
	}
}

func (state *StartRoundState) Run(ctx context.Context, args ...any) error {
	state.Logger().Debug("starting round",
		zap.Object("gameInfo", state.gameInfo),
	)

	for uid := range state.playerGroup.Data {
		state.statsCacheGroup.Data[uid] = state.statsGroup.Data[uid]
	}
	state.Logger().Info("stats cache updated", zap.Object("statsCacheGroup", state.statsCacheGroup))

	if err := state.placeStartRoundBet(); err != nil {
		state.Logger().Error(
			"failed to place start round bet",
			zap.Error(err),
			zap.String("roundId", state.gameInfo.RoundId),
			zap.Object("players", state.playerGroup),
			zap.Object("gameSetting", state.gameSetting),
		)
		state.GameController().GoErrorState()
		return nil
	}

	state.Logger().Info("starting round, all set",
		zap.String("roundId", state.gameInfo.RoundId),
		zap.Object("playerGroup", state.playerGroup),
		zap.Object("gameInfo", state.gameInfo),
	)

	if err := state.table.BetStageFSM.Fire(stage.NextStageTrigger); err != nil {
		state.Logger().Error(
			"failed to fire next stage trigger",
			zap.Error(err),
			zap.Object("table", state.table),
		)
		state.GameController().GoErrorState()
		return nil
	}

	if err := state.gameRepoService.CreateRound(); err != nil {
		state.Logger().Error(
			"failed to create round",
			zap.Error(err),
			zap.Object("gameInfo", state.gameInfo),
		)
		state.GameController().GoErrorState()
		return nil
	}

	state.GameController().GoNextState(txpokerstate.GoDealPocketState)
	return nil
}

func (state *StartRoundState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
		Model: &txpokergrpc.Model{
			SeatStatusGroup: state.seatStatusGroup.ToProto(),
			ActionHintGroup: state.actionHintGroup.ToProto(),
			GameInfo:        state.gameInfo.ToProto(),
			StatsCacheGroup: state.statsCacheGroup.ToProto(),
		},
	})
	return nil
}

func (state *StartRoundState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context:   &txpokergrpc.GameState_StartRoundStateContext{StartRoundStateContext: &txpokergrpc.StartRoundStateContext{}},
	}
}

// Place sb, bb bets
// Do it in this state to prevent race condition such as:
// BB player stand up before placing BB bet.
// Oh wait, we don't have stand up feature in zoom tx poker.
// This is just same logic flow copied from original tx poker.
func (state *StartRoundState) placeStartRoundBet() error {
	bbPlayer, ok := lo.Find(lo.Values(state.playerGroup.Data), func(p *model2.Player) bool { return p.Role == role.BB })
	if !ok {
		return fmt.Errorf("bb player not found")
	}

	if err := state.actionHintService.OpenBet(bbPlayer.Uid, action.BB, state.gameSetting.BigBlind); err != nil {
		return fmt.Errorf("failed to open bet for bb: %w", err)
	}

	sbPlayer, ok := lo.Find(lo.Values(state.playerGroup.Data), func(p *model2.Player) bool { return p.Role == role.SB })
	if !ok {
		return fmt.Errorf("sb player not found")
	}

	if err := state.actionHintService.FollowBet(sbPlayer.Uid, action.SB); err != nil {
		return fmt.Errorf("failed to follow bet for sb: %w", err)
	}

	state.Logger().Debug("bets placed", zap.Object("actionHints", state.actionHintGroup))
	return nil
}
