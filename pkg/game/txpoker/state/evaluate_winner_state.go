package state

import (
	"context"
	"fmt"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/common/type/gamemode"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/service"
	action2 "card-game-server-prototype/pkg/game/txpoker/type/action"
	"card-game-server-prototype/pkg/game/txpoker/type/fold"
	hand2 "card-game-server-prototype/pkg/game/txpoker/type/hand"
	pot2 "card-game-server-prototype/pkg/game/txpoker/type/pot"
	"card-game-server-prototype/pkg/game/txpoker/type/stage"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"math"
	"reflect"
	"sort"
	"strconv"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoEvaluateWinnerState = &core.StateTrigger{
	Name:      "GoEvaluateWinnerState",
	ArgsTypes: []reflect.Type{reflect.TypeOf(false)},
}

type EvaluateWinnerState struct {
	core.State

	gameSetting     *model2.GameSetting
	playerGroup     *model2.PlayerGroup
	actionHintGroup *model2.ActionHintGroup
	table           *model2.Table
	replay          *model2.Replay
	roomInfo        *commonmodel.RoomInfo
	jackpotService  service.JackpotService

	hasShowdown bool
}

func ProvideEvaluateWinnerState(
	stateFactory *core.StateFactory,
	gameSetting *model2.GameSetting,
	playerGroup *model2.PlayerGroup,
	actionHintGroup *model2.ActionHintGroup,
	table *model2.Table,
	replay *model2.Replay,
	roomInfo *commonmodel.RoomInfo,
	jackpotService service.JackpotService,
) *EvaluateWinnerState {
	return &EvaluateWinnerState{
		State: stateFactory.Create("EvaluateWinnerState"),

		gameSetting:     gameSetting,
		playerGroup:     playerGroup,
		actionHintGroup: actionHintGroup,
		table:           table,
		replay:          replay,
		roomInfo:        roomInfo,
		jackpotService:  jackpotService,
	}
}

func (state *EvaluateWinnerState) Run(ctx context.Context, args ...any) error {
	state.hasShowdown = args[0].(bool)

	if err := state.evalPlayerHand(); err != nil {
		state.Logger().Error(
			"failed to eval player hand",
			zap.Error(err),
			zap.Object("actionHints", state.actionHintGroup),
			zap.Object("players", state.playerGroup),
			zap.Object("table", state.table),
			zap.Object("replay", state.replay),
		)
		state.GameController().GoErrorState()
		return nil
	}

	if err := state.evalPotWinners(); err != nil {
		state.Logger().Error(
			"failed to eval pot winner",
			zap.Error(err),
			zap.Object("actionHints", state.actionHintGroup),
			zap.Object("players", state.playerGroup),
			zap.Object("table", state.table),
			zap.Object("replay", state.replay),
		)
		state.GameController().GoErrorState()
		return nil
	}
	state.Logger().Info("evaluated pot winners", zap.Object("table", state.table))

	state.recordReplay()
	state.Logger().Debug("replay recorded", zap.Object("replay", state.replay))

	state.GameController().GoNextState(GoDeclareWinnerState, 0)
	return nil
}

func (state *EvaluateWinnerState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
		Model: &txpokergrpc.Model{
			Table:           state.table.ToProto(),
			ActionHintGroup: state.actionHintGroup.ToProto(),
		},
	})
	return nil
}

func (state *EvaluateWinnerState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context:   &txpokergrpc.GameState_EvaluateWinnerStateContext{EvaluateWinnerStateContext: &txpokergrpc.EvaluateWinnerStateContext{}},
	}
}

func (state *EvaluateWinnerState) evalPlayerHand() error {
	if state.hasShowdown {
		for uid, pocketCards := range state.table.ShowdownPocketCards {
			hand, err := hand2.New(pocketCards, state.table.CommunityCards)
			if err != nil {
				return fmt.Errorf("failed to create hand for uid %s: %w", uid.String(), err)
			}

			state.playerGroup.Data[uid].Hand = hand
		}

		state.Logger().Debug("player hands evaluated", zap.Object("players", state.playerGroup))
	}
	return nil
}

