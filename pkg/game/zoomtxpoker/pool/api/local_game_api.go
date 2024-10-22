package api

import (
	commonmodel "card-game-server-prototype/pkg/common/model"
	txpokerapi "card-game-server-prototype/pkg/game/txpoker/api"
	"card-game-server-prototype/pkg/util"

	"go.uber.org/zap"
)

type LocalGameAPI struct {
	roomInfo *commonmodel.RoomInfo

	logger *zap.Logger
}

func ProvideLocalGameAPI(
	roomInfo *commonmodel.RoomInfo,
	loggerFactory *util.LoggerFactory,
) *LocalGameAPI {
	return &LocalGameAPI{
		roomInfo: roomInfo,
		logger:   loggerFactory.Create("LocalGameAPI"),
	}
}

func (api *LocalGameAPI) FetchGameSetting() (*GameSettingResponse, error) {
	return &GameSettingResponse{
		Data: map[string]struct {
			GameMetaUid                  string `json:"gameMetaUid"`
			Game                         int    `json:"game"`
			GameMode                     int    `json:"gameMode"`
			SmallBlind                   int    `json:"smallBlind"`
			BigBlind                     int    `json:"bigBlind"`
			TurnSecond                   int    `json:"turnSecond"`
			InitialExtraTurnSecond       int    `json:"initialExtraTurnSecond"`
			ExtraTurnRefillIntervalRound int    `json:"extraTurnRefillIntervalRound"`
			RefillExtraTurnSecond        int    `json:"refillExtraTurnSecond"`
			MaxExtraTurnSecond           int    `json:"maxExtraTurnSecond"`
			InitialSitOutSecond          int    `json:"initialSitOutSecond"`
			SitOutRefillIntervalSecond   int    `json:"sitOutRefillIntervalSecond"`
			RefillSitOutSecond           int    `json:"refillSitOutSecond"`
			MaxSitOutSecond              int    `json:"maxSitOutSecond"`
			MinEnterLimitBB              int    `json:"minEnterLimitBB"`
			MaxEnterLimitBB              int    `json:"maxEnterLimitBB"`
			WaterPct                     int    `json:"waterPct"`
			TableSize                    int    `json:"tableSize"`
			LeastGamePlayerAmount        int    `json:"leastGamePlayerAmount"`
			MaxWaterLimitBB              int    `json:"maxWaterLimitBB"`
			MaxPlayerCount               int    `json:"maxPlayerCount"`
		}{
			api.roomInfo.GameMetaUid: {
				GameMetaUid:                  api.roomInfo.GameMetaUid,
				SmallBlind:                   10,
				BigBlind:                     20,
				TurnSecond:                   12,
				InitialExtraTurnSecond:       1,
				ExtraTurnRefillIntervalRound: 30,
				RefillExtraTurnSecond:        10,
				MaxExtraTurnSecond:           30,
				InitialSitOutSecond:          420,
				SitOutRefillIntervalSecond:   60,
				RefillSitOutSecond:           300,
				MaxSitOutSecond:              420,
				MinEnterLimitBB:              20,
				MaxEnterLimitBB:              200,
				TableSize:                    9,
				WaterPct:                     5,
			},
		},
	}, nil
}

func (api *LocalGameAPI) FetchGameWater() (*txpokerapi.GetGameWaterResp, error) {
	return &txpokerapi.GetGameWaterResp{
		ErrCode: 0,
		Msg:     "success",
		Data: []*txpokerapi.GameWater{
			{
				Id:     "zoomtxpkr",
				Common: 5,
				Buddy:  5,
			},
		},
	}, nil
}
