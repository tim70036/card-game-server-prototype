package participant

import "card-game-server-prototype/pkg/grpc/txpokergrpc"

type State int

const (
	UndefinedState    State = 0
	ObservingState    State = 1
	MatchingState     State = 2
	PlayingState      State = 3
	BuyingInState     State = 4
	CashingOutState   State = 5
	ExitingMatchState State = 6
)

func (s State) String() string {
	name, ok := stateNames[s]
	if !ok {
		return "UndefinedState"
	}
	return name
}

func (s State) ToProto() txpokergrpc.ParticipantState {
	proto, ok := stateProtos[s]
	if !ok {
		return txpokergrpc.ParticipantState_UNDEFINED
	}
	return proto
}

var stateNames = map[State]string{
	UndefinedState:    "UndefinedState",
	ObservingState:    "ObservingState",
	MatchingState:     "MatchingState",
	PlayingState:      "PlayingState",
	BuyingInState:     "BuyingInState",
	CashingOutState:   "CashingOutState",
	ExitingMatchState: "ExitingMatchState",
}

var stateProtos = map[State]txpokergrpc.ParticipantState{
	UndefinedState:    txpokergrpc.ParticipantState_UNDEFINED,
	ObservingState:    txpokergrpc.ParticipantState_OBSERVING,
	MatchingState:     txpokergrpc.ParticipantState_MATCHING,
	PlayingState:      txpokergrpc.ParticipantState_PLAYING,
	BuyingInState:     txpokergrpc.ParticipantState_BUYING_IN,
	CashingOutState:   txpokergrpc.ParticipantState_CASHING_OUT,
	ExitingMatchState: txpokergrpc.ParticipantState_EXITING_MATCH,
}
