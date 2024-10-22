package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/actor"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/game/darkchess/service"
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"sync"
)

var GoPickFirstState = core.NewStateTrigger("GoPickFirstState")

type PickFirstState struct {
	core.State
	actorGroup      *actor.Group
	playerGroup     *model2.PlayerGroup
	boardService    *service.BoardService
	actionHintGroup *model2.ActionHintGroup
	pickBoard       *model2.PickBoard

	cancelTimer context.CancelFunc
	endOnce     *sync.Once
}

func ProvidePickFirstState(
	stateFactory *core.StateFactory,

	actorGroup *actor.Group,
	playerGroup *model2.PlayerGroup,
	boardService *service.BoardService,
	pickBoard *model2.PickBoard,
	actionHints *model2.ActionHintGroup,
) *PickFirstState {
	return &PickFirstState{
		State: stateFactory.Create("PickFirstState"),

		actorGroup:      actorGroup,
		playerGroup:     playerGroup,
		boardService:    boardService,
		pickBoard:       pickBoard,
		actionHintGroup: actionHints,
	}
}

func (state *PickFirstState) Run(context.Context, ...any) error {
	state.cancelTimer = func() {}
	state.endOnce = &sync.Once{}

	var duration = constant.PickFirstPeriod

	state.cancelTimer = state.GameController().RunTimer(duration, func() {
		state.goNextState()
	})

	for _, curActor := range state.actorGroup.Data {
		actorReqs, err := curActor.DecidePick()
		if err != nil {
			state.Logger().Error("actor cannot decide pick",
				zap.Error(err),
				zap.Object("actor", curActor),
			)
			state.GameController().GoErrorState()
			return nil
		}
		state.GameController().RunActorRequests(actorReqs)
	}

	return nil
}

func (state *PickFirstState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	return nil
}

func (state *PickFirstState) ToProto(_ core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_PickFirstStateContext{
			PickFirstStateContext: &gamegrpc.PickFirstStateContext{
				Duration: durationpb.New(constant.PickFirstPeriod),
			},
		},
	}
}

func (state *PickFirstState) BeforeLeave(uid core.Uid) error {
	state.Logger().Warn("forbid to leave during game", zap.String("uid", uid.String()))
	return status.Errorf(codes.FailedPrecondition, "forbid to leave during game")
}

func (state *PickFirstState) AcceptRequestTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(&gamegrpc.PickRequest{}),
	}
}

func (state *PickFirstState) HandleRequest(req *core.Request) error {
	if state.isPicked(req.Uid) {
		state.Logger().Warn("player already picked",
			zap.Object("req", req),
		)
		return status.Errorf(codes.FailedPrecondition, "player already picked")
	}

	switch req.Msg.(type) {

	case *gamegrpc.PickRequest:

		state.pick(req.Uid)

		state.Logger().Info("req:pick",
			zap.Object("req", req),
			util.DebugField(zap.Object("pickBoard", state.pickBoard)),
		)

	default:
		state.Logger().Warn("not supported request", zap.Object("req", req))
		return status.Errorf(codes.Unimplemented, "not supported request")
	}

	if state.isAllPicked() {
		state.Logger().Debug("player all picked")
		state.goNextState()
		return nil
	}

	return nil
}

func (state *PickFirstState) pick(uid core.Uid) {
	if state.isPicked(uid) {
		return
	}

	pick := state.pickBoard.Pieces[0]
	state.pickBoard.PickedPieces[uid] = pick

	// remove picked
	state.pickBoard.Pieces = state.pickBoard.Pieces[1:]

	if !state.isAllPicked() {
		state.MsgBus().Broadcast(core.ModelTopic, &gamegrpc.Model{
			PickResult: state.pickBoard.ToProto(),
		})
		return
	}

	var tmpPicked piece.Piece

	for u, picked := range state.pickBoard.PickedPieces {
		if state.pickBoard.FirstPlayer == "" {
			state.pickBoard.FirstPlayer = u
			tmpPicked = picked
			continue
		}

		if state.boardService.IsAllowToCapture(picked, tmpPicked) {
			state.pickBoard.FirstPlayer = u
		}
	}

	state.MsgBus().Broadcast(core.ModelTopic, &gamegrpc.Model{
		PickResult: state.pickBoard.ToProto(),
	})
}

func (state *PickFirstState) isPicked(uid core.Uid) bool {
	_, isPicked := state.pickBoard.PickedPieces[uid]
	return isPicked
}

func (state *PickFirstState) isAllPicked() bool {
	return len(state.pickBoard.PickedPieces) == constant.MaxUserCount
}

func (state *PickFirstState) goNextState() {
	state.endOnce.Do(func() {
		state.cancelTimer()

		for uid := range state.playerGroup.Data {
			if state.isPicked(uid) {
				continue
			}

			state.pick(uid)

			// auto pick
			state.Logger().Info("timeout, autoPick",
				zap.String("uid", uid.String()),
				zap.Object("pick_board", state.pickBoard),
			)
		}

		// Set first player
		firstPlayer := state.pickBoard.FirstPlayer
		for uid := range state.actionHintGroup.Data {
			if uid == firstPlayer {
				state.actionHintGroup.Data[uid].TurnCount = 1
				break
			}
		}

		state.Logger().Debug("end pick first",
			zap.Object("pick_board", state.pickBoard),
		)

		state.GameController().RunTimer(constant.AfterPickingPeriod, func() {
			state.GameController().GoNextState(GoStartTurnState)
		})
	})
}
