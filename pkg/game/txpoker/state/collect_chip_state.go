package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/service"
	"card-game-server-prototype/pkg/game/txpoker/type/action"
	pot2 "card-game-server-prototype/pkg/game/txpoker/type/pot"
	"card-game-server-prototype/pkg/game/txpoker/type/seatstatus"
	stage2 "card-game-server-prototype/pkg/game/txpoker/type/stage"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoCollectChipState = &core.StateTrigger{
	Name:      "GoCollectChipState",
	ArgsTypes: []reflect.Type{},
}

type CollectChipState struct {
	core.State

	seatStatusGroup   *model2.SeatStatusGroup
	chipCacheGroup    *model2.ChipCacheGroup
	table             *model2.Table
	actionHintGroup   *model2.ActionHintGroup
	replay            *model2.Replay
	actionHintService *service.ActionHintService

	collectedChips map[core.Uid]int
}

func ProvideCollectChipState(
	stateFactory *core.StateFactory,
	seatStatusGroup *model2.SeatStatusGroup,
	chipCacheGroup *model2.ChipCacheGroup,
	table *model2.Table,
	actionHintGroup *model2.ActionHintGroup,
	replay *model2.Replay,
	actionHintService *service.ActionHintService,
) *CollectChipState {
	return &CollectChipState{
		State: stateFactory.Create("CollectChipState"),

		seatStatusGroup:   seatStatusGroup,
		chipCacheGroup:    chipCacheGroup,
		table:             table,
		actionHintGroup:   actionHintGroup,
		replay:            replay,
		actionHintService: actionHintService,

		collectedChips: map[core.Uid]int{},
	}
}

func (state *CollectChipState) Run(ctx context.Context, args ...any) error {
	state.collectedChips = map[core.Uid]int{}

	state.collect()
	state.returnOverBet()

	potChipSum := lo.Sum(lo.Map(state.table.Pots, func(p *pot2.Pot, _ int) int { return lo.Sum(lo.Values(p.Chips)) }))
	curStage := state.table.BetStageFSM.MustState().(stage2.Stage)
	state.replay.StagePotChip[curStage] = potChipSum

	state.actionHintService.ClearAction()

	if len(state.collectedChips) <= 0 {
		state.GameController().RunTimer(250*time.Millisecond, state.endCollectChip)
	} else {
		state.GameController().RunTimer(1250*time.Millisecond, state.endCollectChip)
	}
	return nil

}

func (state *CollectChipState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
		Model: &txpokergrpc.Model{
			ActionHintGroup: state.actionHintGroup.ToProto(),
			ChipCacheGroup:  state.chipCacheGroup.ToProto(),
		},
	})
	return nil
}

func (state *CollectChipState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &txpokergrpc.GameState_CollectChipStateContext{CollectChipStateContext: &txpokergrpc.CollectChipStateContext{
			CollectedChips: lo.MapEntries(state.collectedChips, func(k core.Uid, v int) (string, int32) { return k.String(), int32(v) }),
			Pots:           state.table.Pots.ToProto(),
		}},
	}
}

func (state *CollectChipState) Cleanup(ctx context.Context, args ...any) error {
	// Frontend need a little delay to render collect chip animation
	// and then update pot chip in the end of state.
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		Model: &txpokergrpc.Model{
			Table: state.table.ToProto(),
		},
	})
	return nil
}

