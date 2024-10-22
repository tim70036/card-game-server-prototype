package service

import (
	"fmt"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/api"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	action2 "card-game-server-prototype/pkg/game/txpoker/type/action"
	"card-game-server-prototype/pkg/game/txpoker/type/card"
	"card-game-server-prototype/pkg/game/txpoker/type/hand"
	"card-game-server-prototype/pkg/game/txpoker/type/stage"
	"card-game-server-prototype/pkg/grpc/commongrpc"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type BaseGameRepoService struct {
	roomInfo          *commonmodel.RoomInfo
	gameInfo          *model2.GameInfo
	gameSetting       *model2.GameSetting
	playerGroup       *model2.PlayerGroup
	seatStatusGroup   *model2.SeatStatusGroup
	tableProfitsGroup *model2.TableProfitsGroup
	replay            *model2.Replay
	table             *model2.Table
	gameAPI           api.GameAPI
	msgBus            core.MsgBus
	logger            *zap.Logger
}

func ProvideBaseGameRepoService(
	roomInfo *commonmodel.RoomInfo,
	gameInfo *model2.GameInfo,
	gameSetting *model2.GameSetting,
	playerGroup *model2.PlayerGroup,
	seatStatusGroup *model2.SeatStatusGroup,
	tableProfitsGroup *model2.TableProfitsGroup,
	replay *model2.Replay,
	table *model2.Table,
	gameAPI api.GameAPI,
	msgBus core.MsgBus,
	loggerFactory *util.LoggerFactory,
) *BaseGameRepoService {
	return &BaseGameRepoService{
		roomInfo:          roomInfo,
		gameInfo:          gameInfo,
		gameSetting:       gameSetting,
		playerGroup:       playerGroup,
		seatStatusGroup:   seatStatusGroup,
		tableProfitsGroup: tableProfitsGroup,
		replay:            replay,
		table:             table,
		gameAPI:           gameAPI,
		msgBus:            msgBus,
		logger:            loggerFactory.Create("BaseGameRepoService"),
	}
}

func (s *BaseGameRepoService) UpdateRoomInfo() error {
	vpip := s.gameInfo.CalculateVPIP()
	emptySeatNum := s.gameSetting.TableSize - len(s.seatStatusGroup.TableUids)
	go s.gameAPI.UpdateRoom(s.roomInfo.RoomId, s.roomInfo.GameMetaUid, vpip, emptySeatNum)
	return nil
}

func (s *BaseGameRepoService) CreateRound() error {
	vpip := s.gameInfo.CalculateVPIP()
	emptySeatNum := s.gameSetting.TableSize - len(s.seatStatusGroup.TableUids)

	// Don't need care about the result of start game, since it's
	// only updating txpoker room status.
	go s.gameAPI.StartGame(s.gameInfo.RoundId, s.roomInfo.RoomId, s.roomInfo.GameMetaUid, vpip, emptySeatNum)
	return nil
}

func (s *BaseGameRepoService) TriggerJackpot() (map[core.Uid]int, error) {
	// Unit test: pkg/poker/txpoker/type/hand/jackpot_test.go
	jackpotResult := map[core.Uid]int{}
	jackpotPlayers := lo.Filter(lo.Values(s.playerGroup.Data), func(p *model2.Player, _ int) bool {
		if p.Hand == nil {
			return false
		}

		// 完成牌型時玩家必須用到兩張手牌，且鐵支時必須為手裡對

		if p.Hand.Type() == hand.RoyalFlush || p.Hand.Type() == hand.StraightFlush {
			return card.MatchAllPocketCards(p.Hand.Cards(), p.PocketCards)
		}

		if p.Hand.Type() == hand.FourOfAKind {
			return card.MatchAllPocketCards(p.Hand.Cards(), p.PocketCards) &&
				p.PocketCards[0].Face == p.PocketCards[1].Face
		}

		return false
	})

	if (len(jackpotPlayers)) <= 0 {
		return jackpotResult, nil
	}

	s.logger.Info("jackpot hit", zap.Array("jackpotPlayers", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		for _, p := range jackpotPlayers {
			enc.AppendObject(p)
		}
		return nil
	})))

	resp, err := s.gameAPI.TriggerJackpot(jackpotPlayers, s.roomInfo.GameMetaUid)
	if err != nil {
		return nil, fmt.Errorf("failed request trigger jackpot api: %w", err)
	}

	for _, data := range resp.Data {
		jackpotResult[core.Uid(data.Uid)] = data.Amount
	}

	s.logger.Debug("trigger jackpot response", zap.Any("resp", resp), zap.Any("jackpotResult", jackpotResult))
	return jackpotResult, nil
}

