package state

import (
	"context"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/game/darkchess/service"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoEndRoundState = core.NewStateTrigger("GoEndRoundState")

type EndRoundState struct {
	core.State

	userGroup        *commonmodel.UserGroup
	playerGroup      *model2.PlayerGroup
	eventGroup       *model2.EventGroup
	eventService     service.EventService
	gameInfo         *model2.GameInfo
	playSettingGroup *model2.PlaySettingGroup
	replayGroup      *model2.ReplayGroup
}

func ProvideEndRoundState(
	stateFactory *core.StateFactory,

	userGroup *commonmodel.UserGroup,
	playerGroup *model2.PlayerGroup,
	eventGroup *model2.EventGroup,
	eventService service.EventService,
	gameInfo *model2.GameInfo,
	playSettingGroup *model2.PlaySettingGroup,
	replayGroup *model2.ReplayGroup,
) *EndRoundState {
	return &EndRoundState{
		State: stateFactory.Create("EndRoundState"),

		userGroup:        userGroup,
		playerGroup:      playerGroup,
		eventGroup:       eventGroup,
		eventService:     eventService,
		gameInfo:         gameInfo,
		playSettingGroup: playSettingGroup,
		replayGroup:      replayGroup,
	}
}

func (state *EndRoundState) Run(context.Context, ...any) error {
	for uid, p := range state.playerGroup.Data {
		if !state.userGroup.Data[uid].IsConnected {
			p.DisconnectedRoundCount++
			state.Logger().Info("user is not connected, incr disconnect round count",
				zap.String("uid", uid.String()),
			)
		}
	}

	if err := state.eventService.EvalRoundEvents(); err != nil {
		state.Logger().Error("eval round events failed", zap.Error(err))
		state.GameController().GoErrorState()
		return nil
	}

	state.Logger().Debug("endRound",
		zap.Object("evaluatedEvent", state.eventGroup),
		util.DebugField(zap.Array("replay", state.replayGroup)),
	)

	if err := state.eventService.Submit(); err != nil {
		state.Logger().Error("submit round events failed", zap.Error(err))
		state.GameController().GoErrorState()
		return nil
	}

	state.goNextState()
	return nil
}

func (state *EndRoundState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	return nil
}

func (state *EndRoundState) ToProto(core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_EndRoundStateContext{
			EndRoundStateContext: &gamegrpc.EndRoundStateContext{},
		},
	}
}

func (state *EndRoundState) BeforeLeave(uid core.Uid) error {
	state.Logger().Warn("forbid to leave during game", zap.String("uid", uid.String()))
	return status.Errorf(codes.FailedPrecondition, "forbid to leave during game")
}

func (state *EndRoundState) isToEndGame() bool {
	for _, player := range state.playerGroup.Data {
		// is Someone Bankrupt
		if player.IsBankrupt {
			state.Logger().Info("someone bankrupt, ending game")
			return true
		}

		// is Someone Disconnect Too Long
		if player.DisconnectedRoundCount >= constant.EndGameDisconnectedRoundCount {
			state.Logger().Info("someone disconnect for too long, ending game")
			return true
		}
	}

	if state.gameInfo.RoundCount == state.gameInfo.Setting.TotalRound {
		state.Logger().Debug("endingGame")
		return true
	}

	return false
}

func (state *EndRoundState) goNextState() {
	if state.isToEndGame() {
		state.GameController().GoNextState(GoGameScoreboardState)
		return
	}

	state.gameInfo.RoundCount++
	state.GameController().GoNextState(GoResetRoundState)
}
