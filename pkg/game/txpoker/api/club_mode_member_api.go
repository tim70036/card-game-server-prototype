package api

import (
	"fmt"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/util"
	"github.com/imroc/req/v3"
	"github.com/samber/lo"
)

type ClubModeMemberApi struct {
	httpClient *req.Client
	cfg        *config.Config
}

func ProvideClubModeMemberAPI(httpClient *req.Client, apiCFG *config.APIConfig, cfg *config.Config) *ClubModeMemberApi {
	return &ClubModeMemberApi{
		cfg: cfg,
		httpClient: httpClient.Clone().
			SetBaseURL("https://"+*apiCFG.MainServerHost).
			SetCommonHeader("jtoken", *apiCFG.MainServerAPIKey),
	}
}

func (api *ClubModeMemberApi) FetchDetail(clubId string, uids ...core.Uid) (*ClubMemberDetailResponse, error) {
	uidQueries := util.JoinStrings(lo.Map(uids, func(uid core.Uid, i int) string {
		return fmt.Sprintf("&uid[%v]=%s", i, uid.String())
	}))

	path := fmt.Sprintf("/game/game/lmj-member-info?clubid=%s%s", clubId, uidQueries)

	request := api.httpClient.Get(path)
	resp := &ClubMemberDetailResponse{}
	err := request.Do().Into(resp)

	return resp, err
}
