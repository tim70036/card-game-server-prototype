package server

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 接受前端(client)傳來的動作。

type ChatServiceServer struct {
	commongrpc.UnimplementedChatServiceServer

	gameNotifier core.GameNotifier
	msgBus       core.MsgBus
	logger       *zap.Logger
}

func ProvideChatServiceServer(
	gameNotifier core.GameNotifier,
	msgBus core.MsgBus,
	loggerFactory *util.LoggerFactory,
) *ChatServiceServer {
	return &ChatServiceServer{
		gameNotifier: gameNotifier,
		msgBus:       msgBus,
		logger:       loggerFactory.Create("ChatServiceServer"),
	}
}

func (s *ChatServiceServer) Subscribe(_ *grpc.Empty, stream commongrpc.ChatService_SubscribeServer) error {
	if err := HandleSubscribe(
		core.ChatTopic,
		s.msgBus,
		stream,
		func(uid core.Uid, msg *commongrpc.ChatMessage) (*commongrpc.ChatMessage, error) {
			return msg, nil
		},
		s.logger,
	); err != nil {
		s.logger.Error("failed when handling stream", zap.Error(err))
		return err
	}

	return nil
}

func (s *ChatServiceServer) Send(ctx context.Context, req *commongrpc.ChatRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ChatServiceServer) handleFireForgetRequest(ctx context.Context, req proto.Message) (*grpc.Empty, error) {
	resp := HandleRequest(ctx, req, s.gameNotifier)
	if resp != nil && resp.Err != nil {
		s.logger.Error("handleFireForgetRequest failed", zap.Error(resp.Err))
		return nil, resp.Err
	}

	return &grpc.Empty{}, nil
}
