package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	action2 "card-game-server-prototype/pkg/game/txpoker/type/action"
	"card-game-server-prototype/pkg/game/txpoker/type/stage"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoShowdownState = &core.StateTrigger{
	Name:      "GoShowdownState",
	ArgsTypes: []reflect.Type{},
}

type ShowdownState struct {
	core.State

	actionHintGroup *model2.ActionHintGroup
	playerGroup     *model2.PlayerGroup
	table           *model2.Table
	replay          *model2.Replay
}

func ProvideShowdownState(
	stateFactory *core.StateFactory,
	actionHintGroup *model2.ActionHintGroup,
	playerGroup *model2.PlayerGroup,
	table *model2.Table,
	replay *model2.Replay,
) *ShowdownState {
	return &ShowdownState{
		State: stateFactory.Create("ShowdownState"),

		actionHintGroup: actionHintGroup,
		playerGroup:     playerGroup,
		table:           table,
		replay:          replay,
	}
}

func (state *ShowdownState) Run(ctx context.Context, args ...any) error {
	var showdownUids core.UidList = lo.Reject(lo.Keys(state.actionHintGroup.Hints), func(uid core.Uid, _ int) bool {
		return state.actionHintGroup.Hints[uid].Action == action2.Fold
	})

	curBetStage := state.table.BetStageFSM.MustState().(stage.Stage)
	for _, uid := range showdownUids {
		state.table.ShowdownPocketCards[uid] = state.playerGroup.Data[uid].PocketCards

		state.replay.ActionLog[curBetStage] = append(
			state.replay.ActionLog[curBetStage],
			&action2.ShowdownRecord{
				Uid:         uid,
				Role:        state.playerGroup.Data[uid].Role,
				PocketCards: state.playerGroup.Data[uid].PocketCards,
			},
		)
	}

	state.Logger().Debug("showdown!", zap.Array("showdownUids", showdownUids))
	state.GameController().RunTimer(1*time.Second, func() {
		state.GameController().GoNextState(GoDealRemainCommunityState)
	})
	return nil
}

func (state *ShowdownState) Publish(ctx context.Context, args ...any) error {
	playerGroupProto := &txpokergrpc.PlayerGroup{Players: make(map[string]*txpokergrpc.Player)}
	for uid, player := range state.playerGroup.Data {
		playerGroupProto.Players[uid.String()] = player.ToProto()
		_, playerGroupProto.Players[uid.String()].HasShowdown = state.table.ShowdownPocketCards[uid]
	}

	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
		Model: &txpokergrpc.Model{
			Table:       state.table.ToProto(),
			PlayerGroup: playerGroupProto,
		},
	})
	return nil
}

func (state *ShowdownState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context:   &txpokergrpc.GameState_ShowdownStateContext{ShowdownStateContext: &txpokergrpc.ShowdownStateContext{}},
	}
}
