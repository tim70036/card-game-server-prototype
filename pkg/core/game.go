package core

import (
	"context"
	"fmt"
	"card-game-server-prototype/pkg/common/constant"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"card-game-server-prototype/pkg/util"
	"github.com/qmuntal/stateless"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"reflect"
	"sync/atomic"
	"time"
)

type GameNotifier interface {
	NotifyConnect() chan<- *ConnectNotification
	NotifyEnter() chan<- *EnterNotification
	NotifyLeave() chan<- *LeaveNotification
	NotifyDisconnect() chan<- *DisconnectNotification
	NotifyRequest() chan<- *RequestNotification
}

type GameController interface {
	GoNextState(trigger *StateTrigger, args ...any)
	GoErrorState()
	CurrentState() State

	// RunTimer Run delayed task that needs to be thread safe with main runner
	// loop. This is useful for case such as modifying model data or
	// transitioning state machine.
	RunTimer(timeout time.Duration, callback func()) context.CancelFunc

	// RunTicker Run task periodically that needs to be thread safe with main
	// runner loop. This is useful for case such as modifying model
	// data or transitioning state machine.
	RunTicker(interval time.Duration, callback func()) context.CancelFunc

	// RunTask Run task that needs to be thread safe with main runner loop.
	// This is useful for case such as modifying model data or
	// transitioning state machine.
	RunTask(task func())

	// RunActorRequests Run actor requests that needs to be thread safe with main runner.
	RunActorRequests(actorReqs ActorRequestList)
	Shutdown()
}

type Game interface {
	GameNotifier
	GameController

	// Run Main runner loop of game. It will process notification(event) 1
	// by 1.
	Run()

	ConfigInitState(initState State) *stateless.StateConfiguration
	ConfigErrorState(errorHandleState State) *stateless.StateConfiguration
	ConfigState(state State) *stateless.StateConfiguration
	ConfigTriggerParamsType(trigger *StateTrigger)
	ConfigHandler(handler Handler)

	OnConnect(func(uid Uid) error)
	OnEnter(func(uid Uid) error)
	OnLeave(func(uid Uid) error)
	OnDisconnect(func(uid Uid) error)
	OnRequest(func(req *Request) *Response)

	WaitShutdown() <-chan struct{}

	mustEmbedBaseGame()
}

func ProvideGame(loggerFactory *util.LoggerFactory) Game {
	g := &baseGame{
		fsm:            stateless.NewStateMachine("Nil"),
		onConnect:      make(chan *ConnectNotification, constant.CoreGameRunnerBufferSize),
		onEnter:        make(chan *EnterNotification, constant.CoreGameRunnerBufferSize),
		onLeave:        make(chan *LeaveNotification, constant.CoreGameRunnerBufferSize),
		onDisconnect:   make(chan *DisconnectNotification, constant.CoreGameRunnerBufferSize),
		onRequest:      make(chan *RequestNotification, constant.CoreGameRunnerBufferSize),
		onTask:         make(chan *TaskNotification, constant.CoreGameRunnerBufferSize),
		onNextState:    make(chan *NextStateNotification, constant.CoreGameRunnerBufferSize),
		notifyShutdown: make(chan struct{}, constant.PerNotificationBufferSize),

		isRunning:          &atomic.Bool{},
		connectHandlers:    []func(uid Uid) error{},
		enterHandlers:      []func(uid Uid) error{},
		leaveHandlers:      []func(uid Uid) error{},
		disconnectHandlers: []func(uid Uid) error{},
		requestHandlers:    []func(req *Request) *Response{},

		initState:  nil,
		errorState: nil,
		errorTrigger: &StateTrigger{
			Name:      "GoErrorState",
			ArgsTypes: []reflect.Type{reflect.TypeOf(commongrpc.KickoutReason_GAME_EXCEPTION)},
		},
		logger: loggerFactory.Create("Game"),
	}

	return g
}

