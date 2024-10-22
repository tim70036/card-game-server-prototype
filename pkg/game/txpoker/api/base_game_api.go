package api

import (
	"encoding/json"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/common/type/rawevent"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/type/card"
	"card-game-server-prototype/pkg/game/txpoker/type/hand"
	"strconv"

	"github.com/imroc/req/v3"
	"github.com/samber/lo"
)

type BaseGameAPI struct {
	httpClient *req.Client
	roomInfo   *commonmodel.RoomInfo
}

func ProvideBaseGameAPI(httpClient *req.Client, apiCFG *config.APIConfig, roomInfo *commonmodel.RoomInfo) *BaseGameAPI {
	return &BaseGameAPI{
		httpClient: httpClient.Clone().
			SetBaseURL("https://"+*apiCFG.MainServerHost).
			SetCommonHeader("jtoken", *apiCFG.MainServerAPIKey),
		roomInfo: roomInfo,
	}
}

func (api *BaseGameAPI) FetchGameSetting() (*GameSettingResponse, error) {
	resp := &GameSettingResponse{}
	err := api.httpClient.Get("/game/txpkr/setting").
		Do().
		Into(resp)
	return resp, err
}

func (api *BaseGameAPI) FetchGameWater() (*GetGameWaterResp, error) {
	resp := &GetGameWaterResp{}
	err := api.httpClient.Get("/game/game-setting/game-water").
		Do().
		Into(resp)
	return resp, err
}

func (api *BaseGameAPI) FetchRoom(roomId string, gameMetaUid string) (*RoomResponse, error) {
	resp := &RoomResponse{}
	err := api.httpClient.Get("/game/txpkr/room").
		AddQueryParam("roomid", roomId).
		AddQueryParam("gamemetauid", gameMetaUid).
		Do().
		Into(resp)
	return resp, err
}

func (api *BaseGameAPI) FetchGameResult(gameId string) (*GameResultResponse, error) {
	resp := &GameResultResponse{}
	err := api.httpClient.Get("/game/txpkr/game").
		AddQueryParam("gameid", gameId).
		Do().
		Into(resp)
	return resp, err
}

func (api *BaseGameAPI) UpdateRoom(roomId string, gameMetaUid string, vpip float64, emptySeatNum int) error {
	req := &startGameRequest{
		RoomId:      roomId,
		GameMetaUid: gameMetaUid,
		Vpip:        strconv.FormatFloat(vpip, 'f', 6, 64),
		EmptySeats:  strconv.Itoa(emptySeatNum),
	}

	return api.httpClient.Put("/game/txpkr/room/session").
		SetBodyJsonMarshal(req).
		Do().
		Err
}

func (api *BaseGameAPI) StartGame(gameId, roomId string, gameMetaUid string, vpip float64, emptySeatNum int) error {
	req := &startGameRequest{
		RoomId:      roomId,
		GameMetaUid: gameMetaUid,
		Vpip:        strconv.FormatFloat(vpip, 'f', 6, 64),
		EmptySeats:  strconv.Itoa(emptySeatNum),
		GameId:      gameId,
	}

	return api.httpClient.Put("/game/txpkr/room/start-session").
		SetBodyJsonMarshal(req).
		Do().
		Err
}

func (api *BaseGameAPI) TriggerJackpot(jackpotPlayers []*model2.Player, gameMetaUid string) (*TriggerJackpotResponse, error) {
	req := &triggerJackpotRequest{
		Players: lo.Map(jackpotPlayers, func(p *model2.Player, _ int) struct {
			Uid       string `json:"uid"`
			PrizeType string `json:"prizeType"`
		} {
			prizeType := lo.If(p.Hand.Type() == hand.RoyalFlush, "0").
				ElseIf(p.Hand.Type() == hand.StraightFlush, "1").
				ElseIf(p.Hand.Type() == hand.FourOfAKind, "2").
				Else("-1")

			return struct {
				Uid       string `json:"uid"`
				PrizeType string `json:"prizeType"`
			}{
				Uid:       p.Uid.String(),
				PrizeType: prizeType,
			}
		}),

		GameMetaUid: gameMetaUid,
		GameMode:    api.roomInfo.GameMode.Val(),
	}

	resp := &TriggerJackpotResponse{}
	err := api.httpClient.Post("/game/txpkr/jackpot").
		SetBodyJsonMarshal(req).
		Do().
		Into(resp)

	return resp, err
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
	scores := map[core.Uid]*PlayerScore{}
	for uid, cards := range pocketCards {
		scores[uid] = &PlayerScore{
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

	req := &endGameRequest{
		RoomId:         roomId,
		GameId:         gameId,
		CommunityCards: communityCards.ToHexStr(),
		Replay:         string(rawReplay),
		Players:        lo.Values(scores),
	}

	return api.httpClient.Post("/game/txpkr/game").
		SetBodyJsonMarshal(req).
		Do().
		Err
}

func (api *BaseGameAPI) SubmitWatchEvents(uid core.Uid, rawEvents rawevent.RawEventList) error {
	req := &eventsRequest{
		Uid:    uid.String(),
		Events: rawEvents,
	}

	return api.httpClient.Post("/game/event/events/txpoker").
		SetBodyJsonMarshal(req).
		Do().
		Err
}

func (api *BaseGameAPI) SubmitUserEvents(uid core.Uid, rawEvents rawevent.RawEventList) error {
	req := &eventsRequest{
		Uid:    uid.String(),
		Events: rawEvents,
	}

	return api.httpClient.Post("/game/event/user-event").
		SetBodyJsonMarshal(req).
		Do().
		Err
}
