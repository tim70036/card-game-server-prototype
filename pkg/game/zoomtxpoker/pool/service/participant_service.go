package service

import (
	"context"
	commonapi "card-game-server-prototype/pkg/common/api"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/common/type/gametype"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/constant"
	model2 "card-game-server-prototype/pkg/game/zoomtxpoker/pool/model"
	participant2 "card-game-server-prototype/pkg/game/zoomtxpoker/pool/type/participant"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"
	"github.com/qmuntal/stateless"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
	"time"
)

type ParticipantService struct {
	gameController core.GameController
	msgBus         core.MsgBus

	userGroup         *commonmodel.UserGroup
	participantGroup  *model2.ParticipantGroup
	userAPI           commonapi.UserAPI
	gameSetting       *model2.GameSetting
	forceBuyInGroup   *model2.ForceBuyInGroup
	tableProfitsGroup *model2.TableProfitsGroup

	chipSnapshotService *ChipSnapshotService

	isToppingUp map[core.Uid]struct{}
	logger      *zap.Logger
}

func ProvideParticipantService(
	gameController core.GameController,
	msgBus core.MsgBus,
	userGroup *commonmodel.UserGroup,
	participantGroup *model2.ParticipantGroup,
	userAPI commonapi.UserAPI,
	gameSetting *model2.GameSetting,
	forceBuyInGroup *model2.ForceBuyInGroup,
	tableProfitsGroup *model2.TableProfitsGroup,
	chipSnapshotService *ChipSnapshotService,
	loggerFactory *util.LoggerFactory) *ParticipantService {
	return &ParticipantService{
		gameController:      gameController,
		msgBus:              msgBus,
		userGroup:           userGroup,
		participantGroup:    participantGroup,
		userAPI:             userAPI,
		gameSetting:         gameSetting,
		forceBuyInGroup:     forceBuyInGroup,
		tableProfitsGroup:   tableProfitsGroup,
		chipSnapshotService: chipSnapshotService,
		isToppingUp:         make(map[core.Uid]struct{}),
		logger:              loggerFactory.Create("ParticipantService"),
	}
}

func (s *ParticipantService) BuyIn(uid core.Uid, buyInChip int) error {
	part, ok := s.participantGroup.Data[uid]
	if !ok {
		return status.Errorf(codes.NotFound, "cannot found uid %v in participant group with uids: %v", uid, lo.Keys(s.participantGroup.Data))
	}

	if err := part.FSM.Fire(participant2.BuyInTrigger, buyInChip); err != nil {
		return err
	}

	s.logger.Info(
		"buying in",
		zap.String("uid", uid.String()),
		zap.Int("buyInChip", buyInChip),
		zap.Object("participant", part),
	)

	return nil
}

func (s *ParticipantService) ExitMatch(uid core.Uid, isLeaving bool) (<-chan struct{}, error) {
	part, ok := s.participantGroup.Data[uid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "cannot found uid %v in participant group with uids: %v", uid, lo.Keys(s.participantGroup.Data))
	}

	cashOutDone := make(chan struct{})
	if err := part.FSM.Fire(participant2.ExitMatchTrigger, cashOutDone, isLeaving); err != nil {
		return nil, err
	}

	s.logger.Info("exiting match",
		zap.String("uid", uid.String()),
		zap.Object("participant", part),
		zap.Bool("isLeaving", isLeaving),
	)
	return cashOutDone, nil
}

// EnterGame will change participant state to PlayingState. It will also transfer participant's chip to game.
// It's designed as transaction on purpose. We need make sure all
// participants done enter game together. In case of some participants failed to
// enter game, we need to revert all participants' state.
func (s *ParticipantService) EnterGame(chipIntoGame map[core.Uid]int) error {
	for uid := range chipIntoGame {
		if _, ok := s.participantGroup.Data[uid]; !ok {
			return status.Errorf(codes.NotFound, "cannot found uid %v in participant group with uids: %v", uid, lo.Keys(s.participantGroup.Data))
		}
	}

	s.logger.Info("entering game", zap.Object("chipIntoGame", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for uid, chip := range chipIntoGame {
			enc.AddInt(uid.String(), chip)
		}
		return nil
	})))

	// Fire enter game trigger for each participant.
	for uid := range chipIntoGame {
		if err := s.participantGroup.Data[uid].FSM.Fire(participant2.EnterGameTrigger); err != nil {
			s.rollbackEnterGame(lo.Keys(chipIntoGame))
			return err
		}
	}

	for uid := range chipIntoGame {
		s.participantGroup.Data[uid].Chip -= chipIntoGame[uid]
	}

	return nil
}

