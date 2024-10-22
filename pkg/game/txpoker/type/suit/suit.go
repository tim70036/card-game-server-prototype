package suit

import (
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
)

type Suit int

const (
	Undefined Suit = 0
	Clubs     Suit = 1
	Diamonds  Suit = 2
	Hearts    Suit = 3
	Spades    Suit = 4
)

func (s Suit) ShortString() string {
	if name, ok := names[s]; ok {
		return name
	}
	return "X"
}

func (s Suit) ToProto() txpokergrpc.Suit {
	proto, ok := protos[s]
	if !ok {
		return txpokergrpc.Suit_UNDEFINED
	}
	return proto
}

var names = map[Suit]string{
	Spades:   "S",
	Hearts:   "H",
	Diamonds: "D",
	Clubs:    "C",
}

var protos = map[Suit]txpokergrpc.Suit{
	Undefined: txpokergrpc.Suit_UNDEFINED,
	Clubs:     txpokergrpc.Suit_CLUBS,
	Diamonds:  txpokergrpc.Suit_DIAMONDS,
	Hearts:    txpokergrpc.Suit_HEARTS,
	Spades:    txpokergrpc.Suit_SPADES,
}
