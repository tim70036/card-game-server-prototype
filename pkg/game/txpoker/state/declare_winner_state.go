package state

import (
	"context"
	commonapi "card-game-server-prototype/pkg/common/api"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/common/type/gamemode"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/type/fold"
	"card-game-server-prototype/pkg/game/txpoker/type/hand"
	"card-game-server-prototype/pkg/game/txpoker/type/seatstatus"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoDeclareWinnerState = &core.StateTrigger{
	Name:      "GoDeclareWinnerState",
	ArgsTypes: []reflect.Type{reflect.TypeOf(0)},
}

type DeclareWinnerState struct {
	core.State

	roomInfo        *commonmodel.RoomInfo
	seatStatusGroup *model2.SeatStatusGroup
	playerGroup     *model2.PlayerGroup
	chipCacheGroup  *model2.ChipCacheGroup
	table           *model2.Table
	userAPI         commonapi.UserAPI

	potIdx                   int
	cancelTimer              context.CancelFunc
	startTime                time.Time
	hasSomeOneShowFoldBefore bool
	isLastPot                bool
	defaultDuration          time.Duration

	endOnce *sync.Once
}

func ProvideDeclareWinnerState(
	stateFactory *core.StateFactory,
	roomInfo *commonmodel.RoomInfo,
	seatStatusGroup *model2.SeatStatusGroup,
	playerGroup *model2.PlayerGroup,
	chipCacheGroup *model2.ChipCacheGroup,
	table *model2.Table,
	userAPI commonapi.UserAPI,
) *DeclareWinnerState {
	return &DeclareWinnerState{
		State: stateFactory.Create("DeclareWinnerState"),

		roomInfo:        roomInfo,
		seatStatusGroup: seatStatusGroup,
		playerGroup:     playerGroup,
		chipCacheGroup:  chipCacheGroup,
		table:           table,
		userAPI:         userAPI,

		defaultDuration: 2 * time.Second,
	}
}

func (state *DeclareWinnerState) Run(ctx context.Context, args ...any) error {
	state.cancelTimer = func() {}
	state.endOnce = &sync.Once{}
	state.potIdx = args[0].(int)
	if state.potIdx < 0 || state.potIdx >= len(state.table.Pots) {
		state.Logger().Error("potIdx out of range", zap.Int("potIdx", state.potIdx), zap.Int("len(state.table.Pots)", len(state.table.Pots)))
		state.GameController().GoErrorState()
		return nil
	}

	state.isLastPot = state.potIdx+1 >= len(state.table.Pots)
	state.hasSomeOneShowFoldBefore = lo.ContainsBy(
		lo.Values(state.playerGroup.Data),
		func(p *model2.Player) bool { return p.HasShowFoldSet() },
	)

	// Will record chip for winners that standup (could sit down in
	// another seat, and become reserving state). These will be cash
	// out once declare winner is done. For sitting out winners, still
	// send cheap back to the seat status.
	pot := state.table.Pots[state.potIdx]
	state.Logger().Debug(
		"declare pot winner",
		zap.Object("pot", pot),
		zap.Bool("isLastPot", state.isLastPot),
		zap.Bool("hasSomeOneShowFoldBefore", state.hasSomeOneShowFoldBefore),
	)
	for uid, winner := range pot.Winners {
		afterWaterWinChip := winner.Chip - winner.Water - winner.JackpotWater
		state.chipCacheGroup.SeatStatusChips[uid] += afterWaterWinChip

		seatStatus, ok := state.seatStatusGroup.Status[uid]
		if !ok {
			state.chipCacheGroup.CashOutChips[uid] += afterWaterWinChip
			continue
		}

		seatStatusState := seatStatus.FSM.MustState().(seatstatus.SeatStatusState)
		if seatStatusState == seatstatus.PlayingState || seatStatusState == seatstatus.SittingOutState {
			seatStatus.Chip += afterWaterWinChip
		} else {
			state.chipCacheGroup.CashOutChips[uid] += afterWaterWinChip
		}
	}

	// If last pot and someone show fold before, then delay 2 more seconds.
	duration := lo.Ternary(
		state.isLastPot && state.hasSomeOneShowFoldBefore,
		state.defaultDuration+2*time.Second,
		state.defaultDuration,
	)

	state.startTime = time.Now()
	state.cancelTimer = state.GameController().RunTimer(duration, state.endDeclareWinner)
	return nil
}

func (state *DeclareWinnerState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
	})
	return nil
}

func (state *DeclareWinnerState) Cleanup(ctx context.Context, args ...any) error {
	// Frontend cannot update chip when winner declared, they need a little delay for chip pushing.
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		Model: &txpokergrpc.Model{
			ChipCacheGroup: state.chipCacheGroup.ToProto(),
		},
	})
	return nil
}

