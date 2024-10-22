package state

import "card-game-server-prototype/pkg/core"

type Initiator interface{}

func Init(
	game core.Game,
	closedState *ClosedState,
	initState *InitState,
	matchingState *MatchingState,
) Initiator {
	game.ConfigTriggerParamsType(GoClosedState)
	game.ConfigErrorState(closedState)

	game.ConfigTriggerParamsType(GoInitState)
	game.ConfigInitState(initState).
		Permit(GoMatchingState, matchingState)

	game.ConfigTriggerParamsType(GoMatchingState)
	game.ConfigState(matchingState).
		Permit(GoClosedState, closedState)

	return nil
}
