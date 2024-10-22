package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/actor"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	service2 "card-game-server-prototype/pkg/game/darkchess/service"
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"sync"
	"time"
)

var GoWaitActionState = core.NewStateTrigger("GoWaitActionState",
	core.UidType,
)

type WaitActionState struct {
	core.State

	boardService     *service2.BoardService
	playerService    *service2.PlayerService
	actorGroup       *actor.Group
	playerGroup      *model2.PlayerGroup
	gameInfo         *model2.GameInfo
	playSettingGroup *model2.PlaySettingGroup
	board            *model2.Board
	actionHintGroup  *model2.ActionHintGroup
	replayGroup      *model2.ReplayGroup
	capturedPieces   *model2.CapturedPieces

	actUid        core.Uid
	startWaitTime time.Time

	cancelTimer context.CancelFunc
	endOnce     *sync.Once
}

func ProvideWaitActionState(
	stateFactory *core.StateFactory,

	boardService *service2.BoardService,
	playerService *service2.PlayerService,
	actorGroup *actor.Group,
	playerGroup *model2.PlayerGroup,
	gameInfo *model2.GameInfo,
	playSettingGroup *model2.PlaySettingGroup,
	board *model2.Board,
	actionHintGroup *model2.ActionHintGroup,
	replayGroup *model2.ReplayGroup,
	capturedPieces *model2.CapturedPieces,
) *WaitActionState {
	return &WaitActionState{
		State: stateFactory.Create("WaitActionState"),

		boardService:     boardService,
		playerService:    playerService,
		playerGroup:      playerGroup,
		gameInfo:         gameInfo,
		playSettingGroup: playSettingGroup,
		actorGroup:       actorGroup,
		board:            board,
		actionHintGroup:  actionHintGroup,
		replayGroup:      replayGroup,
		capturedPieces:   capturedPieces,
	}
}

func (state *WaitActionState) Run(_ context.Context, args ...any) error {
	state.cancelTimer = func() {}
	state.endOnce = &sync.Once{}

	if state.board.IsDraw {
		state.Logger().Debug("draw, end game")
		state.endWaitAction(GoDrawState)
		return nil
	}

	state.actUid = args[0].(core.Uid)

	state.Logger().Debug("waitAction",
		zap.String("actor_uid", state.actUid.String()))

	// random action
	curActor := state.actorGroup.Data[state.actUid]
	actorReqs, err := curActor.DecideAction()
	if err != nil {
		state.Logger().Error("actor cannot decide action", zap.Error(err), zap.Object("actor", curActor))
		state.GameController().GoErrorState()
		return nil
	}
	state.GameController().RunActorRequests(actorReqs)

	state.cancelTimer = state.GameController().RunTimer(state.gameInfo.Setting.TurnSecond, state.autoPlayForIdle)
	state.startWaitTime = time.Now()
	return nil
}

func (state *WaitActionState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	state.MsgBus().Broadcast(core.ModelTopic, &gamegrpc.Model{
		PlayerGroup:     state.playerGroup.ToProto(),
		ActionHintGroup: state.actionHintGroup.ToProto(),
	})
	return nil
}

func (state *WaitActionState) ToProto(core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_WaitActionStateContext{
			WaitActionStateContext: &gamegrpc.WaitActionStateContext{
				ActorUid:  state.actUid.String(),
				TurnCount: int32(state.board.TurnCount),
			}},
	}
}

func (state *WaitActionState) BeforeLeave(uid core.Uid) error {
	state.Logger().Warn("forbid to leave during game", zap.String("uid", uid.String()))
	return status.Errorf(codes.FailedPrecondition, "forbid to leave during game")
}