func (s *ParticipantService) rollbackEnterGame(uids core.UidList) {
	for _, uid := range uids {
		part := s.participantGroup.Data[uid]
		if err := part.FSM.Fire(participant2.ExitGameTrigger); err != nil {
			s.logger.Error("Failed to fire exit game trigger during rollback", zap.Error(err), zap.String("uid", uid.String()))
		}
	}
}

func (s *ParticipantService) ExitGame(uid core.Uid, chipFromGame int) error {
	part, ok := s.participantGroup.Data[uid]
	if !ok {
		return status.Errorf(codes.NotFound, "cannot found uid %v in participant group with uids: %v", uid, lo.Keys(s.participantGroup.Data))
	}

	s.logger.Info("exiting game",
		zap.String("uid", uid.String()),
		zap.Object("participant", part),
		zap.Int("chipFromGame", chipFromGame),
	)

	part.Chip += chipFromGame
	s.chipSnapshotService.TakeSnapshot(uid, part.Chip)

	if err := part.FSM.Fire(participant2.ExitGameTrigger); err != nil {
		s.logger.Error(
			"Failed to fire exit game trigger",
			zap.Error(err),
			zap.Int("chipFromGame", chipFromGame),
			zap.Object("participant", part),
		)
		return err
	}

	return nil
}

func (s *ParticipantService) TopUp(uid core.Uid, topUpChip int) (<-chan struct{}, error) {
	if topUpChip <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "topUpChip must > 0, got %v", topUpChip)
	}

	if _, ok := s.isToppingUp[uid]; ok {
		return nil, status.Errorf(codes.AlreadyExists, "uid %v is already topping up", uid)
	}

	part, ok := s.participantGroup.Data[uid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "cannot found uid %v in participant group with uids: %v", uid, lo.Keys(s.participantGroup.Data))
	}

	user, ok := s.userGroup.Data[uid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "cannot found uid %v in user group with uids: %v", uid, lo.Keys(s.userGroup.Data))
	}

	if topUpChip > user.Cash {
		return nil, status.Errorf(codes.InvalidArgument, "uid %v has not enough cash %v to top up %v", uid, user.Cash, topUpChip)
	}

	if !lo.Contains(
		[]participant2.State{participant2.MatchingState, participant2.PlayingState},
		part.FSM.MustState().(participant2.State),
	) {
		return nil, status.Errorf(codes.FailedPrecondition, "uid %v is not in matching or playing state, got %v", uid, part.FSM.MustState().(participant2.State).String())
	}

	// check top up chip is valid
	maxEnterLimitChip := s.gameSetting.MaxEnterLimitBB * s.gameSetting.BigBlind
	if part.Chip >= maxEnterLimitChip {
		return nil, status.Errorf(codes.FailedPrecondition, "uid %v already has chip %v that >= maxEnterLimitChip %v", uid, part.Chip, maxEnterLimitChip)
	}

	topUpChip = lo.Ternary(
		part.Chip+topUpChip > maxEnterLimitChip,
		maxEnterLimitChip-part.Chip,
		topUpChip,
	)

	topUpDone := make(chan struct{})
	s.isToppingUp[uid] = struct{}{}
	s.logger.Info("start top up",
		zap.Object("participant", part),
		zap.Int("topUpChip", topUpChip),
	)

	go func() {
		err := s.userAPI.ExchangeChip(uid, gametype.ZoomTXPoker, topUpChip*-1)
		go s.fetchUserCash(uid)

		s.gameController.RunTask(func() {
			defer func() {
				delete(s.isToppingUp, uid)
				close(topUpDone)
			}()

			if err != nil {
				s.logger.Error("failed top up", zap.Error(err), zap.Object("participant", part), zap.Int("topUpChip", topUpChip))
				return
			}

			// After top up, need to check if current status of participant.
			revertTopUp := func() {
				if err := s.userAPI.ExchangeChip(uid, gametype.ZoomTXPoker, topUpChip); err != nil {
					s.logger.Error("failed to revert top up chips",
						zap.Error(err),
						zap.Object("participant", part),
						zap.Int("revertTopUpChip", topUpChip),
					)
				}
				return
			}

			if _, ok := s.userGroup.Data[uid]; !ok {
				s.logger.Warn("revert top up chips since user has left the room",
					zap.Object("participant", part),
					zap.Int("revertTopUpChip", topUpChip),
				)
				go revertTopUp()
				return
			}

			if !lo.Contains(
				[]participant2.State{participant2.MatchingState, participant2.PlayingState},
				part.FSM.MustState().(participant2.State),
			) {
				s.logger.Warn("revert top up chips since user is not in matching or playing state",
					zap.Object("participant", part),
					zap.Int("revertTopUpChip", topUpChip),
				)

				go revertTopUp()
				return
			}

			s.gameController.RunTask(func() {
				if tableProfits, ok := s.tableProfitsGroup.Get(uid); ok {
					tableProfits.SumBuyInChips += topUpChip
					s.tableProfitsGroup.Save(tableProfits)
				}
			})

			if err := s.AddChip(uid, topUpChip); err != nil {
				s.logger.Error("failed to add chip after top up", zap.Error(err), zap.Object("participant", part), zap.Int("topUpChip", topUpChip))
			}

			s.logger.Info("done top up", zap.Object("participant", part))

			s.msgBus.Unicast(uid, core.MessageTopic, &txpokergrpc.Message{
				Model: &txpokergrpc.Model{
					Participant: part.ToProto(),
				},
			})

			s.msgBus.Unicast(uid, core.EventTopic, &txpokergrpc.Event{
				TopupSuccess: &txpokergrpc.TopUpSuccess{
					Amount: int32(part.Chip),
				},
			})
		})
	}()

	return topUpDone, nil
}

