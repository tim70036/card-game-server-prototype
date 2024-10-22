package api

import (
	"card-game-server-prototype/pkg/common/type/rawevent"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
)

var _ GameAPI = new(LocalGameApi)

type LocalGameApi struct {
	cfg *config.Config
}

func ProvideLocalGameAPI(cfg *config.Config) *LocalGameApi {
	return &LocalGameApi{
		cfg: cfg,
	}
}

func (api *LocalGameApi) FetchGameSetting() (*GameSetting, error) {
	return &GameSetting{
		GameMetaUid:             *api.cfg.GameMetaUid,
		TotalRound:              constant.TotalRound,
		TurnSecond:              constant.TurnSecond,
		ExtraTurnSecond:         constant.ExtraTurnSecond,
		AnteAmount:              constant.AnteAmount,
		IsCaptureRevealPieces:   constant.IsCaptureRevealPieces,
		IsCaptureUnrevealPieces: constant.IsCaptureUnrevealPieces,
		IsCaptureUnrevealPiece:  constant.IsCaptureUnrevealPiece,
		HasRookRules:            constant.HasRookRules,
		HasBishopRules:          constant.HasBishopRules,
		MaxRepeatMoves:          constant.MaxRepeatMoves,
		MaxChaseSamePiece:       constant.MaxChaseSamePiece,
		EnterLimit:              constant.EnterLimit,
		WaterPct:                constant.WaterPct,
	}, nil
}

func (api *LocalGameApi) FetchRoundResult(gameId string, round int) (*RoundResultResponse, error) {
	return nil, nil
}

func (api *LocalGameApi) FetchGameResult(gameId string) (*GameResultResponse, error) {
	return nil, nil
}

func (api *LocalGameApi) StartGame(roomId string, gameId string, uids core.UidList) error {
	return nil
}

func (api *LocalGameApi) StartRound(gameId string, round int) error {
	return nil
}

func (api *LocalGameApi) EndRound(gameId string, round int, roundScoreboard *model2.RoundScoreboard, gameInfo *model2.GameInfo) error {
	return nil
}

func (api *LocalGameApi) EndGame(gameId string, round int, gameScoreboard *model2.GameScoreboard) error {
	return nil
}

func (api *LocalGameApi) SubmitWatchEvents(uid core.Uid, rawEvents rawevent.RawEventList) error {
	return nil
}

func (api *LocalGameApi) SubmitUserEvents(uid core.Uid, rawEvents rawevent.RawEventList) error {
	return nil
}
