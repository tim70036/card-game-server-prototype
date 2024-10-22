package api

import (
	"errors"
	"card-game-server-prototype/pkg/common/type/gametype"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"strconv"

	"github.com/imroc/req/v3"
)

type BaseUserAPI struct {
	httpClient *req.Client
}

func ProvideBaseUserAPI(httpClient *req.Client, apiCFG *config.APIConfig, config *config.Config) *BaseUserAPI {
	return &BaseUserAPI{
		httpClient: httpClient.Clone().
			SetBaseURL("https://"+*apiCFG.MainServerHost+"/game").
			SetCommonHeader("jtoken", *apiCFG.MainServerAPIKey),
	}
}

func (api *BaseUserAPI) FetchUserDetail(uid core.Uid) (*UserDetailResponse, error) {
	resp := &UserDetailResponse{}
	err := api.httpClient.Get("user/detail").
		AddQueryParam("uid", uid.String()).
		Do().
		Into(resp)

	return resp, err
}

func (api *BaseUserAPI) ExchangeChip(uid core.Uid, gameType gametype.GameType, amount int) error {
	req := &exchangeChipRequest{
		Uid:      uid.String(),
		GameType: string(gameType),
		Amount:   strconv.Itoa(amount),
	}
	return api.httpClient.Put("inventory/chip").
		SetBodyJsonMarshal(req).
		Do().
		Err
}

func (api *BaseUserAPI) GetIdleAIs() ([]string, error) {
	type Resp struct {
		Data struct {
			IdleCpuList []string `json:"idleCpuUids"`
		} `json:"data"`
	}
	resp := &Resp{}

	if err := api.httpClient.Get("user/idle-cpu").
		Do().
		Into(resp); err != nil {
		return nil, err
	}

	if len(resp.Data.IdleCpuList) == 0 {
		return nil, errors.New("no idle AI")
	}

	return resp.Data.IdleCpuList, nil
}
