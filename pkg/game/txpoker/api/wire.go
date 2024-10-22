package api

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	ProvideBaseGameAPI,
	ProvideLocalGameAPI,
	ProvideClubModeUserAPI,
	ProvideClubModeMemberAPI,
	ProvideClubModeGameAPI,
)
