package api

import (
	"card-game-server-prototype/pkg/config"
	"github.com/imroc/req/v3"
)

type ClubModeGameAPI struct {
	cfg        *config.Config
	httpClient *req.Client
}

func ProvideClubModeGameAPI(httpClient *req.Client, apiCFG *config.APIConfig, cfg *config.Config) *ClubModeGameAPI {
	return &ClubModeGameAPI{
		cfg: cfg,
		httpClient: httpClient.Clone().
			SetBaseURL("https://"+*apiCFG.MainServerHost).
			SetCommonHeader("jtoken", *apiCFG.MainServerAPIKey),
	}
}

func (api *ClubModeGameAPI) FetchGameSetting() (*ClubGameSettingResponse, error) {
	resp := &ClubGameSettingResponse{}
	err := api.httpClient.Get("/game/game/club/otherGameSetting").
		AddQueryParam("gameMetaUid", *api.cfg.GameMetaUid).
		Do().
		Into(resp)
	return resp, err
}
