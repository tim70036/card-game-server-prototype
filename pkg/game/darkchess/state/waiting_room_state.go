package state

import (
	"context"
	"card-game-server-prototype/pkg/common/api"
	commonmodel "card-game-server-prototype/pkg/common/model"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/common/type/chat"
	"card-game-server-prototype/pkg/common/type/kicker"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	"card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"github.com/samber/lo"
	"reflect"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoWaitingRoomState = core.NewStateTrigger("GoWaitingRoomState")

type WaitingRoomState struct {
	core.State

	buddyGroup  *commonmodel.BuddyGroup
	userGroup   *commonmodel.UserGroup
	userService commonservice.UserService
	gameInfo    *model.GameInfo
	roomInfo    *commonmodel.RoomInfo
	userApi     api.UserAPI

	kickedList         *kicker.KickedList
	chatHistory        *chat.History
	cancelTimer        context.CancelFunc
	cancelPremiumTimer context.CancelFunc
	endOnce            *sync.Once
}

func ProvideWaitingRoomState(
	stateFactory *core.StateFactory,

	buddyGroup *commonmodel.BuddyGroup,
	userGroup *commonmodel.UserGroup,
	userService commonservice.UserService,
	gameInfo *model.GameInfo,
	roomInfo *commonmodel.RoomInfo,
	userApi api.UserAPI,
) *WaitingRoomState {
	return &WaitingRoomState{
		State: stateFactory.Create("WaitingRoomState"),

		buddyGroup:  buddyGroup,
		userGroup:   userGroup,
		gameInfo:    gameInfo,
		userService: userService,
		roomInfo:    roomInfo,
		userApi:     userApi,
	}
}

func (state *WaitingRoomState) Run(context.Context, ...any) error {
	state.cancelTimer = func() {}
	state.cancelPremiumTimer = func() {}
	state.endOnce = &sync.Once{}

	state.kickedList = kicker.NewKickedList()
	state.chatHistory = chat.NewHistory()

	// If all users already disconnect during game, close game.
	if len(state.userGroup.Data) > 0 && !state.userGroup.HasAnyConnected() {
		state.Logger().Info("all users disconnected, closing game",
			zap.Object("users", state.userGroup),
		)
		state.endWaitingRoom(false)
		return nil
	}

	// 計算方式因遊戲而異，建立新遊戲需要特別確認這裡。
	// 這裡是：幾局就幾張

	for uid, info := range state.userGroup.Data {
		// AI 不需處理 cash、roomCards
		if info.IsAI {
			continue
		}

		if info.Cash < state.gameInfo.Setting.EnterLimit {
			if err := state.kickoutCannotAffordUser(uid, commongrpc.KickoutReason_NOT_ENOUGH_CASH); err != nil {
				state.GameController().GoErrorState()
				return nil
			}
		}

		// 房長卡：不需檢查 roomCards
		if state.roomInfo.IsPremium {
			continue
		}

		if info.RoomCards < state.gameInfo.Setting.TotalRound {
			if err := state.kickoutCannotAffordUser(uid, commongrpc.KickoutReason_NOT_ENOUGH_ROOM_CARD); err != nil {
				state.GameController().GoErrorState()
				return nil
			}
		}
	}

	if len(state.buddyGroup.Data) > 0 && lo.NoneBy(
		lo.Values(state.buddyGroup.Data),
		func(buddy *commonmodel.Buddy) bool { return buddy.IsOwner },
	) {
		state.buddyGroup.Data[lo.Keys(state.buddyGroup.Data)[0]].IsOwner = true
		state.Logger().Debug("init owner",
			zap.Object("buddy", state.buddyGroup),
		)
	}

	for _, buddy := range state.buddyGroup.Data {
		// AI 不能重置 ready 狀態
		if state.userGroup.Data[buddy.Uid].IsAI {
			continue
		}
		buddy.IsReady = false
	}

	// After the game is created, if no user has join the game for too
	// long. We should close the game since there might be some error.
	if len(state.userGroup.Data) == 0 {
		state.cancelTimer = state.GameController().RunTimer(constant.EmptyWaitingRoomGracefulPeriod, func() {
			state.Logger().Debug("no one enter for too long, closing game")
			state.endWaitingRoom(false)
		})
	}

	if state.roomInfo.IsPremium {
		if isStop := state.checkPremiumExpiredAndCloseRoom(); isStop {
			return nil
		}

		state.cancelPremiumTimer = state.GameController().RunTicker(time.Minute, func() {
			_ = state.checkPremiumExpiredAndCloseRoom()
		})
	}

	return nil
}

func (state *WaitingRoomState) checkPremiumExpiredAndCloseRoom() bool {
	if state.roomInfo.PremiumEndTimestamp <= time.Now().Unix() {
		if err := state.userService.Destroy(core.EventTopic, &gamegrpc.Event{
			Kickout: &commongrpc.Kickout{
				Reason: commongrpc.KickoutReason_ROOM_EXPIRED,
			},
		}, lo.Keys(state.userGroup.Data)...); err != nil {
			state.Logger().Error(
				"failed to destroy user when premium room card is expired",
				zap.Error(err),
				zap.Object("roomInfo", state.roomInfo),
				zap.Object("userGroup", state.userGroup),
			)
			state.GameController().GoErrorState()
			return true
		}

		state.GameController().RunTask(func() {
			state.endWaitingRoom(false)

			state.Logger().Info("premium room card is expired, game closed",
				zap.Object("roomInfo", state.roomInfo),
				zap.Object("userGroup", state.userGroup),
			)
		})
		return true
	}

	return false
}

func (state *WaitingRoomState) Publish(context.Context, ...any) error {
	state.MsgBus().Broadcast(core.GameStateTopic,
		state.ToProto("").(*gamegrpc.GameState),
	)
	state.MsgBus().Broadcast(core.ModelTopic, &gamegrpc.Model{
		UserGroup:  state.userGroup.ToProto(),
		BuddyGroup: state.buddyGroup.ToProto(),
	})
	return nil
}

func (state *WaitingRoomState) ToProto(_ core.Uid) proto.Message {
	return &gamegrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &gamegrpc.GameState_WaitingRoomStateContext{
			WaitingRoomStateContext: &gamegrpc.WaitingRoomStateContext{},
		},
	}
}

