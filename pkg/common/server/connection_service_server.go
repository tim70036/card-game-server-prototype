package server

import (
	"context"
	"fmt"
	"card-game-server-prototype/pkg/common/constant"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/grpc"
	"card-game-server-prototype/pkg/grpc/coregrpc"
	"card-game-server-prototype/pkg/util"

	"sync"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ConnectionServiceServer struct {
	coregrpc.UnimplementedConnectionServiceServer

	gameNotifier core.GameNotifier
	testCFG      *config.TestConfig
	cfg          *config.Config

	// map[uid]*peer
	peers *sync.Map

	// map[idToken]uid
	// Session will be created when client first connect to server.
	// Each id token maps to a uid.
	// Session will then be used to authenticate subsequent request.
	sessions *sync.Map

	peerFactory *PeerFactory
	logger      *zap.Logger
}

func ProvideConnectionServiceServer(
	gameNotifier core.GameNotifier,
	cfg *config.Config,
	testCFG *config.TestConfig,
	peerFactory *PeerFactory,
	loggerFactory *util.LoggerFactory,
) *ConnectionServiceServer {

	server := &ConnectionServiceServer{
		gameNotifier: gameNotifier,
		cfg:          cfg,
		testCFG:      testCFG,
		peers:        &sync.Map{},
		sessions:     &sync.Map{},
		peerFactory:  peerFactory,
		logger:       loggerFactory.Create("ConnectionServiceServer"),
	}

	return server
}

// client rpc request Leave() → server deletes ms session → server send
// kickout → client rpc request Close() → client leave room / server deletes gs
// session

// 2 Cases that peer is consider closed:
// 1. Hasn't sent ping for a while
// 2. Client actively close by calling Close() rpc
// For both cases, 3 actions are needed:
// - Stop running peer
// - Notify app that peer disconnect
// - Remove peer from maps so that subsequent request from that peer will not success.
func (c *ConnectionServiceServer) Connect(ctx context.Context, req *coregrpc.ConnectRequest) (*grpc.Empty, error) {
	var uid core.Uid
	if *c.testCFG.NoAuth {
		uid = core.Uid(req.IdToken)
	} else {
		parsedToken, err := jwt.Parse(req.IdToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected jwt signing method: %v", token.Header["alg"])
			}
			return []byte(*c.cfg.JWTKey), nil
		})

		if err != nil || !parsedToken.Valid {
			c.logger.Warn("failed to parse id token", zap.Error(err))
			return nil, status.Errorf(codes.Unauthenticated, "invalid id_token")
		}

		uid = core.Uid(parsedToken.Claims.(jwt.MapClaims)["uid"].(string))
	}

	peer := c.peerFactory.Create(uid, req.IdToken, func(p *Peer) { c.removePeer(p) })
	c.logger.Debug("peer created", zap.Object("peer", peer))
	if actual, loaded := c.peers.LoadOrStore(uid, peer); loaded {
		oldPeer := actual.(*Peer)
		c.sessions.Delete(oldPeer.idToken) // TODO: What if oldPeer.idToken == req.IdToken? this could cause a race condition.
		oldPeer.Close()

		c.peers.Store(uid, peer)
		c.logger.Debug("old peer replaced", zap.Object("peer", peer), zap.Object("oldPeer", oldPeer))
	} else {
		done := make(chan error, constant.PerNotificationBufferSize)
		c.gameNotifier.NotifyConnect() <- &core.ConnectNotification{
			Uid:  uid,
			Done: done,
		}

		if err := <-done; err != nil {
			c.peers.Delete(uid)
			return nil, err
		}
	}

	c.sessions.Store(req.IdToken, uid)
	go peer.Run()
	c.logger.Debug("peer running", zap.Object("peer", peer))
	return &grpc.Empty{}, nil
}

func (c *ConnectionServiceServer) Enter(ctx context.Context, req *grpc.Empty) (*grpc.Empty, error) {
	uid, err := GetUidFromContext(ctx)
	if err != nil {
		return nil, err
	}

	done := make(chan error, constant.PerNotificationBufferSize)
	c.gameNotifier.NotifyEnter() <- &core.EnterNotification{
		Uid:  uid,
		Done: done,
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return &grpc.Empty{}, nil
}

func (c *ConnectionServiceServer) Leave(ctx context.Context, req *grpc.Empty) (*grpc.Empty, error) {
	uid, err := GetUidFromContext(ctx)
	if err != nil {
		return nil, err
	}

	done := make(chan error, constant.PerNotificationBufferSize)
	c.gameNotifier.NotifyLeave() <- &core.LeaveNotification{
		Uid:  uid,
		Done: done,
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return &grpc.Empty{}, nil
}

func (c *ConnectionServiceServer) Ping(ctx context.Context, req *grpc.Empty) (*coregrpc.PingResponse, error) {
	uid, err := GetUidFromContext(ctx)
	if err != nil {
		return nil, err
	}

	value, ok := c.peers.Load(uid)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "peer not found")
	}

	peer := value.(*Peer)
	peer.Ping()

	return &coregrpc.PingResponse{
		ServerTime: timestamppb.Now(),
	}, nil
}

func (c *ConnectionServiceServer) Close(ctx context.Context, req *grpc.Empty) (*grpc.Empty, error) {
	uid, err := GetUidFromContext(ctx)
	if err != nil {
		return nil, err
	}

	value, ok := c.peers.Load(uid)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "peer not found")
	}

	peer := value.(*Peer)
	if err := c.removePeer(peer); err != nil {
		return nil, err
	}

	return &grpc.Empty{}, nil
}

func (c *ConnectionServiceServer) removePeer(peer *Peer) error {
	c.logger.Debug("removing peer", zap.Object("peer", peer))
	peer.Close()
	c.peers.Delete(peer.Uid)
	c.sessions.Delete(peer.idToken)

	done := make(chan error, constant.PerNotificationBufferSize)
	c.gameNotifier.NotifyDisconnect() <- &core.DisconnectNotification{
		Uid:  peer.Uid,
		Done: done,
	}

	if err := <-done; err != nil {
		return err
	}

	return nil
}
