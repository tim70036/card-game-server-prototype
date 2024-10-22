package api

import (
	"encoding/json"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/common/type/rawevent"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/api"
	txpokerapi "card-game-server-prototype/pkg/game/txpoker/api"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/type/card"
	"github.com/imroc/req/v3"
	"github.com/samber/lo"
)

type BaseGameAPI struct {
	*api.BaseGameAPI
	httpClient *req.Client
	roomInfo   *commonmodel.RoomInfo
}

func ProvideBaseGameAPI(
	baseGameAPI *api.BaseGameAPI,
	httpClient *req.Client,
	apiCFG *config.APIConfig,
	roomInfo *commonmodel.RoomInfo,
) *BaseGameAPI {
	return &BaseGameAPI{
		BaseGameAPI: baseGameAPI,
		httpClient: httpClient.Clone().
			SetBaseURL("https://"+*apiCFG.MainServerHost).
			SetCommonHeader("jtoken", *apiCFG.MainServerAPIKey),
		roomInfo: roomInfo,
	}
}

func (api *BaseGameAPI) UpdateRoom(roomId string, gameMetaUid string, vpip float64, emptySeatNum int) error {
	req := &struct {
		RoomId      string  `json:"roomId"`
		GameMetaUid string  `json:"gameMetaUid"`
		Vpip        float64 `json:"vpip"`
		EmptySeats  int     `json:"emptySeats"`
		GameId      string  `json:"gameId"`
	}{
		RoomId:      roomId,
		GameMetaUid: gameMetaUid,
		Vpip:        vpip,
		EmptySeats:  emptySeatNum,
	}

	return api.httpClient.Put("/game/ztxpkr/room/session").
		SetBodyJsonMarshal(req).
		Do().
		Err
}

func (api *BaseGameAPI) StartGame(gameId, roomId string, gameMetaUid string, vpip float64, emptySeatNum int) error {
	r := &struct {
		RoomId      string  `json:"roomId"`
		GameMetaUid string  `json:"gameMetaUid"`
		Vpip        float64 `json:"vpip"`
		EmptySeats  int     `json:"emptySeats"`
		GameId      string  `json:"gameId"`
	}{
		RoomId:      roomId,
		GameMetaUid: gameMetaUid,
		Vpip:        vpip,
		EmptySeats:  emptySeatNum,
		GameId:      gameId,
	}

	return api.httpClient.Put("/game/ztxpkr/room/start-session").
		SetBodyJsonMarshal(r).
		Do().
		Err
}

func (api *BaseGameAPI) EndGame(
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
) error {
	scores := map[core.Uid]*txpokerapi.PlayerScore{}
	for uid, cards := range pocketCards {
		scores[uid] = &txpokerapi.PlayerScore{
			Uid:               uid.String(),
			BetAmount:         betAmount[uid],
			Profit:            profits[uid],
			Water:             waters[uid],
			JackpotWater:      jackpotWaters[uid],
			Cards:             cards.ToHexStr(),
			DidEnterFlopStage: didEnterFlopStage[uid],
		}
	}

	rawReplay, err := json.Marshal(replay)
	if err != nil {
		return err
	}

	req := &struct {
		RoomId         string                    `json:"roomId"`
		GameId         string                    `json:"gameId"`
		CommunityCards string                    `json:"communityCards"`
		Players        []*txpokerapi.PlayerScore `json:"players"`
		Replay         string                    `json:"replay"`
	}{
		RoomId:         roomId,
		GameId:         gameId,
		CommunityCards: communityCards.ToHexStr(),
		Players:        lo.Values(scores),
		Replay:         string(rawReplay),
	}

	return api.httpClient.Post("/game/ztxpkr/game").
		SetBodyJsonMarshal(req).
		Do().
		Err
}

func (api *BaseGameAPI) FetchGameResult(gameId string) (*api.GameResultResponse, error) {
	resp := &txpokerapi.GameResultResponse{}
	err := api.httpClient.Get("/game/ztxpkr/game").
		AddQueryParam("gameId", gameId).
		Do().
		Into(resp)
	return resp, err
}

// FetchRoom only use CreationId for close room.
func (api *BaseGameAPI) FetchRoom(roomId string, gameMetaUid string) (*api.RoomResponse, error) {
	return &txpokerapi.RoomResponse{
		Data: struct {
			RoomId       string  `json:"roomId"`
			GameMetaUid  string  `json:"gameMetaUid"`
			CreationId   int     `json:"creationId"`
			Vpip         float64 `json:"vpip"`
			EmptySeats   int     `json:"emptySeats"`
			LastGameTime string  `json:"lastGameTime"`
		}{}}, nil
}

// Todo: After v1.36.0
func (api *BaseGameAPI) TriggerJackpot(jackpotPlayers []*model2.Player, gameMetaUid string) (*api.TriggerJackpotResponse, error) {
	return &txpokerapi.TriggerJackpotResponse{
		Data: []struct {
			Uid    string `json:"uid"`
			Amount int    `json:"amount"`
		}{},
	}, nil
}

func (api *BaseGameAPI) SubmitWatchEvents(uid core.Uid, rawEvents rawevent.RawEventList) error {
	req := &struct {
		Uid    string                `json:"uid"`
		Events rawevent.RawEventList `json:"events"`
	}{
		Uid:    uid.String(),
		Events: rawEvents,
	}

	return api.httpClient.Post("/game/event/events/zoom-txpoker").
		SetBodyJsonMarshal(req).
		Do().
		Err
}

func (api *BaseGameAPI) SubmitUserEvents(uid core.Uid, rawEvents rawevent.RawEventList) error {
	req := &struct {
		Uid    string                `json:"uid"`
		Events rawevent.RawEventList `json:"events"`
	}{
		Uid:    uid.String(),
		Events: rawEvents,
	}

	return api.httpClient.Post("/game/event/user-event").
		SetBodyJsonMarshal(req).
		Do().
		Err
}
