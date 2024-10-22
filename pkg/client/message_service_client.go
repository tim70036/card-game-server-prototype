package client

import (
	"context"
	"card-game-server-prototype/pkg/grpc"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"
	"io"
	"log"

	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

type MessageServiceClient struct {
	txpokergrpc.MessageServiceClient
	logger *zap.Logger
}

func ProvideMessageServiceClient(
	connectionClient *ConnectionServiceClient,
	loggerFactory *util.LoggerFactory,
) *MessageServiceClient {
	return &MessageServiceClient{
		MessageServiceClient: txpokergrpc.NewMessageServiceClient(connectionClient.conn),
		logger:               loggerFactory.Create("MessageServiceClient"),
	}
}

func (c *MessageServiceClient) Run() {
	stream, err := c.Subscribe(context.Background(), &grpc.Empty{})
	if err != nil {
		c.logger.Error("subscribe error", zap.Error(err))
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			c.logger.Error("recv failed", zap.Error(err))
			break
		}

		log.Printf("recv %s", protojson.Format(msg))
	}
}