type baseGame struct {
	fsm *stateless.StateMachine

	onConnect      chan *ConnectNotification
	onEnter        chan *EnterNotification
	onLeave        chan *LeaveNotification
	onDisconnect   chan *DisconnectNotification
	onRequest      chan *RequestNotification
	onTask         chan *TaskNotification
	onNextState    chan *NextStateNotification
	notifyShutdown chan struct{}

	isRunning          *atomic.Bool
	connectHandlers    []func(uid Uid) error
	enterHandlers      []func(uid Uid) error
	leaveHandlers      []func(uid Uid) error
	disconnectHandlers []func(uid Uid) error
	requestHandlers    []func(req *Request) *Response

	initState State

	// A state that handles error. If sth bad happens in a state,
	// then it should go to this state.
	errorState   State
	errorTrigger *StateTrigger

	logger *zap.Logger
}

// ----- TTT ConnectionService -----
func (g *baseGame) NotifyConnect() chan<- *ConnectNotification       { return g.onConnect }
func (g *baseGame) NotifyEnter() chan<- *EnterNotification           { return g.onEnter }
func (g *baseGame) NotifyLeave() chan<- *LeaveNotification           { return g.onLeave }
func (g *baseGame) NotifyDisconnect() chan<- *DisconnectNotification { return g.onDisconnect }

// ----- TTT ActionService -----
func (g *baseGame) NotifyRequest() chan<- *RequestNotification { return g.onRequest }

// ----- TTT in games/{game} -----
func (g *baseGame) WaitShutdown() <-chan struct{} { return g.notifyShutdown }

func (g *baseGame) CurrentState() State { return g.fsm.MustState().(State) }
func (g *baseGame) mustEmbedBaseGame()  {}

