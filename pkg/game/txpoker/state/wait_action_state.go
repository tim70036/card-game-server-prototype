package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/actor"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	service2 "card-game-server-prototype/pkg/game/txpoker/service"
	"card-game-server-prototype/pkg/game/txpoker/type/action"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"
	"reflect"
	"sync"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoWaitActionState = &core.StateTrigger{
	Name:      "GoWaitActionState",
	ArgsTypes: []reflect.Type{core.UidType},
}

type WaitActionState struct {
	core.State

	seatStatusGroup *model2.SeatStatusGroup
	chipCacheGroup  *model2.ChipCacheGroup
	actionHintGroup *model2.ActionHintGroup
	actorGroup      *actor.ActorGroup
	playerGroup     *model2.PlayerGroup
	gameInfo        *model2.GameInfo

	seatStatusService service2.SeatStatusService
	actionHintService *service2.ActionHintService

	actionUid        core.Uid
	stateTriggerMap  map[action.ActionType]*core.StateTrigger
	requestActionMap map[reflect.Type]action.ActionType
	startWaitTime    time.Time
	duration         time.Duration

	cancelTimer context.CancelFunc
	endOnce     *sync.Once
}

func ProvideWaitActionState(
	stateFactory *core.StateFactory,
	seatStatusGroup *model2.SeatStatusGroup,
	chipCacheGroup *model2.ChipCacheGroup,
	actionHintGroup *model2.ActionHintGroup,
	actorGroup *actor.ActorGroup,
	playerGroup *model2.PlayerGroup,
	gameInfo *model2.GameInfo,

	seatStatusService service2.SeatStatusService,
	actionHintService *service2.ActionHintService,
) *WaitActionState {
	return &WaitActionState{
		State: stateFactory.Create("WaitActionState"),

		seatStatusGroup: seatStatusGroup,
		chipCacheGroup:  chipCacheGroup,
		actionHintGroup: actionHintGroup,
		actorGroup:      actorGroup,
		playerGroup:     playerGroup,
		gameInfo:        gameInfo,

		seatStatusService: seatStatusService,
		actionHintService: actionHintService,

		stateTriggerMap: map[action.ActionType]*core.StateTrigger{
			action.Fold:  GoFoldState,
			action.Check: GoCheckState,
			action.Bet:   GoBetState,
			action.Call:  GoCallState,
			action.Raise: GoRaiseState,
			action.AllIn: GoAllInState,
		},

		requestActionMap: map[reflect.Type]action.ActionType{
			reflect.TypeOf(&txpokergrpc.FoldRequest{}):  action.Fold,
			reflect.TypeOf(&txpokergrpc.CheckRequest{}): action.Check,
			reflect.TypeOf(&txpokergrpc.BetRequest{}):   action.Bet,
			reflect.TypeOf(&txpokergrpc.CallRequest{}):  action.Call,
			reflect.TypeOf(&txpokergrpc.RaiseRequest{}): action.Raise,
			reflect.TypeOf(&txpokergrpc.AllInRequest{}): action.AllIn,
		},
	}
}

func (state *WaitActionState) Run(ctx context.Context, args ...any) error {
	state.cancelTimer = func() {}
	state.endOnce = &sync.Once{}

	state.actionUid = args[0].(core.Uid)
	actionHint := state.actionHintGroup.Hints[state.actionUid]

	state.duration = actionHint.Duration
	if seatStatus, ok := state.seatStatusGroup.Status[state.actionUid]; ok {
		state.duration += seatStatus.ActionExtraDuration
	}

	// Edge case:
	// user request and timeout happen at the same time.
	// It comes to both action are executed.
	state.duration += time.Millisecond * 500

	state.Logger().Debug("start waiting", zap.Object("actionHint", actionHint))
	state.cancelTimer = state.GameController().RunTimer(state.duration, func() {
		act := action.Fold
		if lo.Contains(state.actionHintGroup.Hints[state.actionUid].AvailableActions, action.Check) {
			act = action.Check
		}

		if err := state.actionHintService.Pass(state.actionUid, act); err != nil {
			state.Logger().Error(
				"timeout auto fold/check, but pass failed",
				zap.Error(err),
				zap.String("actionUid", state.actionUid.String()),
				zap.String("action", act.String()),
				zap.Object("seatStatus", state.seatStatusGroup),
				zap.Object("actionHints", state.actionHintGroup),
			)
			state.GameController().GoErrorState()
			return
		}

		state.Logger().Info("timeout, auto act",
			zap.String("actionUid", state.actionUid.String()),
			zap.String("action", act.String()),
		)
		state.endWaitAction(act, false)
	})

	actorReqs, err := state.actorGroup.Data[state.actionUid].DecideAction()
	if err != nil {
		state.Logger().Error(
			"cannot decide action",
			zap.Error(err),
			zap.String("actionUid", state.actionUid.String()),
			zap.Object("actionHint", actionHint),
		)
		state.GameController().GoErrorState()
		return nil
	}

	state.GameController().RunActorRequests(actorReqs)

	state.startWaitTime = time.Now()
	return nil
}

