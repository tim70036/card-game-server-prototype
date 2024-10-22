package model

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProvidePlaySettingGroup,
	ProvideGameScoreboard,
	ProvideEventGroup,
	ProvideGameInfo,
	ProvidePlayerGroup,
	ProvideRoundScoreboard,
	ProvideRoundScoreboardRecords,
	ProvidePickBoard,
	ProvideBoard,
	ProvideActionHintGroup,
	ProvideCapturedPieces,
	ProvideReplayGroup,
)
