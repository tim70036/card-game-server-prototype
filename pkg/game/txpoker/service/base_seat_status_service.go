package service

import (
	"context"
	"fmt"
	commonapi "card-game-server-prototype/pkg/common/api"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/common/type/gametype"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/constant"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/type/action"
	"card-game-server-prototype/pkg/game/txpoker/type/event"
	role2 "card-game-server-prototype/pkg/game/txpoker/type/role"
	seatstatus2 "card-game-server-prototype/pkg/game/txpoker/type/seatstatus"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"
	"reflect"
	"time"

	"github.com/qmuntal/stateless"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BaseSeatStatusService struct {
	gameController core.GameController
	msgBus         core.MsgBus

	userGroup         *commonmodel.UserGroup
	seatStatusGroup   *model2.SeatStatusGroup
	actionHintGroup   *model2.ActionHintGroup
	playerGroup       *model2.PlayerGroup
	playSettingGroup  *model2.PlaySettingGroup
	statsGroup        *model2.StatsGroup
	gameSetting       *model2.GameSetting
	gameInfo          *model2.GameInfo
	userAPI           commonapi.UserAPI
	forceBuyInGroup   *model2.ForceBuyInGroup
	eventGroup        *model2.EventGroup
	tableProfitsGroup *model2.TableProfitsGroup

	roundStarted bool
	isToppingUp  map[core.Uid]struct{}
	logger       *zap.Logger
}

func ProvideBaseSeatStatusService(
	gameController core.GameController,
	msgBus core.MsgBus,
	userGroup *commonmodel.UserGroup,
	seatStatusGroup *model2.SeatStatusGroup,
	actionHintGroup *model2.ActionHintGroup,
	playerGroup *model2.PlayerGroup,
	playSettingGroup *model2.PlaySettingGroup,
	statsGroup *model2.StatsGroup,
	gameSetting *model2.GameSetting,
	gameInfo *model2.GameInfo,
	userAPI commonapi.UserAPI,
	loggerFactory *util.LoggerFactory,
	forceBuyInGroup *model2.ForceBuyInGroup,
	eventGroup *model2.EventGroup,
	tableProfitsGroup *model2.TableProfitsGroup,
) *BaseSeatStatusService {
	return &BaseSeatStatusService{
		gameController:    gameController,
		msgBus:            msgBus,
		userGroup:         userGroup,
		seatStatusGroup:   seatStatusGroup,
		actionHintGroup:   actionHintGroup,
		playerGroup:       playerGroup,
		playSettingGroup:  playSettingGroup,
		statsGroup:        statsGroup,
		gameSetting:       gameSetting,
		gameInfo:          gameInfo,
		userAPI:           userAPI,
		forceBuyInGroup:   forceBuyInGroup,
		eventGroup:        eventGroup,
		tableProfitsGroup: tableProfitsGroup,

		roundStarted: false,
		isToppingUp:  make(map[core.Uid]struct{}),
		logger:       loggerFactory.Create("BaseSeatStatusService"),
	}
}

func (s *BaseSeatStatusService) SitDown(uid core.Uid, seatId int) error {
	if _, ok := s.seatStatusGroup.Status[uid]; !ok {
		return status.Errorf(codes.NotFound, "cannot found uid %v in seat status group with uids: %v", uid, lo.Keys(s.seatStatusGroup.Status))
	}

	seatStatus := s.seatStatusGroup.Status[uid]
	if err := seatStatus.FSM.Fire(seatstatus2.SitDownTrigger, seatId); err != nil {
		return err
	}

	s.logger.Debug("sit down", zap.Object("seatStatus", seatStatus))
	return nil
}

func (s *BaseSeatStatusService) BuyIn(uid core.Uid, buyInChip int) error {
	if _, ok := s.seatStatusGroup.Status[uid]; !ok {
		return status.Errorf(codes.NotFound, "cannot found uid %v in seat status group with uids: %v", uid, lo.Keys(s.seatStatusGroup.Status))
	}

	seatStatus := s.seatStatusGroup.Status[uid]
	if err := seatStatus.FSM.Fire(seatstatus2.BuyInTrigger, buyInChip); err != nil {
		return err
	}

	s.logger.Debug("buying in",
		zap.String("uid", uid.String()),
		zap.Int("Chip", seatStatus.Chip),
		util.DebugField(zap.Object("seatStatus", seatStatus)),
	)
	return nil
}

func (s *BaseSeatStatusService) StandUp(uid core.Uid) (<-chan struct{}, error) {
	if _, ok := s.seatStatusGroup.Status[uid]; !ok {
		return nil, status.Errorf(codes.NotFound, "cannot found uid %v in seat status group with uids: %v", uid, lo.Keys(s.seatStatusGroup.Status))
	}

	seatStatus := s.seatStatusGroup.Status[uid]
	cashOutDone := make(chan struct{})

	if err := seatStatus.FSM.Fire(seatstatus2.StandUpTrigger, cashOutDone); err != nil {
		return nil, err
	}

	s.logger.Debug("standing up",
		zap.String("uid", uid.String()),
		zap.Int("Chip", seatStatus.Chip),
		util.DebugField(zap.Object("seatStatus", seatStatus)),
	)
	return cashOutDone, nil
}

func (s *BaseSeatStatusService) SitOut(uid core.Uid) error {
	if _, ok := s.seatStatusGroup.Status[uid]; !ok {
		return status.Errorf(codes.NotFound, "cannot found uid %v in seat status group with uids: %v", uid, lo.Keys(s.seatStatusGroup.Status))
	}

	seatStatus := s.seatStatusGroup.Status[uid]

	if err := seatStatus.FSM.Fire(seatstatus2.SitOutTrigger); err != nil {
		return err
	}

	s.logger.Debug("sitting out",
		zap.String("uid", uid.String()),
		zap.Int("Chip", seatStatus.Chip),
		util.DebugField(zap.Object("seatStatus", seatStatus)),
	)
	return nil
}

func (s *BaseSeatStatusService) TopUp(uid core.Uid, topUpChip int) (<-chan struct{}, error) {
	if topUpChip <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "topUpChip must > 0, got %v", topUpChip)
	}

	if _, ok := s.isToppingUp[uid]; ok {
		return nil, status.Errorf(codes.AlreadyExists, "uid %v is already topping up", uid)
	}

	seatStatus, ok := s.seatStatusGroup.Status[uid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "cannot found uid %v in seat status group with uids: %v", uid, lo.Keys(s.seatStatusGroup.Status))
	}

	user, ok := s.userGroup.Data[uid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "cannot found uid %v in user group with uids: %v", uid, lo.Keys(s.userGroup.Data))
	}

	if topUpChip > user.Cash {
		return nil, status.Errorf(codes.InvalidArgument, "uid %v has not enough cash %v to top up %v", uid, user.Cash, topUpChip)
	}

	if !lo.Contains(
		[]seatstatus2.SeatStatusState{seatstatus2.JoiningState, seatstatus2.SittingOutState, seatstatus2.WaitingState, seatstatus2.PlayingState},
		seatStatus.FSM.MustState().(seatstatus2.SeatStatusState),
	) {
		return nil, status.Errorf(codes.FailedPrecondition, "uid %v is not in joining, sitting out, waiting or playing state, got %v", uid, seatStatus.FSM.MustState().(seatstatus2.SeatStatusState).String())
	}

	maxEnterLimitChip := s.gameSetting.MaxEnterLimitBB * s.gameSetting.BigBlind
	if seatStatus.Chip >= maxEnterLimitChip {
		return nil, status.Errorf(codes.FailedPrecondition, "uid %v already has chip %v that >= maxEnterLimitChip %v", uid, seatStatus.Chip, maxEnterLimitChip)
	}

	topUpChip = lo.Ternary(
		seatStatus.Chip+topUpChip > maxEnterLimitChip,
		maxEnterLimitChip-seatStatus.Chip,
		topUpChip,
	)

	topUpDone := make(chan struct{})
	s.isToppingUp[uid] = struct{}{}

	s.logger.Info("start top up",
		zap.String("uid", uid.String()),
		zap.Int("Chip", seatStatus.Chip),
		util.DebugField(zap.Object("seatStatus", seatStatus)),
	)

	go func() {
		err := s.userAPI.ExchangeChip(uid, gametype.TXPoker, topUpChip*-1)
		go s.fetchUserCash(uid)

		s.gameController.RunTask(func() {
			defer func() {
				delete(s.isToppingUp, uid)
				close(topUpDone)
			}()

			if err != nil {
				s.logger.Error("failed top up", zap.Error(err), zap.Object("seatStatus", seatStatus), zap.Int("topUpChip", topUpChip))
				return
			}

			// Edge case: top up close to game start, user is in playing state and not fold yet.
			// https://game-soul-technology.atlassian.net/browse/GCS-2899
			revertTopUp := func() {
				if err := s.userAPI.ExchangeChip(uid, gametype.TXPoker, topUpChip); err != nil {
					s.logger.Error("failed to revert top up chips",
						zap.Error(err),
						zap.Object("seatStatus", seatStatus),
						zap.Int("revertTopUpChip", topUpChip),
					)
				}
				return
			}

			if _, ok := s.userGroup.Data[uid]; !ok {
				s.logger.Warn("revert top up chips since user has left the room",
					zap.Object("seatStatus", seatStatus),
					zap.Int("revertTopUpChip", topUpChip),
				)
				go revertTopUp()
				return
			}

			if !lo.Contains(
				[]seatstatus2.SeatStatusState{seatstatus2.JoiningState, seatstatus2.SittingOutState, seatstatus2.WaitingState, seatstatus2.PlayingState},
				seatStatus.FSM.MustState().(seatstatus2.SeatStatusState),
			) {

				s.logger.Warn("revert top up chips since user is not in joining, sitting out, waiting or playing state",
					zap.Object("seatStatus", seatStatus),
					zap.Int("revertTopUpChip", topUpChip),
				)

				go revertTopUp()
				return
			}

			if seatStatus.FSM.MustState().(seatstatus2.SeatStatusState) == seatstatus2.PlayingState &&
				s.actionHintGroup.Hints[uid].Action != action.Fold {

				s.logger.Warn("revert top up chips since user is in playing state and not fold yet",
					zap.Object("seatStatus", seatStatus),
					zap.Int("revertTopUpChip", topUpChip),
				)

				go revertTopUp()
				return
			}

			seatStatus.Chip += topUpChip

			s.gameController.RunTask(func() {
				if tableProfits, ok := s.tableProfitsGroup.Get(uid); ok {
					tableProfits.SumBuyInChips += topUpChip
					s.tableProfitsGroup.Save(tableProfits)
				}
			})

			s.logger.Info("done top up",
				zap.String("uid", uid.String()),
				zap.Int("Chip", seatStatus.Chip),
				util.DebugField(zap.Object("seatStatus", seatStatus)),
			)
			s.msgBus.Broadcast(core.MessageTopic, &txpokergrpc.Message{
				Model: &txpokergrpc.Model{
					SeatStatusGroup: s.seatStatusGroup.ToProto(),
				},
			})
		})
	}()

	return topUpDone, nil
}

