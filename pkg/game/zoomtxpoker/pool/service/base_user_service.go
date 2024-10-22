package service

import (
	"errors"
	commonapi "card-game-server-prototype/pkg/common/api"
	commonmodel "card-game-server-prototype/pkg/common/model"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/type/event"
	"card-game-server-prototype/pkg/game/txpoker/type/role"
	model3 "card-game-server-prototype/pkg/game/zoomtxpoker/pool/model"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type BaseUserService struct {
	*commonservice.BaseUserService

	msgBus             core.MsgBus
	testCFG            *config.TestConfig
	roomInfo           *commonmodel.RoomInfo
	gameSetting        *model3.GameSetting
	userGroup          *commonmodel.UserGroup
	roomAPI            commonapi.RoomAPI
	playSettingGroup   *model3.PlaySettingGroup
	statsGroup         *model3.StatsGroup
	participantGroup   *model3.ParticipantGroup
	playedHistoryGroup *model3.PlayedHistoryGroup
	forceBuyInGroup    *model3.ForceBuyInGroup
	tableProfitsGroup  *model3.TableProfitsGroup
	participantService *ParticipantService
}

func ProvideBaseUserService(
	baseUserService *commonservice.BaseUserService,
	msgBus core.MsgBus,
	testCFG *config.TestConfig,
	roomInfo *commonmodel.RoomInfo,
	gameSetting *model3.GameSetting,
	userGroup *commonmodel.UserGroup,
	roomAPI commonapi.RoomAPI,
	playSettingGroup *model3.PlaySettingGroup,
	statsGroup *model3.StatsGroup,
	participantGroup *model3.ParticipantGroup,
	playedHistoryGroup *model3.PlayedHistoryGroup,
	forceBuyInGroup *model3.ForceBuyInGroup,
	tableProfitsGroup *model3.TableProfitsGroup,
	participantService *ParticipantService,
) *BaseUserService {
	return &BaseUserService{
		BaseUserService:    baseUserService,
		msgBus:             msgBus,
		testCFG:            testCFG,
		roomInfo:           roomInfo,
		gameSetting:        gameSetting,
		userGroup:          userGroup,
		roomAPI:            roomAPI,
		playSettingGroup:   playSettingGroup,
		statsGroup:         statsGroup,
		participantGroup:   participantGroup,
		playedHistoryGroup: playedHistoryGroup,
		forceBuyInGroup:    forceBuyInGroup,
		tableProfitsGroup:  tableProfitsGroup,
		participantService: participantService,
	}
}

func (s *BaseUserService) Init(uids ...core.Uid) error {
	results := make(chan lo.Tuple2[core.Uid, error], len(uids))
	for _, uid := range uids {
		go func(uid core.Uid) {
			err := s.roomAPI.EnterRoom(s.roomInfo.RoomId, uid)
			results <- lo.T2(uid, err)
		}(uid)
	}

	var errs error = nil
	for i := 0; i < len(uids); i++ {
		result := <-results
		uid, err := result.A, result.B

		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		if err := s.BaseUserService.Init(uid); err != nil {
			// Revert enter room if cannot init user. Otherwise, user
			// will be stuck.
			go s.roomAPI.LeaveRoom(s.roomInfo.RoomId, uid)
			errs = errors.Join(errs, err)
			continue
		}

		s.playSettingGroup.Data[uid] = &model2.PlaySetting{
			Uid:                       uid,
			WaitBB:                    true,
			AutoTopUp:                 false,
			AutoTopUpThresholdPercent: 0.0,
			AutoTopUpChipPercent:      0.0,
		}

		s.playedHistoryGroup.Data[uid] = &model3.PlayedHistory{
			Uid:            uid,
			LastPlayedRole: role.Undefined,
			CountRoles:     map[role.Role]int{},
		}

		s.statsGroup.Data[uid] = &model3.Stats{
			Uid:            uid,
			EventAmountSum: map[event.EventType]int{},
		}

		s.tableProfitsGroup.Save(&model2.TableProfits{
			Uid:             uid,
			Name:            "",
			CountGames:      0,
			SumBuyInChips:   0,
			SumWinLoseChips: 0,
		})

		s.participantGroup.Data[uid] = s.participantService.NewParticipant(uid)
	}

	close(results)
	return errs
}

func (s *BaseUserService) FetchFromRepo(uids ...core.Uid) error {
	if err := s.BaseUserService.FetchFromRepo(uids...); err != nil {
		return err
	}

	for _, uid := range uids {
		if tableProfits, ok := s.tableProfitsGroup.Get(uid); ok {
			tableProfits.Name = s.userGroup.Data[uid].Name
			s.tableProfitsGroup.Save(tableProfits)
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
		delete(s.statsGroup.Data, uid)
		delete(s.participantGroup.Data, uid)
		delete(s.playedHistoryGroup.Data, uid)
	}

	return nil
}

func (s *BaseUserService) CanConnect(uid core.Uid) error {
	if err := s.BaseUserService.CanConnect(uid); err != nil {
		return err
	}

	if _, ok := s.userGroup.Data[uid]; !ok && len(s.userGroup.Data) >= s.gameSetting.MaxUserAmount {
		return status.Errorf(codes.FailedPrecondition, "exceed max user count %v", s.gameSetting.MaxUserAmount)
	}

	return nil
}
