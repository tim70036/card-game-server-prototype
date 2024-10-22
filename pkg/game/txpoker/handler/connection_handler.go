package handler

import (
	commonmodel "card-game-server-prototype/pkg/common/model"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/constant"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	service2 "card-game-server-prototype/pkg/game/txpoker/service"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConnectionHandler struct {
	core.BaseHandler
	game             core.Game
	userGroup        *commonmodel.UserGroup
	seatStatusGroup  *model2.SeatStatusGroup
	actionHintGroup  *model2.ActionHintGroup
	playSettingGroup *model2.PlaySettingGroup
	eventGroup       *model2.EventGroup

	userService       commonservice.UserService
	resyncService     *service2.ResyncService
	seatStatusService service2.SeatStatusService
	msgBus            core.MsgBus
	logger            *zap.Logger
}

func ProvideConnectionHandler(
	game core.Game,
	userGroup *commonmodel.UserGroup,
	userService commonservice.UserService,
	resyncService *service2.ResyncService,
	seatStatusService service2.SeatStatusService,
	msgBus core.MsgBus,
	loggerFactory *util.LoggerFactory,
) *ConnectionHandler {
	return &ConnectionHandler{
		game:              game,
		userGroup:         userGroup,
		userService:       userService,
		resyncService:     resyncService,
		seatStatusService: seatStatusService,
		msgBus:            msgBus,
		logger:            loggerFactory.Create("ConnectionHandler"),
	}
}

func (handler *ConnectionHandler) HandleConnect(uid core.Uid) error {
	if err := handler.userService.CanConnect(uid); err != nil {
		handler.logger.Warn("user OnConnect but not allow to connect", zap.Error(err), zap.String("uid", uid.String()))
		return err
	}

	if _, ok := handler.userGroup.Data[uid]; !ok {
		if err := handler.userService.Init(uid); err != nil {
			handler.logger.Error("user OnConnect but failed to init user data", zap.Error(err), zap.String("uid", uid.String()))
			return status.Errorf(codes.Internal, "failed to init user data")
		}

		handler.logger.Debug("user OnConnect created new user", zap.Object("user", handler.userGroup.Data[uid]))
	}

	if err := handler.userService.FetchFromRepo(uid); err != nil {
		handler.logger.Error("user OnConnect but failed to refresh user data", zap.Error(err), zap.String("uid", uid.String()))
		return status.Errorf(codes.Internal, "failed to refresh user data")
	}

	user := handler.userGroup.Data[uid]
	user.IsConnected = true

	handler.logger.Info(
		"user OnConnect",
		zap.Object("user", user),
		zap.Array("users", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
			for _, v := range handler.userGroup.Data {
				enc.AppendString(v.Uid.String())
			}
			return nil
		})),
	)
	handler.userService.BroadcastUpdate()
	return nil
}

func (handler *ConnectionHandler) HandleEnter(uid core.Uid) error {
	user, ok := handler.userGroup.Data[uid]
	if !ok {
		handler.logger.Error("user OnEnter but does not exist", zap.String("uid", uid.String()), zap.Object("users", handler.userGroup))
		return status.Errorf(codes.NotFound, "user does not exist")
	}

	user.HasEntered = true

	handler.logger.Debug("user OnEnter",
		zap.String("uid", user.Uid.String()),
		zap.Object("users", handler.userGroup),
	)
	handler.userService.BroadcastUpdate()

	handler.resyncService.Send(uid)
	return nil
}

func (handler *ConnectionHandler) HandleLeave(uid core.Uid) error {
	handler.logger.Info("user OnLeave", zap.String("uid", uid.String()))
	return handler.userLeaveProcess(uid, commongrpc.KickoutReason_USER_REQUESTED)
}

