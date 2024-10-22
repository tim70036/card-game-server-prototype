package api

import (
	commonapi "card-game-server-prototype/pkg/common/api"
	"card-game-server-prototype/pkg/common/type/gametype"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"github.com/imroc/req/v3"
	"strconv"
)

type ClubModeUserAPI struct {
	httpClient *req.Client
	cfg        *config.Config
}

func ProvideClubModeUserAPI(httpClient *req.Client, apiCFG *config.APIConfig, cfg *config.Config) *ClubModeUserAPI {
	return &ClubModeUserAPI{
		cfg: cfg,
		httpClient: httpClient.Clone().
			SetBaseURL("https://"+*apiCFG.MainServerHost).
			SetCommonHeader("jtoken", *apiCFG.MainServerAPIKey),
	}
}

func (api *ClubModeUserAPI) FetchUserDetail(uid core.Uid) (*commonapi.UserDetailResponse, error) {
	resp := &commonapi.UserDetailResponse{}
	err := api.httpClient.Get("/game/user/detail").
		AddQueryParam("uid", uid.String()).
		Do().
		Into(resp)

	return resp, err
}

func (api *ClubModeUserAPI) ExchangeChip(uid core.Uid, gameType gametype.GameType, amount int) error {
	req := &exchangeChipForClubRequest{
		Uid:         uid.String(),
		GameType:    string(gameType),
		Amount:      strconv.Itoa(amount),
		GameMetaUid: *api.cfg.GameMetaUid,
	}

	return api.httpClient.Put("/game/club/chip").
		SetBodyJsonMarshal(req).
		Do().
		Err
}

func (api *ClubModeUserAPI) GetIdleAIs() ([]string, error) {
	return nil, nil
}
