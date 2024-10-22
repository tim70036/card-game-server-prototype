package constant

import "time"

const (
	ShutdownGracefulPeriod         = 30 * time.Second
	EmptyWaitingRoomGracefulPeriod = 60 * time.Second
	RoomPingInterval               = 60 * time.Second
	RoomPingFailThreshold          = 3
	PerConnectionBufferSize        = 16
	PerNotificationBufferSize      = 4
	CoreGameRunnerBufferSize       = 256
	ChatHistoryLimit               = 5
)
