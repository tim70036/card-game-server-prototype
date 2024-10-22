package service

import (
	"card-game-server-prototype/pkg/core"
	role2 "card-game-server-prototype/pkg/game/txpoker/type/role"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/model"
	"github.com/samber/lo"
	"math"
	"sort"
)

type RoleService struct {
	playedHistoryGroup *model.PlayedHistoryGroup
}

func ProvideRoleService(
	playedHistoryGroup *model.PlayedHistoryGroup,
) *RoleService {
	return &RoleService{
		playedHistoryGroup: playedHistoryGroup,
	}
}

// Assign 分派角色(要改用%數計算，用次數計算會有高機率出現【晚進來的玩家一定會一直跳小盲】的情況)
// 1. 前一輪擔任 BB or SB 的人，此輪不會再被分配 BB or SB。
// 2. 原本 shuffle users 後直接分配就可以達到平均值了，但這裡還有使用 natural8 的規則。（以下稱 N8 rules）
func (s *RoleService) Assign(matchedUidGroup core.UidList) (map[core.Uid]role2.Role, error) {
	// N8 rules - Rush & Cash 的分派角色規則目的是要 BB 和 SB 的數量達到平均值，作法如下：
	// 每局玩家將隨機就座，但盲注位將給較少坐在盲注位置的人。
	// 如果所有玩家坐在大盲%數相同，則將根據坐在小盲%數最少的玩家入座。
	// 如果所有玩家坐在小盲%數相同，則隨機選擇座位。
	// 這樣，隨著時間的推移，每個人都將支付相同百分比的盲注。

	players := make(core.UidList, len(matchedUidGroup))
	copy(players, matchedUidGroup)

	roles, err := role2.GetRoles(len(players))
	if err != nil {
		return nil, err
	}

	var playedRoles []*model.PlayedHistory

	// Init data
	for _, uid := range players {
		if _, ok := s.playedHistoryGroup.Data[uid]; !ok {
			playedRoles = append(playedRoles, &model.PlayedHistory{
				Uid:            uid,
				CountRoles:     map[role2.Role]int{},
				LastPlayedRole: role2.Undefined,
			})
			continue
		}

		copyCountRoles := map[role2.Role]int{}
		for k, v := range s.playedHistoryGroup.Data[uid].CountRoles {
			copyCountRoles[k] = v
		}

		playedRoles = append(playedRoles, &model.PlayedHistory{
			Uid:            uid,
			CountRoles:     copyCountRoles,
			LastPlayedRole: s.playedHistoryGroup.Data[uid].LastPlayedRole,
		})
	}

	assignedRoles := map[core.Uid]role2.Role{}
	rolePercents := s.getRolePercent(players)
	hasSamePercentBB, hasSamePercentSB := s.hasSamePercentBBorSB(rolePercents)

	if hasSamePercentBB && hasSamePercentSB {
		var playersLastPlayedBBAndPickNothing core.UidList
		playersLastPlayedBBorSB := lo.Keys(lo.PickBy(s.playedHistoryGroup.Data, func(_ core.Uid, played *model.PlayedHistory) bool {
			return played.LastPlayedRole == role2.BB || played.LastPlayedRole == role2.SB
		}))

		if len(playersLastPlayedBBorSB) > 0 {
			rolesToBeAssigned := lo.Shuffle(
				lo.Reject(lo.Without(roles, lo.Values(assignedRoles)...), func(r role2.Role, _ int) bool {
					return r != role2.BB && r != role2.SB
				}),
			)

			for _, uid := range playersLastPlayedBBorSB {
				// edge case：當過 BB 但只剩 BB,SB 可以選，沒得選，放到下一輪處理。
				if len(rolesToBeAssigned) == 0 {
					if r, ok := s.playedHistoryGroup.Data[uid]; ok && r.LastPlayedRole == role2.BB {
						playersLastPlayedBBAndPickNothing = append(playersLastPlayedBBAndPickNothing, uid)
					}
					continue
				}

				pickRole := rolesToBeAssigned[0]
				assignedRoles[uid] = pickRole
				s.updateHistory(uid, pickRole)
				rolesToBeAssigned = lo.Reject(rolesToBeAssigned, func(r role2.Role, _ int) bool {
					return r == pickRole
				})
			}
		}

		// 這裡開始處理 last played 是 BB or SB 但選不出結果的 edge case
		// 處理過後至少不會是 BB->BB, 頂多就是 BB->SB or SB->BB

		// 優先讓上一輪是 BB 但前面選不出結果的人先選非 BB 的角色
		if len(playersLastPlayedBBAndPickNothing) > 0 {
			rolesToBeAssigned := lo.Shuffle(
				lo.Reject(lo.Without(roles, lo.Values(assignedRoles)...), func(r role2.Role, _ int) bool {
					return r != role2.BB
				}),
			)

			for _, uid := range playersLastPlayedBBAndPickNothing {
				// 終極 edge case：當過 BB 但只剩 BB 可以選，沒得選，放到下一輪處理。
				if len(rolesToBeAssigned) == 0 {
					continue
				}

				pickRole := rolesToBeAssigned[0]
				assignedRoles[uid] = pickRole
				s.updateHistory(uid, pickRole)
				rolesToBeAssigned = lo.Reject(rolesToBeAssigned, func(r role2.Role, _ int) bool {
					return r == pickRole
				})
			}
		}

		rolesToBeAssigned := lo.Without(roles, lo.Values(assignedRoles)...)
		lo.Shuffle(rolesToBeAssigned)

		// 還沒分完的角色就讓剩下的人就隨機分配
		for _, uid := range players {
			if len(rolesToBeAssigned) == 0 {
				break
			}

			if _, ok := assignedRoles[uid]; ok {
				continue
			}

			pickRole := rolesToBeAssigned[0]
			assignedRoles[uid] = pickRole
			s.updateHistory(uid, pickRole)
			rolesToBeAssigned = lo.Reject(rolesToBeAssigned, func(r role2.Role, _ int) bool {
				return r == pickRole
			})
		}

		return assignedRoles, nil
	}

	rolesToBeAssigned := make([]role2.Role, len(roles))
	copy(rolesToBeAssigned, roles)

	// Pick BB

	bb := s.pickByLessPlayedFirst([]role2.Role{role2.BB, role2.SB}, lo.Ternary(hasSamePercentBB, role2.SB, role2.BB),
		rolePercents, playedRoles, assignedRoles)

	assignedRoles[bb] = role2.BB
	s.updateHistory(bb, assignedRoles[bb])

	// Pick rest roles randomly if rest users have same SB percent
	if hasSamePercentSB {
		remainRoles, remainUids := s.getRemains(rolesToBeAssigned, players, assignedRoles)
		lo.Shuffle(remainUids)

		for i, uid := range remainUids {
			assignedRoles[uid] = remainRoles[i]
			s.updateHistory(uid, assignedRoles[uid])
		}

		return assignedRoles, nil
	}

	// Pick SB

	sb := s.pickByLessPlayedFirst([]role2.Role{role2.BB, role2.SB}, role2.SB, rolePercents, playedRoles, assignedRoles)

	assignedRoles[sb] = role2.SB
	s.updateHistory(sb, assignedRoles[sb])

	remainRoles, remainUids := s.getRemains(rolesToBeAssigned, players, assignedRoles)
	lo.Shuffle(remainUids)

	for i, uid := range remainUids {
		assignedRoles[uid] = remainRoles[i]
		s.updateHistory(uid, assignedRoles[uid])
	}

	return assignedRoles, nil
}

