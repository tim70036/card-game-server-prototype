package state

import (
	"card-game-server-prototype/pkg/core"
)

type Initiator interface{}

func Init(
	game core.Game,
	closedState *ClosedState,
	initState *InitState,
	resetState *ResetState,

	waitUserState *WaitUserState,
	startRoundState *StartRoundState,
	dealPocketState *DealPocketState,
	evaluateActionState *EvaluateActionState,
	collectChipState *CollectChipState,
	dealCommunityState *DealCommunityState,

	waitActionState *WaitActionState,
	foldState *FoldState,
	checkState *CheckState,
	betState *BetState,
	callState *CallState,
	raiseState *RaiseState,
	allInState *AllInState,

	declareShowdownState *DeclareShowdownState,
	showdownState *ShowdownState,
	dealRemainCommunityState *DealRemainCommunityState,
	evaluateWinnerState *EvaluateWinnerState,
	declareWinnerState *DeclareWinnerState,
	jackpotState *JackpotState,
	endRoundState *EndRoundState,
) Initiator {

	game.ConfigTriggerParamsType(GoClosedState)
	game.ConfigErrorState(closedState)

	game.ConfigTriggerParamsType(GoInitState)
	game.ConfigInitState(initState).
		Permit(GoResetState, resetState)

	game.ConfigTriggerParamsType(GoResetState)
	game.ConfigState(resetState).
		Permit(GoWaitUserState, waitUserState)

	game.ConfigTriggerParamsType(GoWaitUserState)
	game.ConfigState(waitUserState).
		Permit(GoStartRoundState, startRoundState).
		Permit(GoClosedState, closedState)

	game.ConfigTriggerParamsType(GoStartRoundState)
	game.ConfigState(startRoundState).
		Permit(GoDealPocketState, dealPocketState).
		Permit(GoWaitUserState, waitUserState)

	game.ConfigTriggerParamsType(GoDealPocketState)
	game.ConfigState(dealPocketState).
		Permit(GoEvaluateActionState, evaluateActionState)

	game.ConfigTriggerParamsType(GoEvaluateActionState)
	game.ConfigState(evaluateActionState).
		Permit(GoWaitActionState, waitActionState).
		Permit(GoCollectChipState, collectChipState)

	game.ConfigTriggerParamsType(GoCollectChipState)
	game.ConfigState(collectChipState).
		Permit(GoDealCommunityState, dealCommunityState).
		Permit(GoDeclareShowdownState, declareShowdownState).
		Permit(GoShowdownState, showdownState).
		Permit(GoEvaluateWinnerState, evaluateWinnerState)

	game.ConfigTriggerParamsType(GoDealCommunityState)
	game.ConfigState(dealCommunityState).
		Permit(GoEvaluateActionState, evaluateActionState)

	game.ConfigTriggerParamsType(GoWaitActionState)
	game.ConfigState(waitActionState).
		Permit(GoFoldState, foldState).
		Permit(GoCheckState, checkState).
		Permit(GoBetState, betState).
		Permit(GoCallState, callState).
		Permit(GoRaiseState, raiseState).
		Permit(GoAllInState, allInState)

	game.ConfigTriggerParamsType(GoFoldState)
	game.ConfigState(foldState).
		Permit(GoEvaluateActionState, evaluateActionState)

	game.ConfigTriggerParamsType(GoCheckState)
	game.ConfigState(checkState).
		Permit(GoEvaluateActionState, evaluateActionState)

	game.ConfigTriggerParamsType(GoBetState)
	game.ConfigState(betState).
		Permit(GoEvaluateActionState, evaluateActionState)

	game.ConfigTriggerParamsType(GoCallState)
	game.ConfigState(callState).
		Permit(GoEvaluateActionState, evaluateActionState)

	game.ConfigTriggerParamsType(GoRaiseState)
	game.ConfigState(raiseState).
		Permit(GoEvaluateActionState, evaluateActionState)

	game.ConfigTriggerParamsType(GoAllInState)
	game.ConfigState(allInState).
		Permit(GoEvaluateActionState, evaluateActionState)

	game.ConfigTriggerParamsType(GoDeclareShowdownState)
	game.ConfigState(declareShowdownState).
		Permit(GoShowdownState, showdownState)

	game.ConfigTriggerParamsType(GoShowdownState)
	game.ConfigState(showdownState).
		Permit(GoDealRemainCommunityState, dealRemainCommunityState)

	game.ConfigTriggerParamsType(GoDealRemainCommunityState)
	game.ConfigState(dealRemainCommunityState).
		Permit(GoEvaluateWinnerState, evaluateWinnerState)

	game.ConfigTriggerParamsType(GoEvaluateWinnerState)
	game.ConfigState(evaluateWinnerState).
		Permit(GoDeclareWinnerState, declareWinnerState)

	game.ConfigTriggerParamsType(GoDeclareWinnerState)
	game.ConfigState(declareWinnerState).
		PermitReentry(GoDeclareWinnerState).
		Permit(GoJackpotState, jackpotState).
		Permit(GoEndRoundState, endRoundState)

	game.ConfigTriggerParamsType(GoJackpotState)
	game.ConfigState(jackpotState).
		Permit(GoEndRoundState, endRoundState)

	game.ConfigTriggerParamsType(GoEndRoundState)
	game.ConfigState(endRoundState).
		Permit(GoResetState, resetState)

	return nil
}