func (s *BaseGameRepoService) SubmitRoundScore() error {
	pocketCards := lo.MapValues(s.playerGroup.Data, func(p *model2.Player, _ core.Uid) card.CardList { return p.PocketCards })

	// Every player needs an entry even if he/she didn't win.
	betAmounts := lo.MapValues(s.playerGroup.Data, func(_ *model2.Player, _ core.Uid) int { return 0 })
	waters := lo.MapValues(s.playerGroup.Data, func(_ *model2.Player, _ core.Uid) int { return 0 })
	jackpotWaters := lo.MapValues(s.playerGroup.Data, func(_ *model2.Player, _ core.Uid) int { return 0 })
	profits := lo.MapValues(s.playerGroup.Data, func(_ *model2.Player, _ core.Uid) int { return 0 })
	winners := map[core.Uid]struct{}{}

	// fix: https://game-soul-technology.atlassian.net/browse/CS-40
	// Pot 有 4 人入池，CHIP 共 295136。中途兩人 fold（各 bet 500）
	// 最後兩個玩家平手，各分: 295136 / 2 = 147568

	// 玩家 betChip: 147068，winChip: 147568，被抽水 water: 7378
	// RawProfit = winChip - betChip = 147568 - 147068 = 500
	// 註：CHIP 算法在 pkg/poker/txpoker/state/evaluate_winner_state.go:225

	// Profit = RawProfit - water = 500 - 7378 = -6878
	// 這裡的 -6878 是合理的，平手的確會有機會扣掉 water 後是倒賠的。
	// bug 在 eval loser 這裡直接用 Profit <= 0 判斷是 lose
	// 在這個情況下會把 win 誤判成 lose
	// 所以 db 的 profit 記錄才會是 -betChip = -147068

	for _, uid := range s.seatStatusGroup.TableUids {
		if tableProfits, ok := s.tableProfitsGroup.Get(uid); ok {
			tableProfits.CountGames++
			s.tableProfitsGroup.Save(tableProfits)
		}
	}

	// eval winner
	for _, pot := range s.table.Pots {
		for uid, chip := range pot.Chips {
			betAmounts[uid] += chip
		}

		for uid, winner := range pot.Winners {
			winners[uid] = struct{}{}
			profits[uid] += winner.RawProfit - winner.Water - winner.JackpotWater
			waters[uid] += winner.Water
			jackpotWaters[uid] += winner.JackpotWater

			if tableProfits, ok := s.tableProfitsGroup.Get(uid); ok {
				tableProfits.SumWinLoseChips += winner.RawProfit
				s.tableProfitsGroup.Save(tableProfits)
			}
		}
	}

	// eval winner but lose game in other pots.
	// Fix Issue: https://game-soul-technology.atlassian.net/browse/GCS-4398
	for uid := range winners {
		for _, p := range s.table.Pots {
			// winner's raw profits has calculated in "evalPotWinners"
			if _, ok := p.Winners[uid]; ok {
				continue
			}

			betAmount, ok := p.Chips[uid]
			if !ok || betAmount <= 0 {
				continue
			}

			profits[uid] -= betAmount
			if tableProfits, ok := s.tableProfitsGroup.Get(uid); ok {
				tableProfits.SumWinLoseChips -= betAmount
				s.tableProfitsGroup.Save(tableProfits)
			}
		}
	}

	// eval looser
	// loser = 入池的人(in betAmounts 的人)排除 winner
	for uid, betAmount := range betAmounts {
		if _, ok := winners[uid]; ok {
			continue
		}
		profits[uid] = -betAmount

		if tableProfits, ok := s.tableProfitsGroup.Get(uid); ok {
			tableProfits.SumWinLoseChips -= betAmount
			s.tableProfitsGroup.Save(tableProfits)
		}
	}

	didEnterFlopStage := lo.MapValues(s.playerGroup.Data, func(_ *model2.Player, _ core.Uid) bool { return false })
	if s.table.BetStageFSM.MustState().(stage.Stage) >= stage.FlopStage {
		for uid := range didEnterFlopStage {
			didEnterFlopStage[uid] = !lo.ContainsBy(s.replay.ActionLog[stage.PreFlopStage], func(r action2.ActionRecord) bool {
				return r.GetUid() == uid && r.GetType() == action2.Fold
			})
		}
	}

	if err := s.gameAPI.EndGame(
		s.roomInfo.RoomId,
		s.gameInfo.RoundId,
		s.table.CommunityCards,
		s.replay,
		pocketCards,
		betAmounts,
		profits,
		waters,
		jackpotWaters,
		didEnterFlopStage,
	); err != nil {
		return fmt.Errorf("failed to call end game api: %w", err)
	}

	// After submitting round score, retreive the result. No need to
	// wait for the result, since it has nothing to do with game
	// logic. It's just for display, so we do it in background
	// (async).
	go func(roundId string) {
		// Need delay... since the result is eventual consistent.
		time.Sleep(2000 * time.Millisecond)
		result, err := s.gameAPI.FetchGameResult(roundId)
		if err != nil {
			s.logger.Error("failed to fetch game result", zap.Error(err))
			return
		}

		for _, expInfo := range result.Data {
			s.msgBus.Unicast(core.Uid(expInfo.Uid), core.MessageTopic, &txpokergrpc.Message{
				Event: &txpokergrpc.Event{
					ExpInfo: &commongrpc.ExpInfo{
						BeforeLevel:    int32(expInfo.BeforeExp),
						BeforeExp:      int32(expInfo.BeforeExp),
						LevelUpExp:     int32(expInfo.LevelUpExp),
						NextLevel:      int32(expInfo.NextLevel),
						NextLevelUpExp: int32(expInfo.NextLevelUpExp),
						IncreaseExp:    int32(expInfo.IncreaseExp),
					},
				},
			})
		}
	}(s.gameInfo.RoundId)

	return nil
}

