package session

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/constant"
	"time"
)

type ForceBuyInSession struct {
	buyInChip  int
	expireTime int64
}

func (f *ForceBuyInSession) GetBuyInChip() int {
	return f.buyInChip
}

func (f *ForceBuyInSession) GetExpireTime() time.Time {
	return time.Unix(f.expireTime, 0)
}

type ForceBuyInSessionGroup struct {
	data map[core.Uid]*ForceBuyInSession
}

func NewForceBuyInSessionGroup() *ForceBuyInSessionGroup {
	return &ForceBuyInSessionGroup{
		data: make(map[core.Uid]*ForceBuyInSession),
	}
}

func (m *ForceBuyInSessionGroup) Set(uid core.Uid, buyInChip int, now time.Time) {
	if _, ok := m.data[uid]; !ok {
		m.data[uid] = &ForceBuyInSession{}
	}

	m.data[uid].buyInChip = buyInChip
	m.data[uid].expireTime = now.Add(constant.CacheExpirationDuration).Unix()
}

func (m *ForceBuyInSessionGroup) Get(uid core.Uid) (*ForceBuyInSession, bool) {
	d, ok := m.data[uid]
	if !ok {
		return nil, false
	}

	if m.isExpired(time.Unix(d.expireTime, 0), time.Now()) {
		return nil, false
	}

	return d, true
}

func (m *ForceBuyInSessionGroup) CleanExpired() {
	now := time.Now()

	for uid, forceBuyIn := range m.data {
		expireTime := time.Unix(forceBuyIn.expireTime, 0)
		if m.isExpired(expireTime, now) || expireTime.IsZero() {
			delete(m.data, uid)
		}
	}
}

func (m *ForceBuyInSessionGroup) isExpired(t, now time.Time) bool {
	return t.Equal(now) || t.Before(now)
}
