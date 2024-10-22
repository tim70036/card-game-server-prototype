package service

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/api"
	"time"
)

type ClubModeGameRepoService struct {
	*BaseGameRepoService
	clubGameApi api.ClubGameAPI
}

func ProvideClubModeGameRepoService(
	baseGameRepoService *BaseGameRepoService,
	clubGameApi api.ClubGameAPI,
) *ClubModeGameRepoService {
	return &ClubModeGameRepoService{
		BaseGameRepoService: baseGameRepoService,
		clubGameApi:         clubGameApi,
	}
}

func (s *ClubModeGameRepoService) TriggerJackpot() (map[core.Uid]int, error) {
	// CLUB MODE skip for now
	return nil, nil
}

func (s *ClubModeGameRepoService) FetchGameInfo() error {
	rawGameSettings, err := s.clubGameApi.FetchGameSetting()
	if err != nil {
		return err
	}

	rawRoom, err := s.gameAPI.FetchRoom(s.roomInfo.RoomId, s.roomInfo.GameMetaUid)
	if err != nil {
		return err
	}

	s.roomInfo.GameMetaUid = rawGameSettings.Data.GameMetaUid
	s.roomInfo.ClubId = rawGameSettings.Data.ClubId

	s.gameInfo.CreationId = rawRoom.Data.CreationId
	s.gameSetting.GameMetaUid = rawGameSettings.Data.GameMetaUid
	s.gameSetting.SmallBlind = rawGameSettings.Data.GameSettingMeta.SmallBlind
	s.gameSetting.BigBlind = rawGameSettings.Data.GameSettingMeta.BigBlind
	s.gameSetting.TurnDuration = time.Duration(rawGameSettings.Data.GameSettingMeta.TurnSecond) * time.Second
	s.gameSetting.InitialExtraTurnDuration = time.Duration(rawGameSettings.Data.GameSettingMeta.InitialExtraTurnSecond) * time.Second
	s.gameSetting.ExtraTurnRefillIntervalRound = rawGameSettings.Data.GameSettingMeta.ExtraTurnRefillIntervalRound
	s.gameSetting.RefillExtraTurnDuration = time.Duration(rawGameSettings.Data.GameSettingMeta.RefillExtraTurnSecond) * time.Second
	s.gameSetting.MaxExtraTurnDuration = time.Duration(rawGameSettings.Data.GameSettingMeta.MaxExtraTurnSecond) * time.Second
	s.gameSetting.InitialSitOutDuration = time.Duration(rawGameSettings.Data.GameSettingMeta.InitialSitOutSecond) * time.Second
	s.gameSetting.SitOutRefillIntervalDuration = time.Duration(rawGameSettings.Data.GameSettingMeta.SitOutRefillIntervalSecond) * time.Second
	s.gameSetting.RefillSitOutDuration = time.Duration(rawGameSettings.Data.GameSettingMeta.RefillSitOutSecond) * time.Second
	s.gameSetting.MaxSitOutDuration = time.Duration(rawGameSettings.Data.GameSettingMeta.MaxSitOutSecond) * time.Second
	s.gameSetting.MinEnterLimitBB = rawGameSettings.Data.GameSettingMeta.MinEnterLimitBB
	s.gameSetting.MaxEnterLimitBB = rawGameSettings.Data.GameSettingMeta.MaxEnterLimitBB
	s.gameSetting.WaterPct = rawGameSettings.Data.GameSettingMeta.WaterPct
	s.gameSetting.TableSize = rawGameSettings.Data.GameSettingMeta.TableSize
	s.gameSetting.LeastGamePlayerAmount = rawGameSettings.Data.GameSettingMeta.LeastGamePlayerAmount
	s.gameSetting.MaxWaterLimitBB = rawGameSettings.Data.GameSettingMeta.MaxWaterLimitBB

	return nil
}
