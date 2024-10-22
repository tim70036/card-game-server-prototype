package service

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	ProvideBaseUserService,
	ProvideBaseGameRepoService,
	ProvideBaseEventService,
	ProvideResyncService,
	ProvideBaseSeatStatusService,
	ProvideActionHintService,
	ProvideClubModeUserService,
	ProvideClubModeGameRepoService,
	ProvideBuddyModeGameRepoService,
	ProvideClubModeSeatStatusService,
	ProvideBaseJackpotService,
	ProvideClubModeJackpotService,
)