func (state *WaitActionState) AcceptRequestTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(&gamegrpc.RevealRequest{}),
		reflect.TypeOf(&gamegrpc.MoveRequest{}),
		reflect.TypeOf(&gamegrpc.CaptureRequest{}),
		reflect.TypeOf(&gamegrpc.AskExtraSecondsRequest{}),
		reflect.TypeOf(&gamegrpc.ClaimDrawRequest{}),
		reflect.TypeOf(&gamegrpc.AnswerDrawRequest{}),
		reflect.TypeOf(&gamegrpc.SurrenderRequest{}),
		reflect.TypeOf(&gamegrpc.UpdatePlaySettingRequest{}),
	}
}

// 驗證 move, reveal, capture 錯誤時，僅阻擋後續變更資料，不回傳錯誤讓遊戲繼續進行

func (state *WaitActionState) HandleRequest(req *core.Request) error {
	if _, ok := req.Msg.(*gamegrpc.AnswerDrawRequest); ok {
		// 處理 request 在 darkchess.go 的 HandleRequest，這裡處理 go next state。
		if state.board.IsDraw {
			state.Logger().Info("req:accept draw, end game",
				zap.Object("req", req),
			)
			state.endWaitAction(GoDrawState)
		}
		return nil
	}

	// idle 啟用的自動代打需要累計 idle 次數。
	// 如果代打有設置 delay，這裡會等 delay 結束後才判斷，可能會有「此時 player disable auto mode」的狀況。
	if s, ok := state.playSettingGroup.Data[req.Uid]; ok && s.IsAuto && s.IsTriggerAutoByIdle {
		if winner, isEndGame := state.incrIdleCount(); isEndGame {
			state.endWaitAction(GoShowRoundResultState, false, winner)
			return nil
		}
	}

	// Other players forbidden to action, only allow to update PlaySetting.
	if _, ok := req.Msg.(*gamegrpc.UpdatePlaySettingRequest); ok && req.Uid != state.actUid {
		state.Logger().Warn("other player update play setting, ignored", zap.Object("req", req))
		return nil
	}

	if _, ok := req.Msg.(*gamegrpc.RevealRequest); ok {
		if req.Uid != state.actUid {
			msg := "player raise a RevealRequest in opposite's turn"
			state.Logger().Warn(msg, zap.Object("req", req))
			return status.Error(codes.PermissionDenied, msg)
		}

		if state.playerGroup.Data[req.Uid].HasActed {
			state.Logger().Warn("player has raised a RevealRequest, ignored", zap.Object("req", req))
			return nil
		}
	}

	if _, ok := req.Msg.(*gamegrpc.MoveRequest); ok {
		if req.Uid != state.actUid {
			msg := "player raise a MoveRequest in opposite's turn"
			state.Logger().Warn(msg, zap.Object("req", req))
			return status.Error(codes.PermissionDenied, msg)
		}

		if state.playerGroup.Data[req.Uid].HasActed {
			state.Logger().Warn("player has raised a MoveRequest, ignored", zap.Object("req", req))
			return nil
		}
	}

	if _, ok := req.Msg.(*gamegrpc.CaptureRequest); ok {
		if req.Uid != state.actUid {
			msg := "player raise a CaptureRequest in opposite's turn"
			state.Logger().Warn(msg, zap.Object("req", req))
			return status.Error(codes.PermissionDenied, msg)
		}

		if state.playerGroup.Data[req.Uid].HasActed {
			state.Logger().Warn("player has raised a CaptureRequest, ignored", zap.Object("req", req))
			return nil
		}
	}

	if _, ok := req.Msg.(*gamegrpc.AskExtraSecondsRequest); ok {
		if req.Uid != state.actUid {
			msg := "player raise a AskExtraSecondsRequest in opposite's turn"
			state.Logger().Warn(msg, zap.Object("req", req))
			return status.Error(codes.PermissionDenied, msg)
		}

		if _, ok := state.actionHintGroup.Data[req.Uid]; !ok {
			state.Logger().Error("get actionHintGroup failed", zap.Object("req", req))
			return status.Error(codes.NotFound, "get actionHintGroup failed")
		}

		if actionHint, ok := state.actionHintGroup.Data[req.Uid]; ok {
			if extend, err := lo.Last(actionHint.TimeExtendTurns); err == nil && extend.Turn == state.board.TurnCount {
				state.Logger().Warn("player has raised a AskExtraSecondsRequest, ignored", zap.Object("req", req))
				return nil
			}
		}

		if (constant.TimeExtendCount - len(state.actionHintGroup.Data[req.Uid].TimeExtendTurns)) <= 0 {
			msg := "player run out of TimeExtend times"
			state.Logger().Warn(msg, zap.Object("req", req))
			return status.Error(codes.ResourceExhausted, msg)
		}
	}

	if _, ok := req.Msg.(*gamegrpc.ClaimDrawRequest); ok {
		if state.board.TurnCount < constant.TurnToAllowClaimDraw {
			// 當總回合數≧10，輪到自己行動回合時可提出談和。
			msg := "player raise a ClaimDrawRequest before turn 20"
			state.Logger().Warn(msg, zap.Object("req", req))
			return status.Error(codes.PermissionDenied, msg)
		}

		if req.Uid != state.actUid {
			msg := "player raise a ClaimDrawRequest in opposite's turn"
			state.Logger().Warn(msg, zap.Object("req", req))
			return status.Error(codes.PermissionDenied, msg)
		}

		if _, err := lo.Last(lo.Filter(state.actionHintGroup.ClaimDraws, func(item *model2.ClaimDraw, _ int) bool {
			return item.ClaimUid == req.Uid && item.ClaimTurn == state.board.TurnCount
		})); err == nil {
			state.Logger().Warn("player has raised a ClaimDrawRequest, ignored", zap.Object("req", req))
			return nil
		}

		if claims := lo.Filter(state.actionHintGroup.ClaimDraws, func(item *model2.ClaimDraw, _ int) bool {
			return item.ClaimUid == req.Uid
		}); (constant.ClaimDrawCount - len(claims)) <= 0 {
			msg := "player run out of ClaimDrawRequest times"
			state.Logger().Warn(msg, zap.Object("req", req))
			return status.Error(codes.ResourceExhausted, msg)
		}
	}

	// AnswerDrawRequest 這裡不用驗證重複 request 行為，darkchess.go 的 HandleRequest 已經阻擋了。

	if _, ok := req.Msg.(*gamegrpc.SurrenderRequest); ok {
		if req.Uid != state.actUid {
			msg := "player raise a SurrenderRequest in opposite's turn"
			state.Logger().Warn(msg, zap.Object("req", req))
			return status.Error(codes.PermissionDenied, msg)
		}

		if actionHint, ok := state.actionHintGroup.Data[req.Uid]; ok {
			if surrender, err := lo.Last(actionHint.SurrenderTurns); err == nil && surrender.Turn == state.board.TurnCount {
				state.Logger().Warn("player has raised a SurrenderRequest, ignored", zap.Object("req", req))
				return nil
			}
		}
	}

	switch msg := req.Msg.(type) {

	case *gamegrpc.UpdatePlaySettingRequest:

		if _, ok := state.playSettingGroup.Data[req.Uid]; !ok {
			state.Logger().Error("playSetting not found", zap.String("uid", req.Uid.String()))
			return status.Errorf(codes.NotFound, "playSetting not found")
		}

		state.Logger().Info("req:UpdatePlaySetting",
			zap.Object("req", req),
		)

		state.autoPlay()

		return nil

		// Reveal, Move, Capture, ClaimDraw, Surrender
		// 補這幾個驗證時，記得檢查 actor 需不需要同步補上。
	case *gamegrpc.RevealRequest:

		if msg == nil || msg.GridPosition == nil {
			state.Logger().Warn("invalid reveal request", zap.Object("req", req))
			return status.Error(codes.InvalidArgument, "invalid reveal request")
		}

		x, y := int(msg.GridPosition.X), int(msg.GridPosition.Y)

		if statusErr := state.boardService.ValidRevealRules(x, y); statusErr != nil {
			state.Logger().Warn("invalid reveal position",
				zap.Error(statusErr),
				zap.Object("req", req),
			)
			return statusErr
		}

		state.Logger().Info("req:reveal",
			zap.Object("req", req),
		)

		state.playerGroup.Data[req.Uid].HasActed = true
		state.endWaitAction(GoRevealState, req.Uid, model2.Pos{
			X: x,
			Y: y,
		})
		return nil

	case *gamegrpc.MoveRequest:

		if msg == nil || msg.GoTo == nil {
			state.Logger().Warn("invalid move request", zap.Object("req", req))
			return status.Error(codes.InvalidArgument, "invalid reveal request")
		}

		moveP := piece.New(msg.MovePiece)

		if v, ok := state.actionHintGroup.Data[req.Uid]; ok &&
			v.FreezeCell != nil && moveP.IsSame(v.FreezeCell.Piece) {
			msg := "not allow to move freeze piece"
			state.Logger().Error(msg,
				zap.String("FreezeCell", v.FreezeCell.Piece.GetName()),
				zap.String("uid", req.Uid.String()),
			)
			return status.Error(codes.PermissionDenied, msg)
		}

		if lo.Contains(state.capturedPieces.Pieces, moveP) {
			msg := "not allow to move captured piece"
			state.Logger().Error(msg,
				zap.String("move_piece", moveP.GetName()),
			)
			return status.Error(codes.PermissionDenied, msg)
		}

		x, y, ok := state.boardService.GetCellPos(moveP)
		if !ok {
			msg := "invalid move piece, piece not found"
			state.Logger().Warn(msg,
				zap.Object("req", req),
			)
			return status.Error(codes.NotFound, msg)
		}

		if moveP.GetColor() != state.playerGroup.Data[req.Uid].Color {
			msg := "invalid move piece, player use opposite color"
			state.Logger().Warn(msg,
				zap.Object("req", req),
			)
			return status.Error(codes.PermissionDenied, msg)
		}

		if !moveP.IsSame(state.board.Cells[x][y].Piece) {
			msg := "invalid move piece, not same piece"
			state.Logger().Warn(msg,
				zap.Object("req", req),
				zap.String("real_piece", state.board.Cells[x][y].Piece.GetName()),
			)
			return status.Error(codes.PermissionDenied, msg)
		}

		toX, toY := int(msg.GoTo.X), int(msg.GoTo.Y)

		if statusErr := state.boardService.ValidMoveRules(x, y, toX, toY); statusErr != nil {
			state.Logger().Warn("invalid move",
				zap.Error(statusErr),
				zap.Object("req", req),
			)
			return statusErr
		}

		state.Logger().Info("req:move",
			zap.Object("req", req),
		)

		state.playerGroup.Data[req.Uid].HasActed = true
		from, to := model2.Pos{X: x, Y: y}, model2.Pos{X: toX, Y: toY}
		state.endWaitAction(GoMoveState, req.Uid, from, to)
		return nil

	case *gamegrpc.CaptureRequest:

		if msg == nil || msg.GoTo == nil {
			state.Logger().Warn("invalid capture request", zap.Object("req", req))
			return status.Error(codes.InvalidArgument, "invalid capture request")
		}

		moveP := piece.New(msg.MovePiece)

		x, y, ok := state.boardService.GetCellPos(moveP)
		if !ok {
			msg := "invalid capture piece, piece not found"
			state.Logger().Warn(msg,
				zap.Object("req", req),
			)
			return status.Error(codes.NotFound, msg)
		}

		if lo.Contains(state.capturedPieces.Pieces, moveP) {
			msg := "not allow to move captured piece"
			state.Logger().Error(msg,
				zap.String("move_piece", moveP.GetName()),
			)
			return status.Error(codes.PermissionDenied, msg)
		}

		if moveP.GetColor() != state.playerGroup.Data[req.Uid].Color {
			msg := "invalid hunter piece, player use opposite color"
			state.Logger().Warn(msg,
				zap.Object("req", req),
			)
			return status.Error(codes.PermissionDenied, msg)
		}

		if !moveP.IsSame(state.board.Cells[x][y].Piece) {
			msg := "invalid hunter piece, not same piece"
			state.Logger().Warn(msg,
				zap.Object("req", req),
				zap.String("real_piece", state.board.Cells[x][y].Piece.GetName()),
			)
			return status.Error(codes.PermissionDenied, msg)
		}

		toX, toY := int(msg.GoTo.X), int(msg.GoTo.Y)

		capturedP := piece.New(msg.CapturedPiece)

		if lo.Contains(state.capturedPieces.Pieces, capturedP) {
			msg := "not allow to capture captured piece"
			state.Logger().Error(msg,
				zap.String("captured_piece", moveP.GetName()),
			)
			return status.Error(codes.PermissionDenied, msg)
		}

		if !piece.New(msg.CapturedPiece).IsSame(state.board.Cells[toX][toY].Piece) {
			msg := "invalid captured piece, not same piece"
			state.Logger().Warn(msg,
				zap.Object("req", req),
				zap.String("real_piece", state.board.Cells[toX][toY].Piece.GetName()),
			)
			return status.Error(codes.PermissionDenied, msg)
		}

		if capturedP.GetColor() == state.playerGroup.Data[req.Uid].Color {
			msg := "invalid target piece, player capture his own piece"
			state.Logger().Warn(msg,
				zap.Object("req", req),
			)
			return status.Error(codes.PermissionDenied, msg)
		}

		if statusErr := state.boardService.ValidCaptureRules(x, y, toX, toY); statusErr != nil {
			state.Logger().Warn("invalid capture",
				zap.Error(statusErr),
				zap.Object("req", req),
			)
			return statusErr
		}

		state.Logger().Info("req:capture",
			zap.Object("req", req),
		)

		state.playerGroup.Data[req.Uid].HasActed = true
		from, to := model2.Pos{X: x, Y: y}, model2.Pos{X: toX, Y: toY}
		state.endWaitAction(GoCaptureState, req.Uid, from, to)
		return nil

	case *gamegrpc.AskExtraSecondsRequest:

		state.actionHintGroup.Data[req.Uid].TurnDuration += state.gameInfo.Setting.ExtraTurnSecond
		state.actionHintGroup.Data[req.Uid].TimeExtendTurns = append(state.actionHintGroup.Data[req.Uid].TimeExtendTurns,
			model2.RaiseData{
				Turn:     state.board.TurnCount,
				RaisedAt: time.Now(),
			})

		state.Logger().Info("req:extraSec",
			zap.Object("req", req),
			zap.Object("actionHint", state.actionHintGroup.Data[req.Uid]),
		)

		state.MsgBus().Broadcast(core.ModelTopic, &gamegrpc.Model{
			PlayerGroup:     state.playerGroup.ToProto(),
			ActionHintGroup: state.actionHintGroup.ToProto(),
		})

		state.cancelTimer()
		actualDuration := time.Since(state.startWaitTime)
		remain := state.actionHintGroup.Data[req.Uid].TurnDuration - actualDuration
		state.cancelTimer = state.GameController().RunTimer(remain, state.autoPlayForIdle) // restart timer
		state.startWaitTime = time.Now()

		return nil

	case *gamegrpc.ClaimDrawRequest:

		// 有未回應的請求，拒絕提出新請求
		if lastClaim, err := lo.Last(state.actionHintGroup.ClaimDraws); err == nil && !lastClaim.IsAnswered {
			msg := "a ClaimDrawRequest is not handled"
			state.Logger().Warn(msg, zap.Object("actionHints", state.actionHintGroup))
			return status.Error(codes.FailedPrecondition, msg)
		}

		claimDraw := &model2.ClaimDraw{
			ClaimUid:  req.Uid,
			ClaimTurn: state.board.TurnCount,
		}

		state.actionHintGroup.ClaimDraws = append(state.actionHintGroup.ClaimDraws, claimDraw)

		state.Logger().Info("req:claimDraw",
			zap.Object("req", req),
			zap.Object("actionHints", state.actionHintGroup),
		)

		state.MsgBus().Broadcast(core.ModelTopic, &gamegrpc.Model{
			ClaimDraw: &gamegrpc.ClaimDraw{
				ClaimUid:  req.Uid.String(),
				ClaimTurn: int32(claimDraw.ClaimTurn),
			},
			ActionHintGroup: state.actionHintGroup.ToProto(),
		})

		return nil

	case *gamegrpc.SurrenderRequest:

		// 第一回合結束後開放投降，亦即 turn >= 3
		if state.board.TurnCount < 3 {
			msg := "not allow to surrender in the beginning"
			state.Logger().Warn(msg, zap.Object("req", req))
			return status.Error(codes.PermissionDenied, msg)
		}

		state.actionHintGroup.Data[req.Uid].SurrenderTurns = append(state.actionHintGroup.Data[req.Uid].SurrenderTurns,
			model2.RaiseData{
				Turn:     state.board.TurnCount,
				RaisedAt: time.Now(),
			})

		state.MsgBus().Broadcast(core.ModelTopic, &gamegrpc.Model{
			Surrender: &gamegrpc.Surrender{
				Uid: req.Uid.String(),
			},
		})

		state.Logger().Info("req:surrender",
			zap.Object("req", req),
		)

		state.endWaitAction(GoSurrenderState)
		return nil

	default:
		state.Logger().Debug("not supported request", zap.Object("req", req))
		return status.Errorf(codes.Unimplemented, "not supported request")
	}
}