func (s *BaseSeatStatusService) IsReadyToStartRound() bool {
	joiningCnt := lo.CountBy(lo.Values(s.seatStatusGroup.Status), func(status *model2.SeatStatus) bool {
		return seatstatus2.JoiningState.IsEqual(status.FSM.MustState())
	})

	return joiningCnt >= constant.StartGameUserCount
}

func (s *BaseSeatStatusService) RoundStarted() error {
	s.roundStarted = true

	candidateUids := lo.FilterMap(lo.Values(s.seatStatusGroup.Status), func(s *model2.SeatStatus, _ int) (core.Uid, bool) {
		return s.Uid, s.FSM.MustState().(seatstatus2.SeatStatusState) == seatstatus2.JoiningState
	})

	lastBBSeatId := -1
	if s.playerGroup.LastBBPlayer != nil {
		lastBBSeatId = s.playerGroup.LastBBPlayer.SeatId
	}

	roleAssignment, err := role2.EvalRoleAssignment(candidateUids, s.seatStatusGroup.TableUids, lastBBSeatId)
	if err != nil {
		return err
	}

	var playingUids core.UidList = lo.Filter(
		lo.Keys(roleAssignment),
		func(uid core.Uid, _ int) bool {
			seatStatus := s.seatStatusGroup.Status[uid]
			return s.gameInfo.StreakRoundCount <= 1 || // For the first round, all players are playing.
				!seatStatus.ShouldPlaceBB || // All players who have placed BB are playing.
				roleAssignment[uid] == role2.BB || // BB is playing.
				(!s.playSettingGroup.Data[uid].WaitBB && roleAssignment[uid] != role2.SB) // If playSetting wait bb off and is not sb, then playing.
		},
	)

	// Special case: If only BB player is able to play, then let SB
	// play too (don't care if he has not placed BB yet and is waiting
	// BB). Otherwise, the game will have insufficient player.
	if len(playingUids) < 2 {
		sbPlayerUid := lo.Invert(roleAssignment)[role2.SB]
		playingUids = append(playingUids, sbPlayerUid)
	}

	for _, uid := range playingUids {
		seatStatus := s.seatStatusGroup.Status[uid]
		if err := seatStatus.FSM.Fire(seatstatus2.RoundStartedTrigger); err != nil {
			return fmt.Errorf("failed to fire round started trigger for uid %v: %w", uid, err)
		}

		totalRoundCount := s.statsGroup.Data[uid].EventAmountSum[event.Game]
		if totalRoundCount%s.gameSetting.ExtraTurnRefillIntervalRound == 0 &&
			totalRoundCount > 0 {
			seatStatus.ActionExtraDuration = lo.Ternary(
				seatStatus.ActionExtraDuration+s.gameSetting.RefillExtraTurnDuration > s.gameSetting.MaxExtraTurnDuration,
				s.gameSetting.MaxExtraTurnDuration,
				seatStatus.ActionExtraDuration+s.gameSetting.RefillExtraTurnDuration,
			)
			s.logger.Debug("refill action extra duration", zap.Object("seatStatus", seatStatus))
		}
	}

	s.logger.Debug("start round done", zap.Array("playingUids", playingUids), zap.Object("seatStatus", s.seatStatusGroup))
	return nil
}

