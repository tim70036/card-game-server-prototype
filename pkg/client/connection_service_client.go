package client

import (
	"context"
	"crypto/tls"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/grpc/coregrpc"
	"card-game-server-prototype/pkg/util"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ConnectionServiceClient struct {
	coregrpc.ConnectionServiceClient

	clientCFG       *config.ClientConfig
	conn            *grpc.ClientConn
	authInterceptor *AuthInterceptor
	logger          *zap.Logger
}

func ProvideConnectionServiceClient(
	clientCFG *config.ClientConfig,
	authInterceptor *AuthInterceptor,
	loggerFactory *util.LoggerFactory,
) (*ConnectionServiceClient, error) {

	dialOptions := []grpc.DialOption{
		grpc.WithUnaryInterceptor(authInterceptor.unaryInterceptor),
		grpc.WithStreamInterceptor(authInterceptor.streamInterceptor),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
			Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
			PermitWithoutStream: true,             // send pings even without active streams
		}),
	}

	if *clientCFG.WithTLS {
		tlsCreds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(tlsCreds))
	} else {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(*clientCFG.Address, dialOptions...)
	if err != nil {
		return nil, err
	}

	return &ConnectionServiceClient{
		ConnectionServiceClient: coregrpc.NewConnectionServiceClient(conn),

		clientCFG:       clientCFG,
		conn:            conn,
		authInterceptor: authInterceptor,
		logger:          loggerFactory.Create("ConnectionServiceClient"),
	}, nil
}

func (c *ConnectionServiceClient) Init() error {
	c.logger.Info("connecting...", zap.Object("config", c.clientCFG))
	_, err := c.Connect(context.Background(), &coregrpc.ConnectRequest{IdToken: *c.clientCFG.Uid})
	if err != nil {
		c.logger.Error("connect failed", zap.Error(err))
		return err
	}

	c.authInterceptor.hasAuth = true
	c.authInterceptor.idToken = *c.clientCFG.Uid
	c.logger.Info("connect successfully")
	return nil
}

func (c *ConnectionServiceClient) Run() {
	pingFailedCount := 0
	for {
		if _, err := c.Ping(context.Background(), &emptypb.Empty{}); err != nil {
			pingFailedCount++
			c.logger.Error("ping failed", zap.Error(err))
			if pingFailedCount >= 3 {
				return
			}
		}
		time.Sleep(30 * time.Second)
	}
}