func (state *WaitingRoomState) kickoutCannotAffordUser(uid core.Uid, reason commongrpc.KickoutReason) error {
	kickoutReason := &gamegrpc.Event{
		Kickout: &commongrpc.Kickout{
			Reason: reason,
		},
	}

	if err := state.userService.Destroy(core.EventTopic, kickoutReason, uid); err != nil {
		state.Logger().Error(
			"failed to destroy cannot afford user",
			zap.Error(err),
			zap.String("uid", uid.String()),
			zap.String("reason", reason.String()),
			zap.Object("gameInfo", state.gameInfo),
		)
		return err
	}

	state.Logger().Debug("destroyed cannot afford user",
		zap.String("uid", uid.String()),
		zap.String("reason", reason.String()),
		zap.Object("gameInfo", state.gameInfo),
	)

	return nil
}

func (state *WaitingRoomState) BeforeConnect(uid core.Uid) error {
	if state.kickedList.Has(uid) {
		state.Logger().Warn("forbid user to connect since was kicked", zap.Array("kickedUids", state.kickedList.Get()), zap.String("uid", uid.String()))
		return status.Errorf(codes.PermissionDenied, "forbid user to connect since was kicked")
	}
	return nil
}

func (state *WaitingRoomState) BeforeEnter(uid core.Uid) error {
	if !state.buddyGroup.HasAnyOwner() {
		if _, ok := state.buddyGroup.Data[uid]; !ok {
			state.Logger().Error("set owner failed before enter waiting room",
				zap.String("uid", uid.String()),
				zap.Object("buddy", state.buddyGroup),
			)
			return status.Errorf(codes.NotFound, "cannot find user in buddyGroup")
		}

		state.buddyGroup.Data[uid].IsOwner = true
		state.Logger().Debug("set entered user as owner", zap.Object("buddy", state.buddyGroup))
	}

	state.cancelTimer()
	return nil
}

func (state *WaitingRoomState) BeforeLeave(leaveUid core.Uid) error {
	owner, hasOwner := lo.Find(
		lo.Values(state.buddyGroup.Data),
		func(b *commonmodel.Buddy) bool { return b.IsOwner },
	)

	if !hasOwner {
		state.Logger().Error("owner not found", zap.Object("buddy", state.buddyGroup), zap.String("uid", leaveUid.String()))
		return status.Errorf(codes.Internal, "owner not found")
	}

	if owner.Uid == leaveUid {
		if err := state.ownerLeaving(leaveUid); err != nil {
			state.GameController().GoErrorState()
			return err
		}
		return nil
	}

	return nil
}

