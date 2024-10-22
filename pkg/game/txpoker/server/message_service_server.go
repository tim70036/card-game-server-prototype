package server

import (
	commonserver "card-game-server-prototype/pkg/common/server"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/type/fold"
	"card-game-server-prototype/pkg/grpc"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type MessageServiceServer struct {
	txpokergrpc.UnimplementedMessageServiceServer

	msgBus core.MsgBus
	logger *zap.Logger
}

func ProvideMessageServiceServer(
	msgBus core.MsgBus,
	loggerFactory *util.LoggerFactory,
) *MessageServiceServer {
	return &MessageServiceServer{
		msgBus: msgBus,
		logger: loggerFactory.Create("MessageServiceServer"),
	}
}

func (s *MessageServiceServer) Subscribe(req *grpc.Empty, stream txpokergrpc.MessageService_SubscribeServer) error {
	if err := commonserver.HandleSubscribe(
		core.MessageTopic,
		s.msgBus,
		stream,
		func(uid core.Uid, msg *txpokergrpc.Message) (*txpokergrpc.Message, error) {
			finalMsg := msg
			if finalMsg.Model != nil && (finalMsg.Model.PlayerGroup != nil) {
				finalMsg = proto.Clone(msg).(*txpokergrpc.Message)
				if finalMsg.Model.PlayerGroup != nil {
					s.maskPocketCards(uid, finalMsg.Model.PlayerGroup)
				}

			}
			return finalMsg, nil
		},
		s.logger,
	); err != nil {
		s.logger.Error("failed when handling stream", zap.Error(err))
		return err
	}

	return nil
}

func (s *MessageServiceServer) maskPocketCards(targetUid core.Uid, msg *txpokergrpc.PlayerGroup) {
	for uid, player := range msg.Players {
		if player.ShowFoldType != int32(fold.ShowNone) {
			continue
		}

		if uid != targetUid.String() {
			player.PocketCards = player.PocketCards[:0]
		}
	}
}
