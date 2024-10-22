package api

import (
	"card-game-server-prototype/pkg/core"
)

type LocalRoomAPI struct {
}

func ProvideLocalRoomAPI() *LocalRoomAPI {
	return &LocalRoomAPI{}
}

func (api *LocalRoomAPI) GetDetail(roomId string) (*GetRoomDetailResp, error) {
	return &GetRoomDetailResp{}, nil
}

func (api *LocalRoomAPI) Heartbeat(roomId string) error {
	return nil
}

func (api *LocalRoomAPI) EnterRoom(roomId string, uid core.Uid) error {
	return nil
}

func (api *LocalRoomAPI) LeaveRoom(roomId string, uid core.Uid) error {
	return nil
}

func (api *LocalRoomAPI) CloseRoom(roomId string) error {
	return nil
}