func (state *CollectChipState) collect() {
	if state.actionHintService.OnlyOneNotFold() {
		// Edge case: only one player not fold.
		// In this case, the minCommitChip is the min bet chip except player not fold.

		// -- Case 1: all-in
		// A $10000 Allin <- winner
		// B $50 Fold
		// C $25 Fold
		// commitChip = Max($50, $25) = $50
		// Collect chip = $50(B) + $25(C) + $50(A) = $125
		// Return chip to A $9950($10000 - $50)

		// -- Case 2: only BB
		// A $20 (Bet BB) <- winner
		// B $10 (Bet 0.5 BB) Fold
		// C $0 Fold
		// commitChip = Max($10, $0) = $10
		// Collect chip = $10(B) + $0(C) + $10(A) = $20
		// Return chip to A $10($20 - $10)

		remainHints := lo.PickBy(state.actionHintGroup.Hints, func(k core.Uid, v *model2.ActionHint) bool {
			return v.BetChip > 0
		})

		foldHints := lo.PickBy(remainHints, func(k core.Uid, v *model2.ActionHint) bool {
			return v.Action == action.Fold
		})

		maxFoldChip := lo.Max(lo.Map(lo.Values(foldHints), func(h *model2.ActionHint, _ int) int { return h.BetChip }))

		if len(state.table.Pots) <= 0 {
			state.table.Pots = append(state.table.Pots, &pot2.Pot{
				Chips:   map[core.Uid]int{},
				Winners: map[core.Uid]*pot2.Winner{},
			})
		}

		curPot, _ := lo.Last(state.table.Pots)
		for _, remainHint := range remainHints {
			commitChip := lo.Min([]int{maxFoldChip, remainHint.BetChip})

			curPot.Chips[remainHint.Uid] += commitChip
			state.collectedChips[remainHint.Uid] += commitChip
			remainHint.BetChip -= commitChip
		}

		state.Logger().Debug("chip collected(others fold)",
			zap.Object("table", state.table),
			zap.Object("collected", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
				for uid, chip := range state.collectedChips {
					enc.AddInt(uid.String(), chip)
				}
				return nil
			})),
		)

		return
	}

	// https://www.quora.com/What-are-some-algorithms-for-handling-side-pots-in-Texas-Hold-Em
	// https://stackoverflow.com/questions/62672186/distributing-side-pots-in-poker

	// Normally the action (bets, raises, and calls) for each player
	// is kept separate until the betting round concludes. Let  x be
	// the least amount of action in front of a player who has not
	// folded. Pull x from each action pile into the current pot (or
	// the entire action from a folded player with a lesser amount of
	// action). In the common case that no players are all-in, all
	// unfolded players will have the same amount, x , in front of
	// them, so the dealer just sweeps it all into the main pot. If
	// any action remains, repeat the process with the remaining money
	// and the new side-pot as the current pot.
	for {
		remainHints := lo.PickBy(state.actionHintGroup.Hints, func(k core.Uid, v *model2.ActionHint) bool {
			return v.BetChip > 0
		})

		// Only 1 player left, his remain bet chip cannot be
		// collected into pot. These over bet chip is returned to
		// player before the end of game.
		//
		// EX:
		// x 100 all in
		// y 10 all in
		// z 15 all in
		// k 10 fold
		// 2 pot will be created, 1st pot is 10 * 4, 2nd pot is 5 * 2.
		// x will have remain 85 bet chip on the table that will not go into pot.
		if len(remainHints) <= 1 {
			break
		}

		// err != nil
		//   -> Create first pot.
		// (err == nil && curPot.HasAllIn)
		//   -> Create side pot when cur pot has all-in player's chips.
		if curPot, err := lo.Last(state.table.Pots); err != nil || (err == nil && curPot.HasAllIn) {
			state.table.Pots = append(state.table.Pots, &pot2.Pot{
				Chips:   map[core.Uid]int{},
				Winners: map[core.Uid]*pot2.Winner{},
			})
		}

		remainLiveHints := lo.PickBy(remainHints, func(k core.Uid, v *model2.ActionHint) bool {
			return v.Action != action.Fold
		})

		minCommitChip := lo.Min(lo.Map(lo.Values(remainLiveHints), func(h *model2.ActionHint, _ int) int { return h.BetChip }))

		curPot, _ := lo.Last(state.table.Pots)
		curPot.HasAllIn = lo.ContainsBy(lo.Values(remainLiveHints), func(h *model2.ActionHint) bool {
			return h.Action == action.AllIn ||
				(h.Action == action.BB && h.IsBBAllIn)
		})

		for _, hint := range remainHints {
			// In case of fold, some may have less bet chip than minCommitChip.
			commitChip := lo.Min([]int{minCommitChip, hint.BetChip})

			curPot.Chips[hint.Uid] += commitChip
			state.collectedChips[hint.Uid] += commitChip
			hint.BetChip -= commitChip
		}
	}

	state.Logger().Info("chip collected",
		zap.Object("table", state.table),
		zap.Object("collected", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			for uid, chip := range state.collectedChips {
				enc.AddInt(uid.String(), chip)
			}
			return nil
		})),
	)
}

