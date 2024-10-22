package api

import (
	"fmt"
	"card-game-server-prototype/pkg/common/type/rawevent"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	"card-game-server-prototype/pkg/util"
	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"strconv"
)

var _ GameAPI = new(BaseGameApi)

type BaseGameApi struct {
	httpClient *req.Client
	cfg        *config.Config
}

func ProvideBaseGameAPI(httpClient *req.Client, apiCFG *config.APIConfig, cfg *config.Config) *BaseGameApi {
	return &BaseGameApi{
		httpClient: httpClient.Clone().
			SetBaseURL("https://"+*apiCFG.MainServerHost).
			SetCommonHeader("jtoken", *apiCFG.MainServerAPIKey),
		cfg: cfg,
	}
}

// FetchGameSetting get All GameSettingType Setting
func (api *BaseGameApi) FetchGameSetting() (*GameSetting, error) {
	type data map[string]struct {
		GameMetaUid             string `json:"gameMetaUid"`
		TotalRound              int    `json:"totalRound"`
		TurnSecond              int    `json:"turnSecond"`
		ExtraTurnSecond         int    `json:"extraTurnSecond"`
		AnteAmount              int    `json:"anteAmount"`
		IsCaptureRevealPieces   bool   `json:"isCaptureRevealPieces"`
		IsCaptureUnrevealPieces bool   `json:"isCaptureUnrevealPieces"`
		IsCaptureUnrevealPiece  bool   `json:"isCaptureUnrevealPiece"`
		HasRookRules            bool   `json:"hasRookRules"`
		HasBishopRules          bool   `json:"hasBishopRules"`
		MaxRepeatMoves          int    `json:"maxRepeatMoves"`
		MaxChaseSamePiece       int    `json:"maxChaseSamePiece"`
		EnterLimit              int    `json:"enterLimit"`
		WaterPct                int    `json:"waterPct"`
	}

	type resp struct {
		ErrCode int    `json:"errCode"`
		Msg     string `json:"msg"`
		Data    data   `json:"data"`
	}

	r := &resp{}
	request := api.httpClient.Get("/game/dc/setting")
	if err := request.Do().Into(r); err != nil {
		return nil, err
	}

	if r == nil || (r != nil && r.Data == nil) {
		return nil, fmt.Errorf("null base game setting, gameMetaUid: %v ", *api.cfg.GameMetaUid)
	}

	v, ok := r.Data[*api.cfg.GameMetaUid]
	if !ok {
		return nil, fmt.Errorf("null base game setting, gameMetaUid: %v ", *api.cfg.GameMetaUid)
	}

	return &GameSetting{
		GameMetaUid:             v.GameMetaUid,
		TotalRound:              v.TotalRound,
		TurnSecond:              v.TurnSecond,
		ExtraTurnSecond:         v.ExtraTurnSecond,
		AnteAmount:              v.AnteAmount,
		IsCaptureRevealPieces:   v.IsCaptureRevealPieces,
		IsCaptureUnrevealPieces: v.IsCaptureUnrevealPieces,
		IsCaptureUnrevealPiece:  v.IsCaptureUnrevealPiece,
		HasRookRules:            v.HasRookRules,
		HasBishopRules:          v.HasBishopRules,
		MaxRepeatMoves:          v.MaxRepeatMoves,
		MaxChaseSamePiece:       v.MaxChaseSamePiece,
		EnterLimit:              v.EnterLimit,
		WaterPct:                v.WaterPct,
	}, nil
}

func (api *BaseGameApi) StartGame(roomId string, gameId string, uids core.UidList) error {
	var uidStrings = make([]string, len(uids))
	for i, uid := range uids {
		uidStrings[i] = uid.String()
	}

	r := &startGameRequest{
		RoomId:  roomId,
		GameId:  gameId,
		Players: uidStrings,
	}

	return api.httpClient.Post("/game/dc/game").
		SetBodyJsonMarshal(r).
		Do().
		Err
}

