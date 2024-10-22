package darkchess

import (
	commonmodel "card-game-server-prototype/pkg/common/model"
	commonserver "card-game-server-prototype/pkg/common/server"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	server2 "card-game-server-prototype/pkg/game/darkchess/server"
	service2 "card-game-server-prototype/pkg/game/darkchess/service"
	state2 "card-game-server-prototype/pkg/game/darkchess/state"
	event2 "card-game-server-prototype/pkg/game/darkchess/type/event"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"card-game-server-prototype/pkg/grpc/coregrpc"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"card-game-server-prototype/pkg/util"
	"github.com/samber/lo"
	"math/rand"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DarkChess struct {
	grpcServer *commonserver.GrpcServer

	userService      commonservice.UserService
	resyncService    *service2.ResyncService
	userGroup        *commonmodel.UserGroup
	playSettingGroup *model2.PlaySettingGroup
	actionHintGroup  *model2.ActionHintGroup
	eventGroup       *model2.EventGroup
	board            *model2.Board
	playerService    *service2.PlayerService
	playerGroup      *model2.PlayerGroup
	msgBus           core.MsgBus
	logger           *zap.Logger

	game core.Game
}

func ProvideGame(
	grpcServer *commonserver.GrpcServer,
	connectionServiceServer *commonserver.ConnectionServiceServer,
	chatServiceServer *commonserver.ChatServiceServer,
	emoteRpcServiceServer *commonserver.EmoteRpcServiceServer,
	actionServiceServer *server2.ActionServiceServer,
	messageServiceServer *server2.MessageServiceServer,

	userService commonservice.UserService,
	resyncService *service2.ResyncService,
	userGroup *commonmodel.UserGroup,
	playSettingGroup *model2.PlaySettingGroup,
	actionHintGroup *model2.ActionHintGroup,
	eventGroup *model2.EventGroup,
	board *model2.Board,
	playerService *service2.PlayerService,
	playerGroup *model2.PlayerGroup,

	msgBus core.MsgBus,
	loggerFactory *util.LoggerFactory,

	game core.Game,

	// Common State
	initState *state2.InitState,
	resetGameState *state2.ResetGameState,
	waitUserState *state2.WaitUserState,
	waitingRoomState *state2.WaitingRoomState,
	startGameState *state2.StartGameState,
	resetRoundState *state2.ResetRoundState,
	startRoundState *state2.StartRoundState,
	endRoundState *state2.EndRoundState,
	roundScoreboardState *state2.RoundScoreboardState,
	endGameState *state2.EndGameState,
	gameScoreboardState *state2.GameScoreboardState,
	closedState *state2.ClosedState,

	// Unique Game State
	pickFirstState *state2.PickFirstState,
	startTurnState *state2.StartTurnState,
	waitActionState *state2.WaitActionState,
	revealState *state2.RevealState,
	moveState *state2.MoveState,
	captureState *state2.CaptureState,
	endTurnState *state2.EndTurnState,
	drawState *state2.DrawState,
	surrenderState *state2.SurrenderState,
	showRoundResultState *state2.ShowRoundResultState,
) *DarkChess {
	coregrpc.RegisterConnectionServiceServer(
		grpcServer,
		connectionServiceServer,
	)

	commongrpc.RegisterChatServiceServer(
		grpcServer,
		chatServiceServer,
	)

	commongrpc.RegisterEmoteRpcServiceServer(
		grpcServer,
		emoteRpcServiceServer,
	)

	gamegrpc.RegisterActionServiceServer(
		grpcServer,
		actionServiceServer,
	)

	gamegrpc.RegisterMessageServiceServer(
		grpcServer,
		messageServiceServer,
	)

	darkChess := &DarkChess{
		game:             game,
		grpcServer:       grpcServer,
		userGroup:        userGroup,
		userService:      userService,
		resyncService:    resyncService,
		playSettingGroup: playSettingGroup,
		actionHintGroup:  actionHintGroup,
		board:            board,
		playerService:    playerService,
		playerGroup:      playerGroup,
		eventGroup:       eventGroup,
		msgBus:           msgBus,
		logger:           loggerFactory.Create("DarkChess"),
	}

	// TTT: 代表不太好追的流程，用 find in files 找。
	// TTT OnXXX 比較像 Register，實際上執行動作在 core 的 chan 負責處理。
	game.OnConnect(darkChess.HandleConnect)
	game.OnEnter(darkChess.HandleEnter)
	game.OnLeave(darkChess.HandleLeave)
	game.OnDisconnect(darkChess.HandleDisconnect)
	game.OnRequest(darkChess.HandleRequest)

	// BEFORE STATE: Register trigger
	game.ConfigTriggerParamsType(state2.GoInitState)
	game.ConfigTriggerParamsType(state2.GoResetGameState)
	game.ConfigTriggerParamsType(state2.GoWaitUserState)
	game.ConfigTriggerParamsType(state2.GoWaitingRoomState)
	game.ConfigTriggerParamsType(state2.GoStartGameState)
	game.ConfigTriggerParamsType(state2.GoResetRoundState)
	game.ConfigTriggerParamsType(state2.GoStartRoundState)
	game.ConfigTriggerParamsType(state2.GoEndRoundState)
	game.ConfigTriggerParamsType(state2.GoRoundScoreboardState)
	game.ConfigTriggerParamsType(state2.GoEndGameState)
	game.ConfigTriggerParamsType(state2.GoGameScoreboardState)
	game.ConfigTriggerParamsType(state2.GoClosedState)

	game.ConfigTriggerParamsType(state2.GoPickFirstState)
	game.ConfigTriggerParamsType(state2.GoStartTurnState)
	game.ConfigTriggerParamsType(state2.GoWaitActionState)
	game.ConfigTriggerParamsType(state2.GoRevealState)
	game.ConfigTriggerParamsType(state2.GoMoveState)
	game.ConfigTriggerParamsType(state2.GoCaptureState)
	game.ConfigTriggerParamsType(state2.GoDrawState)
	game.ConfigTriggerParamsType(state2.GoSurrenderState)
	game.ConfigTriggerParamsType(state2.GoEndTurnState)
	game.ConfigTriggerParamsType(state2.GoShowRoundResultState)

	// Set Common Game State

	game.ConfigErrorState(closedState)

	game.ConfigInitState(initState).
		Permit(state2.GoResetGameState, resetGameState)

	game.ConfigState(resetGameState).
		Permit(state2.GoWaitUserState, waitUserState).
		Permit(state2.GoWaitingRoomState, waitingRoomState)

	game.ConfigState(waitingRoomState).
		Permit(state2.GoWaitUserState, waitUserState)

	game.ConfigState(waitUserState).
		Permit(state2.GoStartGameState, startGameState)

	game.ConfigState(startGameState).
		Permit(state2.GoResetRoundState, resetRoundState)

	game.ConfigState(resetRoundState).
		Permit(state2.GoStartRoundState, startRoundState)

	// Unique Game State Start

	game.ConfigState(startRoundState).
		Permit(state2.GoPickFirstState, pickFirstState)

	game.ConfigState(pickFirstState).
		PermitReentry(state2.GoPickFirstState).
		Permit(state2.GoStartTurnState, startTurnState)

	game.ConfigState(startTurnState).
		Permit(state2.GoWaitActionState, waitActionState)

	game.ConfigState(waitActionState).
		PermitReentry(state2.GoWaitActionState).
		Permit(state2.GoRevealState, revealState).
		Permit(state2.GoMoveState, moveState).
		Permit(state2.GoCaptureState, captureState).
		Permit(state2.GoDrawState, drawState).
		Permit(state2.GoSurrenderState, surrenderState).
		Permit(state2.GoEndRoundState, endRoundState).
		Permit(state2.GoShowRoundResultState, showRoundResultState)

	game.ConfigState(revealState).
		Permit(state2.GoEndTurnState, endTurnState)

	game.ConfigState(moveState).
		Permit(state2.GoEndTurnState, endTurnState)

	game.ConfigState(captureState).
		Permit(state2.GoEndTurnState, endTurnState)

	game.ConfigState(endTurnState).
		Permit(state2.GoStartTurnState, startTurnState).
		Permit(state2.GoShowRoundResultState, showRoundResultState)

	game.ConfigState(drawState).
		Permit(state2.GoShowRoundResultState, showRoundResultState)

	game.ConfigState(surrenderState).
		Permit(state2.GoShowRoundResultState, showRoundResultState)

	game.ConfigState(showRoundResultState).
		Permit(state2.GoRoundScoreboardState, roundScoreboardState)

	// Unique Game State End

	game.ConfigState(roundScoreboardState).
		Permit(state2.GoEndRoundState, endRoundState)

	game.ConfigState(endRoundState).
		Permit(state2.GoResetRoundState, resetRoundState).
		Permit(state2.GoGameScoreboardState, gameScoreboardState)

	game.ConfigState(gameScoreboardState).
		Permit(state2.GoEndGameState, endGameState)

	game.ConfigState(endGameState).
		Permit(state2.GoResetGameState, resetGameState).
		Permit(state2.GoClosedState, closedState)

	return darkChess
}

func (darkChess *DarkChess) Run() {
	darkChess.logger.Debug("start")

	rand.NewSource(time.Now().UnixNano())

	go darkChess.grpcServer.Run()
	go darkChess.game.Run()

	<-darkChess.game.WaitShutdown()
	darkChess.logger.Debug("gameRunner has shutdown, exiting")
}

func (darkChess *DarkChess) HandleConnect(uid core.Uid) error {
	if err := darkChess.userService.CanConnect(uid); err != nil {
		darkChess.logger.Warn("user OnConnect but not allow to connect", zap.Error(err), zap.String("uid", uid.String()))
		return err
	}

	if _, ok := darkChess.userGroup.Data[uid]; !ok {
		if err := darkChess.userService.Init(uid); err != nil {
			darkChess.logger.Error("user OnConnect but failed to init user data", zap.Error(err), zap.String("uid", uid.String()))
			return status.Errorf(codes.Internal, "failed to init user data")
		}

		darkChess.logger.Debug("user OnConnect created new user", zap.Object("user", darkChess.userGroup.Data[uid]))
	}

	if err := darkChess.userService.FetchFromRepo(uid); err != nil {
		darkChess.logger.Error("user OnConnect but failed to refresh user data", zap.Error(err), zap.String("uid", uid.String()))
		return status.Errorf(codes.Internal, "failed to refresh user data")
	}

	// ----- 例外先放這裏 start -----
	// ----- 例外先放這裏 end -----

	user := darkChess.userGroup.Data[uid]
	user.IsConnected = true
	darkChess.logger.Debug("user OnConnect", zap.String("uid", user.Uid.String()), zap.Object("users", darkChess.userGroup))
	darkChess.userService.BroadcastUpdate()
	return nil
}

func (darkChess *DarkChess) HandleEnter(uid core.Uid) error {
	user, ok := darkChess.userGroup.Data[uid]
	if !ok {
		darkChess.logger.Error("user OnEnter but does not exist", zap.String("uid", uid.String()), zap.Object("users", darkChess.userGroup))
		return status.Errorf(codes.NotFound, "user does not exist")
	}

	// ----- 例外先放這裏 start -----
	// ----- 例外先放這裏 end -----

	user.HasEntered = true
	darkChess.logger.Debug("user OnEnter", zap.String("uid", user.Uid.String()), zap.Object("users", darkChess.userGroup))
	darkChess.userService.BroadcastUpdate()
	darkChess.resyncService.Send(uid)
	return nil
}

func (darkChess *DarkChess) HandleLeave(uid core.Uid) error {
	user, ok := darkChess.userGroup.Data[uid]
	if !ok {
		darkChess.logger.Error("user OnLeave but does not exist", zap.String("uid", uid.String()), zap.Object("users", darkChess.userGroup))
		return status.Errorf(codes.NotFound, "user does not exist")
	}

	kickoutMSG := &gamegrpc.Event{
		Kickout: &commongrpc.Kickout{
			Reason: commongrpc.KickoutReason_USER_REQUESTED,
		},
	}

	if err := darkChess.userService.Destroy(core.EventTopic, kickoutMSG, uid); err != nil {
		darkChess.logger.Error("user OnLeave but failed to destroy user data", zap.Error(err), zap.String("uid", uid.String()), zap.Object("users", darkChess.userGroup))
		return status.Errorf(codes.Internal, "failed to destroy user data")
	}

	// ----- 例外先放這裏 start -----

	// ----- 例外先放這裏 end -----

	darkChess.logger.Debug("user OnLeave", zap.Object("user", user), zap.Object("users", darkChess.userGroup))
	darkChess.userService.BroadcastUpdate()
	return nil
}

func (darkChess *DarkChess) HandleDisconnect(uid core.Uid) error {
	user, ok := darkChess.userGroup.Data[uid]
	if !ok {
		// It's normal if user does not exist in user group. EX:
		// client rpc request Leave() → server deletes ms session →
		// server send kickout → client rpc request Close() → client
		// leave room / server deletes gs session
		darkChess.logger.Warn("user OnDisconnect ignored since does not exist", zap.String("uid", uid.String()), zap.Object("users", darkChess.userGroup))
		return nil
	}

	// ----- 例外先放這裏 start -----
	// ----- 例外先放這裏 end -----

	// This is the case that client disconnect due to ping timeout.
	user.IsConnected = false
	darkChess.logger.Debug("user OnDisconnect", zap.Object("user", user), zap.Object("users", darkChess.userGroup))
	darkChess.userService.BroadcastUpdate()
	return nil
}

func (darkChess *DarkChess) HandleRequest(req *core.Request) *core.Response {
	switch msg := req.Msg.(type) {
	case *commongrpc.ResyncRequest:
		darkChess.resyncService.Send(req.Uid)
		return nil

	case *commongrpc.EmotePingRequest:

		darkChess.msgBus.Broadcast(core.EmoteEventTopic, &commongrpc.EmoteEvent{
			EmotePing: &commongrpc.EmotePing{
				ItemId:    msg.ItemId,
				SenderUid: req.Uid.String(),
				TargetUid: msg.TargetUid,
			},
		})
		return nil

	case *commongrpc.StickerRequest:
		darkChess.eventGroup.Data[req.Uid] = append(darkChess.eventGroup.Data[req.Uid], &event2.Event{
			Code:   event2.UseSticker,
			Amount: 1,
		})

		darkChess.msgBus.Broadcast(core.EmoteEventTopic, &commongrpc.EmoteEvent{
			Sticker: &commongrpc.Sticker{
				Uid:       req.Uid.String(),
				StickerId: msg.StickerId,
			},
		})
		return nil

	case *gamegrpc.UpdatePlaySettingRequest:

		if _, ok := darkChess.playSettingGroup.Data[req.Uid]; !ok {
			return &core.Response{Err: status.Errorf(codes.NotFound, "playSetting not found")}
		}

		darkChess.playSettingGroup.Data[req.Uid].IsAuto = msg.IsAuto

		// Player disable autoplay in darkchess.go's HandleRequest.
		// It means player is back to the game.
		// Here to reset data for idle handling.
		if !msg.IsAuto {
			darkChess.playerService.ResetIdleTurn(req.Uid)
			darkChess.playSettingGroup.Data[req.Uid].IsTriggerAutoByIdle = false

			darkChess.logger.Warn("player disable autoplay mode, reset IdleTurns",
				zap.Object("players", darkChess.playerGroup),
				zap.Object("playSetting", darkChess.playSettingGroup),
				zap.Object("req", req),
			)
		}

		darkChess.msgBus.Unicast(req.Uid, core.ModelTopic, &gamegrpc.Model{
			PlaySetting: darkChess.playSettingGroup.Data[req.Uid].ToProto(),
		})

		darkChess.logger.Info("req:UpdatePlaySetting",
			zap.Object("playSetting", darkChess.playSettingGroup),
			zap.Object("req", req),
		)

		return nil

	case *gamegrpc.AnswerDrawRequest:

		lastClaim, err := lo.Last(darkChess.actionHintGroup.ClaimDraws)

		// 沒有，ignored
		if err != nil {
			return nil
		}

		// 處理過了，ignored
		if lastClaim.IsAnswered {
			darkChess.logger.Warn("player has answered ClaimDrawRequest, ignored",
				zap.Object("lastClaim", lastClaim),
			)
			return nil
		}

		// 阻擋提出和棋的人回應自己的提議
		if lastClaim.ClaimUid == req.Uid {
			msg := "player who raised DrawRequest cannot answer from himself"
			darkChess.logger.Warn(msg, zap.String("uid", req.Uid.String()))
			return &core.Response{Err: status.Error(codes.PermissionDenied, msg)}
		}

		for i := range darkChess.actionHintGroup.ClaimDraws {
			if i == len(darkChess.actionHintGroup.ClaimDraws)-1 {
				darkChess.actionHintGroup.ClaimDraws[i].IsAnswered = true
				darkChess.actionHintGroup.ClaimDraws[i].IsAccepted = msg.IsAccept
			}
		}

		darkChess.logger.Info("req:AnswerClaimDraw",
			zap.String("uid", req.Uid.String()),
			zap.Bool("is_accept", msg.IsAccept),
			zap.Object("actionHints", darkChess.actionHintGroup),
		)

		if msg.IsAccept {
			darkChess.logger.Warn("draw accepted, set draw flag for end game",
				zap.String("accept_uid", req.Uid.String()),
				zap.Object("actionHints", darkChess.actionHintGroup),
			)
			darkChess.board.IsDraw = true
		}

		darkChess.msgBus.Broadcast(core.ModelTopic, &gamegrpc.Model{
			ClaimDraw: &gamegrpc.ClaimDraw{
				ClaimUid:  lastClaim.ClaimUid.String(),
				ClaimTurn: int32(lastClaim.ClaimTurn),
				Claimed:   true,
				Answered:  msg.IsAccept,
			},
		})

		return nil

	default:
		darkChess.logger.Debug("ignored not supported request", zap.Object("req", req))
		return nil
	}
}
