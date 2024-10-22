package service

import (
	"card-game-server-prototype/pkg/game/txpoker/constant"
	txpokermodel "card-game-server-prototype/pkg/game/txpoker/model"
	"github.com/samber/lo"
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

// EvalWater 營利扣水前超過 n BB額外抽 1BB
func (s *BaseJackpotService) EvalWater(totalRawProfits int) int {
	return lo.Ternary(totalRawProfits >= constant.DonateJackpotIfBBReaching*s.gameSetting.BigBlind,
		s.gameSetting.BigBlind,
		0,
	)
}
