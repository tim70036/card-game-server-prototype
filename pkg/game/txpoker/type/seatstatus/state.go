package seatstatus

import (
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"github.com/qmuntal/stateless"
)

type SeatStatusState int

const (
	UndefinedState  SeatStatusState = 0
	StandingState   SeatStatusState = 1
	ReservingState  SeatStatusState = 2
	JoiningState    SeatStatusState = 3
	PlayingState    SeatStatusState = 4
	SittingOutState SeatStatusState = 5
	CashingOutState SeatStatusState = 6
	BuyingInState   SeatStatusState = 7
	WaitingState    SeatStatusState = 8
)

func (s SeatStatusState) String() string {
	name, ok := stateNames[s]
	if !ok {
		return "UndefinedState"
	}
	return name
}

func (s SeatStatusState) ToProto() txpokergrpc.SeatStatusState {
	proto, ok := stateProtos[s]
	if !ok {
		return txpokergrpc.SeatStatusState_UNDEFINED
	}
	return proto
}

func (s SeatStatusState) IsEqual(state stateless.State) bool {
	return s == state.(SeatStatusState)
}

var stateNames = map[SeatStatusState]string{
	UndefinedState:  "UndefinedState",
	StandingState:   "StandingState",
	ReservingState:  "ReservingState",
	JoiningState:    "JoiningState",
	PlayingState:    "PlayingState",
	SittingOutState: "SittingOutState",
	CashingOutState: "CashingOutState",
	BuyingInState:   "BuyingInState",
	WaitingState:    "WaitingState",
}

var stateProtos = map[SeatStatusState]txpokergrpc.SeatStatusState{
	UndefinedState:  txpokergrpc.SeatStatusState_UNDEFINED,
	StandingState:   txpokergrpc.SeatStatusState_STANDING,
	ReservingState:  txpokergrpc.SeatStatusState_RESERVING,
	JoiningState:    txpokergrpc.SeatStatusState_JOINING,
	PlayingState:    txpokergrpc.SeatStatusState_PLAYING,
	SittingOutState: txpokergrpc.SeatStatusState_SITTING_OUT,
	CashingOutState: txpokergrpc.SeatStatusState_CASHING_OUT,
	BuyingInState:   txpokergrpc.SeatStatusState_BUYING_IN,
	WaitingState:    txpokergrpc.SeatStatusState_WAITING,
}
