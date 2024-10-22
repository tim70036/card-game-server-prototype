package event

import (
	"card-game-server-prototype/pkg/game/txpoker/type/hand"
	"card-game-server-prototype/pkg/game/txpoker/type/stage"
)

var HandEventTypeMap = map[hand.HandType]EventType{
	hand.Straight:      Straight,
	hand.Flush:         Flush,
	hand.FullHouse:     FullHouse,
	hand.FourOfAKind:   FourOfAKind,
	hand.StraightFlush: StraightFlush,
	hand.RoyalFlush:    RoyalFlush,
}

var StageGameEventTypeMap = map[stage.Stage]EventType{
	stage.PreFlopStage:  PreflopGame,
	stage.FlopStage:     FlopGame,
	stage.TurnStage:     TurnGame,
	stage.RiverStage:    RiverGame,
	stage.ShowdownStage: ShowdownGame,
}

var StageFoldEventTypeMap = map[stage.Stage]EventType{
	stage.PreFlopStage: PreflopFold,
	stage.FlopStage:    FlopFold,
	stage.TurnStage:    TurnFold,
	stage.RiverStage:   RiverFold,
}

var StageBetEventTypeMap = map[stage.Stage]EventType{
	// 因為大盲，Preflop 只能 Raise，沒有 Bet
	stage.FlopStage:  FlopBet,
	stage.TurnStage:  TurnBet,
	stage.RiverStage: RiverBet,
}

var StageRaiseEventTypeMap = map[stage.Stage]EventType{
	stage.PreFlopStage: PreflopRaise,
}

var StageLastRaiseEventTypeMap = map[stage.Stage]EventType{
	stage.PreFlopStage: PreflopLastRaise,
	stage.FlopStage:    FlopLastRaise,
	stage.TurnStage:    TurnLastRaise,
}

var StageReRaiseEventTypeMap = map[stage.Stage]EventType{
	stage.PreFlopStage: PreflopReRaise,
}

var StageContinueBetEventTypeMap = map[stage.Stage]EventType{
	stage.FlopStage:  FlopContinueBet,
	stage.TurnStage:  TurnContinueBet,
	stage.RiverStage: RiverContinueBet,
}
