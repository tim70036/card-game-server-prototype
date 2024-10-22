package service

import commonapi "card-game-server-prototype/pkg/common/api"

type RoomService interface {
	GetDetail() (*commonapi.RoomDetail, error) // todo: move totxpoker
	RunPingLoop() error
	Close() error
	FetchRoomInfo() error
}
