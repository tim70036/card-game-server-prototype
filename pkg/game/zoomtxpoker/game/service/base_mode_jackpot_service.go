package service

import (
	txpokermodel "card-game-server-prototype/pkg/game/txpoker/model"
)

type BaseJackpotService struct {
	gameSetting *txpokermodel.GameSetting
}

func ProvideBaseJackpotService(
	gameSetting *txpokermodel.GameSetting,
) *BaseJackpotService {
	return &BaseJackpotService{
		gameSetting: gameSetting,
	}
}

func (s *BaseJackpotService) EvalWater(_ int) int {
	return 0
}
