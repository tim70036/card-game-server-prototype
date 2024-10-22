package handler

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/constant"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/service"
	"card-game-server-prototype/pkg/game/txpoker/type/action"
	event2 "card-game-server-prototype/pkg/game/txpoker/type/event"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"
)

type RequestHandler struct {
	core.BaseHandler

	seatStatusGroup   *model2.SeatStatusGroup
	playSettingGroup  *model2.PlaySettingGroup
	actionHintGroup   *model2.ActionHintGroup
	seatStatusService service.SeatStatusService
	forceBuyInGroup   *model2.ForceBuyInGroup
	eventGroup        *model2.EventGroup

	msgBus core.MsgBus
	logger *zap.Logger
}

func ProvideRequestHandler(
	seatStatusGroup *model2.SeatStatusGroup,
	playSettingGroup *model2.PlaySettingGroup,
	actionHintGroup *model2.ActionHintGroup,
	seatStatusService service.SeatStatusService,
	forceBuyInGroup *model2.ForceBuyInGroup,
	eventGroup *model2.EventGroup,

	msgBus core.MsgBus,
	loggerFactory *util.LoggerFactory,
) *RequestHandler {
	return &RequestHandler{
		seatStatusGroup:   seatStatusGroup,
		playSettingGroup:  playSettingGroup,
		actionHintGroup:   actionHintGroup,
		seatStatusService: seatStatusService,
		forceBuyInGroup:   forceBuyInGroup,
		eventGroup:        eventGroup,

		msgBus: msgBus,
		logger: loggerFactory.Create("RequestHandler"),
	}
}