func (s *RoleService) hasSamePercentBBorSB(
	rolePercents map[core.Uid]map[role2.Role]float64,
) (bool, bool) {
	hasSameBBCount := true
	hasSameSBCount := true
	var tmpBBCount, tmpSBCount float64

	var isFirstItem = true

	for uid := range rolePercents {
		if isFirstItem {
			tmpBBCount = rolePercents[uid][role2.BB]
			tmpSBCount = rolePercents[uid][role2.SB]
			isFirstItem = false
			continue
		}

		bb := rolePercents[uid][role2.BB]
		sb := rolePercents[uid][role2.SB]

		if bb != tmpBBCount {
			hasSameBBCount = false
		}

		if sb != tmpSBCount {
			hasSameSBCount = false
		}
	}

	return hasSameBBCount, hasSameSBCount
}

func (s *RoleService) getRolePercent(uids core.UidList) map[core.Uid]map[role2.Role]float64 {
	total := map[core.Uid]int{}
	for _, uid := range uids {
		if _, ok := s.playedHistoryGroup.Data[uid]; !ok {
			continue
		}

		for _, countRoles := range s.playedHistoryGroup.Data[uid].CountRoles {
			total[uid] += countRoles
		}
	}

	result := map[core.Uid]map[role2.Role]float64{}

	for _, uid := range uids {
		if _, ok := result[uid]; !ok {
			result[uid] = map[role2.Role]float64{}
		}

		for r, count := range s.playedHistoryGroup.Data[uid].CountRoles {
			result[uid][r] = getPercent(count, total[uid])
		}
	}
	return result
}