func (s *BaseSeatStatusService) RoundEnd() error {
	s.roundStarted = false

	var joiningUids core.UidList = lo.FilterMap(lo.Values(s.seatStatusGroup.Status), func(s *model2.SeatStatus, _ int) (core.Uid, bool) {
		seatStatusState := s.FSM.MustState().(seatstatus2.SeatStatusState)
		return s.Uid, seatStatusState == seatstatus2.WaitingState || seatStatusState == seatstatus2.PlayingState
	})

	for uid, player := range s.playerGroup.Data {
		// user may leave room (destroy user)
		seatStatus, ok := s.seatStatusGroup.Status[uid]
		if !ok {
			continue
		}

		if seatStatus.FSM.MustState().(seatstatus2.SeatStatusState) != seatstatus2.PlayingState {
			continue
		}

		if !player.IsIdle {
			s.seatStatusGroup.Status[uid].CountIdleRounds = 0
			s.logger.Info("reset idle round",
				zap.String("uid", uid.String()),
				zap.Object("seatStatus", seatStatus))
		}

		if player.IsIdle && player.Role == role2.BB {
			s.seatStatusGroup.Status[uid].CountIdleRounds++
			s.logger.Info("increase idle round",
				zap.String("uid", uid.String()),
				zap.Object("seatStatus", seatStatus))
		}
	}

	for _, uid := range joiningUids {
		seatStatus := s.seatStatusGroup.Status[uid]
		if err := seatStatus.FSM.Fire(seatstatus2.RoundEndTrigger); err != nil {
			return fmt.Errorf("failed to fire round end trigger for uid %v: %w", uid, err)
		}
	}

	s.logger.Debug("end round done", zap.Array("joiningUids", joiningUids), zap.Object("seatStatus", s.seatStatusGroup))
	return nil
}