func (state *WaitActionState) autoPlayForIdle() {
	if winner, isEndGame := state.incrIdleCount(); isEndGame {
		state.endWaitAction(GoShowRoundResultState, false, winner)
		return
	}

	mySetting, ok := state.playSettingGroup.Data[state.actUid]
	if !ok {
		state.Logger().Error(
			"try set auto play for idle player but playSetting not found",
			zap.String("uid", state.actUid.String()),
			zap.Object("playSetting", state.playSettingGroup),
		)
		state.GameController().GoErrorState()
		return
	}

	if mySetting.IsAuto {
		return
	}

	// Handle idle - set autoplay
	// 這裡一定要閒置代打時就馬上啟用
	// 不然以目前的架構會沒有辦法判斷什麼時候該 reset IdleCount。
	mySetting.IsTriggerAutoByIdle = true
	mySetting.IsAuto = true

	state.Logger().Info("timeout, autoplay",
		zap.Object("player", state.playerGroup.Data[state.actUid]),
		zap.Object("playSetting", mySetting),
	)

	state.MsgBus().Unicast(state.actUid, core.ModelTopic,
		&gamegrpc.Model{
			PlaySetting: mySetting.ToProto(),
		},
	)

	// Handle idle - autoplay for player
	// 這裡千萬不要使用 actor 送代打 req 的方式，會產生更多判斷要處理。
	// 像是，取消自動代打時要怎麼取消已經送出的代打 TASK！！！！
	state.autoPlay()
}

