package service

import (
	"fmt"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	"strconv"
	"strings"
	"time"
)

type BuddyModeGameRepoService struct {
	*BaseGameRepoService
}

func ProvideBuddyModeGameRepoService(
	baseGameRepoService *BaseGameRepoService,
) *BuddyModeGameRepoService {
	return &BuddyModeGameRepoService{
		BaseGameRepoService: baseGameRepoService,
	}
}

func (s *BuddyModeGameRepoService) FetchGameSetting() error {
	rawGameMetaUid := strings.ReplaceAll(s.roomInfo.GameMetaUid, "-", "")
	if (len(rawGameMetaUid)) < 18 {
		return fmt.Errorf("invalid rawGameMetaUid %s, too short", rawGameMetaUid)
	}

	s.gameInfo.Setting.GameMetaUid = s.roomInfo.GameMetaUid

	totalRound, err := strconv.ParseInt(rawGameMetaUid[0:2], 16, 0)
	if err != nil {
		return fmt.Errorf("failed to parse totalRound from string %s: %w", rawGameMetaUid[4:6], err)
	}
	s.gameInfo.Setting.TotalRound = int(totalRound)

	turnSecond, err := strconv.ParseInt(rawGameMetaUid[2:4], 16, 0)
	if err != nil {
		return fmt.Errorf("failed to parse turnSecond from string %s: %w", rawGameMetaUid[4:6], err)
	}
	s.gameInfo.Setting.TurnSecond = time.Duration(turnSecond) * time.Second

	extraTurnSecond, err := strconv.ParseInt(rawGameMetaUid[4:6], 16, 0)
	if err != nil {
		return fmt.Errorf("failed to parse extraTurnSecond from string %s: %w", rawGameMetaUid[4:6], err)
	}
	s.gameInfo.Setting.ExtraTurnSecond = time.Duration(extraTurnSecond) * time.Second

	isCaptureRevealPieces, err := strconv.ParseInt(rawGameMetaUid[8:9], 16, 0)
	if err != nil {
		return fmt.Errorf("failed to parse isCaptureRevealPieces from string %s: %w", rawGameMetaUid[4:6], err)
	}
	s.gameInfo.Setting.IsCaptureRevealPieces = isCaptureRevealPieces == 1

	isCaptureUnrevealPiece, err := strconv.ParseInt(rawGameMetaUid[9:10], 16, 0)
	if err != nil {
		return fmt.Errorf("failed to parse isCaptureUnrevealPiece from string %s: %w", rawGameMetaUid[4:6], err)
	}
	s.gameInfo.Setting.IsCaptureUnrevealPiece = isCaptureUnrevealPiece == 1

	isCaptureUnrevealPieces, err := strconv.ParseInt(rawGameMetaUid[10:11], 16, 0)
	if err != nil {
		return fmt.Errorf("failed to parse isCaptureUnrevealPieces from string %s: %w", rawGameMetaUid[4:6], err)
	}
	s.gameInfo.Setting.IsCaptureUnrevealPieces = isCaptureUnrevealPieces == 1

	hasRookRules, err := strconv.ParseInt(rawGameMetaUid[11:12], 16, 0)
	if err != nil {
		return fmt.Errorf("failed to parse hasRookRules from string %s: %w", rawGameMetaUid[4:6], err)
	}
	s.gameInfo.Setting.HasRookRules = hasRookRules == 1

	hasBishopRules, err := strconv.ParseInt(rawGameMetaUid[12:13], 16, 0)
	if err != nil {
		return fmt.Errorf("failed to parse hasBishopRules from string %s: %w", rawGameMetaUid[4:6], err)
	}
	s.gameInfo.Setting.HasBishopRules = hasBishopRules == 1

	s.gameInfo.Setting.MaxRepeatMoves = constant.MaxRepeatMoves
	s.gameInfo.Setting.MaxChaseSamePiece = constant.MaxChaseSamePiece

	waterPct, err := strconv.ParseInt(rawGameMetaUid[16:18], 16, 0)
	if err != nil {
		return fmt.Errorf("failed to parse waterPct from string %s: %w", rawGameMetaUid[4:6], err)
	}
	s.gameInfo.Setting.WaterPct = int(waterPct)

	anteAmount, err := strconv.ParseInt(rawGameMetaUid[20:25], 16, 0)
	if err != nil {
		return fmt.Errorf("failed to parse anteAmount from string %s: %w", rawGameMetaUid[4:6], err)
	}
	s.gameInfo.Setting.AnteAmount = int(anteAmount)

	s.gameInfo.Setting.EnterLimit = int(anteAmount) * constant.BuddyEnterLimitWeight

	return nil
}
