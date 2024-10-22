package service

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/model"
)

type GameRepoService interface {
	UpdateRoomInfo() error
	CreateRound() error
	TriggerJackpot() (map[core.Uid]int, error)
	SubmitRoundScore() error

	FetchGameInfo() error
	FetchGameWater() error
}

type EventService interface {
	EvalRoundEvents() error
	Submit() error
}

type SeatStatusService interface {
	NewSeatStatus(uid core.Uid) *model.SeatStatus

	SitDown(uid core.Uid, seatId int) error
	BuyIn(uid core.Uid, buyInChip int) error
	StandUp(uid core.Uid) (<-chan struct{}, error)
	SitOut(uid core.Uid) error
	TopUp(uid core.Uid, topUpChip int) (<-chan struct{}, error)
	RoundStarted() error
	RoundEnd() error

	IsReadyToStartRound() bool
	StartRefillSitOutDurationLoop()

	Idle(uid core.Uid)
	Act(uid core.Uid)
}

type JackpotService interface {
	EvalWater(totalRawProfits int) int
}
