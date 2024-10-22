package event

import (
	"card-game-server-prototype/pkg/common/type/gamemode"
	"card-game-server-prototype/pkg/game/darkchess/constant"
)

type Code int

func (c Code) String() string {
	name, ok := names[c]
	if !ok {
		return "Undefined"
	}
	return name
}

func (c Code) ToWatchEventId(mode gamemode.GameMode) int {
	return int(c) + eventIdOffset[mode]
}

func (c Code) ToUserEventId(mode gamemode.GameMode) int {
	return int(c) + eventIdOffset[mode] + 700000
}

// 參考：https://docs.google.com/spreadsheets/d/1s1FELnExj13RK8gLJQrowKiWSNkjLECvb2Z4iLMKlDE
const (
	Round                     Code = 1
	RoundWin                  Code = 2
	RoundLose                 Code = 3
	RoundDraw                 Code = 4
	HighestProfitInRound      Code = 5
	Game                      Code = 6
	GameWin                   Code = 7
	GameLose                  Code = 8
	GameDraw                  Code = 9
	HighestProfitInGame       Code = 10
	GameDisconnected          Code = 11
	CapturePieceAmount        Code = 12
	WinRoundScoreModifier1    Code = 101
	WinRoundScoreModifier3    Code = 102
	WinRoundScoreModifier5    Code = 103
	CapturePieceGeneralRed    Code = 111
	CapturePieceAdvisorRed    Code = 112
	CapturePieceElephantRed   Code = 113
	CapturePieceChariotRed    Code = 114
	CapturePieceHorseRed      Code = 115
	CapturePieceCannonRed     Code = 116
	CapturePieceSoldierRed    Code = 117
	CapturePieceGeneralBlack  Code = 121
	CapturePieceAdvisorBlack  Code = 122
	CapturePieceElephantBlack Code = 123
	CapturePieceChariotBlack  Code = 124
	CapturePieceHorseBlack    Code = 125
	CapturePieceCannonBlack   Code = 126
	CapturePieceSoldierBlack  Code = 127
	UseSticker                Code = 49
)

var names = map[Code]string{
	Round:                     "參加局一次",
	RoundWin:                  "勝局一次",
	RoundLose:                 "負局一次",
	RoundDraw:                 "和局一次",
	HighestProfitInRound:      "單局最高收益",
	Game:                      "參加遊戲一次",
	GameWin:                   "勝場一次",
	GameLose:                  "負場一次",
	GameDraw:                  "和場一次",
	HighestProfitInGame:       "單場最高收益",
	GameDisconnected:          "斷線場一次",
	CapturePieceAmount:        "吃棋數量",
	WinRoundScoreModifier1:    "1倍數勝局一次",
	WinRoundScoreModifier3:    "3倍數勝局一次",
	WinRoundScoreModifier5:    "5倍數勝局一次",
	CapturePieceGeneralRed:    "吃棋-紅方帥一次",
	CapturePieceAdvisorRed:    "吃棋-紅方仕一次",
	CapturePieceElephantRed:   "吃棋-紅方相一次",
	CapturePieceChariotRed:    "吃棋-紅方俥一次",
	CapturePieceHorseRed:      "吃棋-紅方傌一次",
	CapturePieceCannonRed:     "吃棋-紅方炮一次",
	CapturePieceSoldierRed:    "吃棋-紅方兵一次",
	CapturePieceGeneralBlack:  "吃棋-黑方將一次",
	CapturePieceAdvisorBlack:  "吃棋-黑方士一次",
	CapturePieceElephantBlack: "吃棋-黑方象一次",
	CapturePieceChariotBlack:  "吃棋-黑方車一次",
	CapturePieceHorseBlack:    "吃棋-黑方馬一次",
	CapturePieceCannonBlack:   "吃棋-黑方包一次",
	CapturePieceSoldierBlack:  "吃棋-黑方卒一次",
	UseSticker:                "使用貼圖一次",
}

var eventIdOffset = map[gamemode.GameMode]int{
	gamemode.Common: constant.EventCommonOffset,
	gamemode.Buddy:  constant.EventBuddyOffset,
}
