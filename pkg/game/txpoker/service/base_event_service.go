package service

import (
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/common/type/rawevent"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/api"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	action2 "card-game-server-prototype/pkg/game/txpoker/type/action"
	event2 "card-game-server-prototype/pkg/game/txpoker/type/event"
	"card-game-server-prototype/pkg/game/txpoker/type/stage"
	"card-game-server-prototype/pkg/util"

	"github.com/samber/lo"
	"go.uber.org/zap"
)

type BaseEventService struct {
	roomInfo    *commonmodel.RoomInfo
	gameSetting *model2.GameSetting
	eventGroup  *model2.EventGroup
	playerGroup *model2.PlayerGroup
	replay      *model2.Replay
	table       *model2.Table

	gameAPI api.GameAPI
	logger  *zap.Logger
}

func ProvideBaseEventService(
	roomInfo *commonmodel.RoomInfo,
	gameSetting *model2.GameSetting,
	eventGroup *model2.EventGroup,
	playerGroup *model2.PlayerGroup,
	replay *model2.Replay,
	table *model2.Table,

	gameAPI api.GameAPI,
	loggerFactory *util.LoggerFactory,
) *BaseEventService {
	return &BaseEventService{
		roomInfo:    roomInfo,
		gameSetting: gameSetting,
		eventGroup:  eventGroup,
		playerGroup: playerGroup,
		replay:      replay,
		table:       table,

		gameAPI: gameAPI,
		logger:  loggerFactory.Create("BaseEventService"),
	}
}

