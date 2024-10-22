package state

import (
	"context"
	"fmt"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/service"
	"card-game-server-prototype/pkg/game/txpoker/type/action"
	"card-game-server-prototype/pkg/game/txpoker/type/stage"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"
	"sort"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoEvaluateActionState = &core.StateTrigger{
	Name:      "GoEvaluateActionState",
	ArgsTypes: []reflect.Type{},
}

type EvaluateActionState struct {
	core.State

	actionHintGroup   *model2.ActionHintGroup
	seatStatusGroup   *model2.SeatStatusGroup
	playerGroup       *model2.PlayerGroup
	table             *model2.Table
	actionHintService *service.ActionHintService

	curActionPlayer *model2.Player
}

func ProvideEvaluateActionState(
	stateFactory *core.StateFactory,
	actionHintGroup *model2.ActionHintGroup,
	seatStatusGroup *model2.SeatStatusGroup,
	playerGroup *model2.PlayerGroup,
	table *model2.Table,
	actionHintService *service.ActionHintService,
) *EvaluateActionState {
	return &EvaluateActionState{
		State: stateFactory.Create("EvaluateActionState"),

		actionHintGroup:   actionHintGroup,
		seatStatusGroup:   seatStatusGroup,
		playerGroup:       playerGroup,
		table:             table,
		actionHintService: actionHintService,

		curActionPlayer: nil,
	}
}

func (state *EvaluateActionState) Run(ctx context.Context, args ...any) error {
	state.actionHintService.Eval()
	state.Logger().Debug(
		"action hint evaluated",
		zap.Object("actionHints", state.actionHintGroup),
	)

	// Bet stage is completed when:
	// 1. Only one player is not fold.
	// 2. All players have acted and no available action.
	if state.actionHintService.OnlyOneNotFold() ||
		state.actionHintService.AllHaveActed() {
		state.curActionPlayer = nil
		state.Logger().Debug("betStage completed",
			zap.Object("actionHints", state.actionHintGroup),
		)

		state.GameController().RunTimer(1*time.Second, func() {
			state.GameController().GoNextState(GoCollectChipState)
		})

		return nil
	}

	nextActionUid, err := state.evalNextActionUid()
	if err != nil {
		state.Logger().Error(
			"cannot find next action uid",
			zap.Error(err),
			zap.Object("players", state.playerGroup),
			zap.Object("actionHints", state.actionHintGroup),
		)
		state.GameController().GoErrorState()
		return nil
	}
	state.curActionPlayer = state.playerGroup.Data[nextActionUid]
	state.Logger().Debug(
		"betStage-nextPlayer",
		zap.Object("actionHints", state.actionHintGroup),
		zap.String("nextUid", string(nextActionUid)),
	)

	state.GameController().GoNextState(GoWaitActionState, nextActionUid)
	return nil
}

func (state *EvaluateActionState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
		Model: &txpokergrpc.Model{
			ActionHintGroup: state.actionHintGroup.ToProto(),
		},
	})
	return nil
}

func (state *EvaluateActionState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context:   &txpokergrpc.GameState_EvaluateActionStateContext{EvaluateActionStateContext: &txpokergrpc.EvaluateActionStateContext{}},
	}
}

func (state *EvaluateActionState) evalNextActionUid() (core.Uid, error) {
	if state.actionHintService.AllHaveActed() {
		return core.Uid(""), fmt.Errorf("bet has completed, no next action player")
	}

	playerSeatIdMap := lo.Invert(lo.MapValues(state.playerGroup.Data, func(p *model2.Player, _ core.Uid) int { return p.SeatId }))
	playerSeatIds := lo.Keys(playerSeatIdMap)
	sort.Ints(playerSeatIds)

	var (
		nextPlayer *model2.Player
		ok         bool
	)

	if state.curActionPlayer == nil {
		firstActionRole, err := stage.FirstActionRole(len(state.playerGroup.Data), state.table.BetStageFSM.MustState().(stage.Stage))
		if err != nil {
			return core.Uid(""), fmt.Errorf("cannot find first action role: %w", err)
		}

		nextPlayer, ok = lo.Find(lo.Values(state.playerGroup.Data), func(p *model2.Player) bool {
			return p.Role == firstActionRole
		})

		if !ok {
			return core.Uid(""), fmt.Errorf("next action player not found")
		}
	} else {
		curPlayerIdx := lo.IndexOf(playerSeatIds, state.curActionPlayer.SeatId)
		nextPlayerSeatId := playerSeatIds[(curPlayerIdx+1)%len(playerSeatIds)]
		nextPlayer = state.playerGroup.Data[playerSeatIdMap[nextPlayerSeatId]]
	}

	// Skip player who has actioned already. Note that BB, SB has the
	// option to action in preflop bet stage. They're considered live.
	idx := lo.IndexOf(playerSeatIds, nextPlayer.SeatId)
	for i := 0; i < len(playerSeatIds)-1; i++ {
		seatId := playerSeatIds[(idx+i)%len(playerSeatIds)]
		uid := playerSeatIdMap[seatId]
		if lo.Contains([]action.ActionType{action.Undefined, action.SB, action.BB}, state.actionHintGroup.Hints[uid].Action) {
			return uid, nil
		}
	}

	return core.Uid(""), fmt.Errorf("no other player has undefined bet action")
}