func (state *EvaluateWinnerState) evalPotWinners() error {
	defer state.evalWinnersWater()

	if !state.hasShowdown {
		hint, ok := lo.Find(lo.Values(state.actionHintGroup.Hints), func(hint *model2.ActionHint) bool {
			return hint.Action != action2.Fold
		})

		if !ok {
			return fmt.Errorf("failed to find winner uid which is the only one not fold")
		}

		winnerUid := hint.Uid
		for _, p := range state.table.Pots {
			p.Winners[winnerUid] = &pot2.Winner{
				Chip:         lo.Sum(lo.Values(p.Chips)),
				RawProfit:    lo.Sum(lo.Values(p.Chips)) - p.Chips[winnerUid],
				Water:        0,
				JackpotWater: 0,
			}
		}

		return nil
	}

	handUidMap := lo.Invert(lo.MapValues(
		lo.OmitBy(state.playerGroup.Data, func(k core.Uid, v *model2.Player) bool { return v.Hand == nil }),
		func(v *model2.Player, k core.Uid) hand2.Hand { return v.Hand },
	))

	var sortedHands hand2.HandList = lo.Keys(handUidMap)
	sort.Sort(sortedHands)

	rankMap := make(map[core.Uid]int)
	curRank := 0
	curHand := sortedHands[0]
	for _, hand := range sortedHands {
		uid := handUidMap[hand]
		if curHand.Equal(hand) {
			rankMap[uid] = curRank
		} else {
			curRank++
			rankMap[uid] = curRank
			curHand = hand
		}
	}

	state.Logger().Debug("hand rank evaluated", zap.Object("rank", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for idx, hand := range sortedHands {
			rank := rankMap[handUidMap[hand]]
			enc.AddObject(strconv.Itoa(idx), zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
				enc.AddString("uid", handUidMap[hand].String())
				enc.AddObject("hand", hand)
				enc.AddInt("rank", rank)
				return nil
			}))
		}
		return nil
	})))

	for _, p := range state.table.Pots {
		maxRankUid := lo.MaxBy(lo.Keys(p.Chips), func(uid core.Uid, maxUid core.Uid) bool {
			if _, ok := rankMap[uid]; !ok {
				return false
			}

			// Set maxUid to default, prevent nil panic.
			if _, ok := rankMap[maxUid]; !ok {
				return true
			}

			return rankMap[uid] > rankMap[maxUid]
		})

		winnerUids := lo.Filter(lo.Keys(p.Chips), func(uid core.Uid, _ int) bool {
			if _, ok := rankMap[uid]; !ok {
				return false
			}

			return rankMap[uid] == rankMap[maxRankUid]
		})

		totalWinChips := lo.Sum(lo.Values(p.Chips))
		chipPerWinner := totalWinChips / len(winnerUids)
		for _, uid := range winnerUids {
			p.Winners[uid] = &pot2.Winner{
				Chip:         chipPerWinner,
				RawProfit:    chipPerWinner - p.Chips[uid],
				Water:        0,
				JackpotWater: 0,
			}
		}

		//  一樣大的牌chop平分底池，如果有無法完全平分的部分，由較
		//  前面的位置(SB>)獲得多餘的部分
		if remainChip := totalWinChips % len(winnerUids); remainChip > 0 {
			maxPriorityUid := lo.MaxBy(winnerUids, func(uid core.Uid, maxUid core.Uid) bool {
				return state.playerGroup.Data[uid].Role < state.playerGroup.Data[maxUid].Role
			})
			p.Winners[maxPriorityUid].Chip += remainChip
			p.Winners[maxPriorityUid].RawProfit += remainChip
		}

	}
	return nil
}

func (state *EvaluateWinnerState) evalWinnersWater() {
	totalChips := lo.Reduce(state.table.Pots, func(totalChips map[core.Uid]int, pot *pot2.Pot, _ int) map[core.Uid]int {
		for uid, winner := range pot.Winners {
			totalChips[uid] += winner.Chip
		}
		return totalChips
	}, map[core.Uid]int{})

	totalRawProfits := lo.Reduce(state.table.Pots, func(totalRawProfits map[core.Uid]int, pot *pot2.Pot, _ int) map[core.Uid]int {
		for uid, winner := range pot.Winners {
			totalRawProfits[uid] += winner.RawProfit
		}
		return totalRawProfits
	}, map[core.Uid]int{})

	var hasMaxWaterLimit bool
	var maxProfitLimit int
	if state.roomInfo.GameMode == gamemode.Club && state.gameSetting.MaxWaterLimitBB > 0 {
		hasMaxWaterLimit = true
		maxProfitLimit = state.gameSetting.MaxWaterLimitBB * state.gameSetting.BigBlind
	}

	for uid, totalChip := range totalChips {
		jackpotWater := state.jackpotService.EvalWater(totalRawProfits[uid])
		water := int(math.Round(float64(totalChip-jackpotWater) * (float64(state.gameSetting.WaterPct) * 0.01)))

		if state.roomInfo.GameMode == gamemode.Club {
			if hasMaxWaterLimit && water > maxProfitLimit {
				water = maxProfitLimit
			}
		}

		// The win pots animation render pot 1 by 1. It's hard to
		// match which pot do the water & jackpot water belongs. So,
		// just gonna choose the first pot and put water in there.
		for _, pot := range state.table.Pots {
			if _, ok := pot.Winners[uid]; ok {
				pot.Winners[uid].Water = water
				pot.Winners[uid].JackpotWater = jackpotWater
				break
			}
		}
	}
}

