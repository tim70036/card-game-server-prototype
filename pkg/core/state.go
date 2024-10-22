package core

import (
	"context"
	"card-game-server-prototype/pkg/util"
	"reflect"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
)

type State interface {
	Name() string
	PublishTime() time.Time
	GameController() GameController
	MsgBus() MsgBus
	Logger() *zap.Logger

	// Run State execution order (see stateless machine for more detail):
	//
	// OnEntry: Run() -> beforePublish() -> Publish()
	//
	// OnExit: Cleanup()
	Run(ctx context.Context, args ...any) error
	beforePublish(ctx context.Context, args ...any) error // Private to package.
	Publish(ctx context.Context, args ...any) error
	Cleanup(ctx context.Context, args ...any) error

	BeforeConnect(uid Uid) error
	HandleConnect(uid Uid) error

	BeforeEnter(uid Uid) error
	HandleEnter(uid Uid) error

	BeforeDisconnect(uid Uid) error
	HandleDisconnect(uid Uid) error

	BeforeLeave(uid Uid) error
	HandleLeave(uid Uid) error

	BeforeRequest(req *Request) error
	HandleRequest(req *Request) error
	AcceptRequestTypes() []reflect.Type

	ToProto(uid Uid) proto.Message
	MarshalLogObject(enc zapcore.ObjectEncoder) error
	mustEmbedBaseState()
}

type StateFactory struct {
	gameController GameController
	msgBus         MsgBus
	loggerFactory  *util.LoggerFactory
}

func ProvideStateFactory(
	gameController GameController,
	msgBus MsgBus,
	loggerFactory *util.LoggerFactory,
) *StateFactory {
	return &StateFactory{
		gameController: gameController,
		msgBus:         msgBus,
		loggerFactory:  loggerFactory,
	}
}

func (f *StateFactory) Create(name string) State {
	return &baseState{
		name:           name,
		gameController: f.gameController,
		msgBus:         f.msgBus,
		logger:         f.loggerFactory.Create(name),
	}
}

type baseState struct {
	name           string
	publishTime    time.Time
	gameController GameController
	msgBus         MsgBus
	logger         *zap.Logger
}

func (state *baseState) Name() string                   { return state.name }
func (state *baseState) PublishTime() time.Time         { return state.publishTime }
func (state *baseState) GameController() GameController { return state.gameController }
func (state *baseState) MsgBus() MsgBus                 { return state.msgBus }
func (state *baseState) Logger() *zap.Logger            { return state.logger }

func (state *baseState) Run(ctx context.Context, args ...any) error {
	state.logger.Warn("Run not implemented")
	return nil
}

func (state *baseState) beforePublish(ctx context.Context, args ...any) error {
	state.publishTime = time.Now()
	return nil
}

func (state *baseState) Publish(ctx context.Context, args ...any) error {
	state.logger.Error("Publish not implemented")
	return nil
}

func (state *baseState) Cleanup(ctx context.Context, args ...any) error {
	// state.logger.Debug("Cleanup not implemented")
	return nil
}

func (state *baseState) BeforeConnect(uid Uid) error {
	// state.logger.Debug("BeforeConnect not implemented", zap.String("uid", uid.String()))
	return nil
}

func (state *baseState) HandleConnect(uid Uid) error {
	// state.logger.Debug("HandleConnect not implemented", zap.String("uid", uid.String()))
	return nil
}

func (state *baseState) BeforeEnter(uid Uid) error {
	// state.logger.Debug("BeforeEnter not implemented", zap.String("uid", uid.String()))
	return nil
}

func (state *baseState) HandleEnter(uid Uid) error {
	// state.logger.Debug("HandleEnter not implemented", zap.String("uid", uid.String()))
	return nil
}

func (state *baseState) BeforeDisconnect(uid Uid) error {
	// state.logger.Debug("BeforeDisconnect not implemented", zap.String("uid", uid.String()))
	return nil
}

func (state *baseState) HandleDisconnect(uid Uid) error {
	// state.logger.Debug("BeforeDisconnect not implemented", zap.String("uid", uid.String()))
	return nil
}

func (state *baseState) BeforeLeave(uid Uid) error {
	// state.logger.Debug("BeforeDisconnect not implemented", zap.String("uid", uid.String()))
	return nil
}

func (state *baseState) HandleLeave(uid Uid) error {
	// state.logger.Debug("BeforeDisconnect not implemented", zap.String("uid", uid.String()))
	return nil
}

func (state *baseState) BeforeRequest(req *Request) error {
	// state.logger.Debug("BeforeRequest not implemented", zap.String("uid", uid.String()))
	return nil
}

func (state *baseState) HandleRequest(req *Request) error {
	// state.logger.Debug("HandleRequest not implemented", zap.String("uid", uid.String()))
	return nil
}

func (state *baseState) AcceptRequestTypes() []reflect.Type {
	return []reflect.Type{}
}

func (state *baseState) ToProto(uid Uid) proto.Message {
	state.logger.Error("ToProto not implemented")
	return nil
}

func (state *baseState) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", state.name)
	return nil
}

func (state *baseState) mustEmbedBaseState() {}
