package api

import (
	txpokerapi "card-game-server-prototype/pkg/game/txpoker/api"
)

type GameAPI interface {
	FetchGameSetting() (*GameSettingResponse, error)
	FetchGameWater() (*txpokerapi.GetGameWaterResp, error)
}
