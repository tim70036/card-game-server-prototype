package state

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	Init,
	ProvideClosedState,
	ProvideInitState,
	ProvideMatchingState,
)