func (state *WaitActionState) autoPlay() {
	if x, y, ok := state.boardService.RandomPickUnrevealedPiece(); ok {
		state.endWaitAction(GoRevealState,
			state.actUid,
			model2.Pos{X: x, Y: y},
		)
	} else if x, y, toX, toY, ok := state.boardService.RandomMovePiece(state.actUid); ok {
		state.endWaitAction(GoMoveState,
			state.actUid,
			model2.Pos{X: x, Y: y},
			model2.Pos{X: toX, Y: toY},
		)
	} else if x, y, cX, cY, ok := state.boardService.RandomCapturePiece(state.actUid); ok {
		state.endWaitAction(GoCaptureState,
			state.actUid,
			model2.Pos{X: x, Y: y},
			model2.Pos{X: cX, Y: cY},
		)
	} else {
		state.actionHintGroup.Data[state.actUid].SurrenderTurns = append(state.actionHintGroup.Data[state.actUid].SurrenderTurns,
			model2.RaiseData{
				Turn:     state.board.TurnCount,
				RaisedAt: time.Now(),
			})

		state.MsgBus().Broadcast(core.ModelTopic, &gamegrpc.Model{
			Surrender: &gamegrpc.Surrender{
				Uid: state.actUid.String(),
			},
		})

		state.Logger().Info("idle to auto surrender",
			zap.String("uid", state.actUid.String()),
			zap.Object("players", state.playerGroup),
		)

		state.endWaitAction(GoSurrenderState)
	}
}

