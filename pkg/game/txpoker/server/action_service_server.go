package server

import (
	"context"
	commonserver "card-game-server-prototype/pkg/common/server"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type ActionServiceServer struct {
	txpokergrpc.UnimplementedActionServiceServer

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

func (s *ActionServiceServer) Resync(ctx context.Context, req *commongrpc.ResyncRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) Ready(ctx context.Context, req *commongrpc.ReadyRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) StartGame(ctx context.Context, req *commongrpc.StartGameRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) Kick(ctx context.Context, req *commongrpc.KickRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) ChangeRoom(ctx context.Context, req *commongrpc.ChangeRoomRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) StandUp(ctx context.Context, req *txpokergrpc.StandUpRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) SitDown(ctx context.Context, req *txpokergrpc.SitDownRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) BuyIn(ctx context.Context, req *txpokergrpc.BuyInRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) SitOut(ctx context.Context, req *txpokergrpc.SitOutRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) TopUp(ctx context.Context, req *txpokergrpc.TopUpRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) Fold(ctx context.Context, req *txpokergrpc.FoldRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) Check(ctx context.Context, req *txpokergrpc.CheckRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) Bet(ctx context.Context, req *txpokergrpc.BetRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) Call(ctx context.Context, req *txpokergrpc.CallRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) Raise(ctx context.Context, req *txpokergrpc.RaiseRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) AllIn(ctx context.Context, req *txpokergrpc.AllInRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) ShowFold(ctx context.Context, req *txpokergrpc.ShowFoldRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) UpdateWaitBBSetting(ctx context.Context, req *txpokergrpc.UpdateWaitBBSettingRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) UpdateAutoTopUpSetting(ctx context.Context, req *txpokergrpc.UpdateAutoTopUpSettingRequest) (*grpc.Empty, error) {
	return s.handleFireForgetRequest(ctx, req)
}

func (s *ActionServiceServer) handleFireForgetRequest(ctx context.Context, req proto.Message) (*grpc.Empty, error) {
	resp := commonserver.HandleRequest(ctx, req, s.gameNotifier)
	if resp != nil && resp.Err != nil {
		s.logger.Error("handleFireForgetRequest failed", zap.Error(resp.Err))
		return nil, resp.Err
	}

	return &grpc.Empty{}, nil
}

func (s *ActionServiceServer) ForceBuyIn(ctx context.Context, req *txpokergrpc.ForceBuyInRequest) (*txpokergrpc.ForceBuyInResponse, error) {
	resp := commonserver.HandleRequest(ctx, req, s.gameNotifier)
	if resp != nil {
		if resp.Err != nil {
			s.logger.Error("failed to handle ForceBuyIn request", zap.Error(resp.Err))
		}

		return resp.Msg.(*txpokergrpc.ForceBuyInResponse), resp.Err
	}

	return nil, status.Errorf(codes.Internal, "interal error to handle ForceBuyIn request")
}
