package service

import (
	"context"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	txpokerservice "card-game-server-prototype/pkg/game/txpoker/service"
	seatstatus2 "card-game-server-prototype/pkg/game/txpoker/type/seatstatus"
	session2 "card-game-server-prototype/pkg/game/zoomtxpoker/game/session"
	"card-game-server-prototype/pkg/util"
	"github.com/qmuntal/stateless"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
	"time"
)

type BaseSeatStatusService struct {
	*txpokerservice.BaseSeatStatusService
	seatStatusGroup *model2.SeatStatusGroup
	forceBuyInGroup *model2.ForceBuyInGroup
	gameSetting     *model2.GameSetting
	gameInfo        *model2.GameInfo
	msgBus          core.MsgBus
	logger          *zap.Logger
}

func ProvideBaseSeatStatusService(
	baseSeatStatusService *txpokerservice.BaseSeatStatusService,
	seatStatusGroup *model2.SeatStatusGroup,
	forceBuyInGroup *model2.ForceBuyInGroup,
	gameSetting *model2.GameSetting,
	gameInfo *model2.GameInfo,
	msgBus core.MsgBus,
	loggerFactory *util.LoggerFactory,
) *BaseSeatStatusService {
	return &BaseSeatStatusService{
		BaseSeatStatusService: baseSeatStatusService,
		seatStatusGroup:       seatStatusGroup,
		gameSetting:           gameSetting,
		gameInfo:              gameInfo,
		forceBuyInGroup:       forceBuyInGroup,
		msgBus:                msgBus,
		logger:                loggerFactory.Create("BaseSeatStatusService"),
	}
}

func (s *BaseSeatStatusService) NewSeatStatus(uid core.Uid) *model2.SeatStatus {
	fsm := stateless.NewStateMachine(seatstatus2.StandingState)
	seatStatus := &model2.SeatStatus{
		Uid:                  uid,
		Chip:                 0,
		FSM:                  fsm,
		SitOutStartTime:      time.Now(),
		CancelSitOutTimer:    func() {},
		CancelReservingTimer: func() {},
		ShouldPlaceBB:        false,
	}

	fsm.SetTriggerParameters(seatstatus2.SitDownTrigger, reflect.TypeOf(1))
	fsm.SetTriggerParameters(seatstatus2.BuyInTrigger, reflect.TypeOf(1))
	fsm.SetTriggerParameters(seatstatus2.StandUpTrigger, reflect.TypeOf(make(chan struct{})))

	fsm.Configure(seatstatus2.StandingState).
		Permit(seatstatus2.SitDownTrigger, seatstatus2.ReservingState).
		OnEntry(func(ctx context.Context, args ...any) error {
			s.seatStatusGroup.TableUids = lo.OmitByValues(s.seatStatusGroup.TableUids, core.UidList{uid})
			return nil
		}).
		OnEntryFrom(seatstatus2.StandUpTrigger, func(ctx context.Context, args ...any) error {
			cashOutDone := args[0].(chan struct{})

			countIdleRounds := seatStatus.CountIdleRounds

			gameResult := &session2.GameResultSession{
				GameId:          s.gameInfo.RoundId,
				Uid:             uid,
				Chip:            seatStatus.Chip,
				CountIdleRounds: countIdleRounds,
			}

			s.logger.Debug("stand up, sending game result",
				zap.Object("gameResult", gameResult),
				util.DebugField(zap.Object("seatStatus", seatStatus)),
			)

			// Make user show buy in, not force buy in, GCS-4299
			s.forceBuyInGroup.Delete(uid)

			seatStatus.Chip = 0
			seatStatus.CountIdleRounds = 0

			s.msgBus.Unicast(uid, session2.GameResultTopic, gameResult)
			close(cashOutDone)
			return nil
		})

	fsm.Configure(seatstatus2.ReservingState).
		Permit(seatstatus2.BuyInTrigger, seatstatus2.PlayingState).
		OnEntry(func(ctx context.Context, args ...any) error {
			seatId := args[0].(int)

			if seatId < 0 || seatId > s.gameSetting.TableSize-1 {
				return status.Errorf(codes.InvalidArgument, "uid %v try to sit down with seatId out of range: %v, must <= %v", uid, seatId, s.gameSetting.TableSize-1)
			}

			if _, ok := s.seatStatusGroup.TableUids[seatId]; ok {
				return status.Errorf(codes.ResourceExhausted, "uid %v tries to take seatId %v but already taken by uid %v", uid, seatId, s.seatStatusGroup.TableUids[seatId])
			}

			curSeatIds := lo.Keys(lo.PickByValues(s.seatStatusGroup.TableUids, core.UidList{uid}))
			if len(curSeatIds) > 0 && curSeatIds[0] != seatId {
				return status.Errorf(codes.FailedPrecondition, "uid %v cannot sit down at seatId %v, since already at seatId %v", uid, seatId, curSeatIds[0])
			}

			s.seatStatusGroup.TableUids[seatId] = uid
			return nil
		})

	fsm.Configure(seatstatus2.PlayingState).
		Permit(seatstatus2.StandUpTrigger, seatstatus2.StandingState).
		PermitReentry(seatstatus2.RoundEndTrigger).
		OnEntryFrom(seatstatus2.BuyInTrigger, func(ctx context.Context, args ...any) error {
			// zoom 的 seatStatus 是僅有系統操作，使用者操作只會碰觸到 pool 的 scope。
			// so, 這裡不需要驗證 buyInChip 是否合法，只需要將 chip 加上去即可。
			buyInChip := args[0].(int)

			s.logger.Info("start buying in",
				zap.String("uid", uid.String()),
				zap.Int("buyInChip", buyInChip),
				zap.Int("chip", seatStatus.Chip),
			)

			seatStatus.Chip += buyInChip

			return nil
		})

	return seatStatus
}
