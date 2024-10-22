package api

type GameSettingResponse struct {
	Data map[string]struct {
		GameMetaUid                  string `json:"gameMetaUid"`
		Game                         int    `json:"game"`
		GameMode                     int    `json:"gameMode"`
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
		WaterPct                     int    `json:"waterPct"`
		TableSize                    int    `json:"tableSize"`
		LeastGamePlayerAmount        int    `json:"leastGamePlayerAmount"`
		MaxWaterLimitBB              int    `json:"maxWaterLimitBB"`
		MaxPlayerCount               int    `json:"maxPlayerCount"`
	} `json:"data"`
}
