package api

import (
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"

	"github.com/imroc/req/v3"
)

type BaseRoomAPI struct {
	httpClient *req.Client
}

func ProvideBaseRoomAPI(httpClient *req.Client, apiCFG *config.APIConfig, config *config.Config) *BaseRoomAPI {
	return &BaseRoomAPI{
		httpClient: httpClient.Clone().
			SetBaseURL("https://"+*apiCFG.MainServerHost+"/game/room").
			SetCommonHeader("jtoken", *apiCFG.MainServerAPIKey),
	}
}

func (api *BaseRoomAPI) GetDetail(roomId string) (*GetRoomDetailResp, error) {
	resp := &GetRoomDetailResp{}
	err := api.httpClient.Get("room").
		AddQueryParam("roomid", roomId).
		Do().
		Into(resp)
	return resp, err
}

func (api *BaseRoomAPI) Heartbeat(roomId string) error {
	req := &heartbeatRequest{
		RoomId: roomId,
	}

	return api.httpClient.Post("heartbeat").
		SetBodyJsonMarshal(req).
		Do().
		Err
}

func (api *BaseRoomAPI) EnterRoom(roomId string, uid core.Uid) error {
	req := &enterRoomRequest{
		RoomId: roomId,
		Uid:    uid.String(),
	}

	return api.httpClient.Post("session").
		SetBodyJsonMarshal(req).
		Do().
		Err
}

func (api *BaseRoomAPI) LeaveRoom(roomId string, uid core.Uid) error {
	// TODO: ask main server to use req body instead of query param. It
	// should be same as EnterRoom.

	//  req := &leaveRoomRequest{
	//  RoomId: roomId,
	//  Uid:    uid.String(),
	// }

	return api.httpClient.Delete("session").
		AddQueryParam("roomid", roomId).
		AddQueryParam("uid", uid.String()).
		// SetBodyJsonMarshal(req).
		Do().
		Err
}

func (api *BaseRoomAPI) CloseRoom(roomId string) error {
	req := &closeRoomRequest{
		RoomId: roomId,
	}

	return api.httpClient.Delete("room").
		SetBodyJsonMarshal(req).
		Do().
		Err
}
