package api

import (
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/common/type/rawevent"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/type/card"
	"card-game-server-prototype/pkg/util"

	"github.com/samber/lo"
	"go.uber.org/zap"
)

type LocalGameAPI struct {
	roomInfo    *commonmodel.RoomInfo
	playerGroup *model2.PlayerGroup

	logger *zap.Logger
}

func ProvideLocalGameAPI(
	roomInfo *commonmodel.RoomInfo,
	playerGroup *model2.PlayerGroup,
	loggerFactory *util.LoggerFactory,
) *LocalGameAPI {
	return &LocalGameAPI{
		roomInfo:    roomInfo,
		playerGroup: playerGroup,
		logger:      loggerFactory.Create("LocalGameAPI"),
	}
}

func (api *LocalGameAPI) FetchGameSetting() (*GameSettingResponse, error) {
	return &GameSettingResponse{
		Data: map[string]struct {
			GameMetaUid                  string `json:"gameMetaUid"`
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
			TableSize                    int    `json:"tableSize"`
			WaterPct                     int    `json:"waterPct"`
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

func (api *LocalGameAPI) FetchGameWater() (*GetGameWaterResp, error) {
	return &GetGameWaterResp{
		ErrCode: 0,
		Msg:     "success",
		Data: []*GameWater{
			{
				Id:     "txpkr",
				Common: 5,
				Buddy:  5,
			},
		},
	}, nil
}

func (api *LocalGameAPI) FetchRoom(roomId string, gameMetaUid string) (*RoomResponse, error) {
	return &RoomResponse{
		Data: struct {
			RoomId       string  `json:"roomId"`
			GameMetaUid  string  `json:"gameMetaUid"`
			CreationId   int     `json:"creationId"`
			Vpip         float64 `json:"vpip"`
			EmptySeats   int     `json:"emptySeats"`
			LastGameTime string  `json:"lastGameTime"`
		}{
			RoomId:       roomId,
			GameMetaUid:  gameMetaUid,
			CreationId:   1,
			Vpip:         0.5,
			EmptySeats:   1,
			LastGameTime: "2023-07-18T00:00:00Z",
		},
	}, nil
}

func (api *LocalGameAPI) FetchGameResult(gameId string) (*GameResultResponse, error) {
	return &GameResultResponse{
		Data: lo.Map(
			lo.Keys(api.playerGroup.Data),
			func(uid core.Uid, _ int) struct {
				Uid            string `json:"uid"`
				BeforeLevel    int    `json:"beforeLevel"`
				BeforeExp      int    `json:"beforeExp"`
				LevelUpExp     int    `json:"levelUpExp"`
				NextLevel      int    `json:"nextLevel"`
				NextLevelUpExp int    `json:"nextLevelUpExp"`
				IncreaseExp    int    `json:"increaseExp"`
			} {
				return struct {
					Uid            string `json:"uid"`
					BeforeLevel    int    `json:"beforeLevel"`
					BeforeExp      int    `json:"beforeExp"`
					LevelUpExp     int    `json:"levelUpExp"`
					NextLevel      int    `json:"nextLevel"`
					NextLevelUpExp int    `json:"nextLevelUpExp"`
					IncreaseExp    int    `json:"increaseExp"`
				}{
					Uid:            uid.String(),
					BeforeLevel:    1,
					BeforeExp:      0,
					LevelUpExp:     0,
					NextLevel:      1,
					NextLevelUpExp: 0,
					IncreaseExp:    0,
				}
			},
		),
	}, nil
}

func (api *LocalGameAPI) UpdateRoom(roomId string, gameMetaUid string, vpip float64, emptySeatNum int) error {
	return nil
}

func (api *LocalGameAPI) StartGame(gameId, roomId string, gameMetaUid string, vpip float64, emptySeatNum int) error {
	return nil
}

func (api *LocalGameAPI) TriggerJackpot(jackpotPlayers []*model2.Player, gameMetaUid string) (*TriggerJackpotResponse, error) {
	return &TriggerJackpotResponse{
		Data: lo.Map(jackpotPlayers, func(p *model2.Player, _ int) struct {
			Uid    string `json:"uid"`
			Amount int    `json:"amount"`
		} {
			return struct {
				Uid    string `json:"uid"`
				Amount int    `json:"amount"`
			}{
				Uid:    p.Uid.String(),
				Amount: 123123,
			}
		}),
	}, nil
}

func (api *LocalGameAPI) EndGame(
	roomId string,
	gameId string,
	communityCards card.CardList,
	replay *model2.Replay,
	pocketCards map[core.Uid]card.CardList,
	betAmount map[core.Uid]int,
	profits map[core.Uid]int,
	waters map[core.Uid]int,
	jackpotWaters map[core.Uid]int,
	didEnterFlopStage map[core.Uid]bool,
) error {
	return nil
}

func (api *LocalGameAPI) SubmitWatchEvents(uid core.Uid, rawEvents rawevent.RawEventList) error {
	return nil
}
func (api *LocalGameAPI) SubmitUserEvents(uid core.Uid, rawEvents rawevent.RawEventList) error {
	return nil
}
