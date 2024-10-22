package service

import (
	txpokermodel "card-game-server-prototype/pkg/game/txpoker/model"
)

type ClubModeJackpotService struct {
	gameSetting *txpokermodel.GameSetting
}

func ProvideClubModeJackpotService(
	gameSetting *txpokermodel.GameSetting,
) *ClubModeJackpotService {
	return &ClubModeJackpotService{
		gameSetting: gameSetting,
	}
}

func (s *ClubModeJackpotService) EvalWater(_ int) int {
	return 0
}
