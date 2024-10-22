package api

import (
	"card-game-server-prototype/pkg/common/type/rawevent"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
)

// 取得遊戲的資料 or 發送需記錄的資料，例如：main server。

type GameAPI interface {
	FetchGameSetting() (*GameSetting, error)
	FetchRoundResult(gameId string, round int) (*RoundResultResponse, error)
	FetchGameResult(gameId string) (*GameResultResponse, error)
	StartGame(roomId, gameId string, uids core.UidList) error
	StartRound(gameId string, round int) error
	EndRound(gameId string, round int, roundScoreboard *model2.RoundScoreboard, gameInfo *model2.GameInfo) error
	EndGame(gameId string, round int, gameScoreboard *model2.GameScoreboard) error
	SubmitWatchEvents(uid core.Uid, rawEvents rawevent.RawEventList) error
	SubmitUserEvents(uid core.Uid, rawEvents rawevent.RawEventList) error
}

type ClubMemberAPI interface {
	FetchDetail(clubId string, uids ...core.Uid) ([]*ClubMemberDetail, error)
}
