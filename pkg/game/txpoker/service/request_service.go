package service

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
	"time"
)

type RequestService struct {
	forceBuyInGroup *model.ForceBuyInGroup
	gameController  core.GameController

	logger *zap.Logger
}

func ProvideRequestService(
	forceBuyInGroup *model.ForceBuyInGroup,
	gameController core.GameController,

	loggerFactory *util.LoggerFactory,
) *RequestService {
	return &RequestService{
		forceBuyInGroup: forceBuyInGroup,
		gameController:  gameController,

		logger: loggerFactory.Create("RequestService"),
	}
}

type ForceBuyInResult struct {
	IsExist        bool
	BuyInChip      int
	RemainDuration time.Duration
}

func (s *RequestService) ForceBuyIn(uid core.Uid) (*ForceBuyInResult, error) {
	data, ok := s.forceBuyInGroup.Get(uid)

	if !ok {
		s.logger.Debug("ForceBuyIn data is not exist", zap.String("uid", uid.String()))
		return &ForceBuyInResult{
			IsExist: false,
		}, nil
	}

	buyInChip := data.GetBuyInChip()
	remainDuration := data.GetExpireTime().Sub(time.Now())

	s.logger.Debug("Get ForceBuyIn",
		zap.String("uid", uid.String()),
		zap.Int("Chip", buyInChip),
		zap.Duration("RemainTime", remainDuration),
	)

	return &ForceBuyInResult{
		IsExist:        true,
		BuyInChip:      buyInChip,
		RemainDuration: remainDuration,
	}, nil
}