func (handler *ConnectionHandler) HandleDisconnect(uid core.Uid) error {
	user, ok := handler.userGroup.Data[uid]
	if !ok {
		// It's normal if user does not exist in user group. EX:
		// client rpc request Leave() → server deletes ms session →
		// server send kickout → client rpc request Close() → client
		// leave room / server deletes gs session
		handler.logger.Info("user OnDisconnect ignored since does not exist", zap.String("uid", uid.String()), zap.Object("users", handler.userGroup))
		return nil
	}

	// This is the case that client disconnect due to ping timeout.
	user.IsConnected = false

	handler.logger.Info(
		"user OnDisconnect",
		zap.Object("user", user),
		util.DebugField(zap.Object("users", handler.userGroup)),
	)
	handler.userService.BroadcastUpdate()

	handler.game.RunTimer(
		constant.DisconnectGracefulPeriod,
		func() { handler.handleDisconnectAfterGracefulPeriod(uid) },
	)
	return nil
}

// 玩家有一段的時間可以連線回來，時間過了會被踢出房間，不可以一直
// 留在房間，否則會有 yablon 爆滿問題 (過多斷線玩家累積在房間內)
func (handler *ConnectionHandler) handleDisconnectAfterGracefulPeriod(uid core.Uid) {
	user, ok := handler.userGroup.Data[uid]
	if !ok {
		handler.logger.Info("user OnDisconnect graceful period end, but does not exist already", zap.String("uid", uid.String()))
		return
	}

	if user.IsConnected {
		handler.logger.Info("user OnDisconnect graceful period end, has connected back", zap.String("uid", uid.String()))
		return
	}

	handler.logger.Info("user OnDisconnect graceful period end, still not connected", zap.Object("user", user))

	if err := handler.userLeaveProcess(uid, commongrpc.KickoutReason_USER_REQUESTED); err != nil {
		handler.logger.Error("user OnDisconnect graceful period end, failed to run user leave process", zap.Error(err), zap.Object("user", user))
	}
}

func (handler *ConnectionHandler) userLeaveProcess(uid core.Uid, reason commongrpc.KickoutReason) error {
	user, ok := handler.userGroup.Data[uid]
	if !ok {
		handler.logger.Error("user OnLeave but does not exist", zap.String("uid", uid.String()), zap.Object("users", handler.userGroup))
		return status.Errorf(codes.NotFound, "user does not exist")
	}

	handler.logger.Info(
		"user leaving start",
		zap.String("uid", user.Uid.String()),
		util.DebugField(zap.Object("user", user)),
		util.DebugField(zap.Object("users", handler.userGroup)),
	)

	cashOutDone, err := handler.seatStatusService.StandUp(uid)
	if err != nil {
		return err
	}

	go func() {
		// Only destroy user after cash out done. Client will not actually leave
		// until kickout message is sent.
		<-cashOutDone
		handler.game.RunTask(func() {
			kickoutMSG := &txpokergrpc.Message{
				Event: &txpokergrpc.Event{
					Kickout: &commongrpc.Kickout{
						Reason: reason,
					},
				},
			}

			if err := handler.userService.Destroy(core.MessageTopic, kickoutMSG, uid); err != nil {
				handler.logger.Error("user leaving but failed to destroy user data", zap.Error(err), zap.String("uid", uid.String()), zap.Object("users", handler.userGroup))
			}

			handler.logger.Info(
				"user leaving done",
				zap.String("uid", user.Uid.String()),
				util.DebugField(zap.Object("user", user)),
				util.DebugField(zap.Object("users", handler.userGroup)),
			)
			handler.userService.BroadcastUpdate()
		})
	}()

	return nil
}

func (handler *ConnectionHandler) HandleRequest(req *core.Request) *core.Response {
	switch req.Msg.(type) {
	case *commongrpc.ResyncRequest:
		handler.resyncService.Send(req.Uid)
		return nil

	case *commongrpc.ChangeRoomRequest:
		handler.logger.Info(">changeRoom", zap.String("uid", req.Uid.String()))

		if err := handler.userLeaveProcess(req.Uid, commongrpc.KickoutReason_CHANGE_ROOM); err != nil {
			return &core.Response{Err: err}
		}
		return nil

	default:
		handler.logger.Debug("ignored not supported request", zap.Object("req", req))
		return nil
	}
}
