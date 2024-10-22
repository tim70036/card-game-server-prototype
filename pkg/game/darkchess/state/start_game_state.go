package state

import (
	"context"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	actor2 "card-game-server-prototype/pkg/game/darkchess/actor"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/game/darkchess/service"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"card-game-server-prototype/pkg/util"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 主導初始化遊戲需要的資料

var GoStartGameState = core.NewStateTrigger("GoStartGameState")

type StartGameState struct {
	core.State

	testCFG          *config.TestConfig
	userGroup        *commonmodel.UserGroup
	actorGroup       *actor2.Group
	baseActorFactory *actor2.BaseActorFactory
	aiActorFactory   *actor2.AiActorFactory
	playerGroup      *model2.PlayerGroup
	gameScoreboard   *model2.GameScoreboard
	roundScoreboard  *model2.RoundScoreboard
	gameRepoService  service.GameRepoService
	gameInfo         *model2.GameInfo

	delayLeaveDur time.Duration
}

func ProvideStartGameState(
	stateFactory *core.StateFactory,

	testCFG *config.TestConfig,
	userGroup *commonmodel.UserGroup,
	actorGroup *actor2.Group,
	baseActorFactory *actor2.BaseActorFactory,
	aiActorFactory *actor2.AiActorFactory,
	playerGroup *model2.PlayerGroup,
	gameScoreboard *model2.GameScoreboard,
	roundScoreboard *model2.RoundScoreboard,
	gameRepoService service.GameRepoService,
	gameInfo *model2.GameInfo,
) *StartGameState {
	return &StartGameState{
		State: stateFactory.Create("StartGameState"),

		testCFG:          testCFG,
		userGroup:        userGroup,
		actorGroup:       actorGroup,
		baseActorFactory: baseActorFactory,
		aiActorFactory:   aiActorFactory,
		playerGroup:      playerGroup,
		gameScoreboard:   gameScoreboard,
		roundScoreboard:  roundScoreboard,
		gameRepoService:  gameRepoService,
		gameInfo:         gameInfo,

		delayLeaveDur: constant.StartGamePeriod,
	}
}

func (state *StartGameState) Run(context.Context, ...any) error {
	for uid := range state.userGroup.Data {
		state.playerGroup.Data[uid] = &model2.Player{
			Uid:   uid,
			Color: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_INVALID,
		}

		state.roundScoreboard.Data[uid] = &model2.RoundScore{
			Uid: uid,
		}

		state.gameScoreboard.Data[uid] = &model2.GameScore{
			Uid: uid,
			ExpInfo: &model2.Exp{
				BeforeLevel:    1, // Prevent frontend crash.
				BeforeExp:      0,
				LevelUpExp:     0,
				NextLevel:      1, // Prevent frontend crash.
				NextLevelUpExp: 0,
				IncreaseExp:    0,
			},
		}

		if *state.testCFG.AutopilotMode || state.userGroup.Data[uid].IsAI {
			state.actorGroup.Data[uid] = state.aiActorFactory.Create(uid)
		} else {
			state.actorGroup.Data[uid] = state.baseActorFactory.Create(uid)
		}
	}

	state.Logger().Debug("starting game",
		zap.Object("users", state.userGroup),
		zap.Object("players", state.playerGroup),
		util.DebugField(zap.Object("actors", state.actorGroup)),
		util.DebugField(zap.Object("gameScoreboard", state.gameScoreboard)),
	)

	if err := state.gameRepoService.CreateGame(); err != nil {
		state.Logger().Error("cannot create game", zap.Error(err))
		state.GameController().GoErrorState()
		return nil
	}

	state.GameController().RunTimer(state.delayLeaveDur, func() {
		state.GameController().GoNextState(GoResetRoundState)
	})
	return nil
}

func (state *StartGameState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	state.MsgBus().Broadcast(core.ModelTopic,
		&gamegrpc.Model{
			PlayerGroup: state.playerGroup.ToProto(),
			GameInfo:    state.gameInfo.ToProto(),
		})
	return nil
}

func (state *StartGameState) ToProto(core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_StartGameStateContext{
			StartGameStateContext: &gamegrpc.StartGameStateContext{},
		},
	}
}

func (state *StartGameState) BeforeLeave(uid core.Uid) error {
	state.Logger().Warn("forbid to leave during game", zap.String("uid", uid.String()))
	return status.Errorf(codes.FailedPrecondition, "forbid to leave during game")
}
