package state

import (
	"context"
	commonmodel "card-game-server-prototype/pkg/common/model"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoResetRoundState = core.NewStateTrigger("GoResetRoundState")

type ResetRoundState struct {
	core.State

	userService     commonservice.UserService
	userGroup       *commonmodel.UserGroup
	playerGroup     *model2.PlayerGroup
	gameInfo        *model2.GameInfo
	roundScoreboard *model2.RoundScoreboard
	pickBoard       *model2.PickBoard
	board           *model2.Board
	actionHintGroup *model2.ActionHintGroup
	capturedPieces  *model2.CapturedPieces
}

func ProvideResetRoundState(
	stateFactory *core.StateFactory,

	userService commonservice.UserService,
	userGroup *commonmodel.UserGroup,
	playerGroup *model2.PlayerGroup,
	gameInfo *model2.GameInfo,
	roundScoreboard *model2.RoundScoreboard,
	pickBoard *model2.PickBoard,
	board *model2.Board,
	actionHintGroup *model2.ActionHintGroup,
	capturedPieces *model2.CapturedPieces,
) *ResetRoundState {
	return &ResetRoundState{
		State: stateFactory.Create("ResetRoundState"),

		userService:     userService,
		userGroup:       userGroup,
		playerGroup:     playerGroup,
		gameInfo:        gameInfo,
		roundScoreboard: roundScoreboard,
		pickBoard:       pickBoard,
		board:           board,
		actionHintGroup: actionHintGroup,
		capturedPieces:  capturedPieces,
	}
}

func (state *ResetRoundState) Run(context.Context, ...any) error {
	state.roundScoreboard.Data = make(map[core.Uid]*model2.RoundScore)
	state.roundScoreboard.WinnerId = ""
	state.roundScoreboard.IsDraw = false

	state.playerGroup.Data = make(map[core.Uid]*model2.Player)

	for _, uid := range lo.Keys(state.userGroup.Data) {
		state.playerGroup.Data[uid] = &model2.Player{
			Uid:   uid,
			Color: commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_INVALID,
		}

		state.roundScoreboard.Data[uid] = &model2.RoundScore{
			Uid: uid,
		}
	}

	state.initPickBoard()

	state.initBoard()
	state.Logger().Debug("init board",
		zap.Object("pick", state.pickBoard),
		zap.Object("board", state.board),
	)

	state.board.IsDraw = false
	state.board.TurnCount = 1
	state.capturedPieces.Pieces = []piece.Piece{}

	state.actionHintGroup.RepeatMovesCount = 0
	state.actionHintGroup.LastAction = nil
	state.actionHintGroup.ClaimDraws = []*model2.ClaimDraw{}
	for _, uid := range lo.Keys(state.userGroup.Data) {
		state.actionHintGroup.Data[uid] = &model2.ActionHint{
			Uid:          uid,
			TurnDuration: state.gameInfo.Setting.TurnSecond,
		}
	}

	state.Logger().Debug("reset round state",
		zap.Object("gameInfo", state.gameInfo),
		zap.Object("players", state.playerGroup),
		zap.Object("roundScoreboard", state.roundScoreboard),
		zap.Object("board", state.board),
	)

	if err := state.userService.FetchFromRepo(lo.Keys(state.userGroup.Data)...); err != nil {
		state.Logger().Error("failed to fetch user from repo", zap.Error(err), zap.Object("users", state.userGroup))
		state.GameController().GoErrorState()
		return nil
	}

	state.GameController().GoNextState(GoStartRoundState)
	return nil
}

func (state *ResetRoundState) Publish(context.Context, ...any) error {
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

func (state *ResetRoundState) ToProto(core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_ResetRoundStateContext{
			ResetRoundStateContext: &gamegrpc.ResetRoundStateContext{},
		},
	}
}

func (state *ResetRoundState) BeforeLeave(uid core.Uid) error {
	state.Logger().Warn("forbid to leave during game", zap.String("uid", uid.String()))
	return status.Errorf(codes.FailedPrecondition, "forbid to leave during game")
}

func (state *ResetRoundState) initPickBoard() {
	state.pickBoard.Pieces = []piece.Piece{
		piece.GeneralBlack,
		piece.AdvisorBlack0,
		piece.ElephantBlack0,
		piece.ChariotBlack0,
		piece.HorseBlack0,
	}

	lo.Shuffle(state.pickBoard.Pieces)
	state.pickBoard.PickedPieces = map[core.Uid]piece.Piece{}
	state.pickBoard.FirstPlayer = ""
}

func (state *ResetRoundState) initBoard() {
	pieces := []piece.Piece{
		piece.GeneralRed,
		piece.AdvisorRed0,
		piece.AdvisorRed1,
		piece.ElephantRed0,
		piece.ElephantRed1,
		piece.ChariotRed0,
		piece.ChariotRed1,
		piece.HorseRed0,
		piece.HorseRed1,
		piece.CannonRed0,
		piece.CannonRed1,
		piece.SoldierRed0,
		piece.SoldierRed1,
		piece.SoldierRed2,
		piece.SoldierRed3,
		piece.SoldierRed4,
		piece.GeneralBlack,
		piece.AdvisorBlack0,
		piece.AdvisorBlack1,
		piece.ElephantBlack0,
		piece.ElephantBlack1,
		piece.ChariotBlack0,
		piece.ChariotBlack1,
		piece.HorseBlack0,
		piece.HorseBlack1,
		piece.CannonBlack0,
		piece.CannonBlack1,
		piece.SoldierBlack0,
		piece.SoldierBlack1,
		piece.SoldierBlack2,
		piece.SoldierBlack3,
		piece.SoldierBlack4,
	}

	lo.Shuffle(pieces)

	state.board.Cells = [][]*model2.Cell{}
	var row []*model2.Cell

	var rowCnt int
	for _, p := range pieces {
		row = append(row, &model2.Cell{
			Piece: p,
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