func (s *BaseGameRepoService) FetchGameInfo() error {
	rawGameSettings, err := s.gameAPI.FetchGameSetting()
	if err != nil {
		return err
	}

	setting, ok := rawGameSettings.Data[s.roomInfo.GameMetaUid]
	if !ok {
		return fmt.Errorf("cannot find corresponding gameMetaUid: %v ", s.roomInfo.GameMetaUid)
	}

	s.gameSetting.GameMetaUid = setting.GameMetaUid
	s.gameSetting.SmallBlind = setting.SmallBlind
	s.gameSetting.BigBlind = setting.BigBlind
	s.gameSetting.TurnDuration = time.Duration(setting.TurnSecond) * time.Second
	s.gameSetting.InitialExtraTurnDuration = time.Duration(setting.InitialExtraTurnSecond) * time.Second
	s.gameSetting.ExtraTurnRefillIntervalRound = setting.ExtraTurnRefillIntervalRound
	s.gameSetting.RefillExtraTurnDuration = time.Duration(setting.RefillExtraTurnSecond) * time.Second
	s.gameSetting.MaxExtraTurnDuration = time.Duration(setting.MaxExtraTurnSecond) * time.Second
	s.gameSetting.InitialSitOutDuration = time.Duration(setting.InitialSitOutSecond) * time.Second
	s.gameSetting.SitOutRefillIntervalDuration = time.Duration(setting.SitOutRefillIntervalSecond) * time.Second
	s.gameSetting.RefillSitOutDuration = time.Duration(setting.RefillSitOutSecond) * time.Second
	s.gameSetting.MaxSitOutDuration = time.Duration(setting.MaxSitOutSecond) * time.Second
	s.gameSetting.MinEnterLimitBB = setting.MinEnterLimitBB
	s.gameSetting.MaxEnterLimitBB = setting.MaxEnterLimitBB
	s.gameSetting.TableSize = setting.TableSize

	rawRoom, err := s.gameAPI.FetchRoom(s.roomInfo.RoomId, s.roomInfo.GameMetaUid)
	if err != nil {
		return err
	}

	s.gameInfo.CreationId = rawRoom.Data.CreationId

	// 熱更抽水版本後，Common & Buddy 的 WaterPct 改由 fetchGameWater 取得。
	// 不再從 GameMetaUid 或 get gamesetting API 取得。
	if err := s.FetchGameWater(); err != nil {
		return err
	}

	return nil
}

func (s *BaseGameRepoService) FetchGameWater() error {
	resp, err := s.gameAPI.FetchGameWater()
	if err != nil {
		return err
	}

	for _, v := range resp.Data {
		if v.Id == "txpkr" && s.gameSetting.WaterPct != v.Common {
			waterPctBefore := s.gameSetting.WaterPct
			s.gameSetting.WaterPct = v.Common

			s.logger.Debug("set new game water",
				zap.Int("waterPctBefore", waterPctBefore),
				zap.Int("waterPct", s.gameSetting.WaterPct),
			)
			return err
		}
	}

	return nil
}