func (s *ParticipantService) AddChip(uid core.Uid, chip int) error {
	if chip <= 0 {
		return status.Errorf(codes.InvalidArgument, "chip must > 0, got %v", chip)
	}

	part, ok := s.participantGroup.Data[uid]

	if !ok {
		return status.Errorf(codes.NotFound, "cannot found uid %v in participant group", uid)
	}

	part.Chip += chip
	s.chipSnapshotService.AddToSnapshot(uid, chip)
	s.logger.Debug("add chip", zap.Int("chip", chip), zap.Object("participant", part))

	return nil
}

func (s *ParticipantService) NewParticipant(uid core.Uid) *model2.Participant {
	fsm := stateless.NewStateMachine(participant2.ObservingState)
	part := &model2.Participant{
		Uid:             uid,
		Chip:            0,
		FSM:             fsm,
		QueuedTopUpChip: 0,
		IdleAt:          time.Now().Unix(),
	}

	fsm.SetTriggerParameters(participant2.BuyInTrigger, reflect.TypeOf(1))
	fsm.SetTriggerParameters(participant2.ExitMatchTrigger, reflect.TypeOf(make(chan struct{})), reflect.TypeOf(false))

	fsm.Configure(participant2.ObservingState).
		Permit(participant2.BuyInTrigger, participant2.BuyingInState).
		PermitReentry(participant2.ExitMatchTrigger).
		OnEntry(func(ctx context.Context, args ...any) error {

			if part.IsIdleRoundsReachMax() {
				// 設定一定過期的時間
				part.IdleAt = time.Now().Add(-1 * constant.IdlingTimeoutDuration).Unix()

				s.logger.Info("idle reach max, kick",
					zap.String("uid", uid.String()),
					zap.Object("participant", part),
				)

				return nil
			}

			part.IdleAt = time.Now().Unix()

			s.logger.Info("set idle, enter observing",
				zap.String("uid", uid.String()),
				zap.Object("participant", part),
			)

			return nil
		}).
		OnExit(func(ctx context.Context, args ...any) error {
			part.IdleAt = 0
			s.logger.Info("clear idle, exit observing", zap.String("uid", uid.String()))
			return nil
		})

	fsm.Configure(participant2.MatchingState).
		Permit(participant2.EnterGameTrigger, participant2.PlayingState).
		Permit(participant2.ExitMatchTrigger, participant2.CashingOutState)

	fsm.Configure(participant2.PlayingState).
		Permit(participant2.ExitGameTrigger, participant2.MatchingState).
		Permit(participant2.ExitMatchTrigger, participant2.CashingOutState)

	fsm.Configure(participant2.BuyingInState).
		Permit(participant2.SuccessTrigger, participant2.MatchingState).
		Permit(participant2.FailedTrigger, participant2.ObservingState).
		OnEntry(func(ctx context.Context, args ...any) error {
			buyInChip := args[0].(int)
			_, isForceBuyInExist := s.forceBuyInGroup.Get(uid)

			s.logger.Info("start buying in",
				zap.String("uid", uid.String()),
				zap.Int("buyInChip", buyInChip),
				zap.Object("participant", part),
				zap.Bool("forceBuyInExist", isForceBuyInExist),
			)

			if !isForceBuyInExist &&
				(buyInChip < s.gameSetting.MinEnterLimitBB*s.gameSetting.BigBlind ||
					buyInChip > s.gameSetting.MaxEnterLimitBB*s.gameSetting.BigBlind) {
				s.gameController.RunTask(func() {
					if err := fsm.Fire(participant2.FailedTrigger); err != nil {
						s.logger.Error("failed to fire failed trigger", zap.Error(err), zap.Object("participant", part))
					}
				})
				return status.Errorf(codes.InvalidArgument, "uid %v try to sit down with invalid buyInChip %v", uid, buyInChip)
			}

			// To avoid blocking game main loop, make it async.
			go func() {
				if err := s.userAPI.ExchangeChip(uid, gametype.ZoomTXPoker, buyInChip*-1); err != nil {
					s.gameController.RunTask(func() {
						if err := fsm.Fire(participant2.FailedTrigger); err != nil {
							s.logger.Error("failed to fire failed trigger", zap.Error(err), zap.Object("participant", part))
						}

						s.logger.Error("failed buying in", zap.Error(err), zap.Object("participant", part))

						s.msgBus.Unicast(uid, core.MessageTopic, &txpokergrpc.Message{
							Model: &txpokergrpc.Model{
								Participant: part.ToProto(),
							},
						})
					})
					return
				}

				go s.fetchUserCash(uid)
				s.gameController.RunTask(func() {
					part.Chip += buyInChip

					s.gameController.RunTask(func() {
						if tableProfits, ok := s.tableProfitsGroup.Get(uid); ok {
							tableProfits.SumBuyInChips += buyInChip
							s.tableProfitsGroup.Save(tableProfits)
						}
					})

					s.chipSnapshotService.TakeSnapshot(uid, part.Chip)

					if err := fsm.Fire(participant2.SuccessTrigger); err != nil {
						s.logger.Error("failed to fire success trigger", zap.Error(err), zap.Object("participant", part))
					}

					s.logger.Info("buy in done",
						zap.String("uid", uid.String()),
						zap.Object("participant", part),
					)

					s.msgBus.Unicast(uid, core.MessageTopic, &txpokergrpc.Message{
						Model: &txpokergrpc.Model{
							Participant: part.ToProto(),
						},
					})
				})
			}()

			return nil
		})

	fsm.Configure(participant2.CashingOutState).
		Permit(participant2.BackTrigger, participant2.ObservingState).
		Permit(participant2.ExitingTrigger, participant2.ExitingMatchState).
		OnEntry(func(ctx context.Context, args ...any) error {
			cashOutDone := args[0].(chan struct{})
			isLeaving := args[1].(bool)

			triggerNext := lo.Ternary(isLeaving, participant2.ExitingTrigger, participant2.BackTrigger)

			go func() {
				s.logger.Info("start cashing out",
					zap.String("uid", uid.String()),
					zap.Object("participant", part),
					zap.Int("Chip", part.Chip),
				)

				if part.Chip != 0 {
					if err := s.userAPI.ExchangeChip(uid, gametype.ZoomTXPoker, part.Chip); err != nil {
						go s.fetchUserCash(uid)
						s.gameController.RunTask(func() {
							part.Chip = 0
							if err := fsm.Fire(triggerNext, cashOutDone); err != nil {
								s.logger.Error("failed to fire trigger", zap.Error(err),
									zap.String("trigger", triggerNext.String()),
									zap.Object("participant", part),
								)
							}

							s.logger.Error("failed cashing out", zap.Error(err), zap.Object("participant", part))

							s.msgBus.Unicast(uid, core.MessageTopic, &txpokergrpc.Message{
								Model: &txpokergrpc.Model{
									Participant: part.ToProto(),
								},
							})

							close(cashOutDone)
						})
						return
					}

					// Set force buy in when user is receiving game result, not here.
				}

				go s.fetchUserCash(uid)
				s.gameController.RunTask(func() {
					part.Chip = 0
					s.chipSnapshotService.DeleteSnapshot(uid)

					if err := fsm.Fire(triggerNext, cashOutDone); err != nil {
						s.logger.Error("failed to fire trigger", zap.Error(err),
							zap.String("trigger", triggerNext.String()),
							zap.Object("participant", part),
						)
					}

					s.logger.Info("done cashing out",
						zap.String("uid", uid.String()),
						zap.Object("participant", part),
					)

					s.msgBus.Unicast(uid, core.MessageTopic, &txpokergrpc.Message{
						Model: &txpokergrpc.Model{
							Participant: part.ToProto(),
						},
					})

					close(cashOutDone)
				})
			}()
			return nil
		})

	fsm.Configure(participant2.ExitingMatchState).
		OnEntry(func(ctx context.Context, args ...any) error {
			s.logger.Info("user is leaving now",
				zap.String("uid", uid.String()),
			)
			return nil
		})

	return part
}

func (s *ParticipantService) fetchUserCash(uid core.Uid) {
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

		s.msgBus.Unicast(uid, core.MessageTopic, &txpokergrpc.Message{
			Model: &txpokergrpc.Model{
				User: user.ToProto(),
			},
		})
	})
}
