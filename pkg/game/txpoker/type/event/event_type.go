package event

import (
	"card-game-server-prototype/pkg/common/type/gamemode"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
)

type EventType int

func (t EventType) String() string {
	name, ok := names[t]
	if !ok {
		return "Undefined"
	}
	return name
}

func (t EventType) ToProto() txpokergrpc.StatsEventType {
	proto, ok := protos[t]
	if !ok {
		return txpokergrpc.StatsEventType_UNDEFINED
	}
	return proto
}

func (t EventType) ToWatchEventId(mode gamemode.GameMode) int {
	return int(t) + eventIdOffset[mode]
}

func (t EventType) ToUserEventId(mode gamemode.GameMode, gameOffset int) int {
	return int(t) + eventIdOffset[mode] + gameOffset
}

const (
	Undefined       EventType = 0
	Game            EventType = 1
	BetGame         EventType = 2
	PreflopGame     EventType = 3
	FlopGame        EventType = 4
	TurnGame        EventType = 5
	RiverGame       EventType = 6
	ShowdownGame    EventType = 7
	GameWin         EventType = 8
	ShowdownGameWin EventType = 9
	GameWinAmount   EventType = 10
	UseSticker      EventType = 49

	Straight      EventType = 101
	Flush         EventType = 102
	FullHouse     EventType = 103
	FourOfAKind   EventType = 104
	StraightFlush EventType = 105
	RoyalFlush    EventType = 106

	Bet   EventType = 201
	Raise EventType = 202
	Call  EventType = 203
	Fold  EventType = 204
	AllIn EventType = 205

	PreflopFold      EventType = 301
	PreflopRaise     EventType = 302
	PreflopLastRaise EventType = 303
	PreflopReRaise   EventType = 304

	FlopFold        EventType = 401
	FlopBet         EventType = 402
	FlopLastRaise   EventType = 403
	FlopContinueBet EventType = 404

	TurnFold        EventType = 501
	TurnBet         EventType = 502
	TurnLastRaise   EventType = 503
	TurnContinueBet EventType = 504

	RiverFold        EventType = 601
	RiverBet         EventType = 602
	RiverContinueBet EventType = 603
)

var names = map[EventType]string{
	Undefined:       "Undefined",
	Game:            "完成一局",
	BetGame:         "入池一次",
	PreflopGame:     "翻牌前遊戲一次",
	FlopGame:        "翻盤圈遊戲一次",
	TurnGame:        "轉牌圈遊戲一次",
	RiverGame:       "河牌圈遊戲一次",
	ShowdownGame:    "攤牌遊戲一次",
	GameWin:         "獲勝一次",
	ShowdownGameWin: "攤牌獲勝一次",
	GameWinAmount:   "收取底池金額",
	UseSticker:      "使用貼圖一次",

	Straight:      "順子一次",
	Flush:         "同花一次",
	FullHouse:     "葫蘆一次",
	FourOfAKind:   "鐵支一次",
	StraightFlush: "同花順一次",
	RoyalFlush:    "皇家同花順一次",

	Bet:   "下注一次",
	Raise: "加注一次",
	Call:  "跟注一次",
	Fold:  "棄牌一次",
	AllIn: "AllIn一次",

	PreflopFold:      "翻牌前棄牌一次",
	PreflopRaise:     "翻牌前加注一次",
	PreflopLastRaise: "翻牌前最後加注一次",
	PreflopReRaise:   "翻盤前再加注一次",

	FlopFold:        "翻牌圈棄牌一次",
	FlopBet:         "翻牌圈下注一次",
	FlopLastRaise:   "翻牌圈最後加注一次",
	FlopContinueBet: "翻牌圈 CB 一次",

	TurnFold:        "轉牌圈棄牌一次",
	TurnBet:         "轉牌圈下注一次",
	TurnLastRaise:   "轉牌圈最後加注一次",
	TurnContinueBet: "轉牌圈 CB 一次",

	RiverFold:        "河牌圈棄牌一次",
	RiverBet:         "河牌圈下注一次",
	RiverContinueBet: "河牌圈 CB 一次",
}

var protos = map[EventType]txpokergrpc.StatsEventType{
	Undefined:       txpokergrpc.StatsEventType_UNDEFINED,
	Game:            txpokergrpc.StatsEventType_GAME,
	BetGame:         txpokergrpc.StatsEventType_BET_GAME,
	PreflopGame:     txpokergrpc.StatsEventType_PREFLOP_GAME,
	FlopGame:        txpokergrpc.StatsEventType_FLOP_GAME,
	TurnGame:        txpokergrpc.StatsEventType_TURN_GAME,
	RiverGame:       txpokergrpc.StatsEventType_RIVER_GAME,
	ShowdownGame:    txpokergrpc.StatsEventType_SHOWDOWN_GAME,
	GameWin:         txpokergrpc.StatsEventType_GAME_WIN,
	ShowdownGameWin: txpokergrpc.StatsEventType_SHOWDOWN_GAME_WIN,
	GameWinAmount:   txpokergrpc.StatsEventType_GAME_WIN_AMOUNT,
	UseSticker:      txpokergrpc.StatsEventType_UNDEFINED, // 預防萬一補 key，val前端使用 dictionary 不會處理到。

	Straight:      txpokergrpc.StatsEventType_STRAIGHT,
	Flush:         txpokergrpc.StatsEventType_FLUSH,
	FullHouse:     txpokergrpc.StatsEventType_FULL_HOUSE,
	FourOfAKind:   txpokergrpc.StatsEventType_FOUR_OF_A_KIND,
	StraightFlush: txpokergrpc.StatsEventType_STRAIGHT_FLUSH,
	RoyalFlush:    txpokergrpc.StatsEventType_ROYAL_FLUSH,

	Bet:   txpokergrpc.StatsEventType_BET,
	Raise: txpokergrpc.StatsEventType_RAISE,
	Call:  txpokergrpc.StatsEventType_CALL,
	Fold:  txpokergrpc.StatsEventType_FOLD,

	PreflopFold:      txpokergrpc.StatsEventType_PREFLOP_FOLD,
	PreflopRaise:     txpokergrpc.StatsEventType_PREFLOP_RAISE,
	PreflopLastRaise: txpokergrpc.StatsEventType_PREFLOP_LAST_RAISE,
	PreflopReRaise:   txpokergrpc.StatsEventType_PREFLOP_RE_RAISE,

	FlopFold:        txpokergrpc.StatsEventType_FLOP_FOLD,
	FlopBet:         txpokergrpc.StatsEventType_FLOP_BET,
	FlopLastRaise:   txpokergrpc.StatsEventType_FLOP_LAST_RAISE,
	FlopContinueBet: txpokergrpc.StatsEventType_FLOP_CONTINUE_BET,

	TurnFold:        txpokergrpc.StatsEventType_TURN_FOLD,
	TurnBet:         txpokergrpc.StatsEventType_TURN_BET,
	TurnLastRaise:   txpokergrpc.StatsEventType_TURN_LAST_RAISE,
	TurnContinueBet: txpokergrpc.StatsEventType_TURN_CONTINUE_BET,

	RiverFold:        txpokergrpc.StatsEventType_RIVER_FOLD,
	RiverBet:         txpokergrpc.StatsEventType_RIVER_BET,
	RiverContinueBet: txpokergrpc.StatsEventType_RIVER_CONTINUE_BET,
}

var eventIdOffset = map[gamemode.GameMode]int{
	gamemode.Common:   0,
	gamemode.Buddy:    20000,
	gamemode.Rank:     40000,
	gamemode.Carnival: 60000,
}
