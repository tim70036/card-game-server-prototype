package server

import (
	"crypto/tls"
	"fmt"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/util"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type GrpcServer struct {
	*grpc.Server
	serverCFG *config.ServerConfig
	logger    *zap.Logger
}

func ProvideGrpcServer(
	serverCFG *config.ServerConfig,
	authInterceptor *AuthInterceptor,
	logInterceptor *LogInterceptor,
	loggerFactory *util.LoggerFactory,
) (*GrpcServer, error) {
	logger := loggerFactory.Create("GrpcServer")
	grpcOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(logInterceptor.UnaryInterceptor, authInterceptor.UnaryInterceptor),
		grpc.ChainStreamInterceptor(logInterceptor.StreamInterceptor, authInterceptor.StreamInterceptor),

		// https://github.com/grpc/grpc-go/blob/master/Documentation/keepalive.md
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             1 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
			PermitWithoutStream: true,            // Allow pings even when there are no active streams
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			// MaxConnectionIdle: 15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
			// MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
			// MaxConnectionAgeGrace: 5 * time.Second, // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
			Time:    5 * time.Second, // Ping the client if it is idle for x seconds to ensure the connection is still active
			Timeout: 5 * time.Second, // Wait x second for the ping ack before assuming the connection is dead
		}),
	}

	if *serverCFG.WithTLS {
		cert, err := tls.X509KeyPair([]byte(*serverCFG.TLSCert), []byte(*serverCFG.TLSKey))
		if err != nil {
			logger.Error("failed to read tls credential", zap.Error(err), zap.Object("cfg", serverCFG))
			return nil, err
		}
		grpcOptions = append(grpcOptions, grpc.Creds(credentials.NewServerTLSFromCert(&cert)))
	}

	return &GrpcServer{
		Server:    grpc.NewServer(grpcOptions...),
		serverCFG: serverCFG,
		logger:    logger,
	}, nil
}

func (s *GrpcServer) Run() {
	s.logger.Info("start gRPC server", zap.Object("serverCFG", s.serverCFG))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", *s.serverCFG.Port))
	if err != nil {
		s.logger.Fatal("failed to net listen", zap.Error(err), zap.Object("cfg", s.serverCFG))
	}

	s.logger.Debug("grpc server start serving", zap.Int("port", *s.serverCFG.Port))
	if err := s.Serve(lis); err != nil {
		s.logger.Fatal("failed to serve grpc server", zap.Error(err), zap.Object("cfg", s.serverCFG))
	}
}
