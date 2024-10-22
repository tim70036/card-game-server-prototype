package service

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProvideResyncService,
	ProvideBaseUserService,
	ProvideBaseGameRepoService,
	ProvideBuddyModeGameRepoService,
	ProvideBuddyModeUserService,
	ProvideBaseEventService,
	ProvideBoardService,
	ProvidePlayerService,
)
