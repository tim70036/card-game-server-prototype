package role

import (
	"fmt"
	"card-game-server-prototype/pkg/core"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sort"
)

var rolesTable = map[int][]Role{
	2: {BB, SB},
	3: {BB, BTN, SB},
	4: {BB, UTG, BTN, SB},
	5: {BB, UTG, CO, BTN, SB},
	6: {BB, UTG, MP, CO, BTN, SB},

	7: {BB, UTG, LJ, HJ, CO, BTN, SB},
	8: {BB, UTG, UTG1, LJ, HJ, CO, BTN, SB},
	9: {BB, UTG, UTG1, UTG2, LJ, HJ, CO, BTN, SB},
}

func GetRoles(playerCount int) ([]Role, error) {
	roles, ok := rolesTable[playerCount]
	if !ok {
		return nil, fmt.Errorf("cannot found the correspond roles setting with player count: %d", playerCount)
	}

	return roles, nil
}

func EvalRoleAssignment(uids core.UidList, tableUids map[int]core.Uid, lastBBSeatId int) (map[core.Uid]Role, error) {
	if len(uids) < 2 {
		return nil, status.Errorf(codes.InvalidArgument, "uids must have at least 2 elements but only has: %+v", uids)
	}

	seatIds := lo.Keys(lo.PickByValues(tableUids, uids))
	if len(seatIds) != len(uids) {
		return nil, status.Errorf(codes.NotFound, "cannot found the correspond seat id in seat status group with uids: %+v", uids)
	}

	sort.Ints(seatIds)
	nextBBIdx := 0 // Edge case: when lastBBSeatId is the last seatId, nextBBSeatId will be the first seatId.
	if lastBBSeatId >= 0 {
		for idx, seatId := range seatIds {
			if seatId > lastBBSeatId {
				nextBBIdx = idx
				break
			}
		}
	}

	// Assign roles
	assignment := make(map[core.Uid]Role)
	roles, err := GetRoles(len(seatIds))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot found roles setting with player count: %d: %w", len(seatIds), err)
	}

	nextIdx := nextBBIdx
	for i := 0; i < len(seatIds); i++ {
		uid := tableUids[seatIds[nextIdx]]
		assignment[uid] = roles[i]
		nextIdx = (nextIdx + 1) % len(seatIds)
	}

	return assignment, nil
}
