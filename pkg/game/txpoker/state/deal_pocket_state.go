package state

import (
	"context"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/constant"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/type/cheat"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"go.uber.org/zap/zapcore"
	"math"
	"reflect"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoDealPocketState = &core.StateTrigger{
	Name:      "GoDealPocketState",
	ArgsTypes: []reflect.Type{},
}

type DealPocketState struct {
	core.State

	playerGroup *model2.PlayerGroup
	table       *model2.Table
	testCFG     *config.TestConfig
	roomInfo    *commonmodel.RoomInfo
}

func ProvideDealPocketState(
	stateFactory *core.StateFactory,
	playerGroup *model2.PlayerGroup,
	table *model2.Table,
	testCFG *config.TestConfig,
	roomInfo *commonmodel.RoomInfo,
) *DealPocketState {
	return &DealPocketState{
		State: stateFactory.Create("DealPocketState"),

		playerGroup: playerGroup,
		table:       table,
		testCFG:     testCFG,
		roomInfo:    roomInfo,
	}
}

func (state *DealPocketState) Run(ctx context.Context, args ...any) error {
	pocketCards := state.table.Deck[len(state.table.Deck)-len(state.playerGroup.Data)*constant.PocketCardsPerPlayer:]
	state.table.Deck = state.table.Deck[:len(state.table.Deck)-len(state.playerGroup.Data)*constant.PocketCardsPerPlayer]

	for _, player := range state.playerGroup.Data {
		player.PocketCards = pocketCards[:constant.PocketCardsPerPlayer].Clone()
		pocketCards = pocketCards[constant.PocketCardsPerPlayer:]
	}

	if state.testCFG.EnableCheatMode(string(state.roomInfo.GameType)) {
		cheatData, err := cheat.FromRawCheatData(*state.testCFG.CheatData)
		if err != nil {
			state.Logger().Warn("failed to parse raw cheat data", zap.String("rawCheatData", *state.testCFG.CheatData), zap.Error(err))
		}

		for uid, cheatCards := range cheatData.PlayerPocketCards {
			if _, ok := state.playerGroup.Data[uid]; !ok {
				state.Logger().Warn("cheat mode, overriding pocket cards, but uid not found in player group",
					zap.String("uid", string(uid)),
					zap.Object("players", state.playerGroup),
					zap.Object("cheatData", cheatData),
				)
				continue
			}

			state.playerGroup.Data[uid].PocketCards = cheatCards
		}

		state.Logger().Warn("cheat mode, overriding pocket cards", zap.Object("cheatData", cheatData))
	}

	state.Logger().Debug("pocket cards dealt",
		zap.Object("players", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			for uid, player := range state.playerGroup.Data {
				enc.AddString(uid.String(), player.PocketCards.ToString())
			}
			return nil
		})),
	)

	x := float64(len(state.playerGroup.Data)) / float64(constant.MaxSeatId+1)
	ceiling := 0.6 // ceiling + floor == real ceiling
	floor := 1.65
	duration := time.Duration((ceiling*(1-math.Cos(x*math.Pi/2))+floor)*1000) * time.Millisecond
	state.GameController().RunTimer(duration, func() {
		state.GameController().GoNextState(GoEvaluateActionState)
	})
	return nil
}

func (state *DealPocketState) Publish(ctx context.Context, args ...any) error {
	playerGroupProto := &txpokergrpc.PlayerGroup{Players: make(map[string]*txpokergrpc.Player)}
	for uid, player := range state.playerGroup.Data {
		playerGroupProto.Players[uid.String()] = player.ToProto()
		_, playerGroupProto.Players[uid.String()].HasShowdown = state.table.ShowdownPocketCards[uid]
	}

	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
		Model: &txpokergrpc.Model{
			PlayerGroup: playerGroupProto,
		},
	})
	return nil
}

func (state *DealPocketState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context:   &txpokergrpc.GameState_DealPocketStateContext{DealPocketStateContext: &txpokergrpc.DealPocketStateContext{}},
	}
}
