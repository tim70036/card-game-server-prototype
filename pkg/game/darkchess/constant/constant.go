package constant

import (
	"time"
)

const (
	StartGameUserCount = 2
	MaxUserCount       = 2

	EndGameDisconnectedRoundCount = 1

	IdleThreshold = 6

	EventCommonOffset = 0
	EventBuddyOffset  = 20000

	// BuddyEnterLimitWeight 用來加倍好友房的 enterLimit
	// 作法：ante * BuddyEnterLimitWeight。
	// 參考 Buddy Room Encode
	BuddyEnterLimitWeight = 30
	AnteToAvoidAI         = 500
)

const (
	WaitUserEnterGracefulPeriod    = time.Second * 20
	EmptyWaitingRoomGracefulPeriod = time.Second * 60

	StartGamePeriod       = time.Second * 10
	StartRoundPeriod      = time.Millisecond * 100
	StartTurnPeriod       = time.Millisecond * 25
	PickFirstPeriod       = time.Second * 7
	AfterPickingPeriod    = time.Second * 2
	RevealPeriod          = time.Millisecond * 200
	MovePeriod            = time.Millisecond * 200
	CapturePeriod         = time.Millisecond * 200
	DrawPeriod            = time.Second
	SurrenderPeriod       = time.Second
	ShowRoundResultPeriod = time.Second * 5
)

// Default game setting
const (
	TotalRound              = 1
	TurnSecond              = 60
	ExtraTurnSecond         = 10
	AnteAmount              = 50
	IsCaptureRevealPieces   = false
	IsCaptureUnrevealPieces = false
	IsCaptureUnrevealPiece  = false
	HasRookRules            = false
	HasBishopRules          = false
	MaxRepeatMoves          = 50
	MaxChaseSamePiece       = 6
	EnterLimit              = 50
	WaterPct                = 5
)

const (
	BoardWidth            = 8 // must > 0
	BoardHeight           = 4 // must > 0
	OneColorTotalPieceCnt = 16

	TurnToAllowClaimDraw = 20
	TimeExtendCount      = 1
	ClaimDrawCount       = 2
)
