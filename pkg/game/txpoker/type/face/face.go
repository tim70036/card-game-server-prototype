package face

import (
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
)

type Face int

const (
	Undefined Face = 0
	Two       Face = 2
	Three     Face = 3
	Four      Face = 4
	Five      Face = 5
	Six       Face = 6
	Seven     Face = 7
	Eight     Face = 8
	Nine      Face = 9
	Ten       Face = 10
	Jack      Face = 11
	Queen     Face = 12
	King      Face = 13
	Ace       Face = 14
)

func (f Face) String() string {
	if name, ok := names[f]; ok {
		return name
	}
	return "0"
}

func (f Face) ToProto() txpokergrpc.Face {
	proto, ok := protos[f]
	if !ok {
		return txpokergrpc.Face_UNDEFINED
	}
	return proto
}

var names = map[Face]string{
	Ace:   "A",
	Two:   "2",
	Three: "3",
	Four:  "4",
	Five:  "5",
	Six:   "6",
	Seven: "7",
	Eight: "8",
	Nine:  "9",
	Ten:   "T",
	Jack:  "J",
	Queen: "Q",
	King:  "K",
}

var protos = map[Face]txpokergrpc.Face{
	Undefined: txpokergrpc.Face_UNDEFINED,
	Two:       txpokergrpc.Face_TWO,
	Three:     txpokergrpc.Face_THREE,
	Four:      txpokergrpc.Face_FOUR,
	Five:      txpokergrpc.Face_FIVE,
	Six:       txpokergrpc.Face_SIX,
	Seven:     txpokergrpc.Face_SEVEN,
	Eight:     txpokergrpc.Face_EIGHT,
	Nine:      txpokergrpc.Face_NINE,
	Ten:       txpokergrpc.Face_TEN,
	Jack:      txpokergrpc.Face_JACK,
	Queen:     txpokergrpc.Face_QUEEN,
	King:      txpokergrpc.Face_KING,
	Ace:       txpokergrpc.Face_ACE,
}
