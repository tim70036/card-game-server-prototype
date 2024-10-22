package server

import (
	"context"
	"card-game-server-prototype/pkg/common/constant"
	"card-game-server-prototype/pkg/core"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// grpc.ServerStream does not provide a way to modify its RPC context.
// The streaming interceptor therefore needs to implement the
// grpc.ServerStream interface and return a context with updated
// metadata. The easiest way to do this would be to create a type
// which embeds the grpc.ServerStream interface and overrides only the
// Context() method to return a context with updated metadata. The
// streaming interceptor would then pass this wrapped stream to the
// provided handler.
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *wrappedStream) Context() context.Context {
	return s.ctx
}

func GetUidFromContext(ctx context.Context) (core.Uid, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.InvalidArgument, "no incoming metadata in rpc context")
	}

	v, ok := md["uid"]
	if !ok || len(v) <= 0 {
		return "", status.Errorf(codes.InvalidArgument, "no uid in rpc context")
	}

	return core.Uid(v[0]), nil
}

func HandleRequest(ctx context.Context, req proto.Message, gameNotifier core.GameNotifier) *core.Response {
	uid, err := GetUidFromContext(ctx)
	if err != nil {
		return &core.Response{Err: err}
	}

	done := make(chan *core.Response, constant.PerNotificationBufferSize)
	gameNotifier.NotifyRequest() <- &core.RequestNotification{
		Request: &core.Request{Uid: uid, Msg: req},
		Done:    done,
	}

	return <-done
}

// todo: gpc stream 沒有特別處理 close 的情況。
//  如果 gs 活的夠久 or 人數夠多，memory leaking 是有機會吃掉太多資源。
//  但是目前沒有特別處理的必要，因為遊戲設計不會有太多人，近期也會加「定時重開房間」的機制。

type serverStream[P proto.Message] interface {
	grpc.ServerStream
	Send(P) error
}

func HandleSubscribe[S serverStream[P], P proto.Message](
	topic core.Topic,
	msgBus core.MsgBus,
	stream S,
	msgTransformer func(uid core.Uid, msg P) (P, error),
	logger *zap.Logger,
) error {
	uid, err := GetUidFromContext(stream.Context())
	if err != nil {
		return err
	}

	// It is safe to have a goroutine calling SendMsg and another goroutine
	// calling RecvMsg on the same stream at the same time, but it is not safe
	// to call SendMsg on the same stream in different goroutines.
	// https://pkg.go.dev/google.golang.org/grpc#ServerStream
	// https://github.com/grpc/grpc-go/blob/master/Documentation/concurrency.md#streams
	// https://github.com/grpc/grpc-go/issues/2094
	// https://github.com/grpc/grpc-go/issues/5393
	// https://stackoverflow.com/questions/72143797/grpc-server-blocked-on-sendmsg
	// https://github.com/grpc/grpc-go/issues/2427
	msgQueue := make(chan P, constant.PerConnectionBufferSize)

	// GrpcServerStream.Send() could indefinitely block if the client is not receiving messages or server concurrently
	// calling Send(). topicHandler() Could be called by multiple goroutines. Thus, we used a buffered channel msgQueue to
	// ensure that grpc.ServerStream.Send() will not be concurrently called. More importantly, we should make
	// topicHandler() non-blocking. Otherwise, the whole msgBus will be blocked and game server will be frozen.
	topicHandler := func(topicMsg P) {
		select {
		case msgQueue <- topicMsg:
			// message sent
		default:
			logger.Warn("msg dropped", zap.String("uid", uid.String()), zap.String("topic", string(topic)))
		}
	}

	unsub, err := msgBus.Subscribe(uid, topic, topicHandler)
	if err != nil {
		logger.Error("failed to subscribe msgBus", zap.String("uid", uid.String()), zap.String("topic", string(topic)), zap.Error(err))
		return status.Errorf(codes.Internal, "failed to subscribe msgBus, uid[%s] topic[%s] err[%s]", uid.String(), topic, err.Error())
	}

	defer func() {
		if err := unsub(); err != nil {
			logger.Error("failed to unsubscribe msgBus", zap.String("uid", uid.String()), zap.String("topic", string(topic)), zap.Error(err))
		}
	}()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case msg := <-msgQueue:
			finalMsg, err := msgTransformer(uid, msg)
			if err != nil {
				return err
			}

			if err := stream.Send(finalMsg); err != nil {
				return err
			}
		}
	}
}
