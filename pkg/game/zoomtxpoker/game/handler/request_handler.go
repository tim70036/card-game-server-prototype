package handler

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/model"
	event2 "card-game-server-prototype/pkg/game/txpoker/type/event"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
)

type RequestHandler struct {
	core.BaseHandler

	eventGroup *model.EventGroup

	msgBus core.MsgBus
	logger *zap.Logger
}

func ProvideRequestHandler(
	msgBus core.MsgBus,
	eventGroup *model.EventGroup,
	loggerFactory *util.LoggerFactory,
) *RequestHandler {
	return &RequestHandler{
		msgBus:     msgBus,
		eventGroup: eventGroup,
		logger:     loggerFactory.Create("RequestHandler"),
	}
}

func (handler *RequestHandler) HandleRequest(req *core.Request) *core.Response {
	switch msg := req.Msg.(type) {
	case *commongrpc.EmotePingRequest:

		handler.msgBus.Broadcast(core.EmoteEventTopic, &commongrpc.EmoteEvent{
			EmotePing: &commongrpc.EmotePing{
				ItemId:    msg.ItemId,
				SenderUid: req.Uid.String(),
				TargetUid: msg.TargetUid,
			},
		})
		return nil

	case *commongrpc.StickerRequest:
		handler.eventGroup.Data[req.Uid] = append(handler.eventGroup.Data[req.Uid], &event2.Event{
			Type:   event2.UseSticker,
			Amount: 1,
		})

		handler.msgBus.Broadcast(core.EmoteEventTopic, &commongrpc.EmoteEvent{
			Sticker: &commongrpc.Sticker{
				Uid:       req.Uid.String(),
				StickerId: msg.StickerId,
			},
		})
		return nil

	default:
		handler.logger.Debug("ignored not supported request", zap.Object("req", req))
		return nil
	}
}