func (state *WaitActionState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
	})
	return nil
}

func (state *WaitActionState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &txpokergrpc.GameState_WaitActionStateContext{WaitActionStateContext: &txpokergrpc.WaitActionStateContext{
			ActorUid: state.actionUid.String(),
			Duration: durationpb.New(state.duration),
		}},
	}
}

func (state *WaitActionState) AcceptRequestTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(&txpokergrpc.StandUpRequest{}),
		reflect.TypeOf(&txpokergrpc.SitOutRequest{}),
		reflect.TypeOf(&txpokergrpc.FoldRequest{}),
		reflect.TypeOf(&txpokergrpc.CheckRequest{}),
		reflect.TypeOf(&txpokergrpc.BetRequest{}),
		reflect.TypeOf(&txpokergrpc.CallRequest{}),
		reflect.TypeOf(&txpokergrpc.RaiseRequest{}),
		reflect.TypeOf(&txpokergrpc.AllInRequest{}),
	}
}

func (state *WaitActionState) HandleRequest(req *core.Request) error {
	if state.actionUid != req.Uid {
		switch req.Msg.(type) {
		case *txpokergrpc.StandUpRequest, *txpokergrpc.SitOutRequest:
			return nil
		default:
			state.Logger().Warn("cannot action",
				zap.String("actionUid", state.actionUid.String()),
				zap.Object("req", req),
			)
			return status.Errorf(codes.PermissionDenied, "cannot action")
		}
	}

	actionHint, ok := state.actionHintGroup.Hints[state.actionUid]
	if !ok {
		state.Logger().Warn("cannot find action hint",
			zap.Object("actionHints", state.actionHintGroup),
			zap.String("actionUid", state.actionUid.String()),
			zap.Object("req", req),
		)
		return status.Errorf(codes.NotFound, "cannot find action hint")
	}

	if !lo.Contains([]action.ActionType{action.Undefined, action.BB, action.SB}, actionHint.Action) {
		state.Logger().Warn("already made action",
			zap.Object("actionHints", state.actionHintGroup),
			zap.Object("req", req),
		)
		return status.Errorf(codes.AlreadyExists, "already made action")
	}

	if requestAction, ok := state.requestActionMap[reflect.TypeOf(req.Msg)]; ok {
		if requestAction != action.Fold && // Allow fold under any condition
			!lo.Contains(actionHint.AvailableActions, requestAction) {
			state.Logger().Warn("request action, but not an available action",
				zap.Object("actionHints", state.actionHintGroup),
				zap.Object("req", req),
			)
			return status.Errorf(codes.FailedPrecondition, "%v is not an available action", requestAction.String())
		}
	}

	switch msg := req.Msg.(type) {
	case *txpokergrpc.StandUpRequest:
		if err := state.actionHintService.Pass(state.actionUid, action.Fold); err != nil {
			state.Logger().Warn("standUp, auto fold failed",
				zap.Error(err),
				zap.Object("seatStatus", state.seatStatusGroup),
				zap.Object("actionHints", state.actionHintGroup),
				zap.Object("req", req),
			)
			return err
		}

		state.Logger().Info("req:StandUp, auto fold",
			zap.Object("actionHints", state.actionHintGroup),
			util.DebugField(zap.Object("req", req)),
		)

	case *txpokergrpc.SitOutRequest:
		if err := state.actionHintService.Pass(state.actionUid, action.Fold); err != nil {
			state.Logger().Warn("sitOut failed",
				zap.Error(err),
				zap.Object("seatStatus", state.seatStatusGroup),
				zap.Object("actionHints", state.actionHintGroup),
				zap.Object("req", req),
			)
			return err
		}

		state.Logger().Info("req:SitOut, auto fold",
			zap.Object("actionHints", state.actionHintGroup),
			util.DebugField(zap.Object("req", req)),
		)

	case *txpokergrpc.FoldRequest:
		if err := state.actionHintService.Pass(state.actionUid, action.Fold); err != nil {
			state.Logger().Warn("fold, but pass failed",
				zap.Error(err),
				zap.Object("seatStatus", state.seatStatusGroup),
				zap.Object("actionHints", state.actionHintGroup),
				zap.Object("req", req),
			)
			return err
		}

		state.Logger().Info("req:Fold",
			zap.Object("actionHints", state.actionHintGroup),
			util.DebugField(zap.Object("req", req)),
		)

	case *txpokergrpc.BetRequest:
		if err := state.actionHintService.OpenBet(state.actionUid, action.Bet, int(msg.Chip)); err != nil {
			state.Logger().Warn("bet, but open bet failed",
				zap.Error(err),
				zap.Object("seatStatus", state.seatStatusGroup),
				zap.Object("actionHints", state.actionHintGroup),
				zap.Object("req", req),
			)
			return err
		}

		state.Logger().Info("req:Bet",
			zap.Object("actionHints", state.actionHintGroup),
			zap.Object("req", req),
		)

	case *txpokergrpc.RaiseRequest:
		if err := state.actionHintService.OpenBet(state.actionUid, action.Raise, int(msg.Chip)); err != nil {
			state.Logger().Warn("raise, but open bet failed",
				zap.Error(err),
				zap.Object("seatStatus", state.seatStatusGroup),
				zap.Object("actionHints", state.actionHintGroup),
				zap.Object("req", req),
			)
			return err
		}

		state.Logger().Info("req:Raise",
			zap.Object("actionHints", state.actionHintGroup),
			zap.Object("req", req),
		)

	case *txpokergrpc.CheckRequest:
		if err := state.actionHintService.Pass(state.actionUid, action.Check); err != nil {
			state.Logger().Warn("check, but pass failed",
				zap.Error(err),
				zap.Object("seatStatus", state.seatStatusGroup),
				zap.Object("actionHints", state.actionHintGroup),
				zap.Object("req", req),
			)
			return err
		}

		state.Logger().Info("req:Check",
			zap.Object("actionHints", state.actionHintGroup),
			util.DebugField(zap.Object("req", req)),
		)

	case *txpokergrpc.CallRequest:
		if err := state.actionHintService.FollowBet(state.actionUid, action.Call); err != nil {
			state.Logger().Warn("call, but follow bet failed",
				zap.Error(err),
				zap.Object("seatStatus", state.seatStatusGroup),
				zap.Object("actionHints", state.actionHintGroup),
				zap.Object("req", req),
			)
			return err
		}

		state.Logger().Info("req:Call",
			zap.Object("actionHints", state.actionHintGroup),
			util.DebugField(zap.Object("req", req)),
		)

	case *txpokergrpc.AllInRequest:
		if err := state.actionHintService.AllIn(state.actionUid); err != nil {
			state.Logger().Warn("all in, but failed",
				zap.Error(err),
				zap.Object("seatStatus", state.seatStatusGroup),
				zap.Object("actionHints", state.actionHintGroup),
				zap.Object("req", req),
			)
			return err
		}

		state.Logger().Info("req:AllIn",
			zap.Object("actionHints", state.actionHintGroup),
			util.DebugField(zap.Object("req", req)),
		)

	default:
		state.Logger().Warn("not supported request", zap.Object("req", req))
		return status.Errorf(codes.Unimplemented, "not supported request")

	}

	state.GameController().RunTask(func() {
		state.Logger().Debug("player made action", zap.Object("actionHint", actionHint))
		state.endWaitAction(actionHint.Action, true)
	})
	return nil
}

