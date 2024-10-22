package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/constant"
	"go.uber.org/zap/zapcore"
	"time"
)

// 不是很重要的功能，不必太在意 delete expired data 時機的準確性。

type ForceBuyIn struct {
	buyInChip       int
	expireTimestamp int64
}

func (f *ForceBuyIn) GetBuyInChip() int {
	return f.buyInChip
}

func (f *ForceBuyIn) GetExpireTime() time.Time {
	return time.Unix(f.expireTimestamp, 0)
}

func (f *ForceBuyIn) isExpired(now time.Time) bool {
	return f.expireTimestamp < now.Unix()
}

func (f *ForceBuyIn) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("buyInChip", f.buyInChip)
	enc.AddTime("expireTimestamp", time.Unix(f.expireTimestamp, 0))
	return nil
}

type ForceBuyInGroup struct {
	data map[core.Uid]*ForceBuyIn
}

func ProvideForceBuyInGroup() *ForceBuyInGroup {
	return &ForceBuyInGroup{
		data: make(map[core.Uid]*ForceBuyIn),
	}
}

func (g *ForceBuyInGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for uid, forceBuyIn := range g.data {
		_ = enc.AddObject(uid.String(), forceBuyIn)
	}
	return nil
}

func (g *ForceBuyInGroup) Set(uid core.Uid, buyInChip int, now time.Time) {
	if _, ok := g.data[uid]; !ok {
		g.data[uid] = &ForceBuyIn{}
	}

	g.data[uid].buyInChip = buyInChip
	g.data[uid].expireTimestamp = now.Add(constant.CacheExpirationDuration).Unix()
}

func (g *ForceBuyInGroup) Delete(uid core.Uid) {
	delete(g.data, uid)
}

func (g *ForceBuyInGroup) Get(uid core.Uid) (*ForceBuyIn, bool) {
	d, ok := g.data[uid]
	if !ok {
		return nil, false
	}

	return d, true
}

func (g *ForceBuyInGroup) CleanExpired() {
	now := time.Now()

	for uid, forceBuyIn := range g.data {
		if forceBuyIn.isExpired(now) {
			delete(g.data, uid)
		}
	}
}
