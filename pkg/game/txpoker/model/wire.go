package model

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	ProvideGameInfo,
	ProvideGameSetting,
	ProvideEventGroup,
	ProvidePlayerGroup,
	ProvideSeatStatusGroup,
	ProvideActionHintGroup,
	ProvideTable,
	ProvidePlaySettingGroup,
	ProvideReplay,
	ProvideStatsGroup,
	ProvideChipCacheGroup,
	ProvideStatsCacheGroup,
	ProvideForceBuyInGroup,
	ProvideUserCacheGroup,
	ProvideTableProfitsGroup,
)
