package core

import (
	"card-game-server-prototype/pkg/common/constant"
	"card-game-server-prototype/pkg/util"
)

type Topic string

const (
	MessageTopic    Topic = "Message"
	GameStateTopic  Topic = "GameState"
	ModelTopic      Topic = "Model"
	EventTopic      Topic = "Event"
	ChatTopic       Topic = "Chat"
	EmoteEventTopic Topic = "EmoteEvent"
)

type MsgBus interface {
	Broadcast(topic Topic, args ...interface{})
	Unicast(uid Uid, topic Topic, args ...interface{})
	Subscribe(uid Uid, topic Topic, fn interface{}) (UnsubscribeFunc, error)
}

func ProvideMsgBus(loggerFactory *util.LoggerFactory) MsgBus {
	return &rawMessageBus{
		handlerQueueSize: constant.PerConnectionBufferSize,
		handlers:         make(handlersMap),
		logger:           loggerFactory.Create("MsgBus"),
	}

}
