package handler

import (
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/zoomtxpoker/pool/model"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/service"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/type/participant"
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

	participantGroup   *model2.ParticipantGroup
	playSettingGroup   *model2.PlaySettingGroup
	forceBuyInGroup    *model2.ForceBuyInGroup
	participantService *service.ParticipantService

	game   core.Game
	msgBus core.MsgBus
	logger *zap.Logger
}

func ProvideRequestHandler(
	participantGroup *model2.ParticipantGroup,
	playSettingGroup *model2.PlaySettingGroup,
	forceBuyInGroup *model2.ForceBuyInGroup,
	participantService *service.ParticipantService,

	game core.Game,
	msgBus core.MsgBus,
	loggerFactory *util.LoggerFactory,
) *RequestHandler {
	return &RequestHandler{
		participantGroup:   participantGroup,
		playSettingGroup:   playSettingGroup,
		forceBuyInGroup:    forceBuyInGroup,
		participantService: participantService,

		game:   game,
		msgBus: msgBus,
		logger: loggerFactory.Create("RequestHandler"),
	}
}

func (handler *RequestHandler) HandleRequest(req *core.Request) *core.Response {
	switch msg := req.Msg.(type) {
	case *txpokergrpc.BuyInRequest:
		if msg.BuyInChip <= 0 {
			return &core.Response{Err: status.Errorf(codes.InvalidArgument, "invalid buy in chip: %v", msg.BuyInChip)}
		}

		handler.logger.Info("user request to buyIn", zap.Object("req", req))

		if err := handler.participantService.BuyIn(req.Uid, int(msg.BuyInChip)); err != nil {
			handler.logger.Error("buyIn failed", zap.Error(err), zap.Object("req", req))
			return &core.Response{Err: err}
		}

		handler.msgBus.Unicast(req.Uid, core.MessageTopic, &txpokergrpc.Message{
			Model: &txpokergrpc.Model{
				Participant: handler.participantGroup.Data[req.Uid].ToProto(),
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

	case *txpokergrpc.TopUpRequest:
		if msg.TopUpChip <= 0 {
			return &core.Response{Err: status.Errorf(codes.InvalidArgument, "invalid top up chip: %v", msg.TopUpChip)}
		}

		if part, ok := handler.participantGroup.Data[req.Uid]; ok && part.FSM.MustState().(participant.State) == participant.PlayingState {
			part.QueuedTopUpChip = int(msg.TopUpChip)
			handler.logger.Info("user request to top up, queued", zap.Object("req", req))
			handler.msgBus.Unicast(req.Uid, core.MessageTopic, &txpokergrpc.Message{
				Model: &txpokergrpc.Model{
					Participant: part.ToProto(),
				},
			})
			return nil
		}

		if _, err := handler.participantService.TopUp(req.Uid, int(msg.TopUpChip)); err != nil {
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