func (s *BaseSeatStatusService) NewSeatStatus(uid core.Uid) *model2.SeatStatus {
	fsm := stateless.NewStateMachine(seatstatus2.StandingState)
	seatStatus := &model2.SeatStatus{
		Uid:                  uid,
		Chip:                 0,
		FSM:                  fsm,
		SitOutStartTime:      time.Now(),
		SitOutDuration:       s.gameSetting.InitialSitOutDuration,
		CancelSitOutTimer:    func() {},
		CancelReservingTimer: func() {},
		ActionExtraDuration:  s.gameSetting.InitialExtraTurnDuration,
		ShouldPlaceBB:        true,
		CountIdleRounds:      0,
	}

	fsm.SetTriggerParameters(seatstatus2.SitDownTrigger, reflect.TypeOf(1))
	fsm.SetTriggerParameters(seatstatus2.BuyInTrigger, reflect.TypeOf(1))
	fsm.SetTriggerParameters(seatstatus2.StandUpTrigger, reflect.TypeOf(make(chan struct{})))

	fsm.Configure(seatstatus2.StandingState).
		Permit(seatstatus2.SitDownTrigger, seatstatus2.ReservingState).
		PermitReentry(seatstatus2.StandUpTrigger).
		OnEntry(func(ctx context.Context, args ...any) error {
			s.seatStatusGroup.TableUids = lo.OmitByValues(s.seatStatusGroup.TableUids, core.UidList{uid})
			seatStatus.ShouldPlaceBB = true
			seatStatus.CountIdleRounds = 0

			return nil
		}).
		OnEntryFrom(seatstatus2.StandUpTrigger, func(ctx context.Context, args ...any) error {
			cashOutDone := args[0].(chan struct{})
			close(cashOutDone)
			return nil
		}).
		OnEntryFrom(seatstatus2.FailedTrigger, func(ctx context.Context, args ...any) error {
			s.msgBus.Broadcast(core.MessageTopic, &txpokergrpc.Message{
				Model: &txpokergrpc.Model{
					SeatStatusGroup: s.seatStatusGroup.ToProto(),
				},
			})
			return nil
		})

	fsm.Configure(seatstatus2.ReservingState).
		Permit(seatstatus2.BuyInTrigger, seatstatus2.BuyingInState).
		Permit(seatstatus2.StandUpTrigger, seatstatus2.StandingState).
		Permit(seatstatus2.TimeoutTrigger, seatstatus2.StandingState).
		Permit(seatstatus2.FailedTrigger, seatstatus2.StandingState).
		OnEntry(func(ctx context.Context, args ...any) error {
			seatId := args[0].(int)
			if seatId < 0 || seatId > s.gameSetting.TableSize-1 {
				s.gameController.RunTask(func() {
					if err := fsm.Fire(seatstatus2.FailedTrigger); err != nil {
						s.logger.Error("failed to fire failed trigger", zap.Error(err), zap.Object("seatStatus", seatStatus))
					}
				})
				return status.Errorf(codes.InvalidArgument, "uid %v try to sit down with seatId out of range: %v, must <= %v", uid, seatId, s.gameSetting.TableSize-1)
			}

			if _, ok := s.seatStatusGroup.TableUids[seatId]; ok {
				s.gameController.RunTask(func() {
					if err := fsm.Fire(seatstatus2.FailedTrigger); err != nil {
						s.logger.Error("failed to fire failed trigger", zap.Error(err), zap.Object("seatStatus", seatStatus))
					}
				})
				return status.Errorf(codes.ResourceExhausted, "uid %v tries to take seatId %v but already taken by uid %v", uid, seatId, s.seatStatusGroup.TableUids[seatId])
			}

			curSeatIds := lo.Keys(lo.PickByValues(s.seatStatusGroup.TableUids, core.UidList{uid}))
			if len(curSeatIds) > 0 && curSeatIds[0] != seatId {
				s.gameController.RunTask(func() {
					if err := fsm.Fire(seatstatus2.FailedTrigger); err != nil {
						s.logger.Error("failed to fire failed trigger", zap.Error(err), zap.Object("seatStatus", seatStatus))
					}
				})
				return status.Errorf(codes.FailedPrecondition, "uid %v cannot sit down at seatId %v, since already at seatId %v", uid, seatId, curSeatIds[0])
			}

			s.seatStatusGroup.TableUids[seatId] = uid

			seatStatus.CancelReservingTimer = s.gameController.RunTimer(constant.ReservingGracefulPeriod, func() {
				s.logger.Debug("reserving timeout", zap.Object("seatStatus", seatStatus))
				if err := seatStatus.FSM.Fire(seatstatus2.TimeoutTrigger); err != nil {
					s.logger.Error("failed to fire timeout trigger", zap.Error(err), zap.Object("seatStatus", seatStatus))
				}

				// Notify client that he's force stand up by server.
				s.msgBus.Unicast(uid, core.MessageTopic, &txpokergrpc.Message{
					Event: &txpokergrpc.Event{
						Warning: &txpokergrpc.Warning{Reason: txpokergrpc.WarningReason_RESERVE_TIMEOUT},
					},
				})

				s.msgBus.Broadcast(core.MessageTopic, &txpokergrpc.Message{
					Model: &txpokergrpc.Model{
						SeatStatusGroup: s.seatStatusGroup.ToProto(),
					},
				})
			})
			return nil
		}).
		OnExit(func(ctx context.Context, args ...any) error {
			seatStatus.CancelReservingTimer()
			return nil
		})

	fsm.Configure(seatstatus2.BuyingInState).
		Permit(seatstatus2.SuccessTrigger, seatstatus2.WaitingState, func(ctx context.Context, args ...any) bool {
			return s.roundStarted // game is started
		}).
		Permit(seatstatus2.SuccessTrigger, seatstatus2.JoiningState, func(ctx context.Context, args ...any) bool {
			return !s.roundStarted // game is not started
		}).
		Permit(seatstatus2.FailedTrigger, seatstatus2.StandingState).
		OnEntry(func(ctx context.Context, args ...any) error {
			buyInChip := args[0].(int)
			_, isForceBuyInExist := s.forceBuyInGroup.Get(uid)

			s.logger.Info("start buying in",
				zap.String("uid", uid.String()),
				zap.Int("buyInChip", buyInChip),
				zap.Int("chip", seatStatus.Chip),
				zap.Bool("forceBuyInExist", isForceBuyInExist),
			)

			if !isForceBuyInExist &&
				(buyInChip < s.gameSetting.MinEnterLimitBB*s.gameSetting.BigBlind ||
					buyInChip > s.gameSetting.MaxEnterLimitBB*s.gameSetting.BigBlind) {
				s.gameController.RunTask(func() {
					if err := fsm.Fire(seatstatus2.FailedTrigger); err != nil {
						s.logger.Error("failed to fire failed trigger", zap.Error(err), zap.Object("seatStatus", seatStatus))
					}
				})
				return status.Errorf(codes.InvalidArgument, "uid %v try to sit down with invalid buyInChip %v, not in range %v ~ %v", uid, buyInChip, s.gameSetting.MinEnterLimitBB*s.gameSetting.BigBlind, s.gameSetting.MaxEnterLimitBB*s.gameSetting.BigBlind)
			}

			// To avoid blocking game main loop, make it async.
			go func() {
				if err := s.userAPI.ExchangeChip(uid, gametype.TXPoker, buyInChip*-1); err != nil {
					s.gameController.RunTask(func() {
						if err := fsm.Fire(seatstatus2.FailedTrigger); err != nil {
							s.logger.Error("failed to fire failed trigger", zap.Error(err), zap.Object("seatStatus", seatStatus))
						}

						s.logger.Error("failed buying in", zap.Error(err), zap.Object("seatStatus", seatStatus))
						s.msgBus.Broadcast(core.MessageTopic, &txpokergrpc.Message{
							Model: &txpokergrpc.Model{
								SeatStatusGroup: s.seatStatusGroup.ToProto(),
							},
						})
					})
					return
				}

				go s.fetchUserCash(uid)
				s.gameController.RunTask(func() {
					seatStatus.Chip += buyInChip

					if tableProfits, ok := s.tableProfitsGroup.Get(uid); ok {
						tableProfits.SumBuyInChips += buyInChip
						s.tableProfitsGroup.Save(tableProfits)
					}

					if err := fsm.Fire(seatstatus2.SuccessTrigger); err != nil {
						s.logger.Error("failed to fire success trigger", zap.Error(err), zap.Object("seatStatus", seatStatus))
					}

					s.logger.Info("done buying in",
						zap.String("uid", uid.String()),
						zap.Int("chip", seatStatus.Chip),
						util.DebugField(zap.Object("seatStatus", seatStatus)),
					)
					s.msgBus.Broadcast(core.MessageTopic, &txpokergrpc.Message{
						Model: &txpokergrpc.Model{
							SeatStatusGroup: s.seatStatusGroup.ToProto(),
						},
					})
				})
			}()

			return nil
		})

	fsm.Configure(seatstatus2.CashingOutState).
		Permit(seatstatus2.SuccessTrigger, seatstatus2.StandingState).
		Permit(seatstatus2.FailedTrigger, seatstatus2.StandingState).
		OnEntry(func(ctx context.Context, args ...any) error {
			// To avoid blocking game main loop, make it async.
			go func() {
				s.logger.Info("start cashing out",
					zap.String("uid", seatStatus.Uid.String()),
					zap.Int("Chip", seatStatus.Chip),
					util.DebugField(zap.Object("seatStatus", seatStatus)),
				)

				if seatStatus.Chip != 0 {
					if err := s.userAPI.ExchangeChip(uid, gametype.TXPoker, seatStatus.Chip); err != nil {
						go s.fetchUserCash(uid)
						s.gameController.RunTask(func() {
							seatStatus.Chip = 0
							if err := fsm.Fire(seatstatus2.FailedTrigger); err != nil {
								s.logger.Error("failed to fire failed trigger", zap.Error(err), zap.Object("seatStatus", seatStatus))
							}

							s.logger.Error("failed cashing out", zap.Error(err), zap.Object("seatStatus", seatStatus))
							s.msgBus.Broadcast(core.MessageTopic, &txpokergrpc.Message{
								Model: &txpokergrpc.Model{
									SeatStatusGroup: s.seatStatusGroup.ToProto(),
								},
							})

							if len(args) > 0 {
								cashOutDone := args[0].(chan struct{})
								close(cashOutDone)
							}
						})
						return
					}

					// cash out chip >= min enter limit, allow force buy in
					if seatStatus.Chip >= s.gameSetting.MinEnterLimitBB*s.gameSetting.BigBlind {
						s.gameController.RunTask(func() {
							s.forceBuyInGroup.Set(uid, seatStatus.Chip, time.Now())

							s.logger.Debug("Set ForceBuyIn",
								zap.String("uid", uid.String()),
								zap.Int("Chip", seatStatus.Chip),
							)
						})
					}
				}

				go s.fetchUserCash(uid)
				s.gameController.RunTask(func() {
					seatStatus.Chip = 0
					if err := fsm.Fire(seatstatus2.SuccessTrigger); err != nil {
						s.logger.Error("failed to fire success trigger", zap.Error(err), zap.Object("seatStatus", seatStatus))
					}

					s.logger.Info("done cashing out",
						zap.Int("Chip", seatStatus.Chip),
						zap.String("uid", seatStatus.Uid.String()),
						util.DebugField(zap.Object("seatStatus", seatStatus)),
					)
					s.msgBus.Broadcast(core.MessageTopic, &txpokergrpc.Message{
						Model: &txpokergrpc.Model{
							SeatStatusGroup: s.seatStatusGroup.ToProto(),
						},
					})

					if len(args) > 0 {
						cashOutDone := args[0].(chan struct{})
						close(cashOutDone)
					}
				})
			}()
			return nil
		})

	fsm.Configure(seatstatus2.WaitingState).
		Permit(seatstatus2.StandUpTrigger, seatstatus2.CashingOutState).
		Permit(seatstatus2.RoundEndTrigger, seatstatus2.JoiningState)

	fsm.Configure(seatstatus2.JoiningState).
		Permit(seatstatus2.StandUpTrigger, seatstatus2.CashingOutState).
		Permit(seatstatus2.SitOutTrigger, seatstatus2.SittingOutState).
		Permit(seatstatus2.RoundStartedTrigger, seatstatus2.PlayingState)

	fsm.Configure(seatstatus2.PlayingState).
		Permit(seatstatus2.StandUpTrigger, seatstatus2.CashingOutState).
		Permit(seatstatus2.SitOutTrigger, seatstatus2.SittingOutState).
		Permit(seatstatus2.RoundEndTrigger, seatstatus2.JoiningState)

	fsm.Configure(seatstatus2.SittingOutState).
		Permit(seatstatus2.StandUpTrigger, seatstatus2.CashingOutState).
		Permit(seatstatus2.TimeoutTrigger, seatstatus2.CashingOutState).
		Permit(seatstatus2.SitDownTrigger, seatstatus2.JoiningState).
		OnEntry(func(ctx context.Context, args ...any) error {
			s.logger.Debug("start sit out timer", zap.Object("seatStatus", seatStatus))
			seatStatus.SitOutStartTime = time.Now()
			seatStatus.CancelSitOutTimer = s.gameController.RunTimer(seatStatus.SitOutDuration, func() {
				s.logger.Debug("sit out timeout", zap.Object("seatStatus", seatStatus))
				if err := seatStatus.FSM.Fire(seatstatus2.TimeoutTrigger); err != nil {
					s.logger.Error("failed to fire timeout trigger", zap.Error(err), zap.Object("seatStatus", seatStatus))
				}

				// Notify client that he's force stand up by server.
				s.msgBus.Unicast(uid, core.MessageTopic, &txpokergrpc.Message{
					Event: &txpokergrpc.Event{
						Warning: &txpokergrpc.Warning{Reason: txpokergrpc.WarningReason_IDLE},
					},
				})

				s.msgBus.Broadcast(core.MessageTopic, &txpokergrpc.Message{
					Model: &txpokergrpc.Model{
						SeatStatusGroup: s.seatStatusGroup.ToProto(),
					},
				})
			})
			return nil
		}).
		OnExit(func(ctx context.Context, args ...any) error {
			seatStatus.CancelSitOutTimer()
			seatStatus.SitOutDuration = lo.Ternary(
				seatStatus.SitOutDuration > time.Since(seatStatus.SitOutStartTime),
				seatStatus.SitOutDuration-time.Since(seatStatus.SitOutStartTime),
				0,
			)
			return nil
		})

	return seatStatus
}

