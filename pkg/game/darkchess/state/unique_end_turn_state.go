package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/game/darkchess/service"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"card-game-server-prototype/pkg/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
)

var GoEndTurnState = core.NewStateTrigger("GoEndTurnState")

type EndTurnState struct {
	core.State

	boardService    *service.BoardService
	board           *model2.Board
	actionHintGroup *model2.ActionHintGroup
	playerGroup     *model2.PlayerGroup

	endOnce *sync.Once
}

func ProvideEndTurnState(
	stateFactory *core.StateFactory,

	boardService *service.BoardService,
	board *model2.Board,
	actionHintGroup *model2.ActionHintGroup,
	playerGroup *model2.PlayerGroup,
) *EndTurnState {
	return &EndTurnState{
		State: stateFactory.Create("EndTurnState"),

		boardService:    boardService,
		board:           board,
		actionHintGroup: actionHintGroup,
		playerGroup:     playerGroup,
	}
}

func (state *EndTurnState) Run(context.Context, ...any) error {
	state.endOnce = &sync.Once{}

	state.goNextState()
	return nil
}

func (state *EndTurnState) Publish(context.Context, ...any) error {
	// 這個 state 前端不會 handle，一般有沒有 broadcast 沒差
	// 但這裡是特例，為了壓縮 user 換 turn 的等待時間，所以不 broadcast。
	// 最好的方式是修改遊戲流程，讓 user 可以在非自己的 turn 時做些預操作。
	// 但這是大工程，沒時間改，只好壓縮其他 delay duration。

	// state.MsgBus().Broadcast(core.GameStateTopic,
	// 	state.ToProto("").(*gamegrpc.GameState),
	// )
	return nil
}

func (state *EndTurnState) ToProto(_ core.Uid) proto.Message {
	actorUid, _ := lo.FindKeyBy(state.actionHintGroup.Data, func(_ core.Uid, p *model2.ActionHint) bool {
		return p.TurnCount == state.board.TurnCount
	})

	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_EndTurnStateContext{
			EndTurnStateContext: &gamegrpc.EndTurnStateContext{
				ActorUid: actorUid.String(),
			},
		},
	}
}

func (state *EndTurnState) BeforeLeave(uid core.Uid) error {
	state.Logger().Warn("forbid to leave during game", zap.String("uid", uid.String()))
	return status.Errorf(codes.FailedPrecondition, "forbid to leave during game")
}

func (state *EndTurnState) goNextState() {
	state.endOnce.Do(func() {
		if state.boardService.IsRedAllDead() {
			winner := lo.Keys(lo.PickBy(state.playerGroup.Data, func(uid core.Uid, p *model2.Player) bool {
				return p.Color == commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_BLACK
			}))[0]

			state.playerGroup.Data[winner].IsWinner = true

			state.Logger().Debug("end turn, winner is black",
				zap.String("winner", winner.String()),
			)

			state.GameController().GoNextState(GoShowRoundResultState, false, winner)
			return

		} else if state.boardService.IsBlackAllDead() {
			winner := lo.Keys(lo.PickBy(state.playerGroup.Data, func(uid core.Uid, p *model2.Player) bool {
				return p.Color == commongrpc.CnChessColorType_CN_CHESS_COLOR_TYPE_RED
			}))[0]

			state.playerGroup.Data[winner].IsWinner = true

			state.Logger().Debug("end turn, winner is red",
				zap.String("winner", winner.String()),
			)

			state.GameController().GoNextState(GoShowRoundResultState, false, winner)
			return
		}

		curPlayer := lo.Keys(lo.PickBy(state.actionHintGroup.Data, func(uid core.Uid, actionHint *model2.ActionHint) bool {
			return actionHint.TurnCount == state.board.TurnCount
		}))[0]

		// Next turn count
		state.board.TurnCount++

		// Set Next player
		nextPlayer := lo.Without(lo.Keys(state.actionHintGroup.Data), curPlayer)[0]
		state.actionHintGroup.Data[nextPlayer].TurnCount = state.board.TurnCount

		state.Logger().Debug("endTurn",
			util.DebugField(zap.String("nextPlayer", nextPlayer.String())),
			util.DebugField(zap.Int("nextTurn", state.board.TurnCount)),
			zap.Object("actionHints", state.actionHintGroup),
		)

		state.GameController().GoNextState(GoStartTurnState)
	})
}
