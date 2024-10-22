package api

import (
	"card-game-server-prototype/pkg/common/type/gametype"
	"card-game-server-prototype/pkg/core"
)

type UserAPI interface {
	FetchUserDetail(uid core.Uid) (*UserDetailResponse, error)
	ExchangeChip(uid core.Uid, gameType gametype.GameType, amount int) error
	GetIdleAIs() ([]string, error)
}
