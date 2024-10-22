package api

import (
	"card-game-server-prototype/pkg/core"
)

type ClubMemberAPI interface {
	FetchDetail(clubId string, uids ...core.Uid) (*ClubMemberDetailResponse, error)
}

type ClubMemberDetailResponse struct {
	ErrCode int    `json:"errCode"`
	Msg     string `json:"msg"`
	Data    []struct {
		Uid  string `json:"uid"`
		Gold int    `json:"gold"`
		// Ticket float64 `json:"ticket"`
	} `json:"data"`
}
