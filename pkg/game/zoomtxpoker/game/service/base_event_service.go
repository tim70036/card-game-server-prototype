package service

import (
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/common/type/rawevent"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/api"
	"card-game-server-prototype/pkg/game/txpoker/model"
	txpokerservice "card-game-server-prototype/pkg/game/txpoker/service"
	"card-game-server-prototype/pkg/game/txpoker/type/event"
	"card-game-server-prototype/pkg/util"

	"go.uber.org/zap"
)

type BaseEventService struct {
	*txpokerservice.BaseEventService
	roomInfo   *commonmodel.RoomInfo
	eventGroup *model.EventGroup

	gameAPI api.GameAPI
	logger  *zap.Logger
}

func ProvideBaseEventService(
	baseEventService *txpokerservice.BaseEventService,
	roomInfo *commonmodel.RoomInfo,
	eventGroup *model.EventGroup,

	gameAPI api.GameAPI,
	loggerFactory *util.LoggerFactory,
) *BaseEventService {
	return &BaseEventService{
		BaseEventService: baseEventService,
		roomInfo:         roomInfo,
		eventGroup:       eventGroup,

		gameAPI: gameAPI,
		logger:  loggerFactory.Create("BaseEventService"),
	}
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
				EventId: e.Type.ToUserEventId(s.roomInfo.GameMode, 900000),
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

		s.eventGroup.Data[uid] = event.EventList{}
	}

	return nil
}