func (g *baseGame) Run() {
	g.isRunning.Store(true)
	g.logger.Debug("start")

	if g.initState == nil {
		g.logger.Fatal("init state is not configured, call ConfigInitState() first")
		return
	}

	if g.errorState == nil {
		g.logger.Fatal("error state is not configured, call ConfigErrorState() first")
		return
	}

	// TTT State 起點
	g.fsm.Configure("Nil").
		Permit("Start", g.initState)
	if err := g.fsm.Fire("Start"); err != nil {
		g.logger.Fatal("failed to init fsm", zap.Error(err))
		return
	}

	g.fsm.OnUnhandledTrigger(func(ctx context.Context, state stateless.State, trigger stateless.Trigger, unmetGuards []string) error {
		g.logger.Error(
			"unhandled trigger",
			zap.Object("State", state.(State)),
			zap.Object("Trigger", trigger.(*StateTrigger)),
			zap.Strings("UnmetGuards", unmetGuards),
		)
		return fmt.Errorf("unhandled trigger %v from state %v", trigger.(*StateTrigger).Name, state.(State).Name())
	})

	// TTT 處理一些 rpc 接口
	for {
		select {
		case notification := <-g.onConnect:
			g.logger.Debug("onConnect", zap.String("Uid", notification.Uid.String()))

			if err := g.fsm.MustState().(State).BeforeConnect(notification.Uid); err != nil {
				g.logger.Warn("err from state before connect handler", zap.Error(err), zap.String("Uid", notification.Uid.String()), zap.Object("State", g.fsm.MustState().(State)))
				notification.Done <- err
				close(notification.Done)
				continue
			}

			var err error
			for _, handler := range g.connectHandlers {
				if err = handler(notification.Uid); err != nil {
					break
				}
			}

			if err != nil {
				g.logger.Warn("err from connect handlers", zap.Error(err), zap.String("Uid", notification.Uid.String()))
				notification.Done <- err
				close(notification.Done)
				continue
			}

			if err := g.fsm.MustState().(State).HandleConnect(notification.Uid); err != nil {
				g.logger.Warn("err from state handle connect handler", zap.Error(err), zap.String("Uid", notification.Uid.String()), zap.Object("State", g.fsm.MustState().(State)))
				notification.Done <- err
				close(notification.Done)
				continue
			}

			g.logger.Debug("onConnect done", zap.String("Uid", notification.Uid.String()))
			close(notification.Done)

		case notification := <-g.onEnter:
			g.logger.Debug("onEnter", zap.String("Uid", notification.Uid.String()))

			if err := g.fsm.MustState().(State).BeforeEnter(notification.Uid); err != nil {
				g.logger.Warn("err from state before enter handler", zap.Error(err), zap.String("Uid", notification.Uid.String()), zap.Object("State", g.fsm.MustState().(State)))
				notification.Done <- err
				close(notification.Done)
				continue
			}

			var err error
			for _, handler := range g.enterHandlers {
				if err = handler(notification.Uid); err != nil {
					break
				}
			}

			if err != nil {
				g.logger.Warn("err from enter handler", zap.Error(err), zap.String("Uid", notification.Uid.String()))
				notification.Done <- err
				close(notification.Done)
				continue
			}

			if err := g.fsm.MustState().(State).HandleEnter(notification.Uid); err != nil {
				g.logger.Warn("err from state handle enter handler", zap.Error(err), zap.String("Uid", notification.Uid.String()), zap.Object("State", g.fsm.MustState().(State)))
				notification.Done <- err
				close(notification.Done)
				continue
			}

			g.logger.Debug("onEnter done", zap.String("Uid", notification.Uid.String()))
			close(notification.Done)

		case notification := <-g.onDisconnect:
			g.logger.Debug("onDisconnect", zap.String("Uid", notification.Uid.String()))

			if err := g.fsm.MustState().(State).BeforeDisconnect(notification.Uid); err != nil {
				g.logger.Warn("err from state before disconnect handler", zap.Error(err), zap.String("Uid", notification.Uid.String()), zap.Object("State", g.fsm.MustState().(State)))
				notification.Done <- err
				close(notification.Done)
				continue
			}

			var err error
			for _, handler := range g.disconnectHandlers {
				if err = handler(notification.Uid); err != nil {
					break
				}
			}

			if err != nil {
				g.logger.Warn("err from disconnect handler", zap.Error(err), zap.String("Uid", notification.Uid.String()))
				notification.Done <- err
				close(notification.Done)
				continue
			}

			if err := g.fsm.MustState().(State).HandleDisconnect(notification.Uid); err != nil {
				g.logger.Warn("err from state handle disconnect handler", zap.Error(err), zap.String("Uid", notification.Uid.String()), zap.Object("State", g.fsm.MustState().(State)))
				notification.Done <- err
				close(notification.Done)
				continue
			}

			g.logger.Debug("onDisconnect done", zap.String("Uid", notification.Uid.String()))
			close(notification.Done)

		case notification := <-g.onLeave:
			g.logger.Debug("onLeave", zap.String("Uid", notification.Uid.String()))

			if err := g.fsm.MustState().(State).BeforeLeave(notification.Uid); err != nil {
				g.logger.Warn("err from state before leave handler", zap.Error(err), zap.String("Uid", notification.Uid.String()), zap.Object("State", g.fsm.MustState().(State)))
				notification.Done <- err
				close(notification.Done)
				continue
			}

			var err error
			for _, handler := range g.leaveHandlers {
				if err = handler(notification.Uid); err != nil {
					break
				}
			}

			if err != nil {
				g.logger.Warn("err from leave handler", zap.Error(err), zap.String("Uid", notification.Uid.String()))
				notification.Done <- err
				close(notification.Done)
				continue
			}

			if err := g.fsm.MustState().(State).HandleLeave(notification.Uid); err != nil {
				g.logger.Warn("err from state handle leave handler", zap.Error(err), zap.String("Uid", notification.Uid.String()), zap.Object("State", g.fsm.MustState().(State)))
				notification.Done <- err
				close(notification.Done)
				continue
			}

			g.logger.Debug("onLeave done", zap.String("Uid", notification.Uid.String()))
			close(notification.Done)

		case notification := <-g.onRequest:
			g.logger.Debug("onRequest", zap.Object("Request", notification.Request))

			state := g.fsm.MustState().(State)
			isAcceptByState := lo.ContainsBy(
				state.AcceptRequestTypes(),
				func(t reflect.Type) bool {
					return t == reflect.TypeOf(notification.Request.Msg)
				},
			)

			if isAcceptByState {
				if err := state.BeforeRequest(notification.Request); err != nil {
					g.logger.Warn("err from state before request handler", zap.Error(err), zap.Object("Request", notification.Request), zap.Object("state", state))
					notification.Done <- &Response{Err: err}
					close(notification.Done)
					continue
				}
			}

			// TTT 處理 {game} 屬於 common 的 action: sticker...
			var resp *Response
			for _, handler := range g.requestHandlers {
				if resp = handler(notification.Request); resp != nil {
					break
				}
			}

			if resp != nil {
				if resp.Err != nil {
					err := resp.Err
					g.logger.Warn("err from request handler", zap.Error(err), zap.Object("Request", notification.Request))
				} else {
					g.logger.Debug("onRequest receive response done",
						zap.Object("Request", notification.Request),
						zap.Object("Response", resp),
					)
				}

				notification.Done <- resp
				close(notification.Done)
				continue
			}

			// TTT 處理 {game} 屬於自己獨有的 state action
			if isAcceptByState {
				if err := state.HandleRequest(notification.Request); err != nil {
					g.logger.Warn("err from state handle request handler", zap.Error(err), zap.Object("Request", notification.Request), zap.Object("state", state))
					notification.Done <- &Response{Err: err}
					close(notification.Done)
					continue
				}
			}

			g.logger.Debug("onRequest done", zap.Object("Request", notification.Request))
			close(notification.Done)

		case notification := <-g.onTask:
			// g.logger.Debug("onTask")
			notification.Task()
			// g.logger.Debug("onTask done")

		case notification := <-g.onNextState:
			// g.logger.Debug("onNextState", zap.Object("State", g.fsm.MustState().(State)), zap.Object("Trigger", notification.Trigger), zap.Any("Args", notification.Args))

			g.RunTimer(100*time.Millisecond, func() {
				if err := g.fsm.Fire(notification.Trigger, notification.Args...); err != nil {
					if notification.Trigger == g.errorTrigger {
						g.logger.Error("failed to go error state, force shutting down", zap.Error(err))
						g.Shutdown()
					} else {
						g.logger.Error("failed to go next state, go error state", zap.Error(err))
						g.GoErrorState()
					}
				} else {
					// g.logger.Debug("onNextState done")
				}
			})

		case <-g.notifyShutdown:
			g.logger.Info("shutting down run loop")
			return
		}
	}
}

