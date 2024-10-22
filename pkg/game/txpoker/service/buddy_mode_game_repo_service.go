package service

import (
	"fmt"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type BuddyModeGameRepoService struct {
	*BaseGameRepoService
	logger *zap.Logger
}

func ProvideBuddyModeGameRepoService(
	baseGameRepoService *BaseGameRepoService,
	loggerFactory *util.LoggerFactory,
) *BuddyModeGameRepoService {
	return &BuddyModeGameRepoService{
		BaseGameRepoService: baseGameRepoService,
		logger:              loggerFactory.Create("BuddyGameRepoService"),
	}
}

func (s *BuddyModeGameRepoService) FetchGameInfo() error {
	rawGameMetaUid := strings.ReplaceAll(s.roomInfo.GameMetaUid, "-", "")
	lenRawGameMetaUid := len(rawGameMetaUid)

	if lenRawGameMetaUid < 28 {
		return fmt.Errorf("invalid rawGameMetaUid %s, out of range", rawGameMetaUid)
	}

	s.gameSetting.GameMetaUid = s.roomInfo.GameMetaUid
	s.gameSetting.InitialSitOutDuration = 420 * time.Second
	s.gameSetting.SitOutRefillIntervalDuration = 3600 * time.Second
	s.gameSetting.RefillSitOutDuration = 300 * time.Second
	s.gameSetting.MaxSitOutDuration = 420 * time.Second
	s.gameSetting.MaxWaterLimitBB = 100000

	// 1  char: tableSize
	if lenRawGameMetaUid >= 1 {
		tableSize, err := strconv.ParseInt(rawGameMetaUid[0:1], 16, 0)
		if err != nil {
			return fmt.Errorf("failed to parse table size from string %s: %w", rawGameMetaUid[0:1], err)
		}
		s.gameSetting.TableSize = int(tableSize)
	}

	// 2 char: leastGamePlayerAmount
	if lenRawGameMetaUid >= 2 {
		leastGamePlayerAmount, err := strconv.ParseInt(rawGameMetaUid[1:2], 16, 0)
		if err != nil {
			return fmt.Errorf("failed to parse least game player amount from string %s: %w", rawGameMetaUid[1:2], err)
		}
		s.gameSetting.LeastGamePlayerAmount = int(leastGamePlayerAmount)
	}

	// 3 ~ 5 char: minEnterLimitBB
	if lenRawGameMetaUid >= 5 {
		minEnterLimitBB, err := strconv.ParseInt(rawGameMetaUid[2:5], 16, 0)
		if err != nil {
			return fmt.Errorf("failed to parse min enter limit BB from string %s: %w", rawGameMetaUid[2:5], err)
		}
		s.gameSetting.MinEnterLimitBB = int(minEnterLimitBB)
	}

	// 6 ~ 8 char: maxEnterLimitBB
	if lenRawGameMetaUid >= 8 {
		maxEnterLimitBB, err := strconv.ParseInt(rawGameMetaUid[5:8], 16, 0)
		if err != nil {
			return fmt.Errorf("failed to parse max enter limit BB from string %s: %w", rawGameMetaUid[5:8], err)
		}
		s.gameSetting.MaxEnterLimitBB = int(maxEnterLimitBB)
	}

	// 9 ~ 12 char: smallBlind
	if lenRawGameMetaUid >= 12 {
		smallBlind, err := strconv.ParseInt(rawGameMetaUid[8:12], 16, 0)
		if err != nil {
			return fmt.Errorf("failed to parse small blind from string %s: %w", rawGameMetaUid[8:12], err)
		}
		s.gameSetting.SmallBlind = int(smallBlind)
	}

	// 13 ~ 16 char: bigBlind
	if lenRawGameMetaUid >= 16 {
		bigBlind, err := strconv.ParseInt(rawGameMetaUid[12:16], 16, 0)
		if err != nil {
			return fmt.Errorf("failed to parse big blind from string %s: %w", rawGameMetaUid[12:16], err)
		}
		s.gameSetting.BigBlind = int(bigBlind)
	}

	// 17 ~ 18 char: turnSecond
	if lenRawGameMetaUid >= 18 {
		turnSecond, err := strconv.ParseInt(rawGameMetaUid[16:18], 16, 0)
		if err != nil {
			return fmt.Errorf("failed to parse turn second from string %s: %w", rawGameMetaUid[16:18], err)
		}
		s.gameSetting.TurnDuration = time.Duration(turnSecond) * time.Second
	}

	// 19 ~ 20 char : initialExtraTurnSecond
	if lenRawGameMetaUid >= 20 {
		initialExtraTurnSecond, err := strconv.ParseInt(rawGameMetaUid[18:20], 16, 0)
		if err != nil {
			return fmt.Errorf("failed to parse initial extra turn second from string %s: %w", rawGameMetaUid[18:20], err)
		}
		s.gameSetting.InitialExtraTurnDuration = time.Duration(initialExtraTurnSecond) * time.Second
	}

	// 21 ~ 22 char: extraTurnRefillIntervalRound
	if lenRawGameMetaUid >= 22 {
		extraTurnRefillIntervalRound, err := strconv.ParseInt(rawGameMetaUid[20:22], 16, 0)
		if err != nil {
			return fmt.Errorf("failed to parse extra turn refill interval round from string %s: %w", rawGameMetaUid[20:22], err)
		}
		s.gameSetting.ExtraTurnRefillIntervalRound = int(extraTurnRefillIntervalRound)
	}

	// 23 ~ 24 char: refillExtraTurnSecond
	if lenRawGameMetaUid >= 24 {
		refillExtraTurnSecond, err := strconv.ParseInt(rawGameMetaUid[22:24], 16, 0)
		if err != nil {
			return fmt.Errorf("failed to parse refill extra turn second from string %s: %w", rawGameMetaUid[22:24], err)
		}
		s.gameSetting.RefillExtraTurnDuration = time.Duration(refillExtraTurnSecond) * time.Second
	}

	// 25 ~ 26 char: maxExtraTurnSecond
	if lenRawGameMetaUid >= 26 {
		maxExtraTurnSecond, err := strconv.ParseInt(rawGameMetaUid[24:26], 16, 0)
		if err != nil {
			return fmt.Errorf("failed to parse max extra turn second from string %s: %w", rawGameMetaUid[24:26], err)
		}
		s.gameSetting.MaxExtraTurnDuration = time.Duration(maxExtraTurnSecond) * time.Second
	}

	// 熱更抽水版本後，Common & Buddy 的 WaterPct 改由 fetchGameWater 取得。
	// 不再從 GameMetaUid 或 get gamesetting API 取得。
	if err := s.FetchGameWater(); err != nil {
		return err
	}

	return nil
}

func (s *BuddyModeGameRepoService) FetchGameWater() error {
	resp, err := s.gameAPI.FetchGameWater()
	if err != nil {
		return err
	}

	for _, v := range resp.Data {
		if v.Id == "txpkr" && s.gameSetting.WaterPct != v.Buddy {
			waterPctBefore := s.gameSetting.WaterPct
			s.gameSetting.WaterPct = v.Buddy

			s.logger.Debug("set new game water",
				zap.Int("waterPctBefore", waterPctBefore),
				zap.Int("waterPct", s.gameSetting.WaterPct),
			)
			return nil
		}
	}

	return nil
}
