package service

import (
	commonmodel "card-game-server-prototype/pkg/common/model"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/model"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/type/participant"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"time"
)

type LeaveService struct {
	gameController core.GameController

	participantGroup *model.ParticipantGroup
	userGroup        *commonmodel.UserGroup

	userService        commonservice.UserService
	participantService *ParticipantService

	logger *zap.Logger
}

func ProvideLeaveService(
	gameController core.GameController,

	participantGroup *model.ParticipantGroup,
	userGroup *commonmodel.UserGroup,

	userService commonservice.UserService,
	participantService *ParticipantService,

	loggerFactory *util.LoggerFactory,
) *LeaveService {
	return &LeaveService{
		gameController: gameController,

		participantGroup: participantGroup,
		userGroup:        userGroup,

		userService:        userService,
		participantService: participantService,

		logger: loggerFactory.Create("LeaveService"),
	}
}
func (s *LeaveService) RunIdleChecker() {
	s.gameController.RunTicker(time.Second*10, func() {
		for _, uid := range s.getIdleTimeoutUsers() {
			_ = s.OnLeave(uid, commongrpc.KickoutReason_IDLE_TIMEOUT)
		}
	})
}

func (s *LeaveService) getIdleTimeoutUsers() []core.Uid {
	if s.participantGroup == nil || len(s.participantGroup.Data) == 0 {
		return nil
	}

	now := time.Now()

	return lo.Keys(lo.PickBy(s.participantGroup.Data, func(uid core.Uid, part *model.Participant) bool {
		return part.IsIdlingTimeout(now)
	}))
}

func (s *LeaveService) OnLeave(uid core.Uid, reason commongrpc.KickoutReason) error {
	curState := s.participantGroup.Data[uid].FSM.MustState().(participant.State)

	s.logger.Info("user leaving start",
		zap.String("uid", uid.String()),
		zap.String("reason", reason.String()),
		zap.String("curState", curState.String()),
	)

	// Need cash out
	if lo.Contains([]participant.State{participant.MatchingState, participant.PlayingState}, curState) {
		cashOutDone, err := s.participantService.ExitMatch(uid, true)
		if err != nil {
			s.logger.Error("user leaving but failed to exit match",
				zap.Error(err),
				zap.String("uid", uid.String()),
			)
			return err
		}
		go func() {
			// Only destroy user after cash out done. Client will not actually leave
			// until kickout message is sent.
			<-cashOutDone
			s.gameController.RunTask(func() {
				s.destroyUser(uid, reason)
			})
		}()
		return nil
	}

	go func() {
		s.gameController.RunTask(func() {
			s.destroyUser(uid, reason)
		})
	}()

	return nil
}

func (s *LeaveService) destroyUser(uid core.Uid, reason commongrpc.KickoutReason) {
	kickoutMSG := &txpokergrpc.Message{
		Event: &txpokergrpc.Event{
			Kickout: &commongrpc.Kickout{
				Reason: reason,
			},
		},
	}

	if err := s.userService.Destroy(core.MessageTopic, kickoutMSG, uid); err != nil {
		s.logger.Error("user leaving but failed to destroy user data",
			zap.Error(err),
			zap.String("uid", uid.String()),
		)
	}

	// userAmount is for rolling update. Do not remove this field.
	s.logger.Info("user leaving done",
		zap.String("uid", uid.String()),
		zap.String("reason", reason.String()),
		zap.Int("userAmount", len(s.userGroup.Data)),
	)
}
