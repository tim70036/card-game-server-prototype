package server

import (
	"fmt"
	"card-game-server-prototype/pkg/common/constant"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/util"
	"github.com/google/uuid"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	inactiveTimeout = 40 * time.Second
)

type PeerFactory struct {
	loggerFactory *util.LoggerFactory
}

func ProvidePeerFactory(loggerFactory *util.LoggerFactory) *PeerFactory {
	return &PeerFactory{
		loggerFactory: loggerFactory,
	}
}

func (f *PeerFactory) Create(uid core.Uid, idToken string, onTimeout func(*Peer)) *Peer {
	peerId := uuid.NewString()[:8]
	return &Peer{
		Uid:       uid,
		id:        peerId,
		idToken:   idToken,
		hasPing:   true,
		onTimeout: onTimeout,
		close:     make(chan struct{}, constant.PerNotificationBufferSize),
		logger:    f.loggerFactory.Create(fmt.Sprintf("Peer[%v]", peerId)),
	}
}

type Peer struct {
	// Uid is used to identify which user a peer belongs to.
	// A user can have multiple peers if connect multiple times.
	// But only 1 peer has valid session at a time.
	Uid core.Uid

	// Peer id is used to identify a peer.
	id string

	// idToken is used to authenticate a peer. It's in jwt format.
	idToken string

	// hasPing is used to check if a peer is still active. will trigger onTimeout if not active.
	hasPing bool
	mu      sync.Mutex

	onTimeout func(*Peer)
	closeOnce sync.Once
	close     chan struct{}
	logger    *zap.Logger
}

func (p *Peer) Ping() {
	p.logger.Debug("ping", zap.Object("peer", p))
	p.mu.Lock()
	p.hasPing = true
	p.mu.Unlock()
}

func (p *Peer) Run() {
	ticker := time.NewTicker(inactiveTimeout)
	for {
		select {
		case <-ticker.C:
			p.logger.Debug("checking")
			p.mu.Lock()
			if p.hasPing {
				p.hasPing = false
			} else {
				p.logger.Debug("timeout", zap.Object("peer", p))
				p.onTimeout(p)
			}
			p.mu.Unlock()
			p.logger.Debug("checked")

		case <-p.close:
			return
		}
	}
}

func (p *Peer) Close() {
	p.logger.Debug("closing", zap.Object("peer", p))
	// Do nothing if peer is already in the process of closing.
	p.closeOnce.Do(func() { p.close <- struct{}{} })

}

func (p *Peer) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Uid", p.Uid.String())
	enc.AddString("id", p.id)
	enc.AddBool("hasPing", p.hasPing)
	return nil
}
