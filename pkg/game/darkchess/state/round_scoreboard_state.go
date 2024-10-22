package state

import (
	"context"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/actor"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	service2 "card-game-server-prototype/pkg/game/darkchess/service"
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"sync"
	"time"
)

var GoRoundScoreboardState = core.NewStateTrigger("GoRoundScoreboardState")

type RoundScoreboardState struct {
	core.State

	userGroup              *commonmodel.UserGroup
	playerGroup            *model2.PlayerGroup
	actorGroup             *actor.Group
	gameInfo               *model2.GameInfo
	roundScoreboard        *model2.RoundScoreboard
	roundScoreboardRecords *model2.RoundScoreboardRecords
	gameRepoService        service2.GameRepoService
	boardService           *service2.BoardService
	capturedPieces         *model2.CapturedPieces
	board                  *model2.Board

	duration time.Duration
	// Record the player requests to skip waiting duration for
	// scoreboard.
	requestToSkipPlayers map[core.Uid]struct{}
	cancelTimer          context.CancelFunc
	endOnce              *sync.Once
}

func ProvideRoundScoreboardState(
	stateFactory *core.StateFactory,

	userGroup *commonmodel.UserGroup,
	playerGroup *model2.PlayerGroup,
	actorGroup *actor.Group,
	gameInfo *model2.GameInfo,
	roundScoreboard *model2.RoundScoreboard,
	roundScoreboardRecords *model2.RoundScoreboardRecords,
	gameRepoService service2.GameRepoService,
	boardService *service2.BoardService,
	capturedPieces *model2.CapturedPieces,
	board *model2.Board,
) *RoundScoreboardState {
	return &RoundScoreboardState{
		State: stateFactory.Create("RoundScoreboardState"),

		userGroup:              userGroup,
		playerGroup:            playerGroup,
		actorGroup:             actorGroup,
		gameInfo:               gameInfo,
		roundScoreboard:        roundScoreboard,
		roundScoreboardRecords: roundScoreboardRecords,
		gameRepoService:        gameRepoService,
		boardService:           boardService,
		capturedPieces:         capturedPieces,
		board:                  board,

		duration: 10 * time.Second,
	}
}

func (state *RoundScoreboardState) Run(context.Context, ...any) error {
	state.cancelTimer = func() {}
	state.endOnce = &sync.Once{}
	state.requestToSkipPlayers = map[core.Uid]struct{}{}

	// Update results
	for uid, v := range state.playerGroup.Data {
		var captureList []piece.Piece
		for _, captured := range state.capturedPieces.Pieces {
			if v.Color == captured.GetOppositeColor() {
				captureList = append(captureList, captured)
			}
		}
		state.roundScoreboard.Data[uid].CapturePieces = append(state.roundScoreboard.Data[uid].CapturePieces, captureList...)
		state.roundScoreboard.Data[uid].Color = v.Color
	}

	if state.board.IsDraw {
		state.roundScoreboard.IsDraw = true
	} else {
		// eval winner
		winner, hasWinner := lo.FindKeyBy(state.playerGroup.Data, func(uid core.Uid, p *model2.Player) bool {
			return p.IsWinner
		})

		state.roundScoreboard.WinnerId = winner
		state.roundScoreboard.IsDraw = !hasWinner

		if hasWinner {
			loser := lo.Without(lo.Keys(state.playerGroup.Data), winner)[0]

			var scoreModifier int
			capturePiecesCnt := len(state.roundScoreboard.Data[loser].CapturePieces)

			// 倍數是算輸的人的 capture pieces 數量
			if capturePiecesCnt >= 0 && capturePiecesCnt <= 5 {
				scoreModifier = int(gamegrpc.ScoreModifierType_SCORE_MODIFIER_TYPE_5)
			} else if capturePiecesCnt > 5 && capturePiecesCnt <= 12 {
				scoreModifier = int(gamegrpc.ScoreModifierType_SCORE_MODIFIER_TYPE_3)
			} else if capturePiecesCnt > 12 && capturePiecesCnt <= 16 {
				scoreModifier = int(gamegrpc.ScoreModifierType_SCORE_MODIFIER_TYPE_1)
			}

			state.roundScoreboard.Data[winner].ScoreModifier = scoreModifier
			state.roundScoreboard.Data[loser].ScoreModifier = scoreModifier

			rawProfits := state.gameInfo.Setting.AnteAmount * scoreModifier

			state.roundScoreboard.Data[winner].RawProfit = rawProfits
			state.roundScoreboard.Data[loser].RawProfit = rawProfits * -1
		}
	}

	state.Logger().Debug(
		"round scoreboard evaluated",
		zap.Object("roundScoreboard", state.roundScoreboard),
	)

	state.settleBankrupt()

	state.Logger().Debug("bankrupt evaluated",
		zap.Object("users", state.userGroup),
		zap.Object("players", state.playerGroup),
		zap.Object("roundScoreboard", state.roundScoreboard),
	)

	if err := state.gameRepoService.SubmitRoundScore(); err != nil {
		state.Logger().Error("cannot submit round score", zap.Error(err))
		state.GameController().GoErrorState()
		return nil
	}

	state.cancelTimer = state.GameController().RunTimer(state.duration, state.goNextState)

	for _, curActor := range state.actorGroup.Data {
		reqs, err := curActor.DecideSkipScoreboard()
		if err != nil {
			state.Logger().Error("actor cannot decide call",
				zap.Error(err),
				zap.Object("actor", curActor),
			)
			state.GameController().GoErrorState()
			return nil
		}

		state.GameController().RunActorRequests(reqs)
	}

	return nil
}