func (s *RoleService) getRemains(
	rolesToBeAssigned []role2.Role,
	players core.UidList,
	assignedRoles map[core.Uid]role2.Role,
) ([]role2.Role, core.UidList) {

	remainRoles := lo.Reject(rolesToBeAssigned, func(r role2.Role, _ int) bool {
		return lo.Contains(lo.Values(assignedRoles), r)
	})

	remainUids := lo.Reject(players, func(uid core.Uid, _ int) bool {
		_, isAssigned := assignedRoles[uid]
		return isAssigned
	})

	return remainRoles, remainUids
}

func (s *RoleService) pickByLessPlayedFirst(
	notAllowedPlayedBefore []role2.Role,
	accordingRole role2.Role,
	rolePercents map[core.Uid]map[role2.Role]float64,
	playedRoles []*model.PlayedHistory,
	assignedRoles map[core.Uid]role2.Role,
) core.Uid {

	type percent struct {
		uid     core.Uid
		percent map[role2.Role]float64
	}

	rolePercentsArr := lo.MapToSlice(rolePercents, func(k core.Uid, v map[role2.Role]float64) *percent {
		return &percent{
			uid:     k,
			percent: v,
		}
	})

	sort.Slice(rolePercentsArr, func(i, j int) bool {
		if _, ok := rolePercentsArr[i].percent[accordingRole]; !ok {
			// Put user in first place if user doesn't have played the role.
			return true
		}
		return rolePercentsArr[i].percent[accordingRole] < rolePercentsArr[j].percent[accordingRole]
	})

	chkNotAllow := lo.Associate(notAllowedPlayedBefore, func(r role2.Role) (k role2.Role, v struct{}) {
		return r, struct{}{}
	})

	lastPlayedRoles := lo.Associate(playedRoles, func(r *model.PlayedHistory) (k core.Uid, v role2.Role) {
		return r.Uid, r.LastPlayedRole
	})

	for _, v := range rolePercentsArr {
		if _, ok := assignedRoles[v.uid]; ok {
			continue
		}

		lastPlayedRole := role2.Undefined
		if r, ok := lastPlayedRoles[v.uid]; ok {
			lastPlayedRole = r
		}

		if _, ok := chkNotAllow[lastPlayedRole]; ok {
			continue
		}

		return v.uid
	}

	return lo.Reject(playedRoles, func(r *model.PlayedHistory, _ int) bool {
		_, isPicked := assignedRoles[r.Uid]
		return isPicked
	})[0].Uid
}

func (s *RoleService) updateHistory(uid core.Uid, pickedRole role2.Role) {
	s.playedHistoryGroup.Data[uid].LastPlayedRole = pickedRole
	s.playedHistoryGroup.Data[uid].CountRoles[pickedRole]++
}

func getPercent(numerator, denominator int) float64 {
	if denominator == 0 {
		return 0
	}
	percent := (float64(numerator) / float64(denominator)) * 100
	return math.Round(percent*10000) / 10000
}
