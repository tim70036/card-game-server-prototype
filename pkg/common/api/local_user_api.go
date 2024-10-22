package api

import (
	"fmt"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/common/type/gametype"
	"card-game-server-prototype/pkg/core"
	"github.com/samber/lo"
)

type LocalUserAPI struct {
	userGroup *commonmodel.UserGroup
}

func ProvideLocalUserAPI(userGroup *commonmodel.UserGroup) *LocalUserAPI {
	return &LocalUserAPI{
		userGroup: userGroup,
	}
}

func (api *LocalUserAPI) FetchUserDetail(uid core.Uid) (*UserDetailResponse, error) {
	user, userExists := api.userGroup.Data[uid]
	isAI := false
	if userExists {
		isAI = user.IsAI
	}

	return &UserDetailResponse{
		Data: struct {
			Uid       string "json:\"uid\""
			ShortUid  int    "json:\"shortUid\""
			Name      string "json:\"name\""
			Cash      int    "json:\"cash\""
			Level     int    "json:\"level\""
			RoomCards int    "json:\"roomCards\""
			IsAi      int    "json:\"isAi\""
		}{
			Uid:       uid.String(),
			ShortUid:  312344,
			Name:      fmt.Sprintf("local%s", uid.String()[:4]),
			Cash:      10000,
			Level:     1,
			RoomCards: 100,
			IsAi:      lo.Ternary(isAI, 1, 0),
		},
	}, nil
}

func (api *LocalUserAPI) ExchangeChip(uid core.Uid, gameType gametype.GameType, amount int) error {
	return nil
}

func (api *LocalUserAPI) GetIdleAIs() ([]string, error) {
	return nil, nil
}
