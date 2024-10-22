package constant

import "time"

const (
	PlayerPerMatch           = 6 // gameSetting.TableSize?
	DisconnectGracefulPeriod = 120 * time.Second
	IdlingTimeoutDuration    = 2 * time.Minute
)
