package api

import (
	"card-game-server-prototype/pkg/config"
	txpokerapi "card-game-server-prototype/pkg/game/txpoker/api"
	"github.com/imroc/req/v3"
)

type BaseGameAPI struct {
	httpClient *req.Client
}

func ProvideBaseGameAPI(
	httpClient *req.Client,
	apiCFG *config.APIConfig,
) *BaseGameAPI {
	return &BaseGameAPI{
		httpClient: httpClient.Clone().
			SetBaseURL("https://"+*apiCFG.MainServerHost).
			SetCommonHeader("jtoken", *apiCFG.MainServerAPIKey),
	}
}

func (api *BaseGameAPI) FetchGameSetting() (*GameSettingResponse, error) {
	resp := &GameSettingResponse{}
	err := api.httpClient.Get("/game/ztxpkr/setting").
		Do().
		Into(resp)
	return resp, err
}

func (api *BaseGameAPI) FetchGameWater() (*txpokerapi.GetGameWaterResp, error) {
	resp := &txpokerapi.GetGameWaterResp{}
	err := api.httpClient.Get("/game/game-setting/game-water").
		Do().
		Into(resp)
	return resp, err
}
