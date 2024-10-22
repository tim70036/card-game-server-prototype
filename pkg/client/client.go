package client

import (
	"context"
	"card-game-server-prototype/pkg/grpc"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

type Client struct {
	authInterceptor      *AuthInterceptor
	connectionClient     *ConnectionServiceClient
	messageServiceClient *MessageServiceClient
	actionServiceClient  txpokergrpc.ActionServiceClient

	logger *zap.Logger
}

func ProvideClient(
	authInterceptor *AuthInterceptor,
	connectionClient *ConnectionServiceClient,
	messageServiceClient *MessageServiceClient,
	loggerFactory *util.LoggerFactory,
) *Client {
	return &Client{
		authInterceptor:      authInterceptor,
		connectionClient:     connectionClient,
		actionServiceClient:  txpokergrpc.NewActionServiceClient(connectionClient.conn),
		messageServiceClient: messageServiceClient,
		logger:               loggerFactory.Create("Client"),
	}
}

func (c *Client) Run() {
	if err := c.connectionClient.Init(); err != nil {
		c.logger.Error("init failed", zap.Error(err))
		return
	}

	go c.connectionClient.Run()
	go c.messageServiceClient.Run()

	time.Sleep(2 * time.Second) // Wait for subsribe done before Enter()
	if _, err := c.connectionClient.Enter(context.Background(), &grpc.Empty{}); err != nil {
		c.logger.Error("enter failed", zap.Error(err))
		return
	}

	time.Sleep(1 * time.Second)
	if _, err := c.actionServiceClient.UpdateAutoTopUpSetting(context.Background(), &txpokergrpc.UpdateAutoTopUpSettingRequest{
		AutoTopUp:                 true,
		AutoTopUpThresholdPercent: 0.5,
		AutoTopUpChipPercent:      1,
	}); err != nil {
		c.logger.Error("play setting request failed", zap.Error(err))
		return
	}

	time.Sleep(1 * time.Second)
	for {
		if _, err := c.actionServiceClient.SitDown(context.Background(), &txpokergrpc.SitDownRequest{
			SeatId: int32(rand.Intn(9)),
		}); err != nil {
			c.logger.Error("sit down request failed", zap.Error(err))
			continue
		}
		break
	}

	time.Sleep(3 * time.Second)
	if _, err := c.actionServiceClient.BuyIn(context.Background(), &txpokergrpc.BuyInRequest{
		BuyInChip: 1000,
	}); err != nil {
		c.logger.Error("buy in request failed", zap.Error(err))
		return
	}

	// time.Sleep(3 * time.Second)
	// if _, err := c.actionServiceClient.TopUp(context.Background(), &txpokergrpc.TopUpRequest{
	// 	TopUpChip: 1000,
	// }); err != nil {
	// 	c.logger.Error("top up request failed", zap.Error(err))
	// 	return
	// }

	time.Sleep(10000 * time.Second)
	if _, err := c.connectionClient.Close(context.Background(), &grpc.Empty{}); err != nil {
		c.logger.Error("close failed", zap.Error(err))
		return
	}

	select {}
}

func (c *Client) Close() {
	if _, err := c.connectionClient.Close(context.Background(), &grpc.Empty{}); err != nil {
		c.logger.Error("close failed", zap.Error(err))
	}
}
