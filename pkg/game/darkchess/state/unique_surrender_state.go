package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoSurrenderState = core.NewStateTrigger("GoSurrenderState")

type SurrenderState struct {
	core.State
	playerGroup     *model2.PlayerGroup
	actionHintGroup *model2.ActionHintGroup
	board           *model2.Board

	surrender core.Uid
}

func ProvideSurrenderState(
	stateFactory *core.StateFactory,
	playerGroup *model2.PlayerGroup,
	actionHintGroup *model2.ActionHintGroup,
	board *model2.Board,
) *SurrenderState {
	return &SurrenderState{
		State:           stateFactory.Create("SurrenderState"),
		playerGroup:     playerGroup,
		actionHintGroup: actionHintGroup,
		board:           board,
	}
}

func (state *SurrenderState) Run(context.Context, ...any) error {
	var duration = constant.SurrenderPeriod

	surrender, ok := lo.FindKeyBy(state.actionHintGroup.Data, func(uid core.Uid, hint *model2.ActionHint) bool {
		s, err := lo.Last(hint.SurrenderTurns)
		if err == nil {
			return s.Turn == state.board.TurnCount
		}

		return false
	})

	if !ok {
		state.Logger().Error("surrender not found",
			zap.Object("actionHint", state.actionHintGroup))
		state.GameController().GoErrorState()
		return status.Errorf(codes.Internal, "surrender not found")
	}

	state.surrender = surrender
	winner := lo.Without(lo.Keys(state.playerGroup.Data), surrender)[0]
	state.playerGroup.Data[winner].IsWinner = true

	state.Logger().Debug("player surrendered",
		zap.String("surrender", surrender.String()),
		zap.String("winner", winner.String()),
	)

	state.GameController().RunTimer(duration, func() {
		state.GameController().GoNextState(GoShowRoundResultState, false, winner)
	})
	return nil
}

func (state *SurrenderState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	return nil
}

func (state *SurrenderState) ToProto(_ core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_SurrenderStateContext{
			SurrenderStateContext: &gamegrpc.SurrenderStateContext{
				SurrenderUid: state.surrender.String(),
			},
		},
	}
}

func (state *SurrenderState) BeforeLeave(uid core.Uid) error {
	state.Logger().Warn("forbid to leave during game", zap.String("uid", uid.String()))
	return status.Errorf(codes.FailedPrecondition, "forbid to leave during game")
}
