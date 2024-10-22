package server

import (
	"context"
	"fmt"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/util"
	"github.com/golang-jwt/jwt/v5"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	connectionServiceServer *ConnectionServiceServer
	cfg                     *config.Config
	testCFG                 *config.TestConfig
	logger                  *zap.Logger
}

func ProvideAuthInterceptor(
	connectionServiceServer *ConnectionServiceServer,
	cfg *config.Config,
	testCFG *config.TestConfig,
	loggerFactory *util.LoggerFactory,
) *AuthInterceptor {
	return &AuthInterceptor{
		connectionServiceServer: connectionServiceServer,
		cfg:                     cfg,
		testCFG:                 testCFG,
		logger:                  loggerFactory.Create("AuthInterceptor"),
	}
}

func (i *AuthInterceptor) UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "no incoming metadata in rpc context")
	}

	if info.FullMethod != "/common.ConnectionService/Connect" {
		uid, err := i.authSession(md)
		if err != nil {
			return nil, err
		}

		md.Append("uid", uid.String())
		ctx = metadata.NewIncomingContext(ctx, md)
	}

	return handler(ctx, req)
}

func (i *AuthInterceptor) StreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.InvalidArgument, "no incoming metadata in rpc context")
	}

	uid, err := i.authSession(md)
	if err != nil {
		return err
	}

	md.Append("uid", uid.String())
	ctx := metadata.NewIncomingContext(stream.Context(), md)

	return handler(srv, &wrappedStream{stream, ctx})
}

func (i *AuthInterceptor) authSession(md metadata.MD) (core.Uid, error) {
	if *i.testCFG.NoAuth {
		return core.Uid(md["id_token"][0]), nil
	}

	v, ok := md["id_token"]
	if !ok || len(v) <= 0 {
		return "", status.Errorf(codes.Unauthenticated, "id_token is not provided")
	}

	idToken := v[0]

	// 從 jwt 解出 uid。（不影響原有流程）

	parsedToken, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected jwt signing method: %v", token.Header["alg"])
		}
		return []byte(*i.cfg.JWTKey), nil
	})

	var uidFromJwt core.Uid
	if err != nil || !parsedToken.Valid {
		i.logger.Warn("invalid id_token",
			zap.Error(err),
			zap.String("id_token", idToken),
		)
	} else if err == nil && parsedToken.Valid {
		uidFromJwt = core.Uid(parsedToken.Claims.(jwt.MapClaims)["uid"].(string))
	}

	value, ok := i.connectionServiceServer.sessions.Load(idToken)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "session not found, uid: %s", uidFromJwt)
	}

	return value.(core.Uid), nil
}