func (state *DeclareWinnerState) ToProto(uid core.Uid) proto.Message {
	winnerHands := map[string]*txpokergrpc.PokerHand{}
	for uid := range state.table.Pots[state.potIdx].Winners {
		if state.playerGroup.Data[uid].Hand != nil {
			winnerHands[uid.String()] = hand.ToProto(state.playerGroup.Data[uid].Hand)
		}
	}

	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &txpokergrpc.GameState_DeclareWinnerStateContext{DeclareWinnerStateContext: &txpokergrpc.DeclareWinnerStateContext{
			Pot:         state.table.Pots[state.potIdx].ToProto(),
			PotIndex:    int32(state.potIdx),
			WinnerHands: winnerHands,
		}},
	}
}

func (state *DeclareWinnerState) AcceptRequestTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(&txpokergrpc.ShowFoldRequest{}),
	}
}

func (state *DeclareWinnerState) HandleRequest(req *core.Request) error {
	player, ok := state.playerGroup.Data[req.Uid]
	if !ok {
		state.Logger().Warn("player not found", zap.Object("req", req))
		return nil
	}

	switch msg := req.Msg.(type) {
	case *txpokergrpc.ShowFoldRequest:
		state.Logger().Debug("handle ShowFoldRequest", zap.Object("req", req))

		if _, hasShowdown := state.table.ShowdownPocketCards[req.Uid]; hasShowdown {
			return status.Errorf(codes.InvalidArgument, "player already showdown")
		}

		if player.HasShowFoldSet() {
			return status.Errorf(codes.InvalidArgument, "player already show fold")
		}

		if msg.ShowFoldType == 0 {
			return nil
		}

		switch msg.ShowFoldType {
		case int32(fold.ShowRight):
			player.ShowRightFold()
			state.Logger().Debug("show fold", zap.String("type", "right"))
		case int32(fold.ShowLeft):
			player.ShowLeftFold()
			state.Logger().Debug("show fold", zap.String("type", "left"))
		case int32(fold.ShowBoth):
			player.ShowBothFold()
			state.Logger().Debug("show fold", zap.String("type", "both"))
		}

		// last pot winner show fold, reset state timer if no one showFold before.
		// new delay = passing time + 2sec
		if state.isLastPot && !state.hasSomeOneShowFoldBefore {
			state.hasSomeOneShowFoldBefore = true
			passedDuration := lo.Ternary(
				state.defaultDuration-time.Now().Sub(state.startTime) > 0,
				state.defaultDuration-time.Now().Sub(state.startTime),
				0,
			)
			newDuration := passedDuration + (2 * time.Second)
			state.Logger().Debug("last pot request show fold, reset timer",
				zap.Int64("remain_milliseconds", newDuration.Milliseconds()),
			)

			state.cancelTimer()
			state.cancelTimer = state.GameController().RunTimer(newDuration, state.endDeclareWinner)
		}

		playerGroupProto := &txpokergrpc.PlayerGroup{Players: make(map[string]*txpokergrpc.Player)}
		for uid, player := range state.playerGroup.Data {
			playerGroupProto.Players[uid.String()] = player.ToProto()
			_, playerGroupProto.Players[uid.String()].HasShowdown = state.table.ShowdownPocketCards[uid]
		}

		state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
			Model: &txpokergrpc.Model{
				PlayerGroup: playerGroupProto,
			},
		})

	default:
		state.Logger().Warn("not supported request", zap.Object("req", req))
		return status.Errorf(codes.Unimplemented, "not supported request")

	}
	return nil
}

func (state *DeclareWinnerState) cashoutWinChip() {
	if len(state.chipCacheGroup.CashOutChips) <= 0 {
		state.Logger().Debug("no winner need to cashout")
		return
	}

	for uid, chip := range state.chipCacheGroup.CashOutChips {
		go func(uid core.Uid, chip int) {
			if err := state.userAPI.ExchangeChip(uid, state.roomInfo.GameType, chip); err != nil {
				state.Logger().Error("failed to cashout chip for left winners", zap.Error(err), zap.String("uid", uid.String()), zap.Int("chip", chip))
			}
		}(uid, chip)
	}

	state.Logger().Debug("cashing out winner chip",
		zap.Object("cashOutChips", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			for uid, chip := range state.chipCacheGroup.CashOutChips {
				enc.AddInt(uid.String(), chip)
			}
			return nil
		})),
	)
}

func (state *DeclareWinnerState) endDeclareWinner() {
	state.endOnce.Do(func() {
		state.Logger().Debug("end declare winner")
		state.cancelTimer()
		if state.isLastPot {
			state.cashoutWinChip()

			if state.roomInfo.GameMode == gamemode.Club {
				state.GameController().GoNextState(GoEndRoundState)
			} else {
				state.GameController().GoNextState(GoJackpotState)
			}
		} else {
			state.GameController().GoNextState(GoDeclareWinnerState, state.potIdx+1)
		}
	})
}
