package core

// https://github.com/mborders/artifex
type EventQueue struct {
}

// type Serializable[M proto.Message] interface {
// 	ToProto() M
// 	MarshalLogObject(enc zapcore.ObjectEncoder) error
// }

// func (state *baseState) RegisterModel(model any) {

// }

// func (state *baseState) RegisterIterableModel(model any) {
// 	state.modelPublishers = append(state.modelPublishers,
// 		func() {
// 			for uid, data := range model.Shit() {
// 				msgBus.Unicast(uid, model.Topic(), data.ToProto())
// 			}
// 		})
// }

// func (state *baseState) RegisterTargetModel(model any, getUid func() core.Uid) {
// 	msgBus.Unicast(getUid(), model.Topic(), model.GetData(getUid()))
// }

// type ModelGroupOperator[K comparable, V any, M proto.Message] interface {
// 	Insert(group ModelGroup[K, V, M], key K, value V)
// 	Delete(group ModelGroup[K, V, M], keys ...K)
// 	Clear(group ModelGroup[K, V, M])
// 	OnUpdate(callback func(msg M))
// }

// type modelGroupOperator[K comparable, V any, M proto.Message] struct {
// 	dataMux sync.RWMutex

// 	obMux     sync.Mutex
// 	observers []func(msg M)
// 	logger    *zap.Logger
// }

// func (o *modelGroupOperator[K, V, M]) Insert(group ModelGroup[K, V, M], key K, value V) {
// 	defer o.notify(group)

// 	o.dataMux.Lock()
// 	defer o.dataMux.Unlock()

// 	// // TODO: better way?
// 	// value.OnUpdated(func(msg VM) {
// 	// 	g.logger.Debug("group item OnUpdated", zap.Any("key", key), zap.Object("value", value), zap.String("msg", protojson.Format(msg)))
// 	// 	g.notify()
// 	// })

// 	group.maps()[key] = value
// }

// func (o *modelGroupOperator[K, V, M]) Delete(group ModelGroup[K, V, M], keys ...K) {
// 	defer o.notify(group)

// 	o.dataMux.Lock()
// 	defer o.dataMux.Unlock()
// 	for _, key := range keys {
// 		delete(group.maps(), key)
// 	}
// }

// func (o *modelGroupOperator[K, V, M]) Clear(group ModelGroup[K, V, M]) {
// 	defer o.notify(group)

// 	o.dataMux.Lock()
// 	defer o.dataMux.Unlock()
// 	maps.Clear(group.maps())
// }

// func (o *modelGroupOperator[K, V, M]) OnUpdated(callback func(msg M)) {
// 	o.obMux.Lock()
// 	defer o.obMux.Unlock()
// 	o.observers = append(o.observers, callback)
// }

// func (o *modelGroupOperator[K, V, M]) notify(group ModelGroup[K, V, M]) {
// 	msg := group.ToProto()
// 	// o.logger.Debug("notify", zap.Object("group", o), zap.String("msg", protojson.Format(msg)))

// 	o.obMux.Lock()
// 	defer o.obMux.Unlock()
// 	// TODO: what if cb blocked? go routine?
// 	for _, callback := range o.observers {
// 		callback(msg)
// 	}
// }

// type ModelGroup[K comparable, V any, M proto.Message] interface {
// 	Serializable[M]

// 	Keys() []K
// 	Values() []V
// 	maps() map[K]V
// 	mustEmbedBaseModelGroup()
// }

// type baseModelGroup[K comparable, V any, M proto.Message] struct {
// 	data   map[K]V
// 	logger *zap.Logger
// }

// func (group *baseModelGroup[K, V, M]) Keys() []K {
// 	return maps.Keys(group.data)
// }

// func (group *baseModelGroup[K, V, M]) Values() []V {
// 	return maps.Values(group.data)
// }

// func (group *baseModelGroup[K, V, M]) ToProto() M {
// 	group.logger.Warn("ToProto not implemented")
// 	var dummy M
// 	return dummy
// }

// func (group *baseModelGroup[K, V, M]) MarshalLogObject(enc zapcore.ObjectEncoder) error {
// 	return errors.New("MarshalLogObject not implemented")
// }

