package model

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProvideRoomInfo,
	ProvideUserGroup,
	ProvideBuddyGroup,
)
