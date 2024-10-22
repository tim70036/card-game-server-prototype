package state

import (
	"context"
	commonmodel "card-game-server-prototype/pkg/common/model"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/common/type/gamemode"
	"card-game-server-prototype/pkg/core"
	actor2 "card-game-server-prototype/pkg/game/darkchess/actor"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoResetGameState = core.NewStateTrigger("GoResetGameState")

type ResetGameState struct {
	core.State

	roomInfo               *commonmodel.RoomInfo
	userGroup              *commonmodel.UserGroup
	playerGroup            *model2.PlayerGroup
	actorGroup             *actor2.Group
	gameScoreboard         *model2.GameScoreboard
	gameInfo               *model2.GameInfo
	roundScoreboardRecords *model2.RoundScoreboardRecords
	userService            commonservice.UserService
	board                  *model2.Board
	capturedPieces         *model2.CapturedPieces
	actionHintGroup        *model2.ActionHintGroup
}

func ProvideResetGameState(
	stateFactory *core.StateFactory,

	roomInfo *commonmodel.RoomInfo,
	userGroup *commonmodel.UserGroup,
	playerGroup *model2.PlayerGroup,
	actorGroup *actor2.Group,
	gameScoreboard *model2.GameScoreboard,
	gameInfo *model2.GameInfo,
	roundScoreboardRecords *model2.RoundScoreboardRecords,
	userService commonservice.UserService,
	board *model2.Board,
	capturedPieces *model2.CapturedPieces,
	actionHintGroup *model2.ActionHintGroup,
) *ResetGameState {
	return &ResetGameState{
		State: stateFactory.Create("ResetGameState"),

		roomInfo:               roomInfo,
		userGroup:              userGroup,
		playerGroup:            playerGroup,
		actorGroup:             actorGroup,
		gameScoreboard:         gameScoreboard,
		gameInfo:               gameInfo,
		roundScoreboardRecords: roundScoreboardRecords,
		userService:            userService,
		board:                  board,
		capturedPieces:         capturedPieces,
		actionHintGroup:        actionHintGroup,
	}
}

func (state *ResetGameState) Run(context.Context, ...any) error {
	state.gameInfo.GameId = uuid.NewString()
	state.gameInfo.RoundCount = 1

	state.gameScoreboard.Data = map[core.Uid]*model2.GameScore{}
	state.gameScoreboard.WinnerId = ""
	state.gameScoreboard.IsDraw = false
	for _, uid := range lo.Keys(state.userGroup.Data) {
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
	}

	state.roundScoreboardRecords.Data = make([]*model2.RoundScoreboard, 0)

	state.playerGroup.Data = map[core.Uid]*model2.Player{}
	state.actorGroup.Data = map[core.Uid]actor2.Actor{}

	state.cleanBoard()
	state.board.IsDraw = false
	state.capturedPieces.Pieces = []piece.Piece{}

	state.actionHintGroup.RepeatMovesCount = 0
	state.actionHintGroup.ClaimDraws = []*model2.ClaimDraw{}
	for _, uid := range lo.Keys(state.userGroup.Data) {
		state.actionHintGroup.Data[uid] = &model2.ActionHint{
			Uid:          uid,
			TurnDuration: state.gameInfo.Setting.TurnSecond,
		}
	}

	state.Logger().Debug("reset game done",
		zap.Object("gameInfo", state.gameInfo),
		zap.Object("users", state.userGroup),
	)

	if err := state.userService.FetchFromRepo(lo.Keys(state.userGroup.Data)...); err != nil {
		state.Logger().Error("failed to fetch user from repo", zap.Error(err), zap.Object("users", state.userGroup))
		state.GameController().GoErrorState()
		return nil
	}

	if state.roomInfo.GameMode == gamemode.Buddy {
		state.GameController().GoNextState(GoWaitingRoomState)
		return nil
	}

	state.GameController().GoNextState(GoWaitUserState)
	return nil
}

func (state *ResetGameState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	state.MsgBus().Broadcast(core.ModelTopic, &gamegrpc.Model{
		Board:           state.board.ToProto(),
		CapturedPieces:  state.capturedPieces.ToProto(),
		ActionHintGroup: state.actionHintGroup.ToProto(),
	})
	return nil
}

func (state *ResetGameState) ToProto(core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_ResetGameStateContext{
			ResetGameStateContext: &gamegrpc.ResetGameStateContext{},
		},
	}
}

func (state *ResetGameState) cleanBoard() {
	state.board.Cells = [][]*model2.Cell{}
	var row []*model2.Cell

	var rowCnt int
	for i := 0; i < constant.BoardHeight*constant.BoardWidth; i++ {
		row = append(row, &model2.Cell{
			Piece: piece.InvalidPiece,
		})

		if len(row) == constant.BoardHeight {
			var newRow []*model2.Cell
			newRow = append(newRow, row...)
			state.board.Cells = append(state.board.Cells, newRow)
			row = []*model2.Cell{}
			rowCnt++
		}
	}
}