// func (group *baseModelGroup[K, V, M]) mustEmbedBaseModelGroup() {}

// Experiment...
// type MsgBus[T any] interface {
// 	Subscribe(uid Uid) <-chan T
// 	Unsubscribe(uid Uid)
// 	Close()
// 	Broadcast(msg T)
// 	Unicast(uid Uid, msg T)
// }

// // Generic factory not working in DI...
// func NewMsgBus[T any](name string, loggerFactory *util.LoggerFactory) MsgBus[T] {
// 	return &msgBus[T]{
// 		subscribers: make(map[Uid]chan T),
// 		logger:      loggerFactory.Create(name),
// 	}
// }

// type msgBus[T any] struct {
// 	mu sync.Mutex

// 	subscribers map[Uid]chan T
// 	logger      *zap.Logger
// }

// func (bus *msgBus[T]) Subscribe(uid Uid) <-chan T {
// 	bus.mu.Lock()
// 	defer bus.mu.Unlock()

// 	ch := make(chan T)
// 	bus.subscribers[uid] = ch
// 	bus.logger.Debug("register", zap.String("uid", uid.String()))
// 	return ch
// }

// func (bus *msgBus[T]) Unsubscribe(uid Uid) {
// 	bus.mu.Lock()
// 	defer bus.mu.Unlock()

// 	delete(bus.subscribers, uid)
// 	bus.logger.Debug("deregister", zap.String("uid", uid.String()))
// }

// func (bus *msgBus[T]) Close() {
// 	bus.mu.Lock()
// 	defer bus.mu.Unlock()

// 	for _, ch := range bus.subscribers {
// 		close(ch)
// 	}

// 	maps.Clear(bus.subscribers)
// }

// func (bus *msgBus[T]) Broadcast(msg T) {
// 	bus.mu.Lock()
// 	defer bus.mu.Unlock()

// 	bus.logger.Debug("broadcasting", zap.Any("msg", msg), zap.Any("uids", maps.Keys(bus.subscribers)))
// 	// TODO: what if ch blocked?
// 	for uid, ch := range bus.subscribers {
// 		select {
// 		case ch <- msg:
// 		default:
// 			bus.logger.Warn("sending to blocked channel", zap.String("uid", uid.String()))
// 		}

// 	}
// }

// func (bus *msgBus[T]) Unicast(uid Uid, msg T) {
// 	bus.mu.Lock()
// 	defer bus.mu.Unlock()

// 	bus.logger.Debug("unicasting", zap.Any("msg", msg), zap.String("uid", uid.String()))
// 	ch, ok := bus.subscribers[uid]
// 	if !ok {
// 		bus.logger.Warn("unicast to not registerd uid", zap.Any("msg", msg), zap.String("uid", uid.String()))
// 		return
// 	}

// 	// TODO: what if ch blocked?
// 	select {
// 	case ch <- msg:
// 	default:
// 		bus.logger.Warn("sending to blocked channel", zap.String("uid", uid.String()))
// 	}
// }

// type Model[T Serializable[M], M proto.Message] interface {
// 	OnUpdated(callback func(msg M))
// 	Notify(subject T)
// 	Logger() *zap.Logger

// 	mustEmbedBaseModel()
// }

// func NewModel[T Serializable[M], M proto.Message](name string, loggerFactory *util.LoggerFactory) Model[T, M] {
// 	return &baseModel[T, M]{
// 		observers: make([]func(payload M), 0),
// 		logger:    loggerFactory.Create(name),
// 	}
// }

// type baseModel[T Serializable[M], M proto.Message] struct {
// 	mu        sync.Mutex
// 	observers []func(payload M)
// 	logger    *zap.Logger
// }

// func (model *baseModel[T, M]) Logger() *zap.Logger { return model.logger }

// func (model *baseModel[T, M]) OnUpdated(callback func(msg M)) {
// 	model.mu.Lock()
// 	defer model.mu.Unlock()

// 	model.observers = append(model.observers, callback)
// }