func (g *baseGame) ConfigInitState(initState State) *stateless.StateConfiguration {
	g.initState = initState
	return g.ConfigState(initState)
}

// ConfigErrorState Config a state that will be used to handle error.  All game state
// that config by calling ConfigGameState() will transition to this
// error handle state as soon as invoking HandleError(). Remeber to
// call this before calling any ConfigState().
func (g *baseGame) ConfigErrorState(errorState State) *stateless.StateConfiguration {
	g.errorState = errorState
	g.ConfigTriggerParamsType(g.errorTrigger)
	return g.fsm.Configure(errorState).
		OnEntry(errorState.Run).
		OnEntry(errorState.beforePublish).
		OnEntry(errorState.Publish).
		OnExit(errorState.Cleanup)
}

func (g *baseGame) ConfigState(state State) *stateless.StateConfiguration {
	return g.fsm.Configure(state).
		OnEntry(state.Run).
		OnEntry(state.beforePublish).
		OnEntry(state.Publish).
		OnExit(state.Cleanup).
		Permit(g.errorTrigger, g.errorState)
}

func (g *baseGame) ConfigTriggerParamsType(trigger *StateTrigger) {
	g.fsm.SetTriggerParameters(trigger, trigger.ArgsTypes...)
}

func (g *baseGame) ConfigHandler(handler Handler) {
	if g.isRunning.Load() {
		g.logger.Error("ConfigHandler must be called before Run()")
		return
	}

	g.connectHandlers = append(g.connectHandlers, handler.HandleConnect)
	g.enterHandlers = append(g.enterHandlers, handler.HandleEnter)
	g.disconnectHandlers = append(g.disconnectHandlers, handler.HandleDisconnect)
	g.leaveHandlers = append(g.leaveHandlers, handler.HandleLeave)
	g.requestHandlers = append(g.requestHandlers, handler.HandleRequest)
}

func (g *baseGame) OnConnect(handler func(uid Uid) error) {
	if g.isRunning.Load() {
		g.logger.Error("OnConnect must be called before Run()")
		return
	}

	g.connectHandlers = append(g.connectHandlers, handler)
}

