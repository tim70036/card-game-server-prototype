package service

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	"card-game-server-prototype/pkg/game/darkchess/model"
	"github.com/samber/lo"
)

type PlayerService struct {
	playerGroup *model.PlayerGroup
}

func ProvidePlayerService(
	playerGroup *model.PlayerGroup,
) *PlayerService {
	return &PlayerService{
		playerGroup: playerGroup,
	}
}

func (s *PlayerService) GetLatestIdleTurn(uid core.Uid) (int, bool) {
	data, ok := s.playerGroup.Data[uid]

	if !ok {
		return 0, false
	}

	if len(data.IdleTurns) == 0 {
		return 0, false
	}

	lastTurn, _ := lo.Last(data.IdleTurns)

	return lastTurn, true
}

func (s *PlayerService) AppendIdleTurn(uid core.Uid, turn int) (isUpdate bool) {
	if _, ok := s.playerGroup.Data[uid]; !ok {
		return
	}

	latestIdle, err := lo.Last(s.playerGroup.Data[uid].IdleTurns)

	// empty，直接加
	if err != nil {
		s.playerGroup.Data[uid].IdleTurns = append(s.playerGroup.Data[uid].IdleTurns, turn)
		return
	}

	// 避免重複計算
	if latestIdle == turn {
		return
	}

	s.playerGroup.Data[uid].IdleTurns = append(s.playerGroup.Data[uid].IdleTurns, turn)
	return true
}

func (s *PlayerService) ResetIdleTurn(uid core.Uid) {
	if _, ok := s.playerGroup.Data[uid]; ok {
		s.playerGroup.Data[uid].IdleTurns = []int{}
	}
}

func (s *PlayerService) IsReachIdleThreshold(uid core.Uid) bool {
	if _, ok := s.playerGroup.Data[uid]; ok {
		return len(s.playerGroup.Data[uid].IdleTurns) >= constant.IdleThreshold
	}
	return false
}
