package service

import (
	"fmt"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/api"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/model"
	"card-game-server-prototype/pkg/util"
	"time"

	"go.uber.org/zap"
)

type BaseGameRepoService struct {
	roomInfo    *commonmodel.RoomInfo
	gameSetting *model.GameSetting
	gameAPI     api.GameAPI
	msgBus      core.MsgBus
	logger      *zap.Logger
}

func ProvideBaseGameRepoService(
	roomInfo *commonmodel.RoomInfo,
	gameSetting *model.GameSetting,
	gameAPI api.GameAPI,
	msgBus core.MsgBus,
	loggerFactory *util.LoggerFactory,
) *BaseGameRepoService {
	return &BaseGameRepoService{
		roomInfo:    roomInfo,
		gameSetting: gameSetting,
		gameAPI:     gameAPI,
		msgBus:      msgBus,
		logger:      loggerFactory.Create("BaseGameRepoService"),
	}
}

func (s *BaseGameRepoService) FetchGameInfo() error {
	rawGameSettings, err := s.gameAPI.FetchGameSetting()
	if err != nil {
		return err
	}

	setting, ok := rawGameSettings.Data[s.roomInfo.GameMetaUid]
	if !ok {
		return fmt.Errorf("cannot find corresponding gameMetaUid: %v ", s.roomInfo.GameMetaUid)
	}

	s.gameSetting.GameMetaUid = setting.GameMetaUid
	s.gameSetting.SmallBlind = setting.SmallBlind
	s.gameSetting.BigBlind = setting.BigBlind
	s.gameSetting.TurnDuration = time.Duration(setting.TurnSecond) * time.Second
	s.gameSetting.MinEnterLimitBB = setting.MinEnterLimitBB
	s.gameSetting.MaxEnterLimitBB = setting.MaxEnterLimitBB
	s.gameSetting.TableSize = setting.TableSize
	s.gameSetting.MaxUserAmount = setting.MaxPlayerCount

	// 熱更抽水版本後，Common & Buddy 的 WaterPct 改由 fetchGameWater 取得。
	// 不再從 GameMetaUid 或 get gamesetting API 取得。
	if err := s.FetchGameWater(); err != nil {
		return err
	}

	return nil
}

func (s *BaseGameRepoService) FetchGameWater() error {
	resp, err := s.gameAPI.FetchGameWater()
	if err != nil {
		return err
	}

	for _, v := range resp.Data {
		if v.Id == "zoomtxpkr" && s.gameSetting.WaterPct != v.Common {
			waterPctBefore := s.gameSetting.WaterPct
			s.gameSetting.WaterPct = v.Common

			s.logger.Debug("set new game water",
				zap.Int("waterPctBefore", waterPctBefore),
				zap.Int("waterPct", s.gameSetting.WaterPct),
			)

			return err
		}
	}

	return nil
}
