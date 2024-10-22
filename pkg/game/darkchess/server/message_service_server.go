package server

import (
	commonserver "card-game-server-prototype/pkg/common/server"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 發送給前端的訊息。

type MessageServiceServer struct {
	gamegrpc.UnimplementedMessageServiceServer

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

func (s *MessageServiceServer) SubscribeGameState(_ *grpc.Empty, stream gamegrpc.MessageService_SubscribeGameStateServer) error {
	if err := commonserver.HandleSubscribe(
		core.GameStateTopic,
		s.msgBus,
		stream,
		func(uid core.Uid, msg *gamegrpc.GameState) (*gamegrpc.GameState, error) {
			return msg, nil
		},
		s.logger,
	); err != nil {
		s.logger.Error("failed when handling stream", zap.Error(err))
		return err
	}

	return nil
}

func (s *MessageServiceServer) SubscribeModel(_ *grpc.Empty, stream gamegrpc.MessageService_SubscribeModelServer) error {
	if err := commonserver.HandleSubscribe(
		core.ModelTopic,
		s.msgBus,
		stream,
		func(uid core.Uid, msg *gamegrpc.Model) (*gamegrpc.Model, error) {
			finalMsg := msg

			if finalMsg != nil && finalMsg.Board != nil {
				finalMsg = proto.Clone(msg).(*gamegrpc.Model)

				for i := range finalMsg.Board.Cells {
					if finalMsg.Board.Cells[i].IsRevealed {
						continue
					}

					finalMsg.Board.Cells[i].Piece = commongrpc.CnChessPiece_CN_CHESS_PIECE_INVALID
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
func (s *MessageServiceServer) SubscribeEvent(_ *grpc.Empty, stream gamegrpc.MessageService_SubscribeEventServer) error {
	if err := commonserver.HandleSubscribe(
		core.EventTopic,
		s.msgBus,
		stream,
		func(uid core.Uid, msg *gamegrpc.Event) (*gamegrpc.Event, error) {
			return msg, nil
		},
		s.logger,
	); err != nil {
		s.logger.Error("failed when handling stream", zap.Error(err))
		return err
	}

	return nil
}