func (state *WaitActionState) endWaitAction(playerAction action.ActionType, hasAction bool) {
	state.endOnce.Do(func() {
		state.Logger().Debug("ending wait action",
			zap.String("actionUid", state.actionUid.String()),
			zap.String("playerAction", playerAction.String()),
		)
		state.cancelTimer()
		state.updateActionExtraDuration()

		if !hasAction {
			state.seatStatusService.Idle(state.actionUid)
		} else {
			state.seatStatusService.Act(state.actionUid)
		}

		state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
			Model: &txpokergrpc.Model{
				ActionHintGroup: state.actionHintGroup.ToProto(),
				ChipCacheGroup:  state.chipCacheGroup.ToProto(),
			},
		})

		state.GameController().GoNextState(state.stateTriggerMap[playerAction], state.actionUid)
	})
}

func (state *WaitActionState) updateActionExtraDuration() {
	seatStatus, ok := state.seatStatusGroup.Status[state.actionUid]
	if !ok {
		return
	}

	actionHint := state.actionHintGroup.Hints[state.actionUid]
	actualDuration := time.Since(state.startWaitTime)
	excessDuration := lo.Ternary(actualDuration > actionHint.Duration, actualDuration-actionHint.Duration, 0)

	seatStatus.ActionExtraDuration = lo.Ternary(seatStatus.ActionExtraDuration > excessDuration, seatStatus.ActionExtraDuration-excessDuration, 0)
}
