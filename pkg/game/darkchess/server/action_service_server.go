package server

import (
	"context"
	commonserver "card-game-server-prototype/pkg/common/server"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 接受前端(client)傳來的動作。

type ActionServiceServer struct {
	gamegrpc.UnimplementedActionServiceServer

	gameNotifier core.GameNotifier
	msgBus       core.MsgBus
	logger       *zap.Logger
}

func ProvideActionServiceServer(
	gameNotifier core.GameNotifier,
	msgBus core.MsgBus,
	loggerFactory *util.LoggerFactory,
) *ActionServiceServer {
	return &ActionServiceServer{
		gameNotifier: gameNotifier,
		msgBus:       msgBus,
		logger:       loggerFactory.Create("ActionServiceServer"),
	}
}

func (s *ActionServiceServer) handleFireForgetRequest(ctx context.Context, req proto.Message) (*grpc.Empty, error) {
	resp := commonserver.HandleRequest(ctx, req, s.gameNotifier)
	if resp != nil && resp.Err != nil {
		s.logger.Error("handleFireForgetRequest failed", zap.Error(resp.Err))
		return nil, resp.Err
	}

	return &grpc.Empty{}, nil
}

func (s *ActionServiceServer) Kick(ctx context.Context, req *commongrpc.KickRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) Resync(ctx context.Context, req *commongrpc.ResyncRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) SkipScoreboard(ctx context.Context, req *gamegrpc.SkipScoreboardRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) UpdatePlaySetting(ctx context.Context, req *gamegrpc.UpdatePlaySettingRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) Ready(ctx context.Context, req *commongrpc.ReadyRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) StartGame(ctx context.Context, req *commongrpc.StartGameRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) Pick(ctx context.Context, req *gamegrpc.PickRequest) (*emptypb.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) Reveal(ctx context.Context, req *gamegrpc.RevealRequest) (*emptypb.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) Move(ctx context.Context, req *gamegrpc.MoveRequest) (*emptypb.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) Capture(ctx context.Context, req *gamegrpc.CaptureRequest) (*emptypb.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) Surrender(ctx context.Context, req *gamegrpc.SurrenderRequest) (*emptypb.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) ClaimDraw(ctx context.Context, req *gamegrpc.ClaimDrawRequest) (*emptypb.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) AnswerDraw(ctx context.Context, req *gamegrpc.AnswerDrawRequest) (*emptypb.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) AskExtraSeconds(ctx context.Context, req *gamegrpc.AskExtraSecondsRequest) (*emptypb.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) AddAi(ctx context.Context, req *commongrpc.AddAiRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}