func (state *RoundScoreboardState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	return nil
}

func (state *RoundScoreboardState) ToProto(core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_RoundScoreboardStateContext{
			RoundScoreboardStateContext: &gamegrpc.RoundScoreboardStateContext{
				Duration:   durationpb.New(state.duration),
				Scoreboard: state.roundScoreboard.ToProto(),
			}},
	}
}

func (state *RoundScoreboardState) BeforeLeave(uid core.Uid) error {
	state.Logger().Warn("forbid to leave during game", zap.String("uid", uid.String()))
	return status.Errorf(codes.FailedPrecondition, "forbid to leave during game")
}

func (state *RoundScoreboardState) AcceptRequestTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(&gamegrpc.SkipScoreboardRequest{}),
	}
}

func (state *RoundScoreboardState) HandleRequest(req *core.Request) error {
	myPlayer, ok := state.playerGroup.Data[req.Uid]
	if !ok {
		state.Logger().Warn("HandleRequest failed, player not found",
			zap.Object("req", req),
		)
		return status.Errorf(codes.NotFound, "player not found")
	}

	if _, ok := state.requestToSkipPlayers[myPlayer.Uid]; ok {
		state.Logger().Debug("already request to skip round scoreboard", zap.Object("req", req))
		return status.Errorf(codes.AlreadyExists, "player already request to skip round scoreboard")
	}

	switch req.Msg.(type) {
	case *gamegrpc.SkipScoreboardRequest:
		state.requestToSkipPlayers[myPlayer.Uid] = struct{}{}
	default:
		state.Logger().Warn("not supported request", zap.Object("req", req))
		return status.Errorf(codes.Unimplemented, "not supported request")
	}

	if len(state.requestToSkipPlayers) == len(state.playerGroup.Data) {
		state.Logger().Debug("all players have requested to skip round scoreboard")
		state.GameController().RunTask(state.goNextState)
	}

	return nil
}

func (state *RoundScoreboardState) settleBankrupt() {
	for uid, score := range state.roundScoreboard.Data {

		var cash int
		if v, ok := state.userGroup.Data[uid]; ok {
			cash = v.Cash
		}

		if score.RawProfit < 0 && score.RawProfit+cash <= 0 {
			state.playerGroup.Data[uid].IsBankrupt = true
			state.roundScoreboard.Data[uid].RawProfit = cash * -1

			winner := lo.Without(lo.Keys(state.playerGroup.Data), uid)[0]
			state.roundScoreboard.Data[winner].RawProfit = cash
		}
	}
}

func (state *RoundScoreboardState) goNextState() {
	state.endOnce.Do(func() {
		state.Logger().Debug("ending and updating from round result")
		state.cancelTimer()

		if err := state.gameRepoService.FetchRoundResult(); err != nil {
			state.Logger().Error("cannot update from round result",
				zap.Error(err),
			)
			state.GameController().GoErrorState()
			return
		}

		for uid, roundScore := range state.roundScoreboard.Data {
			// "playerGroup.Get" has checked in "Run" step.
			if _, ok := state.playerGroup.Data[uid]; ok {
				state.playerGroup.Data[uid].RawProfit += roundScore.RawProfit
				state.playerGroup.Data[uid].Profit += roundScore.Profit
			}
		}

		state.Logger().Debug("round scoreboard updated from submit result",
			zap.Object("players", state.playerGroup),
			zap.Object("scoreboard", state.roundScoreboard),
		)

		copiedRoundScoreboard := &model2.RoundScoreboard{}
		if err := copier.CopyWithOption(copiedRoundScoreboard, state.roundScoreboard, copier.Option{DeepCopy: true}); err != nil {
			state.Logger().Error("cannot copy round scoreboard",
				zap.Error(err),
			)
			state.GameController().GoErrorState()
			return
		}

		state.roundScoreboardRecords.Data = append(state.roundScoreboardRecords.Data, copiedRoundScoreboard)

		state.Logger().Debug("round scoreboard records updated",
			zap.Array("records", state.roundScoreboardRecords),
		)

		state.GameController().GoNextState(GoEndRoundState)
	})
}
