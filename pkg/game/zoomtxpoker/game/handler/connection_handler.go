package handler

import (
	commonmodel "card-game-server-prototype/pkg/common/model"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/model"
	service2 "card-game-server-prototype/pkg/game/txpoker/service"
	"card-game-server-prototype/pkg/game/txpoker/type/seatstatus"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
)

type ConnectionHandler struct {
	core.BaseHandler
	game            core.Game
	userGroup       *commonmodel.UserGroup
	seatStatusGroup *model.SeatStatusGroup

	userService       commonservice.UserService
	seatStatusService service2.SeatStatusService
	resyncService     *service2.ResyncService
	msgBus            core.MsgBus
	logger            *zap.Logger
}

func ProvideConnectionHandler(
	game core.Game,
	userGroup *commonmodel.UserGroup,
	seatStatusGroup *model.SeatStatusGroup,

	userService commonservice.UserService,
	seatStatusService service2.SeatStatusService,
	resyncService *service2.ResyncService,
	msgBus core.MsgBus,
	loggerFactory *util.LoggerFactory,
) *ConnectionHandler {
	return &ConnectionHandler{
		game:            game,
		userGroup:       userGroup,
		seatStatusGroup: seatStatusGroup,

		userService:       userService,
		seatStatusService: seatStatusService,
		resyncService:     resyncService,
		msgBus:            msgBus,
		logger:            loggerFactory.Create("ConnectionHandler"),
	}
}

func (handler *ConnectionHandler) HandleEnter(uid core.Uid) error {
	handler.logger.Info("user OnEnter",
		zap.String("uid", uid.String()),
	)

	// Resync user data to user
	handler.resyncService.Send(uid)

	return nil
}

func (handler *ConnectionHandler) HandleLeave(uid core.Uid) error {
	handler.logger.Info("user OnLeave",
		zap.String("uid", uid.String()),
		zap.Int("userCount", len(handler.userGroup.Data)),
	)
	if _, err := handler.seatStatusService.StandUp(uid); err != nil {
		handler.logger.Error("failed to stand up for leave", zap.Error(err))
		return err
	}

	return nil
}

func (handler *ConnectionHandler) HandleDisconnect(uid core.Uid) error {
	user, ok := handler.userGroup.Data[uid]
	if !ok {
		// It's normal if user does not exist in user group. EX:
		// client rpc request Leave() → server deletes ms session →
		// server send kickout → client rpc request Close() → client
		// leave room / server deletes gs session
		handler.logger.Info("user OnDisconnect ignored since does not exist", zap.String("uid", uid.String()), zap.Object("userGroup", handler.userGroup))
		return nil
	}

	// This is the case that client disconnect due to ping timeout.
	user.IsConnected = false

	handler.logger.Info(
		"user OnDisconnect",
		zap.Object("user", user),
		zap.Int("userCount", len(handler.userGroup.Data)),
		util.DebugField(zap.Object("userGroup", handler.userGroup)),
	)

	return nil
}

func (handler *ConnectionHandler) HandleRequest(req *core.Request) *core.Response {
	switch req.Msg.(type) {
	case *commongrpc.ResyncRequest:
		handler.resyncService.Send(req.Uid)
		return nil

	case *txpokergrpc.FoldRequest:
		handler.logger.Info("req:FoldRequest", zap.Object("req", req))

		// 阻擋來自於 txpoker base actor 的 auto fold request 觸發的 standUp;
		if s, ok := handler.seatStatusGroup.Status[req.Uid]; ok && s != nil {
			seatStatusState := s.FSM.MustState().(seatstatus.SeatStatusState)
			if seatStatusState == seatstatus.StandingState {
				handler.logger.Info("ignored fold request since user is standing", zap.Object("status", s))
				return nil
			}
		}

		if _, err := handler.seatStatusService.StandUp(req.Uid); err != nil {
			handler.logger.Error("failed to stand up for fold request", zap.Error(err))
			return &core.Response{Err: err}
		}

		return nil

	default:
		handler.logger.Debug("ignored not supported request", zap.Object("req", req))
		return nil
	}
}
