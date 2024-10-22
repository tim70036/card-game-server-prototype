package action

import (
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
)

type ActionType int

// https://en.wikipedia.org/wiki/Betting_in_poker
const (
	Undefined ActionType = 0

	// Fold: At any point during the betting round, a player can fold.
	// They choose to discard their hand and forfeit any chips they
	// have already invested in the pot.
	Fold ActionType = 1

	// Check: If no bet has been made before their turn, a player can
	// choose to check. This means they don't place any additional
	// chips into the pot, and the action moves to the next player.
	Check ActionType = 2

	// Bet: A player can make the first bet in a betting round. They
	// put chips into the pot to start the betting. Under normal
	// circumstances, all other players still in the pot must either
	// call the full amount of the bet or raise if they wish to remain
	// in, the only exceptions being when a player does not have
	// sufficient stake remaining to call the full amount of the bet
	// (in which case they may either call with their remaining stake
	// to go "all-in" or fold) or when the player is already all-in.
	Bet ActionType = 3

	// Call: If a bet has been made before their turn, a player can
	// choose to call. This means they match the current bet amount to
	// stay in the hand and continue playing. A betting round ends
	// when all active players have bet an equal amount. If no
	// opponents call a player's bet or raise, the player wins the
	// pot.
	Call ActionType = 4

	// Raise: After a bet has been made before their turn, a player
	// can choose to raise. They put more chips into the pot than the
	// current bet, forcing other players to match the higher bet
	// amount to stay in the hand.
	Raise ActionType = 5

	// All-In: A player faced with a current bet who wishes to call
	// but has insufficient remaining stake (folding does not require
	// special rules) may bet the remainder of their stake and declare
	// themselves all-in. They may now hold onto their cards for the
	// remainder of the deal as if they had called every bet, but may
	// not win any more money from any player above the amount of
	// their bet. In no-limit games, a player may also go all in, that
	// is, betting their entire stack at any point during a betting
	// round.
	AllIn ActionType = 6

	SB ActionType = 7

	BB ActionType = 8

	Showdown ActionType = 9

	WinPot ActionType = 10
)

func (t ActionType) String() string {
	name, ok := actionNames[t]
	if !ok {
		return "X"
	}
	return name
}

func (t ActionType) ToProto() txpokergrpc.BetActionType {
	proto, ok := actionProtos[t]
	if !ok {
		return txpokergrpc.BetActionType_UNDEFINED
	}
	return proto
}

var actionNames = map[ActionType]string{
	Undefined: "X",
	Fold:      "Fold",
	Check:     "Check",
	Bet:       "Bet",
	Call:      "Call",
	Raise:     "Raise",
	AllIn:     "AllIn",
	SB:        "SB",
	BB:        "BB",
	Showdown:  "Showdown",
	WinPot:    "WInPot",
}

var actionProtos = map[ActionType]txpokergrpc.BetActionType{
	Undefined: txpokergrpc.BetActionType_UNDEFINED,
	Fold:      txpokergrpc.BetActionType_FOLD,
	Check:     txpokergrpc.BetActionType_CHECK,
	Bet:       txpokergrpc.BetActionType_BET,
	Call:      txpokergrpc.BetActionType_CALL,
	Raise:     txpokergrpc.BetActionType_RAISE,
	AllIn:     txpokergrpc.BetActionType_ALL_IN,
	SB:        txpokergrpc.BetActionType_SB,
	BB:        txpokergrpc.BetActionType_BB,
}
