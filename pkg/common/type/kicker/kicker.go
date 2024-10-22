package kicker

import (
	"card-game-server-prototype/pkg/core"
	"sync"
)

type KickedList struct {
	mu   *sync.RWMutex
	data map[core.Uid]struct{}
}

func NewKickedList() *KickedList {
	return &KickedList{
		mu:   &sync.RWMutex{},
		data: map[core.Uid]struct{}{},
	}
}

func (l *KickedList) Add(uid core.Uid) {
	l.mu.Lock()
	l.data[uid] = struct{}{}
	l.mu.Unlock()
}

func (l *KickedList) Remove(uid core.Uid) {
	l.mu.Lock()
	delete(l.data, uid)
	l.mu.Unlock()
}

func (l *KickedList) Has(uid core.Uid) bool {
	l.mu.RLock()
	_, ok := l.data[uid]
	l.mu.RUnlock()
	return ok
}

func (l *KickedList) Get() core.UidList {
	l.mu.RLock()
	list := make(core.UidList, 0, len(l.data))
	for uid := range l.data {
		list = append(list, uid)
	}
	l.mu.RUnlock()
	return list
}
