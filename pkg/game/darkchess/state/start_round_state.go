package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/game/darkchess/service"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

var GoStartRoundState = core.NewStateTrigger("GoStartRoundState")

type StartRoundState struct {
	core.State

	gameRepoService service.GameRepoService
	gameInfo        *model2.GameInfo
	playerGroup     *model2.PlayerGroup

	delayLeaveDur time.Duration
}

func ProvideStartRoundState(
	stateFactory *core.StateFactory,

	gameRepoService service.GameRepoService,
	gameInfo *model2.GameInfo,
	playerGroup *model2.PlayerGroup,
) *StartRoundState {
	return &StartRoundState{
		State: stateFactory.Create("StartRoundState"),

		gameRepoService: gameRepoService,
		gameInfo:        gameInfo,
		playerGroup:     playerGroup,

		delayLeaveDur: constant.StartRoundPeriod,
	}
}

func (state *StartRoundState) Run(context.Context, ...any) error {
	if err := state.gameRepoService.CreateRound(); err != nil {
		state.Logger().Error("cannot create round", zap.Error(err))
		state.GameController().GoErrorState()
		return nil
	}

	state.Logger().Debug("start round",
		zap.Object("players", state.playerGroup),
	)

	state.GameController().RunTimer(state.delayLeaveDur, func() {
		state.GameController().GoNextState(GoPickFirstState)
	})
	return nil
}

func (state *StartRoundState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	state.MsgBus().Broadcast(core.ModelTopic, &gamegrpc.Model{
		GameInfo:    state.gameInfo.ToProto(),
		PlayerGroup: state.playerGroup.ToProto(),
	})
	return nil
}

func (state *StartRoundState) ToProto(_ core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_StartRoundStateContext{
			StartRoundStateContext: &gamegrpc.StartRoundStateContext{},
		},
	}
}

func (state *StartRoundState) BeforeLeave(uid core.Uid) error {
	state.Logger().Warn("forbid to leave during game", zap.String("uid", uid.String()))
	return status.Errorf(codes.FailedPrecondition, "forbid to leave during game")
}
