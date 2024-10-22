package service

import (
	"context"
	commonapi "card-game-server-prototype/pkg/common/api"
	commonconstant "card-game-server-prototype/pkg/common/constant"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/util"
	"time"

	"go.uber.org/zap"
)

type BaseRoomService struct {
	roomInfo *commonmodel.RoomInfo
	roomAPI  commonapi.RoomAPI
	logger   *zap.Logger

	cancelLoop context.CancelFunc
}

func ProvideBaseRoomService(
	roomInfo *commonmodel.RoomInfo,
	roomAPI commonapi.RoomAPI,
	loggerFactory *util.LoggerFactory,
) *BaseRoomService {
	return &BaseRoomService{
		roomInfo: roomInfo,
		roomAPI:  roomAPI,
		logger:   loggerFactory.Create("BaseRoomService"),
	}
}

func (s *BaseRoomService) GetDetail() (*commonapi.RoomDetail, error) {
	resp, err := s.roomAPI.GetDetail(s.roomInfo.RoomId)
	if err != nil {
		s.logger.Warn("failed to get room detail", zap.Error(err))
		return nil, err
	}

	data := &commonapi.RoomDetail{}
	if resp != nil {
		data = resp.Data
	}

	return data, nil
}

func (s *BaseRoomService) FetchRoomInfo() error {
	resp, err := s.roomAPI.GetDetail(s.roomInfo.RoomId)
	if err != nil {
		s.logger.Warn("failed to get room detail", zap.Error(err))
		return err
	}

	data := &commonapi.RoomDetail{}
	if resp != nil {
		data = resp.Data
	}

	s.roomInfo.IsPremium = data.IsPremium
	s.roomInfo.PremiumUid = data.PremiumUid

	premiumEndTime, err := time.Parse(time.RFC3339, data.PremiumEndTime)
	if err != nil {
		s.logger.Warn("failed to parse premium end time", zap.Error(err))
		return err
	}

	s.roomInfo.PremiumEndTimestamp = premiumEndTime.Unix()

	s.logger.Info("fetched room info", zap.Any("roomInfo", s.roomInfo))

	return nil
}

func (s *BaseRoomService) RunPingLoop() error {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelLoop = cancel

	ticker := time.NewTicker(commonconstant.RoomPingInterval)
	s.logger.Debug("ping loop started", zap.Duration("RoomPingInterval", commonconstant.RoomPingInterval))
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("canceling ping loop")
			ticker.Stop()
			return nil
		case <-ticker.C:
			var err error
			for i := 0; i < commonconstant.RoomPingFailThreshold; i++ {
				if err = s.roomAPI.Heartbeat(s.roomInfo.RoomId); err != nil {
					s.logger.Warn("failed to ping", zap.Error(err))
				} else {
					break
				}
			}

			if err != nil {
				s.logger.Error("failed to ping and reached threshold", zap.Error(err))
				return err
			}
		}
	}
}

func (s *BaseRoomService) Close() error {
	if s.cancelLoop != nil {
		s.cancelLoop()
	}

	if err := s.roomAPI.CloseRoom(s.roomInfo.RoomId); err != nil {
		return err
	}

	return nil
}
