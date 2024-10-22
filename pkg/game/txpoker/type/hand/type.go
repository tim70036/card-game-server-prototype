package hand

import (
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
)

type HandType int

const (
	Undefined     HandType = 0
	HighCard      HandType = 1
	Pair          HandType = 2
	TwoPair       HandType = 3
	ThreeOfAKind  HandType = 4
	Straight      HandType = 5
	Flush         HandType = 6
	FullHouse     HandType = 7
	FourOfAKind   HandType = 8
	StraightFlush HandType = 9
	RoyalFlush    HandType = 10
)

func (t HandType) ToProto() txpokergrpc.PokerHandType {
	proto, ok := protos[t]
	if !ok {
		return txpokergrpc.PokerHandType_UNDEFINED
	}
	return proto
}

func (t HandType) String() string {
	name, ok := names[t]
	if !ok {
		return "Undefined"
	}
	return name
}

var names = map[HandType]string{
	Undefined:     "Undefined",
	HighCard:      "HighCard",
	Pair:          "Pair",
	TwoPair:       "TwoPair",
	ThreeOfAKind:  "ThreeOfAKind",
	Straight:      "Straight",
	Flush:         "Flush",
	FullHouse:     "FullHouse",
	FourOfAKind:   "FourOfAKind",
	StraightFlush: "StraightFlush",
	RoyalFlush:    "RoyalFlush",
}

var protos = map[HandType]txpokergrpc.PokerHandType{
	Undefined:     txpokergrpc.PokerHandType_UNDEFINED,
	HighCard:      txpokergrpc.PokerHandType_HIGH_CARD,
	Pair:          txpokergrpc.PokerHandType_PAIR,
	TwoPair:       txpokergrpc.PokerHandType_TWO_PAIR,
	ThreeOfAKind:  txpokergrpc.PokerHandType_THREE_OF_A_KIND,
	Straight:      txpokergrpc.PokerHandType_STRAIGHT,
	Flush:         txpokergrpc.PokerHandType_FLUSH,
	FullHouse:     txpokergrpc.PokerHandType_FULL_HOUSE,
	FourOfAKind:   txpokergrpc.PokerHandType_FOUR_OF_A_KIND,
	StraightFlush: txpokergrpc.PokerHandType_STRAIGHT_FLUSH,
	RoyalFlush:    txpokergrpc.PokerHandType_ROYAL_FLUSH,
}