func (g *baseGame) OnEnter(handler func(uid Uid) error) {
	if g.isRunning.Load() {
		g.logger.Error("OnEnter must be called before Run()")
		return
	}

	g.enterHandlers = append(g.enterHandlers, handler)
}

func (g *baseGame) OnDisconnect(handler func(uid Uid) error) {
	if g.isRunning.Load() {
		g.logger.Error("OnDisconnect must be called before Run()")
		return
	}

	g.disconnectHandlers = append(g.disconnectHandlers, handler)
}

func (g *baseGame) OnLeave(handler func(uid Uid) error) {
	if g.isRunning.Load() {
		g.logger.Error("OnLeave must be called before Run()")
		return
	}

	g.leaveHandlers = append(g.leaveHandlers, handler)
}

func (g *baseGame) OnRequest(handler func(req *Request) *Response) {
	if g.isRunning.Load() {
		g.logger.Error("OnRequest must be called before Run()")
		return
	}

	g.requestHandlers = append(g.requestHandlers, handler)
}

func (g *baseGame) GoNextState(trigger *StateTrigger, args ...any) {
	g.onNextState <- &NextStateNotification{
		Trigger: trigger,
		Args:    args,
	}
}

func (g *baseGame) GoErrorState() {
	g.onNextState <- &NextStateNotification{
		Trigger: g.errorTrigger,
		Args:    []any{commongrpc.KickoutReason_GAME_EXCEPTION},
	}
}

func (g *baseGame) RunTimer(timeout time.Duration, callback func()) context.CancelFunc {
	ctx, cancelTimer := context.WithCancel(context.Background())

	go (func() {
		// timerId := uuid.NewString()[:6]
		timer := time.NewTimer(timeout)
		// g.logger.Debug(
		//	"timer started",
		//	zap.Duration("timeout", timeout),
		//	zap.String("timerId", timerId),
		// )

		select {
		case <-timer.C:
			// g.logger.Debug(
			//	"timer triggered",
			//	zap.Duration("timeout", timeout),
			//	zap.String("timerId", timerId),
			// )
			g.RunTask(callback)
			return

		case <-ctx.Done():
			// g.logger.Debug(
			//	"timer canceled",
			//	zap.Duration("timeout", timeout),
			//	zap.String("timerId", timerId),
			// )

			if !timer.Stop() {
				<-timer.C
			}
			return
		}
	})()

	return cancelTimer
}

func (g *baseGame) RunTicker(interval time.Duration, callback func()) context.CancelFunc {
	ctx, cancelTicker := context.WithCancel(context.Background())

	go (func() {
		// tickerId := uuid.NewString()[:6]
		ticker := time.NewTicker(interval)
		// g.logger.Debug(
		//	"ticker started",
		//	zap.Duration("interval", interval),
		//	zap.String("tickerId", tickerId),
		// )

		for {
			select {
			case <-ticker.C:
				// g.logger.Debug(
				// 	"ticker triggered",
				// 	zap.Duration("interval", interval),
				// 	zap.String("tickerId", tickerId),
				// )
				g.RunTask(callback)
				continue

			case <-ctx.Done():
				// g.logger.Debug(
				//	"ticker canceled",
				//	zap.Duration("interval", interval),
				//	zap.String("tickerId", tickerId),
				// )

				ticker.Stop()
				return
			}
		}
	})()

	return cancelTicker
}

func (g *baseGame) RunTask(task func()) {
	g.onTask <- &TaskNotification{
		Task: task,
	}
}

func (g *baseGame) RunActorRequests(actorReqs ActorRequestList) {
	for _, actorReq := range actorReqs {
		go (func() {
			time.Sleep(actorReq.Delay)
			g.logger.Debug(
				"running actor request",
				zap.Object("actorReq", actorReq),
			)

			g.onRequest <- &RequestNotification{
				Request: actorReq.Req,
				Done:    make(chan *Response, constant.PerNotificationBufferSize),
			}
		})()
	}
}

func (g *baseGame) Shutdown() {
	g.logger.Info("shutting down")
	close(g.notifyShutdown)
}
