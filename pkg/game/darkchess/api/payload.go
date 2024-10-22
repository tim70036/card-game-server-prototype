package api

import (
	"card-game-server-prototype/pkg/common/type/rawevent"
)

// Add api request response data structure here.

type GameSetting struct {
	GameMetaUid             string
	TotalRound              int
	TurnSecond              int
	ExtraTurnSecond         int
	AnteAmount              int
	IsCaptureRevealPieces   bool
	IsCaptureUnrevealPieces bool
	IsCaptureUnrevealPiece  bool
	HasRookRules            bool
	HasBishopRules          bool
	MaxRepeatMoves          int
	MaxChaseSamePiece       int
	EnterLimit              int
	WaterPct                int
}

type startGameRequest struct {
	RoomId  string   `json:"roomId"`
	GameId  string   `json:"gameId"`
	Players []string `json:"players"`
}

type endGamePlayer struct {
	Uid          string `json:"uid"`
	GameResult   int    `json:"gameResult"`
	IsBankrupt   int    `json:"isBankrupt"`
	IsDisconnect int    `json:"isDisconnect"`
}

type endGameRequest struct {
	GameId     string          `json:"gameId"`
	RoundCount string          `json:"roundCount"`
	Players    []endGamePlayer `json:"players"`
}

type fetchRoundRequest struct {
	GameId string `json:"gameId"`
	Round  string `json:"round"`
}

type startRoundRequest struct {
	GameId string `json:"gameId"`
	Round  string `json:"round"`
}

type endRoundPlayerSet struct {
	Cards     string `json:"cards"`
	Type      int    `json:"type"`
	SetNumber int    `json:"setNumber"`
}

type endRoundPlayerResult struct {
	Type  int `json:"type"`
	Point int `json:"point"`
}

type endRoundPlayer struct {
	Uid           string `json:"uid"`
	Result        int    `json:"result"`
	Profit        int    `json:"profit"`
	ScoreModifier int    `json:"scoreModifier"`
	CapturePieces string `json:"capturePieces"`
	Color         int    `json:"color"`
}

type endRoundRequest struct {
	GameId      string           `json:"gameId"`
	Round       string           `json:"round"`
	JackPotCard string           `json:"jackPotCard"`
	Players     []endRoundPlayer `json:"players"`
}

type GameResultData struct {
	Uid            string `json:"uid"`
	BeforeLevel    int    `json:"beforeLevel"`
	BeforeExp      int    `json:"beforeExp"`
	LevelUpExp     int    `json:"levelUpExp"`
	NextLevel      int    `json:"nextLevel"`
	NextLevelUpExp int    `json:"nextLevelUpExp"`
	IncreaseExp    int    `json:"increaseExp"`
}

type GameResultResponse struct {
	ErrCode int              `json:"errCode"`
	Msg     string           `json:"msg"`
	Data    []GameResultData `json:"data"`
}

type RoundResultBet struct {
	Uid    string `json:"uid"`
	Profit int    `json:"profit"`
	Water  int    `json:"water"`
}

type RoundResultData struct {
	Bets []RoundResultBet `json:"bets"`
}

type RoundResultResponse struct {
	ErrCode int             `json:"errCode"`
	Msg     string          `json:"msg"`
	Data    RoundResultData `json:"data"`
}

type eventsRequest struct {
	Uid    string                `json:"uid"`
	Events rawevent.RawEventList `json:"events"`
}

// ----- CLUB Mode -----

type ClubMemberDetail struct {
	Uid  string
	Gold int
}