func (state *WaitingRoomState) ownerLeaving(leaveUid core.Uid) error {
	// 房長卡: 房長離開房間，房間關閉
	if state.roomInfo.IsPremium {
		state.Logger().Info("owner leaving from premium room, close room",
			zap.Object("buddyGroup", state.buddyGroup),
			zap.Object("userGroup", state.userGroup),
			zap.Object("roomInfo", state.roomInfo),
		)

		if err := state.userService.Destroy(core.EventTopic, &gamegrpc.Event{
			Kickout: &commongrpc.Kickout{
				Reason: commongrpc.KickoutReason_GAME_CLOSED,
			},
		}, lo.Without(lo.Keys(state.userGroup.Data), leaveUid)...); err != nil {
			state.Logger().Error(
				"failed to destroy cannot afford user",
				zap.Error(err),
				zap.Object("userGroup", state.userGroup),
				zap.Object("gameInfo", state.gameInfo),
			)
			return err
		}

		state.GameController().RunTask(func() {
			state.endWaitingRoom(false)
		})
		return nil
	}

	// AI 不能成為 owner
	realUsers := lo.PickBy(state.userGroup.Data, func(_ core.Uid, user *commonmodel.User) bool {
		return !user.IsAI
	})
	newOwnerUids := lo.Without(lo.Keys(realUsers), leaveUid)
	if len(newOwnerUids) > 0 {
		newOwner := newOwnerUids[0]
		state.buddyGroup.Data[newOwner].IsOwner = true
		state.buddyGroup.Data[newOwner].IsReady = false
		state.Logger().Debug("owner leaving, set new owner", zap.Object("buddy", state.buddyGroup), zap.String("uid", leaveUid.String()))
	}

	return nil
}

func (state *WaitingRoomState) HandleDisconnect(uid core.Uid) error {
	if !state.userGroup.HasAnyConnected() {
		state.Logger().Info("all users disconnected, closing game", zap.Object("users", state.userGroup), zap.String("uid", uid.String()))
		state.GameController().RunTask(func() {
			state.endWaitingRoom(false)
		})
	}
	return nil
}

func (state *WaitingRoomState) AcceptRequestTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf(&commongrpc.ReadyRequest{}),
		reflect.TypeOf(&commongrpc.StartGameRequest{}),
		reflect.TypeOf(&commongrpc.KickRequest{}),
		reflect.TypeOf(&commongrpc.ChatRequest{}),
		// todo: AI邏輯未完成，暫不開放
		// 	reflect.TypeOf(&commongrpc.AddAiRequest{}),
	}
}

