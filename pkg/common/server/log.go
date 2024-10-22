package server

import (
	"context"
	"fmt"
	"card-game-server-prototype/pkg/util"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type LogInterceptor struct {
	logger *zap.Logger
}

func ProvideLogInterceptor(
	connectionServiceServer *ConnectionServiceServer,
	loggerFactory *util.LoggerFactory,
) *LogInterceptor {
	return &LogInterceptor{
		logger: loggerFactory.Create("LogInterceptor"),
	}
}

func (i *LogInterceptor) UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	i.logger.Debug(
		fmt.Sprintf("=> %s", info.FullMethod),
		zap.Any("req", req),
	)

	resp, err := handler(ctx, req)
	if err != nil {
		i.logger.Warn(
			fmt.Sprintf("<= %s error", info.FullMethod),
			zap.Error(err),
			zap.Any("req", req),
			zap.Any("resp", resp),
		)
	} else {
		i.logger.Debug(
			fmt.Sprintf("<= %s", info.FullMethod),
			zap.Any("req", req),
			zap.Any("resp", resp),
		)
	}

	return resp, err
}

func (i *LogInterceptor) StreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	i.logger.Debug(
		fmt.Sprintf("=> %s", info.FullMethod),
	)

	err := handler(srv, stream)
	if err != nil {
		i.logger.Warn(
			fmt.Sprintf("<= %s error", info.FullMethod),
			zap.Error(err),
		)
	} else {
		i.logger.Debug(
			fmt.Sprintf("<= %s", info.FullMethod),
		)
	}

	return err
}
