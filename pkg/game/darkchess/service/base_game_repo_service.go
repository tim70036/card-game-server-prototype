package service

import (
	"fmt"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/api"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type BaseGameRepoService struct {
	roomInfo        *commonmodel.RoomInfo
	gameInfo        *model2.GameInfo
	roundScoreboard *model2.RoundScoreboard
	gameScoreboard  *model2.GameScoreboard
	playerGroup     *model2.PlayerGroup
	gameAPI         api.GameAPI
	logger          *zap.Logger
}

func ProvideBaseGameRepoService(
	roomInfo *commonmodel.RoomInfo,
	gameInfo *model2.GameInfo,
	roundScoreboard *model2.RoundScoreboard,
	gameScoreboard *model2.GameScoreboard,
	playerGroup *model2.PlayerGroup,
	gameAPI api.GameAPI,
	loggerFactory *util.LoggerFactory,
) *BaseGameRepoService {
	return &BaseGameRepoService{
		roomInfo:        roomInfo,
		gameInfo:        gameInfo,
		roundScoreboard: roundScoreboard,
		gameScoreboard:  gameScoreboard,
		playerGroup:     playerGroup,
		gameAPI:         gameAPI,
		logger:          loggerFactory.Create("BaseGameRepoService"),
	}
}

func (s *BaseGameRepoService) CreateGame() error {
	var uids core.UidList
	for _, v := range s.playerGroup.Data {
		uids = append(uids, v.Uid)
	}
	return s.gameAPI.StartGame(s.roomInfo.RoomId, s.gameInfo.GameId, uids)
}

func (s *BaseGameRepoService) CreateRound() error {
	return s.gameAPI.StartRound(
		s.gameInfo.GameId,
		s.gameInfo.RoundCount,
	)
}

func (s *BaseGameRepoService) SubmitRoundScore() error {
	return s.gameAPI.EndRound(
		s.gameInfo.GameId,
		s.gameInfo.RoundCount,
		s.roundScoreboard,
		s.gameInfo,
	)
}

func (s *BaseGameRepoService) SubmitGameScore() error {
	return s.gameAPI.EndGame(
		s.gameInfo.GameId,
		s.gameInfo.RoundCount,
		s.gameScoreboard,
	)
}

func (s *BaseGameRepoService) FetchRoundResult() error {
	resp, err := s.gameAPI.FetchRoundResult(s.gameInfo.GameId, s.gameInfo.RoundCount)
	if err != nil {
		return err
	}

	if resp == nil {
		return nil
	}

	for _, bet := range resp.Data.Bets {
		uid := core.Uid(bet.Uid)
		myRoundScore, ok := s.roundScoreboard.Data[uid]
		if !ok {
			return fmt.Errorf("can not find user[%s] in roundScoreboard", uid.String())
		}

		myRoundScore.Profit = bet.Profit
	}

	return nil
}

func (s *BaseGameRepoService) FetchGameResult() error {
	resp, err := s.gameAPI.FetchGameResult(s.gameInfo.GameId)
	if err != nil {
		return err
	}

	if resp == nil {
		return nil
	}

	for _, data := range resp.Data {
		myGameScore, ok := s.gameScoreboard.Data[core.Uid(data.Uid)]
		if !ok {
			return status.Error(codes.NotFound, "gameScoreboard not found")
		}

		myGameScore.ExpInfo.BeforeLevel = data.BeforeLevel
		myGameScore.ExpInfo.BeforeExp = data.BeforeExp
		myGameScore.ExpInfo.LevelUpExp = data.LevelUpExp
		myGameScore.ExpInfo.NextLevel = data.NextLevel
		myGameScore.ExpInfo.NextLevelUpExp = data.NextLevelUpExp
		myGameScore.ExpInfo.IncreaseExp = data.IncreaseExp
	}
	return nil
}

func (s *BaseGameRepoService) FetchGameSetting() error {
	setting, err := s.gameAPI.FetchGameSetting()
	if err != nil {
		return err
	}

	if setting != nil {
		s.gameInfo.Setting.GameMetaUid = setting.GameMetaUid
		s.gameInfo.Setting.TotalRound = setting.TotalRound
		s.gameInfo.Setting.TurnSecond = time.Duration(setting.TurnSecond) * time.Second
		s.gameInfo.Setting.ExtraTurnSecond = time.Duration(setting.ExtraTurnSecond) * time.Second
		s.gameInfo.Setting.AnteAmount = setting.AnteAmount
		s.gameInfo.Setting.IsCaptureRevealPieces = setting.IsCaptureRevealPieces
		s.gameInfo.Setting.IsCaptureUnrevealPieces = setting.IsCaptureUnrevealPieces
		s.gameInfo.Setting.IsCaptureUnrevealPiece = setting.IsCaptureUnrevealPiece
		s.gameInfo.Setting.HasRookRules = setting.HasRookRules
		s.gameInfo.Setting.HasBishopRules = setting.HasBishopRules
		s.gameInfo.Setting.MaxRepeatMoves = setting.MaxRepeatMoves
		s.gameInfo.Setting.MaxChaseSamePiece = setting.MaxChaseSamePiece
		s.gameInfo.Setting.EnterLimit = setting.EnterLimit
		s.gameInfo.Setting.WaterPct = setting.WaterPct
	}

	s.gameInfo.Setting.MaxRepeatMoves = constant.MaxRepeatMoves
	s.gameInfo.Setting.MaxChaseSamePiece = constant.MaxChaseSamePiece

	return nil
}
