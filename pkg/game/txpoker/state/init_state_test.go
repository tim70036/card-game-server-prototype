package state

import (
	"card-game-server-prototype/pkg/game/txpoker/constant"
	"testing"
	"time"
)

func TestClosePeriod(t *testing.T) {
	creationIds := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	now := time.Now()

	for _, id := range creationIds {
		d := id * 10 / constant.CloseRoomAmount * constant.CloseRoomPeriod
		closeAt := now.Add(time.Minute * time.Duration(d))
		t.Logf("id: %v[%v], closeAt: %v", id, d, closeAt.Format("2006-01-02 15:04:05"))
	}
}