func (handler *RequestHandler) HandleRequest(req *core.Request) *core.Response {
	switch msg := req.Msg.(type) {
	case *commongrpc.EmotePingRequest:

		handler.msgBus.Broadcast(core.EmoteEventTopic, &commongrpc.EmoteEvent{
			EmotePing: &commongrpc.EmotePing{
				ItemId:    msg.ItemId,
				SenderUid: req.Uid.String(),
				TargetUid: msg.TargetUid,
			},
		})
		return nil

	case *commongrpc.StickerRequest:
		handler.eventGroup.Data[req.Uid] = append(handler.eventGroup.Data[req.Uid], &event2.Event{
			Type:   event2.UseSticker,
			Amount: 1,
		})

		handler.msgBus.Broadcast(core.EmoteEventTopic, &commongrpc.EmoteEvent{
			Sticker: &commongrpc.Sticker{
				Uid:       req.Uid.String(),
				StickerId: msg.StickerId,
			},
		})
		return nil

	case *txpokergrpc.StandUpRequest:
		handler.logger.Info("req:standUp", zap.String("uid", req.Uid.String()))

		if _, err := handler.seatStatusService.StandUp(req.Uid); err != nil {
			handler.logger.Error("user request to stand up but failed", zap.Error(err), zap.String("uid", req.Uid.String()))
			return &core.Response{Err: err}
		}

		handler.msgBus.Broadcast(core.MessageTopic, &txpokergrpc.Message{
			Model: &txpokergrpc.Model{
				SeatStatusGroup: handler.seatStatusGroup.ToProto(),
			},
		})
		return nil

	case *txpokergrpc.SitDownRequest:
		if msg.SeatId < 0 || msg.SeatId > constant.MaxSeatId {
			return &core.Response{Err: status.Errorf(codes.InvalidArgument, "invalid seat id: %v", msg.SeatId)}
		}

		handler.logger.Info("req:sitDown", zap.Object("req", req))

		if err := handler.seatStatusService.SitDown(req.Uid, int(msg.SeatId)); err != nil {
			handler.logger.Error("sitDown failed", zap.Error(err), zap.String("uid", req.Uid.String()))
			return &core.Response{Err: err}
		}

		handler.msgBus.Broadcast(core.MessageTopic, &txpokergrpc.Message{
			Model: &txpokergrpc.Model{
				SeatStatusGroup: handler.seatStatusGroup.ToProto(),
			},
		})
		return nil

	case *txpokergrpc.BuyInRequest:
		if msg.BuyInChip <= 0 {
			return &core.Response{Err: status.Errorf(codes.InvalidArgument, "invalid buy in chip: %v", msg.BuyInChip)}
		}

		handler.logger.Info("req:buyIn", zap.Object("req", req))

		if err := handler.seatStatusService.BuyIn(req.Uid, int(msg.BuyInChip)); err != nil {
			handler.logger.Error("buyIn failed", zap.Error(err), zap.String("uid", req.Uid.String()))
			return &core.Response{Err: err}
		}

		handler.msgBus.Broadcast(core.MessageTopic, &txpokergrpc.Message{
			Model: &txpokergrpc.Model{
				SeatStatusGroup: handler.seatStatusGroup.ToProto(),
			},
		})
		return nil

	case *txpokergrpc.ForceBuyInRequest:
		handler.logger.Info("req:getForceBuyIn", zap.String("uid", req.Uid.String()))

		if data, isExist := handler.forceBuyInGroup.Get(req.Uid); isExist {
			buyInChip := data.GetBuyInChip()
			remainDuration := data.GetExpireTime().Sub(time.Now())

			handler.logger.Info("got forceBuyIn",
				zap.Bool("isExist", isExist),
				zap.Int("buyInChip", buyInChip),
				zap.Duration("remainDuration", remainDuration),
			)

			return &core.Response{
				Msg: &txpokergrpc.ForceBuyInResponse{
					IsBuyIn:    isExist,
					BuyInChip:  int32(buyInChip),
					RemainTime: durationpb.New(remainDuration),
				},
			}
		}

		return &core.Response{
			Msg: &txpokergrpc.ForceBuyInResponse{},
		}

	case *txpokergrpc.SitOutRequest:
		handler.logger.Debug("req:sitOut", zap.String("uid", req.Uid.String()))

		if err := handler.seatStatusService.SitOut(req.Uid); err != nil {
			handler.logger.Error("user request to sit out but failed", zap.Error(err), zap.String("uid", req.Uid.String()))
			return &core.Response{Err: err}
		}

		handler.msgBus.Broadcast(core.MessageTopic, &txpokergrpc.Message{
			Model: &txpokergrpc.Model{
				SeatStatusGroup: handler.seatStatusGroup.ToProto(),
			},
		})
		return nil

	case *txpokergrpc.TopUpRequest:
		if msg.TopUpChip <= 0 {
			return &core.Response{Err: status.Errorf(codes.InvalidArgument, "invalid top up chip: %v", msg.TopUpChip)}
		}

		hint, isPlaying := handler.actionHintGroup.Hints[req.Uid]
		if isPlaying && hint.Action != action.Fold {
			handler.seatStatusGroup.TopUpQueue[req.Uid] = int(msg.TopUpChip)

			handler.logger.Info("req:topUp,queued", zap.Object("req", req))
			handler.msgBus.Broadcast(core.MessageTopic, &txpokergrpc.Message{
				Model: &txpokergrpc.Model{
					SeatStatusGroup: handler.seatStatusGroup.ToProto(),
				},
			})
			return nil
		}

		if _, err := handler.seatStatusService.TopUp(req.Uid, int(msg.TopUpChip)); err != nil {
			handler.logger.Error("user request to top up but failed", zap.Error(err), zap.String("uid", req.Uid.String()))
			return &core.Response{Err: err}
		}

		handler.logger.Info("user request to top up", zap.String("uid", req.Uid.String()), zap.Object("req", req))
		return nil

	case *txpokergrpc.UpdateWaitBBSettingRequest:
		if _, ok := handler.playSettingGroup.Data[req.Uid]; !ok {
			return &core.Response{Err: status.Errorf(codes.NotFound, "play setting does not exist for uid: %v", req.Uid.String())}
		}

		handler.playSettingGroup.Data[req.Uid].WaitBB = msg.WaitBB

		handler.logger.Info("req:waitBB setting", zap.Object("playSetting", handler.playSettingGroup.Data[req.Uid]))
		handler.msgBus.Unicast(req.Uid, core.MessageTopic, &txpokergrpc.Message{
			Model: &txpokergrpc.Model{
				PlaySetting: handler.playSettingGroup.Data[req.Uid].ToProto(),
			},
		})
		return nil

	case *txpokergrpc.UpdateAutoTopUpSettingRequest:
		if _, ok := handler.playSettingGroup.Data[req.Uid]; !ok {
			return &core.Response{Err: status.Errorf(codes.NotFound, "play setting does not exist for uid: %v", req.Uid.String())}
		}

		if msg.AutoTopUpThresholdPercent < 0.0 || msg.AutoTopUpThresholdPercent > 1.0 {
			return &core.Response{Err: status.Errorf(codes.InvalidArgument, "invalid auto top up threshold percent: %v", msg.AutoTopUpThresholdPercent)}
		}

		if msg.AutoTopUpChipPercent < 0.0 || msg.AutoTopUpChipPercent > 1.0 {
			return &core.Response{Err: status.Errorf(codes.InvalidArgument, "invalid auto top up chip percent: %v", msg.AutoTopUpChipPercent)}
		}

		handler.playSettingGroup.Data[req.Uid].AutoTopUp = msg.AutoTopUp
		handler.playSettingGroup.Data[req.Uid].AutoTopUpThresholdPercent = msg.AutoTopUpThresholdPercent
		handler.playSettingGroup.Data[req.Uid].AutoTopUpChipPercent = msg.AutoTopUpChipPercent

		handler.logger.Info("req:autoTopUp setting", zap.Object("playSetting", handler.playSettingGroup.Data[req.Uid]))
		handler.msgBus.Unicast(req.Uid, core.MessageTopic, &txpokergrpc.Message{
			Model: &txpokergrpc.Model{
				PlaySetting: handler.playSettingGroup.Data[req.Uid].ToProto(),
			},
		})
		return nil

	default:
		handler.logger.Debug("ignored not supported request", zap.Object("req", req))
		return nil
	}
}
