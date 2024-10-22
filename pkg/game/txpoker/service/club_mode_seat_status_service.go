package service

import (
	"card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/type/seatstatus"
	"github.com/samber/lo"
)

type ClubModeSeatStatusService struct {
	*BaseSeatStatusService
}

func ProvideClubModeSeatStatusService(
	baseSeatStatusService *BaseSeatStatusService,
) *ClubModeSeatStatusService {
	return &ClubModeSeatStatusService{
		BaseSeatStatusService: baseSeatStatusService,
	}
}

func (s *ClubModeSeatStatusService) IsReadyToStartRound() bool {
	joiningCnt := lo.CountBy(lo.Values(s.seatStatusGroup.Status), func(status *model.SeatStatus) bool {
		return seatstatus.JoiningState.IsEqual(status.FSM.MustState())
	})

	if s.gameSetting.LeastGamePlayerAmount > 0 {
		return joiningCnt >= s.gameSetting.LeastGamePlayerAmount
	}

	return s.BaseSeatStatusService.IsReadyToStartRound()
}