func (state *EvaluateWinnerState) recordReplay() {
	curBetStage := state.table.BetStageFSM.MustState().(stage.Stage)
	for _, pot := range state.table.Pots {
		for uid, winner := range pot.Winners {
			pocketCardsMask := fold.ShowNone

			if state.hasShowdown {
				pocketCardsMask = fold.ShowBoth
			} else {
				if player, ok := state.playerGroup.Data[uid]; ok {
					pocketCardsMask = player.GetShowFold()
				}
			}

			state.replay.ActionLog[curBetStage] = append(
				state.replay.ActionLog[curBetStage],
				&action2.WinPotRecord{
					Uid:             uid,
					Role:            state.playerGroup.Data[uid].Role,
					Chip:            winner.Chip,
					PocketCardsMask: int(pocketCardsMask),
					PocketCards:     state.playerGroup.Data[uid].PocketCards,
					Hand:            state.playerGroup.Data[uid].Hand,
				},
			)
		}
	}

	playerTotalBetChip := lo.Reduce(state.table.Pots, func(totalBetChip map[core.Uid]int, pot *pot2.Pot, _ int) map[core.Uid]int {
		for uid, betChip := range pot.Chips {
			totalBetChip[uid] += betChip
		}
		return totalBetChip
	}, map[core.Uid]int{})

	playerTotalWinChips := lo.Reduce(state.table.Pots, func(totalWinChips map[core.Uid]int, pot *pot2.Pot, _ int) map[core.Uid]int {
		for uid, winner := range pot.Winners {
			totalWinChips[uid] += winner.Chip
		}
		return totalWinChips
	}, map[core.Uid]int{})

	for uid, playerRecord := range state.replay.PlayerRecords {
		playerRecord.PocketCards = state.playerGroup.Data[uid].PocketCards
		playerRecord.BeforeWinSeatStatusChip = playerRecord.InitSeatStatusChip - playerTotalBetChip[uid]
		playerRecord.WinChip = playerTotalWinChips[uid]
		playerRecord.IsWinner = playerTotalWinChips[uid] > 0
		_, playerRecord.HasShowdown = state.table.ShowdownPocketCards[uid]
	}

	state.replay.CommunityCards = state.table.CommunityCards
	state.replay.PotChipSum = lo.Sum(lo.Map(state.table.Pots, func(p *pot2.Pot, _ int) int { return lo.Sum(lo.Values(p.Chips)) }))

	// Special case: In ante stage there will be bb sb chip. Even
	// though they are not collected into pot, we still need to record
	// them in ante stage pot chip for replay rendering.
	state.replay.StagePotChip[stage.AnteStage] = lo.SumBy(state.replay.ActionLog[stage.AnteStage], func(r action2.ActionRecord) int {
		if r.GetType() == action2.BB {
			return r.(*action2.BBRecord).Chip
		}

		if r.GetType() == action2.SB {
			return r.(*action2.SBRecord).Chip
		}
		return 0
	})

	// Fill the empty stage pot chip for frontend rendering.
	// EX:
	// ante stage: 30
	// preflop stage: 100
	// flop stage: 200 (every one all in and proceed to showdown)
	// turn stage: 0
	// river stage: 0
	// showdown stage: 0
	// in this case, turn, river, showdown stage should have chip 200
	lastStageWithChip := stage.AnteStage
	for stage := stage.AnteStage; stage <= curBetStage; stage++ {
		if state.replay.StagePotChip[stage] > 0 {
			lastStageWithChip = stage
		} else {
			state.replay.StagePotChip[stage] = state.replay.StagePotChip[lastStageWithChip]
		}
	}
}
