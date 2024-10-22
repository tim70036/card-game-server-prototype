package api

import "time"

type RoomDetail struct {
	RoomId            string    `json:"roomId"`
	ShortRoomId       string    `json:"shortRoomId"`
	State             int       `json:"state"`
	Game              int       `json:"game"`
	GameMode          int       `json:"gameMode"`
	GameMetaUid       string    `json:"gameMetaUid"`
	PlayersAmount     int       `json:"playersAmount"`
	PlayerWhiteList   []string  `json:"playerWhiteList"`
	GameServerAddress string    `json:"gameServerAddress"`
	LastHeartBeat     time.Time `json:"lastHeartBeat"`

	IsPremium      bool   `json:"isPremium"`
	PremiumUid     string `json:"premiumUid"`
	PremiumEndTime string `json:"premiumEndTime"`
}

type GetRoomDetailResp struct {
	ErrCode int         `json:"errCode"`
	Msg     string      `json:"msg"`
	Data    *RoomDetail `json:"data"`
}

type heartbeatRequest struct {
	RoomId string `json:"roomId"`
}

type enterRoomRequest struct {
	RoomId string `json:"roomId"`
	Uid    string `json:"uid"`
}

type closeRoomRequest struct {
	RoomId string `json:"roomId"`
}

type exchangeChipRequest struct {
	Uid      string `json:"uid"`
	GameType string `json:"gameType"`
	Amount   string `json:"amount"`
}

type UserDetailResponse struct {
	Data struct {
		Uid       string `json:"uid"`
		ShortUid  int    `json:"shortUid"`
		Name      string `json:"name"`
		Cash      int    `json:"cash"`
		Level     int    `json:"level"`
		RoomCards int    `json:"roomCards"`
		IsAi      int    `json:"isAi"`
	} `json:"data"`
}
