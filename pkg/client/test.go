package client

import (
	"context"
	"io"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SubscribeStream[M proto.Message] interface {
	Recv() (M, error)
}

type SubscribeClient[S SubscribeStream[M], M proto.Message] interface {
	Subscribe(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (S, error)
	Logger() *zap.Logger
}

func RunSubscribeClient[M proto.Message](client SubscribeClient[SubscribeStream[M], M]) {
	stream, err := client.Subscribe(context.Background(), &emptypb.Empty{})
	if err != nil {
		client.Logger().Error("subscribe error", zap.Error(err))
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			client.Logger().Error("recv failed", zap.Error(err))
			break
		}

		client.Logger().Info("recv", zap.Any("msg", msg))
	}
}