func (state *WaitingRoomState) HandleRequest(req *core.Request) error {
	buddy, ok := state.buddyGroup.Data[req.Uid]
	if !ok {
		state.Logger().Error("cannot find in buddyGroup", zap.Object("buddy", state.buddyGroup), zap.Object("req", req))
		return status.Errorf(codes.NotFound, "cannot find in buddyGroup")
	}

	switch msg := req.Msg.(type) {
	case *commongrpc.ChatRequest:
		if msg.Message == "" {
			return nil
		}

		chatMsg := chat.NewMsg(req.Uid, msg.Message)

		state.chatHistory.Add(chatMsg)

		state.MsgBus().Broadcast(core.ChatTopic, &commongrpc.ChatMessage{
			LatestMessage: chatMsg.ToProto(),
			History:       state.chatHistory.ToProto(),
		})

		return nil

	case *commongrpc.ReadyRequest:
		if buddy.IsOwner {
			state.Logger().Warn("owner cannot request ready", zap.Object("buddy", state.buddyGroup), zap.Object("req", req))
			return status.Errorf(codes.PermissionDenied, "owner cannot request ready")
		}

		buddy.IsReady = msg.Value
		state.Logger().Debug("req:ready",
			zap.Object("req", req),
			zap.Object("buddy", state.buddyGroup),
		)

		state.MsgBus().Broadcast(core.ModelTopic, &gamegrpc.Model{
			BuddyGroup: state.buddyGroup.ToProto(),
		})
		return nil

	case *commongrpc.StartGameRequest:
		if !buddy.IsOwner {
			state.Logger().Warn("only owner can request start game", zap.Object("buddy", state.buddyGroup), zap.Object("req", req))
			return status.Errorf(codes.PermissionDenied, "only owner can request start game")
		}

		readyCount := lo.CountBy(
			lo.Without(lo.Values(state.buddyGroup.Data), buddy),
			func(b *commonmodel.Buddy) bool { return b.IsReady },
		)

		if readyCount < constant.StartGameUserCount-1 {
			state.Logger().Warn("req:StartGame but not enough ready user", zap.Object("buddy", state.buddyGroup), zap.Object("req", req))
			return status.Errorf(codes.FailedPrecondition, "not enough ready user")
		}

		state.Logger().Debug("req:StartGame", zap.Object("req", req))
		state.GameController().RunTask(func() {
			state.endWaitingRoom(true)
		})
		return nil

	case *commongrpc.KickRequest:
		if !buddy.IsOwner {
			state.Logger().Warn("only owner can request kick", zap.Object("buddy", state.buddyGroup), zap.Object("req", req))
			return status.Errorf(codes.PermissionDenied, "only owner can request kick")
		}

		if state.kickedList.Has(core.Uid(msg.Uid)) {
			state.Logger().Warn("already raised kicked", zap.Object("req", req))
			return nil
		}

		state.kickedList.Add(core.Uid(msg.Uid))

		kickoutMsg := &gamegrpc.Event{
			Kickout: &commongrpc.Kickout{
				Reason: commongrpc.KickoutReason_OWNER_REQUESTED,
			},
		}

		if err := state.userService.Destroy(core.EventTopic, kickoutMsg, core.Uid(msg.Uid)); err != nil {
			state.Logger().Error("failed to destroy kicked user", zap.Error(err), zap.Object("users", state.userGroup), zap.Object("req", req))
			state.GameController().GoErrorState()
			return status.Errorf(codes.Internal, "failed to destroy kicked user")
		}

		state.Logger().Debug("req:Kick", zap.Object("req", req))

		state.MsgBus().Broadcast(core.ModelTopic, &gamegrpc.Model{
			UserGroup:  state.userGroup.ToProto(),
			BuddyGroup: state.buddyGroup.ToProto(),
		})
		return nil

	case *commongrpc.AddAiRequest:
		if !buddy.IsOwner {
			state.Logger().Warn("only owner can request AddAI",
				zap.Object("buddy", state.buddyGroup),
				zap.Object("req", req),
			)
			return status.Errorf(codes.PermissionDenied, "only owner can request AddAI")
		}

		if state.gameInfo.Setting.AnteAmount > constant.AnteToAvoidAI {
			return status.Errorf(codes.FailedPrecondition, "ante is too high to add AI")
		}

		state.GameController().RunTask(func() {
			if len(state.buddyGroup.Data) == constant.StartGameUserCount {
				state.Logger().Warn("table is full, cannot add AI.",
					zap.Object("buddy", state.buddyGroup),
					zap.Object("req", req),
				)
				return
			}

			aiUid, err := state.userService.GetAI(state.gameInfo.Setting.EnterLimit, state.gameInfo.Setting.TotalRound)
			if err != nil {
				state.Logger().Warn("failed to add AI", zap.Error(err))
				return
			}

			if err := state.userService.Init(aiUid); err != nil {
				state.Logger().Error("failed to init AI user",
					zap.Error(err),
					zap.String("ai", aiUid.String()),
				)
				return
			}

			if err := state.userService.FetchFromRepo(aiUid); err != nil {
				state.Logger().Error("failed to fetch AI user from repo",
					zap.Error(err),
					zap.String("ai", aiUid.String()),
				)
				return
			}

			state.userGroup.Data[aiUid].IsConnected = true
			state.userGroup.Data[aiUid].HasEntered = true
			state.buddyGroup.Data[aiUid].IsReady = true

			// Todo: 應該要用 ModelTopic，但神奇的是功能正常，直接改怕會影響到現有功能，要等前端有空再看看能不能直接修正這個錯誤。
			// https://shorturl.at/S3Rpj
			state.MsgBus().Broadcast(core.MessageTopic, &gamegrpc.Model{
				UserGroup:  state.userGroup.ToProto(),
				BuddyGroup: state.buddyGroup.ToProto(),
			})
		})

		return nil

	default:
		state.Logger().Warn("not supported request", zap.Object("req", req))
		return status.Errorf(codes.Unimplemented, "not supported request")
	}
}

func (state *WaitingRoomState) endWaitingRoom(isStartGame bool) {
	state.endOnce.Do(func() {
		state.Logger().Debug("end waiting room", zap.Bool("startGame", isStartGame))
		state.cancelTimer()
		state.cancelPremiumTimer()

		if isStartGame {
			state.GameController().GoNextState(GoWaitUserState)
		} else {
			state.GameController().GoNextState(GoClosedState)
		}
	})
}
