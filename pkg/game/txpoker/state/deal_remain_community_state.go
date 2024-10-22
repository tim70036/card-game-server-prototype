package state

import (
	"context"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/type/cheat"
	stage2 "card-game-server-prototype/pkg/game/txpoker/type/stage"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoDealRemainCommunityState = &core.StateTrigger{
	Name:      "GoDealRemainCommunityState",
	ArgsTypes: []reflect.Type{},
}

type DealRemainCommunityState struct {
	core.State

	table    *model.Table
	testCFG  *config.TestConfig
	roomInfo *commonmodel.RoomInfo
}

func ProvideDealRemainCommunityState(
	stateFactory *core.StateFactory,
	table *model.Table,
	testCFG *config.TestConfig,
	roomInfo *commonmodel.RoomInfo,
) *DealRemainCommunityState {
	return &DealRemainCommunityState{
		State: stateFactory.Create("DealRemainCommunityState"),

		table:    table,
		testCFG:  testCFG,
		roomInfo: roomInfo,
	}
}

func (state *DealRemainCommunityState) Run(ctx context.Context, args ...any) error {
	duration := 500 * time.Millisecond
	for state.table.BetStageFSM.MustState().(stage2.Stage) != stage2.ShowdownStage {
		betStage := state.table.BetStageFSM.MustState().(stage2.Stage)
		dealCnt := lo.Ternary(betStage == stage2.FlopStage, 3, 1)
		duration += time.Duration(lo.Ternary(betStage == stage2.FlopStage, 1500, 900)) * time.Millisecond

		dealCards := state.table.Deck[len(state.table.Deck)-dealCnt:]
		state.table.Deck = state.table.Deck[:len(state.table.Deck)-dealCnt]
		state.table.CommunityCards = append(state.table.CommunityCards, dealCards...)
		state.Logger().Debug("community cards dealt", zap.String("cards", state.table.CommunityCards.ToString()))

		if err := state.table.BetStageFSM.Fire(stage2.NextStageTrigger); err != nil {
			state.Logger().Error(
				"failed to fire next stage trigger",
				zap.Error(err),
				zap.Object("table", state.table),
			)
			state.GameController().GoErrorState()
			return nil
		}
	}

	if state.testCFG.EnableCheatMode(string(state.roomInfo.GameType)) {
		cheatData, err := cheat.FromRawCheatData(*state.testCFG.CheatData)
		if err != nil {
			state.Logger().Warn("failed to parse raw cheat data", zap.String("rawCheatData", *state.testCFG.CheatData), zap.Error(err))
		}

		state.table.CommunityCards = cheatData.CommunityCards[:len(state.table.CommunityCards)]
		state.Logger().Warn("cheat mode, overriding community cards", zap.Object("cheatData", cheatData))
	}

	state.GameController().RunTimer(duration, func() {
		state.GameController().GoNextState(GoEvaluateWinnerState, true)
	})
	return nil
}

func (state *DealRemainCommunityState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
		Model: &txpokergrpc.Model{
			Table: state.table.ToProto(),
		},
	})
	return nil
}

func (state *DealRemainCommunityState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context:   &txpokergrpc.GameState_DealRemainCommunityStateContext{DealRemainCommunityStateContext: &txpokergrpc.DealRemainCommunityStateContext{}},
	}
}