// 觸發 idle 累計 +1，到達 max 值時結束遊戲並回傳是否結束遊戲的 flag。
func (state *WaitActionState) incrIdleCount() (winner core.Uid, isEndGame bool) {
	if isUpdate := state.playerService.AppendIdleTurn(state.actUid, state.board.TurnCount); !isUpdate {
		return "", false
	}

	state.Logger().Info("idleCount++",
		zap.String("idlePlayer", state.actUid.String()),
		zap.Object("players", state.playerGroup),
		zap.Object("actionHints", state.actionHintGroup),
	)

	// New rule:
	// Idle 1st turn 代打 by autoplay，IdleCount++ (enable autoplay mode first)
	// Idle 2nd turn 代打 by autoplay，IdleCount++
	// Idle 3rd turn 代打 by autoplay，IdleCount++
	// Idle 4th turn 代打 by autoplay，IdleCount++
	// Idle 5th turn Game Over, idle player lose.
	// Reset IdleCount when player back to the game. (disable autoplay mode)

	// 若連續3回合玩家未連回並執行走棋，會直接判定另一位玩家獲勝。
	// 並記為斷線1場(EndRoundState 的 DisconnectedRoundCount 處理)。
	if state.playerService.IsReachIdleThreshold(state.actUid) {
		winner := lo.Without(lo.Keys(state.playerGroup.Data), state.actUid)[0]
		state.playerGroup.Data[winner].IsWinner = true
		state.playerGroup.Data[state.actUid].DisconnectedRoundCount++

		state.Logger().Info("game set,idleCount reached threshold",
			zap.String("idlePlayer", state.actUid.String()),
			zap.String("winner", winner.String()),
			zap.Object("actionHints", state.actionHintGroup),
		)

		return winner, true
	}

	return "", false
}

