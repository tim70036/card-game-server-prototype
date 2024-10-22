package service

import (
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/common/type/rawevent"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/api"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	event2 "card-game-server-prototype/pkg/game/darkchess/type/event"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
)

// 生涯統計用

type BaseEventService struct {
	roomInfo               *commonmodel.RoomInfo
	roundScoreboard        *model2.RoundScoreboard
	gameScoreboard         *model2.GameScoreboard
	eventGroup             *model2.EventGroup
	gameAPI                api.GameAPI
	roundScoreboardRecords *model2.RoundScoreboardRecords
	logger                 *zap.Logger
}

func ProvideBaseEventService(
	roomInfo *commonmodel.RoomInfo,
	roundScoreboard *model2.RoundScoreboard,
	gameScoreboard *model2.GameScoreboard,
	eventGroup *model2.EventGroup,
	gameAPI api.GameAPI,
	roundScoreboardRecords *model2.RoundScoreboardRecords,
	loggerFactory *util.LoggerFactory,
) *BaseEventService {
	return &BaseEventService{
		roomInfo:               roomInfo,
		roundScoreboard:        roundScoreboard,
		gameScoreboard:         gameScoreboard,
		eventGroup:             eventGroup,
		gameAPI:                gameAPI,
		roundScoreboardRecords: roundScoreboardRecords,
		logger:                 loggerFactory.Create("BaseEventService"),
	}
}

func (s *BaseEventService) EvalRoundEvents() error {
	if len(s.roundScoreboard.Data) == 0 {
		s.logger.Warn("roundScoreboard is nil")
	}
	for uid, score := range s.roundScoreboard.Data {
		if score == nil {
			s.logger.Warn("score is nil", zap.Any("uid", uid))
			continue
		}

		events := event2.List{}

		events = append(events, &event2.Event{
			Code:   event2.Round,
			Amount: 1,
		})

		if score.RawProfit > 0 {
			events = append(events,
				&event2.Event{
					Code:   event2.RoundWin,
					Amount: 1,
					Extra:  "W",
				},
				&event2.Event{
					Code:   event2.HighestProfitInRound,
					Amount: score.RawProfit,
				},
			)

			switch score.ScoreModifier {
			case int(gamegrpc.ScoreModifierType_SCORE_MODIFIER_TYPE_1):
				events = append(events, &event2.Event{
					Code:   event2.WinRoundScoreModifier1,
					Amount: 1,
					Extra:  "",
				})
			case int(gamegrpc.ScoreModifierType_SCORE_MODIFIER_TYPE_3):
				events = append(events, &event2.Event{
					Code:   event2.WinRoundScoreModifier3,
					Amount: 1,
					Extra:  "",
				})
			case int(gamegrpc.ScoreModifierType_SCORE_MODIFIER_TYPE_5):
				events = append(events, &event2.Event{
					Code:   event2.WinRoundScoreModifier5,
					Amount: 1,
					Extra:  "",
				})
			}
		} else if score.RawProfit < 0 {
			events = append(events, &event2.Event{
				Code:   event2.RoundLose,
				Amount: 1,
				Extra:  "L",
			})
		} else if s.roundScoreboard.IsDraw {
			events = append(events, &event2.Event{
				Code:   event2.RoundDraw,
				Amount: 1,
				Extra:  "D",
			})
		}

		events = append(events, &event2.Event{
			Code:   event2.CapturePieceAmount,
			Amount: len(score.CapturePieces),
			Extra:  "",
		})

		for _, piece := range score.CapturePieces {
			if piece.IsRed() {
				if piece.IsGeneral() {
					events = append(events, &event2.Event{
						Code:   event2.CapturePieceGeneralRed,
						Amount: 1,
						Extra:  "",
					})

				}
				if piece.IsAdvisor() {
					events = append(events, &event2.Event{
						Code:   event2.CapturePieceAdvisorRed,
						Amount: 1,
						Extra:  "",
					})

				}
				if piece.IsElephant() {
					events = append(events, &event2.Event{
						Code:   event2.CapturePieceElephantRed,
						Amount: 1,
						Extra:  "",
					})

				}
				if piece.IsChariot() {
					events = append(events, &event2.Event{
						Code:   event2.CapturePieceChariotRed,
						Amount: 1,
						Extra:  "",
					})

				}
				if piece.IsHorse() {
					events = append(events, &event2.Event{
						Code:   event2.CapturePieceHorseRed,
						Amount: 1,
						Extra:  "",
					})

				}
				if piece.IsCannon() {
					events = append(events, &event2.Event{
						Code:   event2.CapturePieceCannonRed,
						Amount: 1,
						Extra:  "",
					})

				}
				if piece.IsSoldier() {
					events = append(events, &event2.Event{
						Code:   event2.CapturePieceSoldierRed,
						Amount: 1,
						Extra:  "",
					})
				}
			}

			if piece.IsBlack() {
				if piece.IsGeneral() {
					events = append(events, &event2.Event{
						Code:   event2.CapturePieceGeneralBlack,
						Amount: 1,
						Extra:  "",
					})

				}
				if piece.IsAdvisor() {
					events = append(events, &event2.Event{
						Code:   event2.CapturePieceAdvisorBlack,
						Amount: 1,
						Extra:  "",
					})

				}
				if piece.IsElephant() {
					events = append(events, &event2.Event{
						Code:   event2.CapturePieceElephantBlack,
						Amount: 1,
						Extra:  "",
					})

				}
				if piece.IsChariot() {
					events = append(events, &event2.Event{
						Code:   event2.CapturePieceChariotBlack,
						Amount: 1,
						Extra:  "",
					})

				}
				if piece.IsHorse() {
					events = append(events, &event2.Event{
						Code:   event2.CapturePieceHorseBlack,
						Amount: 1,
						Extra:  "",
					})

				}
				if piece.IsCannon() {
					events = append(events, &event2.Event{
						Code:   event2.CapturePieceCannonBlack,
						Amount: 1,
						Extra:  "",
					})

				}
				if piece.IsSoldier() {
					events = append(events, &event2.Event{
						Code:   event2.CapturePieceSoldierBlack,
						Amount: 1,
						Extra:  "",
					})
				}
			}
		}

		s.eventGroup.Data[uid] = append(s.eventGroup.Data[uid], events...)
	}
	return nil
}

