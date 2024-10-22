package state

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/state"
)

type Initiator interface{}

func Init(
	game core.Game,
	closedState *ClosedState,
	initState *InitState,
	resetState *ResetState,
	startRoundState *StartRoundState,

	dealPocketState *state.DealPocketState,
	evaluateActionState *state.EvaluateActionState,
	collectChipState *state.CollectChipState,
	dealCommunityState *state.DealCommunityState,

	waitActionState *state.WaitActionState,
	foldState *state.FoldState,
	checkState *state.CheckState,
	betState *state.BetState,
	callState *state.CallState,
	raiseState *state.RaiseState,
	allInState *state.AllInState,

	declareShowdownState *state.DeclareShowdownState,
	showdownState *state.ShowdownState,
	dealRemainCommunityState *state.DealRemainCommunityState,
	evaluateWinnerState *state.EvaluateWinnerState,
	declareWinnerState *state.DeclareWinnerState,
	jackpotState *state.JackpotState,
	endRoundState *state.EndRoundState,
) Initiator {

	game.ConfigTriggerParamsType(state.GoClosedState)
	game.ConfigErrorState(closedState)

	game.ConfigTriggerParamsType(state.GoInitState)
	game.ConfigInitState(initState).
		Permit(state.GoResetState, resetState)

	game.ConfigTriggerParamsType(state.GoResetState)
	game.ConfigState(resetState).
		Permit(state.GoStartRoundState, startRoundState).
		Permit(state.GoClosedState, closedState)

	game.ConfigTriggerParamsType(state.GoStartRoundState)
	game.ConfigState(startRoundState).
		Permit(state.GoDealPocketState, dealPocketState)

	game.ConfigTriggerParamsType(state.GoDealPocketState)
	game.ConfigState(dealPocketState).
		Permit(state.GoEvaluateActionState, evaluateActionState)

	game.ConfigTriggerParamsType(state.GoEvaluateActionState)
	game.ConfigState(evaluateActionState).
		Permit(state.GoWaitActionState, waitActionState).
		Permit(state.GoCollectChipState, collectChipState)

	game.ConfigTriggerParamsType(state.GoCollectChipState)
	game.ConfigState(collectChipState).
		Permit(state.GoDealCommunityState, dealCommunityState).
		Permit(state.GoDeclareShowdownState, declareShowdownState).
		Permit(state.GoShowdownState, showdownState).
		Permit(state.GoEvaluateWinnerState, evaluateWinnerState)

	game.ConfigTriggerParamsType(state.GoDealCommunityState)
	game.ConfigState(dealCommunityState).
		Permit(state.GoEvaluateActionState, evaluateActionState)

	game.ConfigTriggerParamsType(state.GoWaitActionState)
	game.ConfigState(waitActionState).
		Permit(state.GoFoldState, foldState).
		Permit(state.GoCheckState, checkState).
		Permit(state.GoBetState, betState).
		Permit(state.GoCallState, callState).
		Permit(state.GoRaiseState, raiseState).
		Permit(state.GoAllInState, allInState)

	game.ConfigTriggerParamsType(state.GoFoldState)
	game.ConfigState(foldState).
		Permit(state.GoEvaluateActionState, evaluateActionState)

	game.ConfigTriggerParamsType(state.GoCheckState)
	game.ConfigState(checkState).
		Permit(state.GoEvaluateActionState, evaluateActionState)

	game.ConfigTriggerParamsType(state.GoBetState)
	game.ConfigState(betState).
		Permit(state.GoEvaluateActionState, evaluateActionState)

	game.ConfigTriggerParamsType(state.GoCallState)
	game.ConfigState(callState).
		Permit(state.GoEvaluateActionState, evaluateActionState)

	game.ConfigTriggerParamsType(state.GoRaiseState)
	game.ConfigState(raiseState).
		Permit(state.GoEvaluateActionState, evaluateActionState)

	game.ConfigTriggerParamsType(state.GoAllInState)
	game.ConfigState(allInState).
		Permit(state.GoEvaluateActionState, evaluateActionState)

	game.ConfigTriggerParamsType(state.GoDeclareShowdownState)
	game.ConfigState(declareShowdownState).
		Permit(state.GoShowdownState, showdownState)

	game.ConfigTriggerParamsType(state.GoShowdownState)
	game.ConfigState(showdownState).
		Permit(state.GoDealRemainCommunityState, dealRemainCommunityState)

	game.ConfigTriggerParamsType(state.GoDealRemainCommunityState)
	game.ConfigState(dealRemainCommunityState).
		Permit(state.GoEvaluateWinnerState, evaluateWinnerState)

	game.ConfigTriggerParamsType(state.GoEvaluateWinnerState)
	game.ConfigState(evaluateWinnerState).
		Permit(state.GoDeclareWinnerState, declareWinnerState)

	game.ConfigTriggerParamsType(state.GoDeclareWinnerState)
	game.ConfigState(declareWinnerState).
		PermitReentry(state.GoDeclareWinnerState).
		Permit(state.GoJackpotState, jackpotState).
		Permit(state.GoEndRoundState, endRoundState)

	game.ConfigTriggerParamsType(state.GoJackpotState)
	game.ConfigState(jackpotState).
		Permit(state.GoEndRoundState, endRoundState)

	game.ConfigTriggerParamsType(state.GoEndRoundState)
	game.ConfigState(endRoundState).
		Permit(state.GoResetState, resetState)

	return nil
}
