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

type EmoteRpcServiceServer struct {
	commongrpc.UnimplementedEmoteRpcServiceServer

	gameNotifier core.GameNotifier
	msgBus       core.MsgBus
	logger       *zap.Logger
}

func ProvideEmoteServiceServer(
	gameNotifier core.GameNotifier,
	msgBus core.MsgBus,
	loggerFactory *util.LoggerFactory,
) *EmoteRpcServiceServer {
	return &EmoteRpcServiceServer{
		gameNotifier: gameNotifier,
		msgBus:       msgBus,
		logger:       loggerFactory.Create("EmoteRpcServiceServer"),
	}
}

func (s *EmoteRpcServiceServer) Subscribe(_ *grpc.Empty, stream commongrpc.EmoteRpcService_SubscribeServer) error {
	if err := HandleSubscribe(
		core.EmoteEventTopic,
		s.msgBus,
		stream,
		func(uid core.Uid, msg *commongrpc.EmoteEvent) (*commongrpc.EmoteEvent, error) {
			return msg, nil
		},
		s.logger,
	); err != nil {
		s.logger.Error("failed when handling stream", zap.Error(err))
		return err
	}

	return nil
}

func (s *EmoteRpcServiceServer) SendSticker(ctx context.Context, req *commongrpc.StickerRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *EmoteRpcServiceServer) SendPing(ctx context.Context, req *commongrpc.EmotePingRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *EmoteRpcServiceServer) handleFireForgetRequest(ctx context.Context, req proto.Message) (*grpc.Empty, error) {
	resp := HandleRequest(ctx, req, s.gameNotifier)
	if resp != nil && resp.Err != nil {
		s.logger.Error("handleFireForgetRequest failed", zap.Error(resp.Err))
		return nil, resp.Err
	}

	return &grpc.Empty{}, nil
}
