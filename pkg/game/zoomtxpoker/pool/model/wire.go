package model

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProvideGameSetting,
	ProvideParticipantGroup,
	ProvidePlaySettingGroup,
	ProvideStatsGroup,
	ProvidePlayedHistoryGroup,
	ProvideForceBuyInGroup,
	ProvideTableProfitsGroup,
)
