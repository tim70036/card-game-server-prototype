package state

import (
	"context"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/type/cheat"
	"card-game-server-prototype/pkg/game/txpoker/type/stage"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoDealCommunityState = &core.StateTrigger{
	Name:      "GoDealCommunityState",
	ArgsTypes: []reflect.Type{},
}

type DealCommunityState struct {
	core.State

	table    *model.Table
	testCFG  *config.TestConfig
	roomInfo *commonmodel.RoomInfo
}

func ProvideDealCommunityState(
	stateFactory *core.StateFactory,
	table *model.Table,
	testCFG *config.TestConfig,
	roomInfo *commonmodel.RoomInfo,
) *DealCommunityState {
	return &DealCommunityState{
		State: stateFactory.Create("DealCommunityState"),

		table:    table,
		testCFG:  testCFG,
		roomInfo: roomInfo,
	}
}

func (state *DealCommunityState) Run(ctx context.Context, args ...any) error {
	betStage := state.table.BetStageFSM.MustState().(stage.Stage)
	dealCnt := lo.Ternary(betStage == stage.FlopStage, 3, 1)
	duration := 500*time.Millisecond + time.Duration(lo.Ternary(betStage == stage.FlopStage, 1500, 900))*time.Millisecond

	dealCards := state.table.Deck[len(state.table.Deck)-dealCnt:]
	state.table.Deck = state.table.Deck[:len(state.table.Deck)-dealCnt]
	state.table.CommunityCards = append(state.table.CommunityCards, dealCards...)

	if state.testCFG.EnableCheatMode(string(state.roomInfo.GameType)) {
		cheatData, err := cheat.FromRawCheatData(*state.testCFG.CheatData)
		if err != nil {
			state.Logger().Warn("failed to parse raw cheat data", zap.String("rawCheatData", *state.testCFG.CheatData), zap.Error(err))
		}

		state.table.CommunityCards = cheatData.CommunityCards[:len(state.table.CommunityCards)]
		state.Logger().Warn("cheat mode, overriding community cards", zap.Object("cheatData", cheatData))
	}

	state.Logger().Debug("community cards dealt", zap.String("cards", state.table.CommunityCards.ToString()))

	state.GameController().RunTimer(duration, func() {
		state.GameController().GoNextState(GoEvaluateActionState)
	})
	return nil
}

func (state *DealCommunityState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
		Model: &txpokergrpc.Model{
			Table: state.table.ToProto(),
		},
	})
	return nil
}

func (state *DealCommunityState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context:   &txpokergrpc.GameState_DealCommunityStateContext{DealCommunityStateContext: &txpokergrpc.DealCommunityStateContext{}},
	}
}
