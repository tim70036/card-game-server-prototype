package core

import (
	"errors"
	"fmt"
	"card-game-server-prototype/pkg/common/constant"
	"card-game-server-prototype/pkg/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type route struct {
	uid           Uid
	innerMsgBus   MsgBus
	innerNotifier GameNotifier
	unsubs        []UnsubscribeFunc
}

type Router struct {
	BaseHandler
	msgBus     MsgBus
	logger     *zap.Logger
	routeTable map[Uid]*route
}

func ProvideRouter(
	msgBus MsgBus,
	loggerFactory *util.LoggerFactory) *Router {
	return &Router{
		msgBus:     msgBus,
		logger:     loggerFactory.Create("Router"),
		routeTable: make(map[Uid]*route),
	}
}

func (r *Router) Bind(uid Uid, innerMsgBus MsgBus, innerNotifier GameNotifier) error {
	if _, ok := r.routeTable[uid]; ok {
		return fmt.Errorf("route already exists for uid: %s", uid.String())
	}

	r.logger.Debug("binding", zap.String("uid", uid.String()))

	topics := []Topic{
		MessageTopic,
		ChatTopic,
		EmoteEventTopic,
	}

	bindRoute := &route{
		uid:           uid,
		innerMsgBus:   innerMsgBus,
		innerNotifier: innerNotifier,
		unsubs:        make([]UnsubscribeFunc, 0),
	}

	handlers := lo.SliceToMap(
		topics,
		func(topic Topic) (Topic, func(proto.Message)) {
			return topic, func(msg proto.Message) {
				r.msgBus.Unicast(uid, topic, msg)
			}
		},
	)

	// Bind msg from inner.
	for topic, h := range handlers {
		unsub, err := bindRoute.innerMsgBus.Subscribe(uid, topic, h)
		if err != nil {
			r.logger.Error("Failed to subscribe to topic", zap.Error(err), zap.String("uid", uid.String()), zap.String("topic", string(topic)))
			return err
		}
		bindRoute.unsubs = append(bindRoute.unsubs, unsub)
	}

	r.routeTable[uid] = bindRoute
	return nil
}

func (r *Router) UnBind(uid Uid) error {
	bindRoute, ok := r.routeTable[uid]
	if !ok {
		// No route for this uid. No need to unbind.
		return nil
	}

	r.logger.Debug("unbinding", zap.String("uid", uid.String()))

	var errs error = nil
	for _, unsub := range bindRoute.unsubs {
		if err := unsub(); err != nil {
			r.logger.Error("Failed to unsubscribe during unbind", zap.Error(err), zap.String("uid", uid.String()))
			errs = errors.Join(errs, err)
		}
	}

	delete(r.routeTable, uid)
	return errs
}

func (r *Router) HandleConnect(uid Uid) error {
	bindRoute, ok := r.routeTable[uid]
	if !ok {
		return nil
	}

	notification := &ConnectNotification{
		Uid:  uid,
		Done: make(chan error, constant.PerNotificationBufferSize),
	}

	// Must be non-blocking to make sure caller will not be blocked.
	select {
	case bindRoute.innerNotifier.NotifyConnect() <- notification:
	default:
		r.logger.Warn("failed to route connect, message dropped", zap.String("uid", uid.String()))
	}

	return nil
}

func (r *Router) HandleEnter(uid Uid) error {
	bindRoute, ok := r.routeTable[uid]
	if !ok {
		return nil
	}
	notification := &EnterNotification{
		Uid:  uid,
		Done: make(chan error, constant.PerNotificationBufferSize),
	}

	// Must be non-blocking to make sure caller will not be blocked.
	select {
	case bindRoute.innerNotifier.NotifyEnter() <- notification:
	default:
		r.logger.Warn("failed to route enter, message dropped", zap.String("uid", uid.String()))
	}

	return nil
}

func (r *Router) HandleLeave(uid Uid) error {
	bindRoute, ok := r.routeTable[uid]
	if !ok {
		return nil
	}
	notification := &LeaveNotification{
		Uid:  uid,
		Done: make(chan error, constant.PerNotificationBufferSize),
	}

	// Must be non-blocking to make sure caller will not be blocked.
	select {
	case bindRoute.innerNotifier.NotifyLeave() <- notification:
	default:
		r.logger.Warn("failed to route leave, message dropped", zap.String("uid", uid.String()))
	}

	return nil
}

func (r *Router) HandleDisconnect(uid Uid) error {
	bindRoute, ok := r.routeTable[uid]
	if !ok {
		return nil
	}

	notification := &DisconnectNotification{
		Uid:  uid,
		Done: make(chan error, constant.PerNotificationBufferSize),
	}

	// Must be non-blocking to make sure caller will not be blocked.
	select {
	case bindRoute.innerNotifier.NotifyDisconnect() <- notification:
	default:
		r.logger.Warn("failed to route disconnect, message dropped", zap.String("uid", uid.String()))
	}

	return nil
}

func (r *Router) HandleRequest(req *Request) *Response {
	bindRoute, ok := r.routeTable[req.Uid]
	if !ok {
		return nil
	}

	notification := &RequestNotification{
		Request: req,
		Done:    make(chan *Response, constant.PerNotificationBufferSize),
	}

	// Must be non-blocking to make sure pool will not be blocked.
	select {
	case bindRoute.innerNotifier.NotifyRequest() <- notification:
	default:
		r.logger.Warn("failed to notify request, message dropped", zap.Object("req", req))
	}

	return nil
}
