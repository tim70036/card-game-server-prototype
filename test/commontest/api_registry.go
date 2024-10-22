package commontest

import (
	commonapi "card-game-server-prototype/pkg/common/api"
)

type apiRegistry struct {
	baseRoomAPI *commonapi.BaseRoomAPI
	baseUserAPI *commonapi.BaseUserAPI
}

func ProvideAPIRegistry(
	baseRoomAPI *commonapi.BaseRoomAPI,
	baseUserAPI *commonapi.BaseUserAPI,
) *apiRegistry {
	return &apiRegistry{
		baseRoomAPI: baseRoomAPI,
		baseUserAPI: baseUserAPI,
	}
}
