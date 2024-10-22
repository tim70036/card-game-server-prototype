package service

import (
	"errors"
	commonapi "card-game-server-prototype/pkg/common/api"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"time"
)

type BuddyModeUserService struct {
	*BaseUserService

	buddyGroup *commonmodel.BuddyGroup
	roomAPI    commonapi.RoomAPI
}

func ProvideBuddyModeUserService(
	baseUserService *BaseUserService,

	buddyGroup *commonmodel.BuddyGroup,
	roomAPI commonapi.RoomAPI,
) *BuddyModeUserService {
	return &BuddyModeUserService{
		BaseUserService: baseUserService,

		buddyGroup: buddyGroup,
		roomAPI:    roomAPI,
	}
}

func (s *BuddyModeUserService) Init(uids ...core.Uid) error {
	type result struct {
		uid core.Uid
		err error
	}

	results := make(chan *result, len(uids))

	for _, uid := range uids {
		go func(uid core.Uid) {
			err := s.roomAPI.EnterRoom(s.roomInfo.RoomId, uid)
			results <- &result{uid, err}
		}(uid)
	}

	var errs error = nil
	for i := 0; i < len(uids); i++ {
		data := <-results
		uid, err := data.uid, data.err
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		if err := s.BaseUserService.Init(uid); err != nil {
			// Revert enter room if cannot init user. Otherwise, user
			// will be stucked.
			go s.roomAPI.LeaveRoom(s.roomInfo.RoomId, uid)
			errs = errors.Join(errs, err)
			continue
		}

		s.buddyGroup.Data[uid] = &commonmodel.Buddy{
			Uid:       uid,
			IsReady:   false,
			IsOwner:   false,
			EnterTime: time.Now(),
		}
	}

	close(results)
	return errs
}

func (s *BuddyModeUserService) Destroy(topic core.Topic, kickoutMSG proto.Message, uids ...core.Uid) error {
	if err := s.BaseUserService.Destroy(topic, kickoutMSG, uids...); err != nil {
		return err
	}

	for _, uid := range uids {
		delete(s.buddyGroup.Data, uid)
	}

	return nil
}

func (s *BuddyModeUserService) CanConnect(uid core.Uid) error {
	if err := s.BaseUserService.BaseUserService.CanConnect(uid); err != nil {
		return err
	}

	if _, ok := s.userGroup.Data[uid]; !ok && len(s.userGroup.Data) >= constant.MaxUserCount {
		return status.Errorf(codes.FailedPrecondition, "exceed max user count %v", constant.MaxUserCount)
	}
	return nil
}

func (s *BuddyModeUserService) BroadcastUpdate() {
	s.BaseUserService.BroadcastUpdate()
	s.msgBus.Broadcast(core.ModelTopic, &gamegrpc.Model{
		BuddyGroup: s.buddyGroup.ToProto(),
	})
}
