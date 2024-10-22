package service

import (
	"fmt"
	"card-game-server-prototype/pkg/core"
	role2 "card-game-server-prototype/pkg/game/txpoker/type/role"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/model"
	"github.com/samber/lo"
	"strings"
	"testing"
)

func TestAssignRoles(t *testing.T) {
	type args struct {
		matchedUids core.UidList
		inst        *RoleService
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "111",
			args: args{
				inst: ProvideRoleService(&model.PlayedHistoryGroup{
					Data: make(map[core.Uid]*model.PlayedHistory),
				}),
				matchedUids: core.UidList{"A", "B", "C", "D", "E", "F", "G"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.args.matchedUids {
				if _, ok := tt.args.inst.playedHistoryGroup.Data[v]; !ok {
					tt.args.inst.playedHistoryGroup.Data[v] = &model.PlayedHistory{
						Uid:            v,
						LastPlayedRole: role2.Undefined,
						CountRoles:     map[role2.Role]int{},
					}
				}
			}

			playTimes := 100000000
			countUnStoppableBBorSB := 0
			previousRoles := make(map[core.Uid]role2.Role)

			for i := 0; i < playTimes; i++ {
				matchedUids := lo.Shuffle(tt.args.matchedUids)[:6]
				roles, err := tt.args.inst.Assign(matchedUids)
				if err != nil {
					t.Fatal(err)
				}

				chkRoles, _ := role2.GetRoles(6)
				var count int
				for _, v2 := range roles {
					for _, v := range chkRoles {
						if v.String() == v2.String() {
							count++
							break
						}
					}
				}

				if count != 6 {
					t.Fatalf("roles: %v", roles)
				}

				for uid, r := range roles {
					if prev, ok := previousRoles[uid]; ok {
						if !(prev == role2.BB || prev == role2.SB || r == role2.BB || r == role2.SB) {
							continue
						}

						if prev == r {
							// for u, v := range roles {
							// 	if u == uid {
							// 		continue
							// 	}
							// 	t.Logf("i:%v, uid: %v, prev: %v, r: %v", i, u, previousRoles[u], v)
							// }
							// t.Logf("i:%v, uid: %v, prev: %v, r: %v", i, uid, prev, r)
							countUnStoppableBBorSB++
						}

						if (prev == role2.BB && r == role2.SB) || (prev == role2.SB && r == role2.BB) {
							// for u, v := range roles {
							// 	if u == uid {
							// 		continue
							// 	}
							// 	t.Logf("i:%v, uid: %v, prev: %v, r: %v", i, u, previousRoles[u], v)
							// }
							// t.Logf("i:%v, uid: %v, prev: %v, r: %v", i, uid, prev, r)
							countUnStoppableBBorSB++
						}
					}
				}

				for uid, r := range roles {
					previousRoles[uid] = r
				}
			}

			var sumCount int
			for uid := range tt.args.inst.playedHistoryGroup.Data {
				var s strings.Builder
				sortRoles, _ := role2.GetRoles(6)
				for i, r := range sortRoles {
					if i > 0 {
						s.WriteString(", ")
					}
					if count, ok := tt.args.inst.playedHistoryGroup.Data[uid].CountRoles[r]; ok {
						s.WriteString(fmt.Sprintf("%s: %v", r.String(), count))
						sumCount += count
					} else {
						s.WriteString(fmt.Sprintf("%s: 0", r.String()))
					}
				}

				t.Log(uid, "-", s.String())
				s.Reset()
			}
			t.Log("--------------------")

			rolePercents := tt.args.inst.getRolePercent(tt.args.matchedUids)
			for u, v := range rolePercents {
				fmt.Printf("User: %v\n", u.String())
				for r, percent := range v {
					fmt.Printf(" role: %v (%v%%)\n", r, percent)
				}
			}
			fmt.Println("-------------")

			t.Logf("每個 user 選出 role 算一次，共 %v 次\n 選不出的情況（配對的 6 人 每個人前一輪是 BB 或 SB）有 %v 次，約 %v%%\n",
				sumCount,
				countUnStoppableBBorSB,
				getPercent(countUnStoppableBBorSB, sumCount))
		})
	}
}
