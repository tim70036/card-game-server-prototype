package state

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	Init,
	ProvideClosedState,
	ProvideInitState,
	ProvideResetState,
	ProvideWaitUserState,

	ProvideStartRoundState,
	ProvideDealPocketState,
	ProvideEvaluateActionState,
	ProvideCollectChipState,
	ProvideDealCommunityState,

	ProvideWaitActionState,
	ProvideFoldState,
	ProvideCheckState,
	ProvideBetState,
	ProvideCallState,
	ProvideRaiseState,
	ProvideAllInState,

	ProvideDeclareShowdownState,
	ProvideShowdownState,
	ProvideDealRemainCommunityState,
	ProvideEvaluateWinnerState,
	ProvideDeclareWinnerState,
	ProvideJackpotState,
	ProvideEndRoundState,
)
