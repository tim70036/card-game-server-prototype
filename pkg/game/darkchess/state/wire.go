package state

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProvideClosedState,
	ProvideInitState,
	ProvideResetGameState,
	ProvideWaitUserState,
	ProvideWaitingRoomState,
	ProvideStartGameState,
	ProvideResetRoundState,
	ProvideStartRoundState,
	ProvideRoundScoreboardState,
	ProvideEndRoundState,
	ProvideGameScoreboardState,
	ProvideEndGameState,

	ProvidePickFirstState,
	ProvideStartTurnState,
	ProvideWaitActionState,
	ProvideRevealState,
	ProvideMoveState,
	ProvideCaptureState,
	ProvideEndTurnState,
	ProvideDrawState,
	ProvideSurrenderState,
	ProvideShowRoundResultState,
)
