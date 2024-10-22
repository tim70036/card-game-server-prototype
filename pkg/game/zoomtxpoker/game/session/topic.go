package session

import "card-game-server-prototype/pkg/core"

const (
	// GameResultTopic is used when:
	// 1. User fold.
	// 2. User request to leave room.
	// 3. Game end normally.
	GameResultTopic      core.Topic = "GameResult"
	CashOutResultTopic   core.Topic = "CashOutResult"
	CloseResultTopic     core.Topic = "CloseResult"
	PoolCloseResultTopic core.Topic = "PoolCloseResult"
)
