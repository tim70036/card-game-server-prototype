package core

import (
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"reflect"
	"sync"
)

type UnsubscribeFunc func() error

type handlersMap map[Topic]map[Uid]map[uuid.UUID]*handler

type handler struct {
	callback reflect.Value
	queue    chan []reflect.Value
}

// rawMessageBus implements publish/subscribe messaging paradigm
type rawMessageBus struct {
	handlerQueueSize int
	mtx              sync.RWMutex
	handlers         handlersMap
	logger           *zap.Logger
}

// Publish publishes arguments to the given topic subscribers
// Publish block only when the buffer of one of the subscribers is full.
func (bus *rawMessageBus) Broadcast(topic Topic, args ...interface{}) {
	rArgs := buildHandlerArgs(args)

	bus.mtx.RLock()
	defer bus.mtx.RUnlock()

	if topicUsers, ok := bus.handlers[topic]; ok {
		for _, userHandlers := range topicUsers {
			for _, h := range userHandlers {
				h.queue <- rArgs
			}
		}
	}
}

func (bus *rawMessageBus) Unicast(uid Uid, topic Topic, args ...interface{}) {
	rArgs := buildHandlerArgs(args)

	bus.mtx.RLock()
	defer bus.mtx.RUnlock()

	if topicUsers, ok := bus.handlers[topic]; ok {
		if userHandlers, ok := topicUsers[uid]; ok {
			for _, h := range userHandlers {
				h.queue <- rArgs
			}
		}
	}
}

// Subscribe subscribes to the given topic
func (bus *rawMessageBus) Subscribe(uid Uid, topic Topic, fn interface{}) (UnsubscribeFunc, error) {
	bus.logger.Debug(
		"subscribe",
		zap.String("uid", uid.String()),
		zap.String("topic", string(topic)),
		zap.Int("fn", int(reflect.ValueOf(fn).Pointer())),
	)

	if err := isValidHandler(fn); err != nil {
		return nil, err
	}

	handlerId, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	h := &handler{
		callback: reflect.ValueOf(fn),
		queue:    make(chan []reflect.Value, bus.handlerQueueSize),
	}

	go func() {
		for args := range h.queue {
			h.callback.Call(args)
		}
	}()

	unsub := func() error {
		bus.logger.Debug("unsubscribe",
			zap.String("uid", uid.String()),
			zap.String("topic", string(topic)),
			zap.Int("fn", int(reflect.ValueOf(fn).Pointer())),
		)

		bus.mtx.Lock()
		defer bus.mtx.Unlock()

		if _, ok := bus.handlers[topic]; !ok {
			return fmt.Errorf("topic %s doesn't exist", topic)
		}

		if _, ok := bus.handlers[topic][uid]; !ok {
			return fmt.Errorf("uid %s doesn't exist on topic %v", uid, topic)
		}

		if _, ok := bus.handlers[topic][uid][handlerId]; !ok {
			return fmt.Errorf("handler %v doesn't exist on topic %v for uid %v", handlerId, topic, uid)
		}

		close(bus.handlers[topic][uid][handlerId].queue)

		if len(bus.handlers[topic][uid]) == 1 {
			if len(bus.handlers[topic]) == 1 {
				delete(bus.handlers[topic], uid)
				delete(bus.handlers, topic)
			} else {
				delete(bus.handlers[topic], uid)
			}
		} else {
			delete(bus.handlers[topic][uid], handlerId)
		}

		return nil
	}

	bus.mtx.Lock()
	defer bus.mtx.Unlock()

	if bus.handlers[topic] == nil {
		bus.handlers[topic] = make(map[Uid]map[uuid.UUID]*handler)
	}

	if bus.handlers[topic][uid] == nil {
		bus.handlers[topic][uid] = make(map[uuid.UUID]*handler)
	}

	bus.handlers[topic][uid][handlerId] = h
	return unsub, nil
}

func isValidHandler(fn interface{}) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("%s is not a reflect.Func", reflect.TypeOf(fn))
	}

	return nil
}

func buildHandlerArgs(args []interface{}) []reflect.Value {
	reflectedArgs := make([]reflect.Value, 0)

	for _, arg := range args {
		reflectedArgs = append(reflectedArgs, reflect.ValueOf(arg))
	}

	return reflectedArgs
}
