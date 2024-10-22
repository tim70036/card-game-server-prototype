package service

import (
	"errors"
	commonapi "card-game-server-prototype/pkg/common/api"
	commonmodel "card-game-server-prototype/pkg/common/model"
	commonservice "card-game-server-prototype/pkg/common/service"
	"card-game-server-prototype/pkg/config"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/constant"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	event2 "card-game-server-prototype/pkg/game/txpoker/type/event"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type BaseUserService struct {
	*commonservice.BaseUserService

	msgBus            core.MsgBus
	testCFG           *config.TestConfig
	roomInfo          *commonmodel.RoomInfo
	userGroup         *commonmodel.UserGroup
	roomAPI           commonapi.RoomAPI
	seatStatusGroup   *model2.SeatStatusGroup
	eventGroup        *model2.EventGroup
	playSettingGroup  *model2.PlaySettingGroup
	statsGroup        *model2.StatsGroup
	userCacheGroup    *model2.UserCacheGroup
	seatStatusService *BaseSeatStatusService
	tableProfitsGroup *model2.TableProfitsGroup
}

func ProvideBaseUserService(
	baseUserService *commonservice.BaseUserService,
	msgBus core.MsgBus,
	testCFG *config.TestConfig,
	roomInfo *commonmodel.RoomInfo,
	userGroup *commonmodel.UserGroup,
	roomAPI commonapi.RoomAPI,
	seatStatusGroup *model2.SeatStatusGroup,
	eventGroup *model2.EventGroup,
	playSettingGroup *model2.PlaySettingGroup,
	statsGroup *model2.StatsGroup,
	userCacheGroup *model2.UserCacheGroup,
	seatStatusService *BaseSeatStatusService,
	tableProfitsGroup *model2.TableProfitsGroup,
) *BaseUserService {
	return &BaseUserService{
		BaseUserService:   baseUserService,
		msgBus:            msgBus,
		testCFG:           testCFG,
		roomInfo:          roomInfo,
		userGroup:         userGroup,
		roomAPI:           roomAPI,
		seatStatusGroup:   seatStatusGroup,
		eventGroup:        eventGroup,
		playSettingGroup:  playSettingGroup,
		statsGroup:        statsGroup,
		userCacheGroup:    userCacheGroup,
		seatStatusService: seatStatusService,
		tableProfitsGroup: tableProfitsGroup,
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

		s.seatStatusGroup.Status[uid] = s.seatStatusService.NewSeatStatus(uid)

		s.playSettingGroup.Data[uid] = &model2.PlaySetting{
			Uid:                       uid,
			WaitBB:                    true,
			AutoTopUp:                 false,
			AutoTopUpThresholdPercent: 0.0,
			AutoTopUpChipPercent:      0.0,
		}

		s.userCacheGroup.Data[uid] = &commonmodel.User{
			Uid:         s.userGroup.Data[uid].Uid,
			ShortUid:    s.userGroup.Data[uid].ShortUid,
			Name:        s.userGroup.Data[uid].Name,
			IsAI:        s.userGroup.Data[uid].IsAI,
			IsConnected: s.userGroup.Data[uid].IsConnected,
			HasEntered:  s.userGroup.Data[uid].HasEntered,
			Cash:        s.userGroup.Data[uid].Cash,
			Level:       s.userGroup.Data[uid].Level,
			RoomCards:   s.userGroup.Data[uid].RoomCards,
		}

		s.eventGroup.Data[uid] = event2.EventList{}

		s.statsGroup.Data[uid] = &model2.Stats{
			Uid:            uid,
			EventAmountSum: map[event2.EventType]int{},
		}

		s.tableProfitsGroup.Save(&model2.TableProfits{
			Uid:             uid,
			Name:            "",
			CountGames:      0,
			SumBuyInChips:   0,
			SumWinLoseChips: 0,
		})
	}

	close(results)
	return errs
}

func (s *BaseUserService) FetchFromRepo(uids ...core.Uid) error {
	if err := s.BaseUserService.FetchFromRepo(uids...); err != nil {
		return err
	}

	for _, uid := range uids {
		if cacheUser, ok := s.userCacheGroup.Data[uid]; ok {
			u := s.userGroup.Data[uid]

			cacheUser.Uid = u.Uid
			cacheUser.ShortUid = u.ShortUid
			cacheUser.Name = u.Name
			cacheUser.IsAI = u.IsAI
			cacheUser.IsConnected = u.IsConnected
			cacheUser.HasEntered = u.HasEntered
			cacheUser.Cash = u.Cash
			cacheUser.Level = u.Level
			cacheUser.RoomCards = u.RoomCards
		}

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
		// TODO: check if user is not in standing state?
		delete(s.seatStatusGroup.Status, uid)
		delete(s.playSettingGroup.Data, uid)
		delete(s.eventGroup.Data, uid)
		delete(s.statsGroup.Data, uid)
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

	return nil
}

func (s *BaseUserService) BroadcastUpdate() {
	s.BaseUserService.BroadcastUpdate()
	s.msgBus.Broadcast(core.MessageTopic, &txpokergrpc.Message{
		Model: &txpokergrpc.Model{
			UserGroup:       s.userGroup.ToProto(),
			SeatStatusGroup: s.seatStatusGroup.ToProto(),
			UserCacheGroup:  s.userCacheGroup.ToProto(),
		},
	})
}
