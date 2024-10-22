package util

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProvideLoggerFactory,
	ProvideHttpClient,
	ProvideRedisClient,
)
