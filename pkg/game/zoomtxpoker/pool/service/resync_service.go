package service

import (
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/zoomtxpoker/pool/model"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/type/participant"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type ResyncService struct {
	userGroup        *commonmodel.UserGroup
	roomInfo         *commonmodel.RoomInfo
	gameSetting      *model2.GameSetting
	participantGroup *model2.ParticipantGroup
	playSettingGroup *model2.PlaySettingGroup
	statsGroup       *model2.StatsGroup

	msgBus core.MsgBus
	logger *zap.Logger
}

func ProvideResyncService(
	userGroup *commonmodel.UserGroup,
	roomInfo *commonmodel.RoomInfo,
	gameSetting *model2.GameSetting,
	participantGroup *model2.ParticipantGroup,
	playSettingGroup *model2.PlaySettingGroup,
	statsGroup *model2.StatsGroup,
	msgBus core.MsgBus,
	loggerFactory *util.LoggerFactory,
) *ResyncService {
	return &ResyncService{
		userGroup:        userGroup,
		roomInfo:         roomInfo,
		gameSetting:      gameSetting,
		participantGroup: participantGroup,
		playSettingGroup: playSettingGroup,
		statsGroup:       statsGroup,
		msgBus:           msgBus,
		logger:           loggerFactory.Create("ResyncService"),
	}
}

func (s *ResyncService) Send(uid core.Uid) {
	modelMsg := &txpokergrpc.Model{
		RoomInfo:    s.roomInfo.ToProto(),
		GameSetting: s.gameSetting.ToProto(),
	}

	if user, ok := s.userGroup.Data[uid]; ok {
		modelMsg.User = user.ToProto()
	}

	if stats, ok := s.statsGroup.Data[uid]; ok {
		modelMsg.Stats = stats.ToProto()
	}

	if playSetting, ok := s.playSettingGroup.Data[uid]; ok {
		modelMsg.PlaySetting = playSetting.ToProto()
	}

	if part, ok := s.participantGroup.Data[uid]; ok {
		modelMsg.Participant = part.ToProto()

		// Frontend needs empty player group to make other players disappear from table.
		// And empty action hints group to hide panel.
		if part.FSM.MustState().(participant.State) != participant.PlayingState {
			modelMsg.PlayerGroup = &txpokergrpc.PlayerGroup{Players: map[string]*txpokergrpc.Player{}}
			modelMsg.Table = &txpokergrpc.Table{
				CommunityCards:      make([]*txpokergrpc.Card, 0),
				ShowdownPocketCards: make(map[string]*txpokergrpc.CardList),
				Pots:                make([]*txpokergrpc.Pot, 0),
			}
			modelMsg.ActionHintGroup = &txpokergrpc.ActionHintGroup{Hints: map[string]*txpokergrpc.ActionHint{}}

			s.logger.Debug("Send resync 1",
				zap.String("uid", uid.String()),
				zap.Object("participantGroup", s.participantGroup),
			)

			s.msgBus.Unicast(uid, core.MessageTopic, &txpokergrpc.Message{
				Model: modelMsg,
				GameState: &txpokergrpc.GameState{
					Timestamp: timestamppb.New(time.Now()),
					Context:   &txpokergrpc.GameState_InitStateContext{InitStateContext: &txpokergrpc.InitStateContext{}},
				},
			})
			return
		}
	}

	s.logger.Debug("Send resync 2",
		zap.String("uid", uid.String()),
		zap.Object("participantGroup", s.participantGroup),
	)

	s.msgBus.Unicast(uid, core.MessageTopic, &txpokergrpc.Message{
		Model: modelMsg,
	})
}