func (s *BaseEventService) EvalGameEvents() error {
	for uid, score := range s.gameScoreboard.Data {
		if score == nil {
			continue
		}

		events := event2.List{}

		events = append(events, &event2.Event{
			Code:   event2.Game,
			Amount: 1,
		})

		if score.IsDisconnected {
			events = append(events, &event2.Event{
				Code:   event2.GameDisconnected,
				Amount: 1,
			})
		}

		if score.RawProfit > 0 {
			events = append(events,
				&event2.Event{
					Code:   event2.GameWin,
					Amount: 1,
					Extra:  "W",
				},
				&event2.Event{
					Code:   event2.HighestProfitInGame,
					Amount: score.RawProfit,
				},
			)
		} else if score.RawProfit < 0 {
			events = append(events, &event2.Event{
				Code:   event2.GameLose,
				Amount: 1,
				Extra:  "L",
			})
		} else if s.gameScoreboard.IsDraw {
			events = append(events, &event2.Event{
				Code:   event2.GameDraw,
				Amount: 1,
				Extra:  "D",
			})
		}

		s.eventGroup.Data[uid] = append(s.eventGroup.Data[uid], events...)
	}
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
				EventId: e.Code.ToWatchEventId(s.roomInfo.GameMode),
				Amount:  e.Amount,
				Extra:   e.Extra,
			})

			rawUserEvents = append(rawUserEvents, &rawevent.RawEvent{
				EventId: e.Code.ToUserEventId(s.roomInfo.GameMode),
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

		s.eventGroup.Data[uid] = event2.List{}
	}

	return nil
}
