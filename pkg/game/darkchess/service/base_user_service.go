package service

import (
	commonmodel "card-game-server-prototype/pkg/common/model"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	model2 "card-game-server-prototype/pkg/game/darkchess/model"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"github.com/samber/lo"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type BaseUserService struct {
	*commonservice.BaseUserService

	msgBus           core.MsgBus
	testCFG          *config.TestConfig
	roomInfo         *commonmodel.RoomInfo
	userGroup        *commonmodel.UserGroup
	playSettingGroup *model2.PlaySettingGroup
	actionHintGroup  *model2.ActionHintGroup
	eventGroup       *model2.EventGroup
}

func ProvideBaseUserService(
	baseUserService *commonservice.BaseUserService,

	msgBus core.MsgBus,
	testCFG *config.TestConfig,
	roomInfo *commonmodel.RoomInfo,
	userGroup *commonmodel.UserGroup,
	playSettingGroup *model2.PlaySettingGroup,
	actionHintGroup *model2.ActionHintGroup,
	eventGroup *model2.EventGroup,
) *BaseUserService {
	return &BaseUserService{
		BaseUserService:  baseUserService,
		msgBus:           msgBus,
		testCFG:          testCFG,
		roomInfo:         roomInfo,
		userGroup:        userGroup,
		playSettingGroup: playSettingGroup,
		actionHintGroup:  actionHintGroup,
		eventGroup:       eventGroup,
	}
}

func (s *BaseUserService) Init(uids ...core.Uid) error {
	if err := s.BaseUserService.Init(uids...); err != nil {
		return err
	}

	for _, uid := range uids {
		s.playSettingGroup.Data[uid] = &model2.PlaySetting{
			Uid:    uid,
			IsAuto: false,
		}
	}

	return nil
}

func (s *BaseUserService) Destroy(topic core.Topic, kickoutMSG proto.Message, uids ...core.Uid) error {
	if err := s.BaseUserService.Destroy(topic, kickoutMSG, uids...); err != nil {
		return err
	}

	for _, uid := range uids {
		delete(s.playSettingGroup.Data, uid)
		delete(s.eventGroup.Data, uid)
		delete(s.actionHintGroup.Data, uid)
	}

	return nil
}

func (s *BaseUserService) CanConnect(uid core.Uid) error {
	if err := s.BaseUserService.CanConnect(uid); err != nil {
		return err
	}

	if _, ok := s.userGroup.Data[uid]; !ok && len(s.userGroup.Data) >= constant.MaxUserCount {
		return status.Errorf(codes.FailedPrecondition, "exceed max user count %v", constant.MaxUserCount)
	}

	if s.testCFG.EnableCheatMode(string(s.roomInfo.GameType)) || *s.testCFG.LocalMode {
		return nil
	}

	if !lo.Contains(s.roomInfo.ValidUsers, uid) {
		return status.Errorf(codes.PermissionDenied, "uid %v is not in valid users: %v", uid, s.roomInfo.ValidUsers)
	}

	return nil
}

func (s *BaseUserService) BroadcastUpdate() {
	s.BaseUserService.BroadcastUpdate()
	s.msgBus.Broadcast(core.ModelTopic, &gamegrpc.Model{
		UserGroup: s.userGroup.ToProto(),
	})
}
