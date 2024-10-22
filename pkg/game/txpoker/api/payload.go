package api

import (
	"card-game-server-prototype/pkg/common/type/rawevent"
)

type GameSettingResponse struct {
	Data map[string]struct {
		GameMetaUid                  string `json:"gameMetaUid"`
		SmallBlind                   int    `json:"smallBlind"`
		BigBlind                     int    `json:"bigBlind"`
		TurnSecond                   int    `json:"turnSecond"`
		InitialExtraTurnSecond       int    `json:"initialExtraTurnSecond"`
		ExtraTurnRefillIntervalRound int    `json:"extraTurnRefillIntervalRound"`
		RefillExtraTurnSecond        int    `json:"refillExtraTurnSecond"`
		MaxExtraTurnSecond           int    `json:"maxExtraTurnSecond"`
		InitialSitOutSecond          int    `json:"initialSitOutSecond"`
		SitOutRefillIntervalSecond   int    `json:"sitOutRefillIntervalSecond"`
		RefillSitOutSecond           int    `json:"refillSitOutSecond"`
		MaxSitOutSecond              int    `json:"maxSitOutSecond"`
		MinEnterLimitBB              int    `json:"minEnterLimitBB"`
		MaxEnterLimitBB              int    `json:"maxEnterLimitBB"`
		TableSize                    int    `json:"tableSize"`
		WaterPct                     int    `json:"waterPct"`
	} `json:"data"`
}

type GetGameWaterResp struct {
	ErrCode int          `json:"errCode"`
	Msg     string       `json:"msg"`
	Data    []*GameWater `json:"data"`
}

type GameWater struct {
	Id     string `json:"id"`
	Common int    `json:"common"`
	Buddy  int    `json:"buddy"`
}

type RoomResponse struct {
	Data struct {
		RoomId       string  `json:"roomId"`
		GameMetaUid  string  `json:"gameMetaUid"`
		CreationId   int     `json:"creationId"`
		Vpip         float64 `json:"vpip"`
		EmptySeats   int     `json:"emptySeats"`
		LastGameTime string  `json:"lastGameTime"`
	} `json:"data"`
}

type startGameRequest struct {
	RoomId      string `json:"roomId"`
	GameMetaUid string `json:"gameMetaUid"`
	Vpip        string `json:"vpip"`
	EmptySeats  string `json:"emptySeats"`
	GameId      string `json:"gameId"`
}

type endGameRequest struct {
	RoomId         string         `json:"roomId"`
	GameId         string         `json:"gameId"`
	CommunityCards string         `json:"communityCards"`
	Replay         string         `json:"replay"`
	Players        []*PlayerScore `json:"players"`
}

type PlayerScore struct {
	Uid               string `json:"uid"`
	BetAmount         int    `json:"betAmount"`
	Profit            int    `json:"profit"`
	Water             int    `json:"water"`
	JackpotWater      int    `json:"jackpotWater"`
	Cards             string `json:"cards"`
	DidEnterFlopStage bool   `json:"didEnterFlopStage"`
}

type eventsRequest struct {
	Uid    string                `json:"uid"`
	Events rawevent.RawEventList `json:"events"`
}

type GameResultResponse struct {
	Data []struct {
		Uid            string `json:"uid"`
		BeforeLevel    int    `json:"beforeLevel"`
		BeforeExp      int    `json:"beforeExp"`
		LevelUpExp     int    `json:"levelUpExp"`
		NextLevel      int    `json:"nextLevel"`
		NextLevelUpExp int    `json:"nextLevelUpExp"`
		IncreaseExp    int    `json:"increaseExp"`
	} `json:"data"`
}

type triggerJackpotRequest struct {
	Players []struct {
		Uid       string `json:"uid"`
		PrizeType string `json:"prizeType"`
	} `json:"players"`
	GameMetaUid string `json:"gameMetaUid"`
	GameMode    string `json:"gameMode"`
}

type TriggerJackpotResponse struct {
	Data []struct {
		Uid    string `json:"uid"`
		Amount int    `json:"amount"`
	} `json:"data"`
}

type exchangeChipForClubRequest struct {
	Uid         string `json:"uid"`
	GameType    string `json:"gameType"`
	Amount      string `json:"amount"`
	GameMetaUid string `json:"gameMetaUid"`
}

// ----- CLUB MODE -----

type ClubGameSettingResponse struct {
	ErrCode int    `json:"errCode"`
	Msg     string `json:"msg"`
	Data    struct {
		GameMetaUid     string `json:"gameMetaUid"`
		ClubId          string `json:"clubId"`
		GameType        int    `json:"gameType"`
		GameSettingMeta struct {
			SmallBlind                   int     `json:"smallBlind"`
			BigBlind                     int     `json:"bigBlind"`
			TurnSecond                   int     `json:"turnSecond"`
			InitialExtraTurnSecond       int     `json:"initialExtraTurnSecond"`
			ExtraTurnRefillIntervalRound int     `json:"extraTurnRefillIntervalRound"`
			RefillExtraTurnSecond        int     `json:"refillExtraTurnSecond"`
			MaxExtraTurnSecond           int     `json:"maxExtraTurnSecond"`
			InitialSitOutSecond          int     `json:"initialSitOutSecond"`
			SitOutRefillIntervalSecond   int     `json:"sitOutRefillIntervalSecond"`
			RefillSitOutSecond           int     `json:"refillSitOutSecond"`
			MaxSitOutSecond              int     `json:"maxSitOutSecond"`
			TableSize                    int     `json:"tableSize"`
			MinEnterLimitBB              int     `json:"minEnterLimitBB"`
			MaxEnterLimitBB              int     `json:"maxEnterLimitBB"`
			IsFixedRoomAmount            bool    `json:"isFixedRoomAmount"`
			RoomAmount                   int     `json:"roomAmount"`
			Duration                     int     `json:"duration"`
			LeastGamePlayerAmount        int     `json:"leastGamePlayerAmount"`
			MaxWaterLimitBB              int     `json:"maxWaterLimitBB"`
			WaterPct                     int     `json:"water"`
			TicketSpend                  float64 `json:"ticketSpend"`
			EnterLimit                   int     `json:"enterLimit"`
		} `json:"gameSettingMeta"`
	} `json:"data"`
}
