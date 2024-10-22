package api

import (
	"card-game-server-prototype/pkg/common/type/rawevent"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/type/card"
)

type GameAPI interface {
	FetchGameSetting() (*GameSettingResponse, error)
	FetchGameWater() (*GetGameWaterResp, error)
	FetchRoom(roomId string, gameMetaUid string) (*RoomResponse, error)
	FetchGameResult(gameId string) (*GameResultResponse, error)
	UpdateRoom(roomId string, gameMetaUid string, vpip float64, emptySeatNum int) error
	StartGame(gameId, roomId string, gameMetaUid string, vpip float64, emptySeatNum int) error

	TriggerJackpot(jackpotPlayers []*model2.Player, gameMetaUid string) (*TriggerJackpotResponse, error)

	EndGame(
		roomId string,
		gameId string,
		communityCards card.CardList,
		replay *model2.Replay,
		pocketCards map[core.Uid]card.CardList,
		betAmount map[core.Uid]int,
		profits map[core.Uid]int,
		waters map[core.Uid]int,
		jackpotWaters map[core.Uid]int,
		didEnterFlopStage map[core.Uid]bool,
	) error

	SubmitWatchEvents(uid core.Uid, rawEvents rawevent.RawEventList) error
	SubmitUserEvents(uid core.Uid, rawEvents rawevent.RawEventList) error
}

type ClubGameAPI interface {
	FetchGameSetting() (*ClubGameSettingResponse, error)
}
