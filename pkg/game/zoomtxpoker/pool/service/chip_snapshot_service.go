package service

import (
	"context"
	"fmt"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/util"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"math/rand"
	"time"
)

type ChipSnapshotService struct {
	roomInfo    *commonmodel.RoomInfo
	redisClient *redis.Client
	logger      *zap.Logger
}

type chipSnapshot struct {
	RoomId     string    `redis:"roomId"`
	Uid        string    `redis:"uid"`
	Chip       int       `redis:"chip"`
	UpdateTime time.Time `redis:"updateTime"`
}

func (s *chipSnapshot) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("roomId", s.RoomId)
	enc.AddString("uid", s.Uid)
	enc.AddInt("chip", s.Chip)
	enc.AddTime("updateTime", s.UpdateTime)
	return nil
}

func ProvideChipSnapshotService(
	roomInfo *commonmodel.RoomInfo,
	redisClient *redis.Client,
	loggerFactory *util.LoggerFactory,
) *ChipSnapshotService {
	return &ChipSnapshotService{
		roomInfo:    roomInfo,
		redisClient: redisClient,
		logger:      loggerFactory.Create("[zoom-pool]ChipSnapshotService"),
	}
}

func (s *ChipSnapshotService) TakeSnapshot(uid core.Uid, chip int) {
	go func() {
		snapshot := chipSnapshot{
			RoomId:     s.roomInfo.RoomId,
			Uid:        uid.String(),
			Chip:       chip,
			UpdateTime: time.Now(),
		}

		if err := s.redisClient.HSet(
			context.TODO(),
			s.genSnapshotKey(uid),
			snapshot,
		).Err(); err != nil {
			s.logger.Error("failed to take snapshot", zap.Error(err),
				zap.String("roomId", snapshot.RoomId),
				zap.String("uid", snapshot.Uid),
				zap.Int("chip", snapshot.Chip),
				zap.Time("updateTime", snapshot.UpdateTime),
			)
			return
		}
		s.logger.Info("take snapshot", zap.String("uid", uid.String()), zap.Int("chip", chip))
	}()
}

func (s *ChipSnapshotService) AddToSnapshot(uid core.Uid, chip int) {
	go func() {
		if chip < 0 {
			s.logger.Error("chip add to snapshot must be positive", zap.Int("chip", chip), zap.String("uid", uid.String()))
			return
		}

		snapshotKey := s.genSnapshotKey(uid)
		trx := func(tx *redis.Tx) error {
			snapshot := &chipSnapshot{}
			if err := tx.HGetAll(context.TODO(), snapshotKey).Scan(snapshot); err != nil && err != redis.Nil {
				return err
			}

			snapshot.Chip += chip
			snapshot.UpdateTime = time.Now()

			_, err := tx.TxPipelined(context.TODO(), func(pipe redis.Pipeliner) error {
				pipe.HSet(context.TODO(), snapshotKey, snapshot)
				return nil
			})
			return err
		}

		for retries := 3; retries > 0; retries-- {
			err := s.redisClient.Watch(context.TODO(), trx, snapshotKey)
			if err == nil {
				s.logger.Info("add", zap.String("uid", uid.String()), zap.Int("chip", chip))
				// success
				return
			}

			if err == redis.TxFailedErr {
				// optimistic lock lost, retry
				time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
				continue
			}

			// other error
			s.logger.Error("failed to add to snapshot", zap.Error(err), zap.String("uid", uid.String()), zap.Int("chip", chip))
			return
		}

		s.logger.Error("failed to add to snapshot, trx max retries reached", zap.String("uid", uid.String()), zap.Int("chip", chip))
	}()
}

func (s *ChipSnapshotService) DeleteSnapshot(uid core.Uid) {
	go func() {
		if err := s.redisClient.Del(context.TODO(), s.genSnapshotKey(uid)).Err(); err != nil {
			s.logger.Error("failed to delete snapshot", zap.Error(err), zap.String("uid", uid.String()))
			return
		}
		s.logger.Info("delete snapshot", zap.String("uid", uid.String()))
	}()
}

func (s *ChipSnapshotService) genSnapshotKey(uid core.Uid) string {
	return fmt.Sprintf("roomId:%v:uid:%v", s.roomInfo.RoomId, uid)
}