// Return the over bet chip that is not collected into the pot.
// EX:
// x 100 all in
// y 10 all in
// z 15 all in
// k 10 fold
// 2 pot will be created, 1st pot is 10 * 4, 2nd pot is 5 * 2.
// x will have remain 85 bet chip on the table that will not go into pot.
// These over bet chip should return to the player before declare winner.
func (state *CollectChipState) returnOverBet() {
	for uid, actionHint := range state.actionHintGroup.Hints {
		if actionHint.BetChip > 0 {
			if seatStatus, ok := state.seatStatusGroup.Status[uid]; ok {
				seatStatusState := seatStatus.FSM.MustState().(seatstatus.SeatStatusState)
				if seatStatusState == seatstatus.PlayingState || seatStatusState == seatstatus.SittingOutState {
					seatStatus.Chip += actionHint.BetChip
				}
			} else {
				state.chipCacheGroup.CashOutChips[uid] += actionHint.BetChip
			}

			state.chipCacheGroup.SeatStatusChips[uid] += actionHint.BetChip
			state.Logger().Debug("return BetChip",
				zap.String("uid", uid.String()),
				zap.Int("returnChip", actionHint.BetChip),
			)

			actionHint.BetChip = 0
		}
	}
}

func (state *CollectChipState) endCollectChip() {
	// The only one player that didn't fold wins the pot without showdown.
	if state.actionHintService.OnlyOneNotFold() {
		state.Logger().Debug(
			"others fold,eval winner",
			zap.String("remain", lo.Keys(lo.PickBy(state.actionHintGroup.Hints, func(_ core.Uid, actionHint *model2.ActionHint) bool {
				return actionHint.Action != action.Fold
			}))[0].String()),
			zap.Object("table", state.table),
			zap.Object("actionHints", state.actionHintGroup),
		)
		state.GameController().GoNextState(GoEvaluateWinnerState, false)
		return
	}

	if err := state.table.BetStageFSM.Fire(stage2.NextStageTrigger); err != nil {
		state.Logger().Error(
			"failed to fire next stage trigger",
			zap.Error(err),
			zap.Object("table", state.table),
		)
		state.GameController().GoErrorState()
		return
	}

	// All bet stage completed, go showdown.
	if state.table.BetStageFSM.MustState().(stage2.Stage) == stage2.ShowdownStage {
		state.Logger().Debug(
			"no more betStage, showdown",
			zap.Object("table", state.table),
			zap.Object("actionHints", state.actionHintGroup),
		)
		state.GameController().GoNextState(GoShowdownState)
		return
	}

	// If only one player is not fold/all-in, go showdown.
	if state.actionHintService.OnlyOneIsLive() {
		state.Logger().Debug("only one is live, go declare showdown",
			zap.Object("table", state.table),
			zap.Object("actionHints", state.actionHintGroup),
		)
		state.GameController().GoNextState(GoDeclareShowdownState)
		return
	}

	state.Logger().Debug("next betStage",
		zap.Object("table", state.table),
		zap.Object("actionHints", state.actionHintGroup),
	)
	state.GameController().GoNextState(GoDealCommunityState)
}