func (s *BaseEventService) EvalRoundEvents() error {
	allActionLog := lo.Flatten(lo.Values(s.replay.ActionLog))
	playerEvents := map[core.Uid]event2.EventList{}

	for uid, player := range s.playerGroup.Data {
		events := event2.EventList{}

		events = append(events, &event2.Event{
			Type:   event2.Game,
			Amount: 1,
		})

		if lo.ContainsBy(allActionLog, func(r action2.ActionRecord) bool {
			return (r.GetUid() == uid && lo.Contains([]action2.ActionType{action2.Bet, action2.Call, action2.Raise, action2.AllIn}, r.GetType())) ||
				(r.GetUid() == uid && r.GetType() == action2.BB && s.table.BetStageFSM.MustState().(stage.Stage) >= stage.FlopStage &&
					lo.Contains([]action2.ActionType{action2.Bet, action2.Call, action2.Raise, action2.AllIn}, r.GetType()))
			// Edge case: is BB and the game has gone into flop stage and do bet any chip. In this case, the BB is counted as betting.
		}) {
			events = append(events, &event2.Event{
				Type:   event2.BetGame,
				Amount: 1,
			})
		}

		winPotRecords := lo.Filter(allActionLog, func(r action2.ActionRecord, _ int) bool {
			return r.GetType() == action2.WinPot && r.GetUid() == uid
		})

		if len(winPotRecords) > 0 {
			winChipSum := lo.SumBy(winPotRecords, func(r action2.ActionRecord) int { return r.(*action2.WinPotRecord).Chip })
			events = append(events, &event2.Event{
				Type:   event2.GameWinAmount,
				Amount: winChipSum,
			})

			// 勝率的定義 edge case: 如果是盲注一直贏 (沒有進到flop)，
			// 那總投注局數不會增加, 收取底池次數會增加，造成勝率超過
			// 100%。讓大盲吃小盲注的時候不計算，這樣可以閃掉這個問
			// 題。
			if winChipSum >= s.gameSetting.BigBlind*2 {
				events = append(events, &event2.Event{
					Type:   event2.GameWin,
					Amount: 1,
				})
			}

			// Has showdown.
			if player.Hand != nil {
				events = append(events, &event2.Event{
					Type:   event2.ShowdownGameWin,
					Amount: 1,
				})
			}
		}

		// Has showdown, record hand type.
		if player.Hand != nil {
			if handEventType, ok := event2.HandEventTypeMap[player.Hand.Type()]; ok { // Not all hand type will be recorded.
				events = append(events, &event2.Event{
					Type:   handEventType,
					Amount: 1,
				})
			}
		}

		// Has played in stage.
		for curStage := stage.AnteStage; curStage <= s.table.BetStageFSM.MustState().(stage.Stage); curStage++ {
			if gameEventType, ok := event2.StageGameEventTypeMap[curStage]; ok {
				events = append(events, &event2.Event{
					Type:   gameEventType,
					Amount: 1,
				})
			}

			if lo.ContainsBy(s.replay.ActionLog[curStage], func(r action2.ActionRecord) bool {
				return r.GetUid() == uid && r.GetType() == action2.Fold
			}) {
				break
			}
		}

		if betCount := lo.CountBy(allActionLog, func(r action2.ActionRecord) bool {
			return r.GetUid() == uid &&
				(r.GetType() == action2.Bet || (r.GetType() == action2.AllIn && r.(*action2.AllInRecord).BetType == action2.Bet))
		}); betCount > 0 {
			events = append(events, &event2.Event{
				Type:   event2.Bet,
				Amount: betCount,
			})
		}

		if raiseCount := lo.CountBy(allActionLog, func(r action2.ActionRecord) bool {
			return r.GetUid() == uid &&
				(r.GetType() == action2.Raise || (r.GetType() == action2.AllIn && r.(*action2.AllInRecord).BetType == action2.Raise))
		}); raiseCount > 0 {
			events = append(events, &event2.Event{
				Type:   event2.Raise,
				Amount: raiseCount,
			})
		}

		if callCount := lo.CountBy(allActionLog, func(r action2.ActionRecord) bool {
			return r.GetUid() == uid &&
				(r.GetType() == action2.Call || (r.GetType() == action2.AllIn && r.(*action2.AllInRecord).BetType == action2.Call))
		}); callCount > 0 {
			events = append(events, &event2.Event{
				Type:   event2.Call,
				Amount: callCount,
			})
		}

		if foldCount := lo.CountBy(allActionLog, func(r action2.ActionRecord) bool {
			return r.GetUid() == uid && r.GetType() == action2.Fold
		}); foldCount > 0 {
			events = append(events, &event2.Event{
				Type:   event2.Fold,
				Amount: foldCount,
			})
		}

		if count := lo.CountBy(allActionLog, func(r action2.ActionRecord) bool {
			return r.GetUid() == uid && r.GetType() == action2.AllIn
		}); count > 0 {
			events = append(events, &event2.Event{
				Type:   event2.AllIn,
				Amount: count,
			})
		}

		playerEvents[uid] = append(playerEvents[uid], events...)
	}

	for curStage, stageLog := range s.replay.ActionLog {
		// Has fold in stage.
		if foldEventType, ok := event2.StageFoldEventTypeMap[curStage]; ok {
			foldUids := lo.Uniq(lo.FilterMap(stageLog, func(r action2.ActionRecord, _ int) (core.Uid, bool) {
				return r.GetUid(), r.GetType() == action2.Fold
			}))

			for _, foldUid := range foldUids {
				playerEvents[foldUid] = append(playerEvents[foldUid], &event2.Event{
					Type:   foldEventType,
					Amount: 1,
				})
			}
		}

		// Has bet in stage.
		if betEventType, ok := event2.StageBetEventTypeMap[curStage]; ok {
			betUids := lo.Uniq(lo.FilterMap(stageLog, func(r action2.ActionRecord, _ int) (core.Uid, bool) {
				return r.GetUid(), r.GetType() == action2.Bet || (r.GetType() == action2.AllIn && r.(*action2.AllInRecord).BetType == action2.Bet)
			}))

			for _, betUid := range betUids {
				playerEvents[betUid] = append(playerEvents[betUid], &event2.Event{
					Type:   betEventType,
					Amount: 1,
				})
			}
		}

		// Has raise in stage.
		if raiseEventType, ok := event2.StageRaiseEventTypeMap[curStage]; ok {
			raiseUids := lo.Uniq(lo.FilterMap(stageLog, func(r action2.ActionRecord, _ int) (core.Uid, bool) {
				return r.GetUid(), r.GetType() == action2.Raise || (r.GetType() == action2.AllIn && r.(*action2.AllInRecord).BetType == action2.Raise)
			}))

			for _, raiseUid := range raiseUids {
				playerEvents[raiseUid] = append(playerEvents[raiseUid], &event2.Event{
					Type:   raiseEventType,
					Amount: 1,
				})
			}
		}

		// The last raise in stage.
		if lastRaiseEventType, ok := event2.StageLastRaiseEventTypeMap[curStage]; ok {
			stageLastBetRecord, _, hasLastBet := lo.FindLastIndexOf(stageLog, func(r action2.ActionRecord) bool {
				return r.GetType() == action2.Bet ||
					r.GetType() == action2.Raise ||
					(r.GetType() == action2.AllIn && lo.Contains([]action2.ActionType{action2.Bet, action2.Raise}, r.(*action2.AllInRecord).BetType))
			})

			if hasLastBet {
				playerEvents[stageLastBetRecord.GetUid()] = append(playerEvents[stageLastBetRecord.GetUid()], &event2.Event{
					Type:   lastRaiseEventType,
					Amount: 1,
				})
			}
		}

		// Has re-raised in stage.
		if reRaiseEventType, ok := event2.StageReRaiseEventTypeMap[curStage]; ok {
			raiseOrderUids := lo.FilterMap(stageLog, func(r action2.ActionRecord, _ int) (core.Uid, bool) {
				return r.GetUid(), r.GetType() == action2.Raise || (r.GetType() == action2.AllIn && r.(*action2.AllInRecord).BetType == action2.Raise)
			})

			if len(raiseOrderUids) > 1 {
				reRaiseUids := lo.Uniq(raiseOrderUids[1:])
				for _, reRaiseUid := range reRaiseUids {
					playerEvents[reRaiseUid] = append(playerEvents[reRaiseUid], &event2.Event{
						Type:   reRaiseEventType,
						Amount: 1,
					})
				}
			}
		}

		// Has continued bet in stage.
		if continueBetEventType, ok := event2.StageContinueBetEventTypeMap[curStage]; ok {
			lastStage := curStage - 1
			lastStageLog := s.replay.ActionLog[lastStage]
			lastStageLastBetRecord, _, hasLastBet := lo.FindLastIndexOf(lastStageLog, func(r action2.ActionRecord) bool {
				return r.GetType() == action2.Bet ||
					r.GetType() == action2.Raise ||
					(r.GetType() == action2.AllIn && lo.Contains([]action2.ActionType{action2.Bet, action2.Raise}, r.(*action2.AllInRecord).BetType))
			})

			if hasLastBet && lo.ContainsBy(stageLog, func(r action2.ActionRecord) bool {
				return r.GetUid() == lastStageLastBetRecord.GetUid() &&
					(r.GetType() == action2.Bet || (r.GetType() == action2.AllIn && r.(*action2.AllInRecord).BetType == action2.Bet))
			}) {
				playerEvents[lastStageLastBetRecord.GetUid()] = append(playerEvents[lastStageLastBetRecord.GetUid()], &event2.Event{
					Type:   continueBetEventType,
					Amount: 1,
				})
			}

		}
	}

	for uid, events := range playerEvents {
		// Note that some user might already leave, we must skip for them.
		if _, ok := s.eventGroup.Data[uid]; !ok {
			continue
		}

		// 不可直接覆蓋 event，可能已經存有 sticker 的 event。
		s.eventGroup.Data[uid] = append(s.eventGroup.Data[uid], events...)
	}

	s.logger.Debug("round events evaluated", zap.Object("events", s.eventGroup))
	return nil
}

func (s *BaseEventService) Submit() error {
	for uid, events := range s.eventGroup.Data {
		if len(events) == 0 {
			continue
		}

		var rawWatchEvents, rawUserEvents []*rawevent.RawEvent

		for _, e := range events {
			rawWatchEvents = append(rawWatchEvents, &rawevent.RawEvent{
				EventId: e.Type.ToWatchEventId(s.roomInfo.GameMode),
				Amount:  e.Amount,
				Extra:   e.Extra,
			})

			rawUserEvents = append(rawUserEvents, &rawevent.RawEvent{
				EventId: e.Type.ToUserEventId(s.roomInfo.GameMode, 600000),
				Amount:  e.Amount,
				Extra:   e.Extra,
			})
		}

		// Don't care result, since game event has nothing to do with game
		// logic.
		go func(u core.Uid, we, ue rawevent.RawEventList) {
			_ = s.gameAPI.SubmitWatchEvents(u, we)
			_ = s.gameAPI.SubmitUserEvents(u, ue)
		}(uid, rawWatchEvents, rawUserEvents)

		s.eventGroup.Data[uid] = event2.EventList{}
	}

	return nil
}
