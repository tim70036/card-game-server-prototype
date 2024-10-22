package client

import (
	"context"
	"card-game-server-prototype/pkg/util"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	hasAuth bool
	idToken string
	logger  *zap.Logger
}

func ProvideAuthInterceptor(loggerFactory *util.LoggerFactory) *AuthInterceptor {
	return &AuthInterceptor{
		logger: loggerFactory.Create("AuthInterceptor"),
	}
}

func (a *AuthInterceptor) unaryInterceptor(
	ctx context.Context,
	method string,
	req,
	reply interface{},
	conn *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	if a.hasAuth {
		ctx = metadata.AppendToOutgoingContext(ctx, "id_token", a.idToken, "dummy", "fuckyou")
	}

	a.logger.Debug("", zap.String("method", method))
	err := invoker(ctx, method, req, reply, conn, opts...)
	return err
}

func (a *AuthInterceptor) streamInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	conn *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	if a.hasAuth {
		ctx = metadata.AppendToOutgoingContext(ctx, "id_token", a.idToken, "dummy", "fuckyou")
	}

	a.logger.Debug("", zap.String("method", method))
	s, err := streamer(ctx, desc, conn, method, opts...)
	if err != nil {
		return nil, err
	}
	return newWrappedStream(s), nil
}

// wrappedStream  wraps around the embedded grpc.ClientStream, and intercepts the RecvMsg and
// SendMsg method action.
type wrappedStream struct {
	grpc.ClientStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	// log.Printf("Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	// log.Printf("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.SendMsg(m)
}

func newWrappedStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStream{s}
}