func (s *BaseSeatStatusService) StartRefillSitOutDurationLoop() {
	s.gameController.RunTicker(s.gameSetting.SitOutRefillIntervalDuration, func() {
		for _, seatStatus := range s.seatStatusGroup.Status {
			seatStatus.SitOutDuration = seatStatus.SitOutDuration + s.gameSetting.RefillSitOutDuration
			if seatStatus.SitOutDuration > s.gameSetting.MaxSitOutDuration {
				seatStatus.SitOutDuration = s.gameSetting.MaxSitOutDuration
			}
		}

		s.logger.Debug("refill sit out duration", zap.Object("seatStatus", s.seatStatusGroup))
		s.msgBus.Broadcast(core.MessageTopic, &txpokergrpc.Message{
			Model: &txpokergrpc.Model{
				SeatStatusGroup: s.seatStatusGroup.ToProto(),
			},
		})
	})
}

func (s *BaseSeatStatusService) Idle(uid core.Uid) {
	if _, ok := s.playerGroup.Data[uid]; !ok {
		return
	}

	s.playerGroup.Data[uid].IsIdle = true
	s.logger.Debug("idle",
		zap.String("uid", uid.String()),
		zap.String("role", s.playerGroup.Data[uid].Role.String()),
	)
}

func (s *BaseSeatStatusService) Act(uid core.Uid) {
	if _, ok := s.playerGroup.Data[uid]; !ok {
		return
	}

	s.playerGroup.Data[uid].IsIdle = false
}

func (s *BaseSeatStatusService) fetchUserCash(uid core.Uid) {
	resp, err := s.userAPI.FetchUserDetail(uid)
	if err != nil {
		s.logger.Error("failed to fetch user cash", zap.Error(err))
		return
	}

	s.gameController.RunTask(func() {
		user, ok := s.userGroup.Data[uid]
		if !ok { // user might be removed already
			return
		}

		user.Cash = resp.Data.Cash
		s.msgBus.Broadcast(core.MessageTopic, &txpokergrpc.Message{
			Model: &txpokergrpc.Model{
				UserGroup: s.userGroup.ToProto(),
			},
		})
	})
}