func (api *BaseGameApi) EndGame(gameId string, round int, gameScoreboard *model2.GameScoreboard) error {
	var players = make([]endGamePlayer, 0)

	for _, score := range gameScoreboard.Data {
		players = append(players, endGamePlayer{
			Uid:          score.Uid.String(),
			IsBankrupt:   lo.Ternary(score.IsBankrupt, 1, 0),
			GameResult:   api.getResult(score.RawProfit),
			IsDisconnect: lo.Ternary(score.IsDisconnected, 1, 0),
		})
	}

	r := &endGameRequest{
		GameId:     gameId,
		RoundCount: strconv.Itoa(round),
		Players:    players,
	}

	return api.httpClient.Put("/game/dc/game").
		SetBodyJsonMarshal(r).
		Do().
		Err
}

func (api *BaseGameApi) FetchGameResult(gameId string) (*GameResultResponse, error) {
	resp := &GameResultResponse{}
	err := api.httpClient.Get("/game/dc/game").
		AddQueryParam("gameid", gameId).
		Do().
		Into(resp)
	return resp, err
}

func (api *BaseGameApi) StartRound(gameId string, round int) error {
	r := &startRoundRequest{
		GameId: gameId,
		Round:  strconv.Itoa(round),
	}

	return api.httpClient.Post("/game/dc/round").
		SetBodyJsonMarshal(r).
		Do().
		Err
}

func (api *BaseGameApi) EndRound(
	gameId string,
	round int,
	roundScoreboard *model2.RoundScoreboard,
	gameInfo *model2.GameInfo,
) error {

	var tmpPlayers = map[core.Uid]endRoundPlayer{}

	for uid, score := range roundScoreboard.Data {
		capturePieces := util.JoinStrings(lo.Map(score.CapturePieces, func(p piece.Piece, _ int) string {
			return p.GetHex()
		}))

		tmpPlayers[uid] = endRoundPlayer{
			Uid:           score.Uid.String(),
			Result:        api.getResult(score.RawProfit),
			Profit:        score.RawProfit,
			ScoreModifier: score.ScoreModifier,
			CapturePieces: capturePieces,
			Color:         int(score.Color),
		}
	}

	var players []endRoundPlayer
	for _, v := range tmpPlayers {
		players = append(players, v)
	}

	r := &endRoundRequest{
		GameId:  gameId,
		Round:   strconv.Itoa(round),
		Players: players,
	}

	return api.httpClient.Put("/game/dc/round").
		SetBodyJsonMarshal(r).
		Do().
		Err
}

func (api *BaseGameApi) FetchRoundResult(gameId string, round int) (*RoundResultResponse, error) {
	resp := &RoundResultResponse{
		Data: RoundResultData{Bets: []RoundResultBet{}},
	}
	err := api.httpClient.Get("/game/dc/round").
		AddQueryParam("gameid", gameId).
		AddQueryParam("round", strconv.Itoa(round)).
		Do().
		Into(resp)
	return resp, err
}

func (api *BaseGameApi) SubmitWatchEvents(uid core.Uid, rawEvents rawevent.RawEventList) error {
	r := &eventsRequest{
		Uid:    uid.String(),
		Events: rawEvents,
	}

	return api.httpClient.Post("/game/event/events/darkchess").
		SetBodyJsonMarshal(r).
		Do().
		Err
}

func (api *BaseGameApi) SubmitUserEvents(uid core.Uid, rawEvents rawevent.RawEventList) error {
	r := &eventsRequest{
		Uid:    uid.String(),
		Events: rawEvents,
	}

	return api.httpClient.Post("/game/event/user-event").
		SetBodyJsonMarshal(r).
		Do().
		Err
}

func (api *BaseGameApi) getResult(rawProfit int) int {
	// win: 2, draw: 1, lose: 0
	return lo.
		If(rawProfit > 0, 2).
		ElseIf(rawProfit == 0, 1).
		Else(0)
}
