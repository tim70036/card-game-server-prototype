package role

import (
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
)

type Role int

const (
	Undefined Role = 0
	SB        Role = 1
	BB        Role = 2
	UTG       Role = 3
	UTG1      Role = 4
	UTG2      Role = 5
	LJ        Role = 6
	HJ        Role = 7
	MP        Role = 8
	CO        Role = 9
	BTN       Role = 10
)

func (r Role) String() string {
	name, ok := names[r]
	if !ok {
		return "Undefined"
	}
	return name
}

func (r Role) ToProto() txpokergrpc.Role {
	proto, ok := protos[r]
	if !ok {
		return txpokergrpc.Role_UNDEFINED
	}
	return proto
}

var names = map[Role]string{
	Undefined: "Undefined",
	SB:        "SB",
	BB:        "BB",
	UTG:       "UTG",
	UTG1:      "UTG1",
	UTG2:      "UTG2",
	LJ:        "LJ",
	HJ:        "HJ",
	MP:        "MP",
	CO:        "CO",
	BTN:       "BTN",
}

var protos = map[Role]txpokergrpc.Role{
	Undefined: txpokergrpc.Role_UNDEFINED,
	SB:        txpokergrpc.Role_SB,
	BB:        txpokergrpc.Role_BB,
	UTG:       txpokergrpc.Role_UTG,
	UTG1:      txpokergrpc.Role_UTG1,
	UTG2:      txpokergrpc.Role_UTG2,
	LJ:        txpokergrpc.Role_LJ,
	HJ:        txpokergrpc.Role_HJ,
	MP:        txpokergrpc.Role_MP,
	CO:        txpokergrpc.Role_CO,
	BTN:       txpokergrpc.Role_BTN,
}
