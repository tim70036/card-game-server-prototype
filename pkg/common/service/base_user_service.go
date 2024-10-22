package service

import (
	"errors"
	"fmt"
	"card-game-server-prototype/pkg/common/api"
	commonapi "card-game-server-prototype/pkg/common/api"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/util"
	"strconv"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type BaseUserService struct {
	userGroup *commonmodel.UserGroup
	msgBus    core.MsgBus
	roomInfo  *commonmodel.RoomInfo
	userAPI   commonapi.UserAPI
	roomAPI   commonapi.RoomAPI
	logger    *zap.Logger
}

func ProvideBaseUserService(
	userGroup *commonmodel.UserGroup,
	msgBus core.MsgBus,
	roomInfo *commonmodel.RoomInfo,
	userAPI commonapi.UserAPI,
	roomAPI commonapi.RoomAPI,
	loggerFactory *util.LoggerFactory,
) *BaseUserService {
	return &BaseUserService{
		userGroup: userGroup,
		msgBus:    msgBus,
		roomInfo:  roomInfo,
		userAPI:   userAPI,
		roomAPI:   roomAPI,
		logger:    loggerFactory.Create("UserService"),
	}
}

func (s *BaseUserService) Init(uids ...core.Uid) error {
	alreadyExistUids := lo.Intersect(uids, lo.Keys(s.userGroup.Data))
	if len(alreadyExistUids) > 0 {
		return fmt.Errorf("these uids has already init: %v", alreadyExistUids)
	}

	for _, uid := range uids {
		s.userGroup.Data[uid] = &commonmodel.User{
			Uid:         uid,
			ShortUid:    "",
			Name:        "Undefined",
			IsAI:        false,
			IsConnected: false,
			HasEntered:  false,
			Cash:        0,
			Level:       1,
			RoomCards:   0,
		}
	}

	return nil
}

func (s *BaseUserService) GetAI(enterLimit, roomCardLimit int) (core.Uid, error) {
	aiUids, err := s.userAPI.GetIdleAIs()
	if err != nil {
		s.logger.Warn("failed to get idle AI", zap.Error(err))
		return "", err
	}

	var aiUid core.Uid

	for _, uid := range aiUids {
		aiUid = core.Uid(uid)

		userDetail, err := s.userAPI.FetchUserDetail(aiUid)
		if err != nil {
			s.logger.Warn("failed to fetch user detail", zap.Error(err))
			return "", err
		}

		if userDetail.Data.Cash < enterLimit {
			s.logger.Warn("AI lack of cash")
			continue
		}

		if userDetail.Data.RoomCards < roomCardLimit {
			s.logger.Warn("AI lack of room cards")
			continue
		}

		break
	}

	if aiUid == "" {
		s.logger.Warn("no suitable AI")
		return "", errors.New("no suitable AI")
	}

	return aiUid, nil
}

func (s *BaseUserService) FetchFromRepo(uids ...core.Uid) error {
	notExistUids, _ := lo.Difference(uids, lo.Keys(s.userGroup.Data))
	if len(notExistUids) > 0 {
		return fmt.Errorf("these uids do not exist in user group: %v", notExistUids)
	}

	results := make(chan lo.Tuple2[*api.UserDetailResponse, error], len(uids))
	for _, uid := range uids {
		go func(uid core.Uid) {
			resp, err := s.userAPI.FetchUserDetail(uid)
			results <- lo.T2(resp, err)
		}(uid)
	}

	var errs error = nil
	for i := 0; i < len(uids); i++ {
		result := <-results
		resp, err := result.A, result.B

		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		uid := core.Uid(resp.Data.Uid)
		if _, ok := s.userGroup.Data[uid]; !ok {
			errs = errors.Join(errs, fmt.Errorf("resp contains uid %v does not exist in user group", uid))
			continue
		}

		s.userGroup.Data[uid].ShortUid = strconv.Itoa(resp.Data.ShortUid)
		s.userGroup.Data[uid].Name = resp.Data.Name
		s.userGroup.Data[uid].Cash = resp.Data.Cash
		s.userGroup.Data[uid].Level = resp.Data.Level
		s.userGroup.Data[uid].RoomCards = resp.Data.RoomCards
		s.userGroup.Data[uid].IsAI = resp.Data.IsAi == 1
	}

	close(results)
	return errs
}

func (s *BaseUserService) Destroy(topic core.Topic, kickoutMSG proto.Message, uids ...core.Uid) error {
	notExistUids, _ := lo.Difference(uids, lo.Keys(s.userGroup.Data))
	if len(notExistUids) > 0 {
		return fmt.Errorf("these uids do not exist in user group: %v", notExistUids)
	}

	// TODO:
	// Need to call leave room api sequentially, since main server
	// don't have atomicity on room data. If called concurrently, room
	// data will be corrupted. (multiple users trying to modify the
	// same room data at the same time)
	var errs error = nil
	for idx, uid := range uids {
		if idx > 0 {
			time.Sleep(500 * time.Millisecond)
		}

		if err := s.roomAPI.LeaveRoom(s.roomInfo.RoomId, uid); err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		s.msgBus.Unicast(uid, topic, kickoutMSG)

		delete(s.userGroup.Data, uid)
	}

	return errs
}

func (s *BaseUserService) CanConnect(uid core.Uid) error { return nil }

func (s *BaseUserService) BroadcastUpdate() {}
