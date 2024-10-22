package main

import (
	"card-game-server-prototype/pkg/util"
	"time"

	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	coresdk "agones.dev/agones/pkg/sdk"
	sdk "agones.dev/agones/sdks/go"
	"go.uber.org/zap"
)

type Agones struct {
	WaitAllocated  chan *coresdk.GameServer
	gameServerInfo *coresdk.GameServer
	agonesSDK      *sdk.SDK
	logger         *zap.Logger
}

func ProvideAgones(loggerFactory *util.LoggerFactory) (*Agones, error) {
	return &Agones{
		WaitAllocated:  make(chan *coresdk.GameServer, 4),
		gameServerInfo: nil,
		agonesSDK:      nil,
		logger:         loggerFactory.Create("AgonesClient"),
	}, nil
}

func (a *Agones) Shutdown() {
	if a.agonesSDK == nil {
		a.logger.Info("skip shutdown because agones sidecar is not connected")
		return
	}

	if err := a.agonesSDK.Shutdown(); err != nil {
		a.logger.Error("fail notify shutdown to agones sidecar", zap.Error(err))
	}
	a.logger.Info("notified shutdown to agones sidecar")
}

func (a *Agones) Run() {
	var err error
	a.agonesSDK, err = sdk.NewSDK()
	if err != nil {
		a.logger.Fatal("cannot connect to agones sidecar", zap.Error(err))
	}
	a.logger.Info("connected to agones sidecar")

	a.gameServerInfo, err = a.agonesSDK.GameServer()
	if err != nil {
		a.logger.Fatal("cannot retrieve agones game server info", zap.Error(err))
	}
	a.logger.Info("retrieved agones game server info", zap.Any("gameServerInfo", a.gameServerInfo))

	if err := a.agonesSDK.WatchGameServer(a.onGameServerInfoUpdated); err != nil {
		a.logger.Fatal("cannot watch events from agones sidecar", zap.Error(err))
	}

	if err := a.agonesSDK.Ready(); err != nil {
		a.logger.Fatal("fail notify ready to agones sidecar", zap.Error(err))
	}
	a.logger.Info("notified ready to agones sidecar")

	go a.pingLoop()
}

func (a *Agones) onGameServerInfoUpdated(gameServerInfo *coresdk.GameServer) {
	if gameServerInfo.Status.State != a.gameServerInfo.Status.State {
		a.logger.Info("game server state changed",
			zap.Any("gameServerInfo", gameServerInfo),
			zap.Any("prevGameServerInfo", a.gameServerInfo),
		)

		if gameServerInfo.Status.State == string(agonesv1.GameServerStateAllocated) {
			a.WaitAllocated <- gameServerInfo
		}
	}

	a.gameServerInfo = gameServerInfo
}

func (a *Agones) pingLoop() {
	for {
		time.Sleep(time.Second * 5)
		err := a.agonesSDK.Health()
		if err != nil {
			a.logger.Error("cannot do health ping to agones sidecar", zap.Error(err))
		}
	}
}