// func (model *baseModel[T, M]) Notify(subject T) {
// 	msg := subject.ToProto()
// 	model.logger.Debug("notify", zap.Object("subject", subject))

// 	model.mu.Lock()
// 	defer model.mu.Unlock()
// 	// TODO: what if cb blocked? go routine?
// 	for _, callback := range model.observers {
// 		callback(msg)
// 	}
// }

// func (model *baseModel[T, M]) mustEmbedBaseModel() {}

// type Observable[M proto.Message] interface {
// 	Serializable[M]
// 	OnUpdated(callback func(msg M))
// }

// type Group[K comparable, V Observable[VM], VM proto.Message, M proto.Message] interface {
// 	Observable[M]

// 	Insert(key K, value V)
// 	Delete(keys ...K)
// 	Clear()

// 	Keys() []K
// 	Values() []V

// 	SetToProto(f func() M)

// 	mustEmbedBaseGroup()
// }

// func NewGroup[K comparable, V Observable[VM], VM proto.Message, M proto.Message](name string, loggerFactory *util.LoggerFactory) Group[K, V, VM, M] {
// 	return &baseGroup[K, V, VM, M]{
// 		data:      make(map[K]V),
// 		observers: make([]func(msg M), 0),
// 		logger:    loggerFactory.Create(name),
// 	}
// }

// type baseGroup[K comparable, V Observable[VM], VM proto.Message, M proto.Message] struct {
// 	data    map[K]V
// 	dataMux sync.RWMutex

// 	obMux     sync.Mutex
// 	observers []func(msg M)
// 	logger    *zap.Logger
// 	toProto   func() M
// }

// func (g *baseGroup[K, V, VM, M]) SetToProto(f func() M) {
// 	g.toProto = f
// }

// func (g *baseGroup[K, V, VM, M]) Insert(key K, value V) {
// 	defer g.notify()

// 	g.dataMux.Lock()
// 	defer g.dataMux.Unlock()

// 	// TODO: better way?
// 	value.OnUpdated(func(msg VM) {
// 		g.logger.Debug("group item OnUpdated", zap.Any("key", key), zap.Object("value", value), zap.String("msg", protojson.Format(msg)))
// 		g.notify()
// 	})

// 	g.data[key] = value
// }

// func (g *baseGroup[K, V, VM, M]) Delete(keys ...K) {
// 	defer g.notify()

// 	g.dataMux.Lock()
// 	defer g.dataMux.Unlock()
// 	for _, key := range keys {
// 		delete(g.data, key)
// 	}
// }

// func (g *baseGroup[K, V, VM, M]) Clear() {
// 	defer g.notify()

// 	g.dataMux.Lock()
// 	defer g.dataMux.Unlock()
// 	maps.Clear(g.data)
// }

// func (g *baseGroup[K, V, VM, M]) Keys() []K {
// 	return maps.Keys(g.data)
// }

// func (g *baseGroup[K, V, VM, M]) Values() []V {
// 	return maps.Values(g.data)
// }

// func (g *baseGroup[K, V, VM, M]) notify() {
// 	msg := g.toProto()
// 	g.logger.Debug("notify", zap.Object("group", g), zap.String("msg", protojson.Format(msg)))

// 	g.obMux.Lock()
// 	defer g.obMux.Unlock()
// 	// TODO: what if cb blocked? go routine?
// 	for _, callback := range g.observers {
// 		callback(msg)
// 	}
// }

// func (g *baseGroup[K, V, VM, M]) OnUpdated(callback func(msg M)) {
// 	g.obMux.Lock()
// 	defer g.obMux.Unlock()
// 	g.observers = append(g.observers, callback)
// }

// func (g *baseGroup[K, V, VM, M]) ToProto() M {
// 	g.logger.Warn("ToProto not implemented")
// 	var dummy M
// 	return dummy
// }

// func (g *baseGroup[K, V, VM, M]) MarshalLogObject(enc zapcore.ObjectEncoder) error {
// 	return errors.New("MarshalLogObject not implemented")
// }

// func (g *baseGroup[K, V, VM, M]) mustEmbedBaseGroup() {}
