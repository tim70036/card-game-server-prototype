package api

import (
	"card-game-server-prototype/pkg/core"
)

type RoomAPI interface {
	GetDetail(roomId string) (*GetRoomDetailResp, error)
	Heartbeat(roomId string) error
	EnterRoom(roomId string, uid core.Uid) error
	LeaveRoom(roomId string, uid core.Uid) error
	CloseRoom(roomId string) error
}