func (state *WaitActionState) endWaitAction(nextStateTrigger *core.StateTrigger, args ...any) {
	state.endOnce.Do(func() {
		state.Logger().Debug("endWaitAction")
		state.cancelTimer()

		// 改為 claimDraw 只能存活到對手的下一個回合結束，因此不需要倒數、不需要 timer、不需要 expire、
		if claimDraw, err := lo.Last(state.actionHintGroup.ClaimDraws); err == nil &&
			!claimDraw.IsAnswered && claimDraw.ClaimUid != state.actUid {
			// Set last claim draw to rejected
			claimDraw.IsAccepted = false
			claimDraw.IsAnswered = true

			state.MsgBus().Broadcast(core.ModelTopic, &gamegrpc.Model{
				ClaimDraw: &gamegrpc.ClaimDraw{
					ClaimUid:  claimDraw.ClaimUid.String(),
					ClaimTurn: int32(claimDraw.ClaimTurn),
					Claimed:   true,
					Answered:  false, // refused
				},
				ActionHintGroup: state.actionHintGroup.ToProto(),
			})

			state.Logger().Info("ignore a claim",
				zap.Object("handleDraw", claimDraw),
				zap.Object("actionHints", state.actionHintGroup),
			)
		}

		// reset
		state.actionHintGroup.Data[state.actUid].TurnDuration = state.gameInfo.Setting.TurnSecond
		state.GameController().GoNextState(nextStateTrigger, args...)
	})
}
