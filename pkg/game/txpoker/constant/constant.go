package constant

import "time"

const (
	StartGameUserCount        = 2
	MaxUserCount              = 50
	MaxSeatId                 = 8
	PocketCardsPerPlayer      = 2
	MaxRoundIdHistoryCount    = 20
	DisconnectGracefulPeriod  = 120 * time.Second
	ReservingGracefulPeriod   = 10 * time.Second
	CacheExpirationDuration   = time.Hour
	ClearCacheInterval        = time.Minute
	DonateJackpotIfBBReaching = 12
	MaxIdleRounds             = 2
	TableProfitsMaxSize       = 10000
	TableProfitRenderMaxSize  = 50
)

const (
	CloseRoomAmount = 20
	CloseRoomPeriod = 3
)
