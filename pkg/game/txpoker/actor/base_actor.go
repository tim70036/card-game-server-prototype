package actor

import (
	"fmt"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/type/seatstatus"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type BaseActorFactory struct {
	userGroup       *commonmodel.UserGroup
	seatStatusGroup *model.SeatStatusGroup
	loggerFactory   *util.LoggerFactory
}

func ProvideBaseActorFactory(
	userGroup *commonmodel.UserGroup,
	seatStatusGroup *model.SeatStatusGroup,
	loggerFactory *util.LoggerFactory,
) *BaseActorFactory {
	return &BaseActorFactory{
		userGroup:       userGroup,
		seatStatusGroup: seatStatusGroup,
		loggerFactory:   loggerFactory,
	}
}

func (f *BaseActorFactory) Create(uid core.Uid) *BaseActor {
	return &BaseActor{
		uid:             uid,
		reqDelay:        1000 * time.Millisecond,
		logger:          f.loggerFactory.Create(fmt.Sprintf("BaseActor[%s]", uid)),
		userGroup:       f.userGroup,
		seatStatusGroup: f.seatStatusGroup,
	}
}

type BaseActor struct {
	uid      core.Uid
	reqDelay time.Duration
	logger   *zap.Logger

	userGroup       *commonmodel.UserGroup
	seatStatusGroup *model.SeatStatusGroup
}

func (a *BaseActor) Uid() core.Uid { return a.uid }
func (a *BaseActor) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("uid", a.uid.String())
	enc.AddString("actorType", "BaseActor")
	return nil
}

func (a *BaseActor) DecideAction() (core.ActorRequestList, error) {
	actorReqs := make([]*core.ActorRequest, 0)

	user, userExists := a.userGroup.Data[a.uid]
	if !userExists {
		// 已離開，執行自動棄牌
		a.logger.Info("user leaved, auto fold")
		actorReqs = append(actorReqs, &core.ActorRequest{
			Delay: a.reqDelay,
			Req: &core.Request{
				Uid: a.uid,
				Msg: &txpokergrpc.FoldRequest{},
			},
		})
	} else if userExists && !user.IsConnected {
		// 斷線，執行自動棄牌
		a.logger.Info("user disconnected, auto fold")
		actorReqs = append(actorReqs, &core.ActorRequest{
			Delay: a.reqDelay,
			Req: &core.Request{
				Uid: a.uid,
				Msg: &txpokergrpc.FoldRequest{},
			},
		})
	}

	// 若某局尚未結束，玩家仍有手牌，執行站起觀戰，並視同自動棄牌
	// 若某局尚未結束，玩家仍有手牌，執行留座離桌，並視同自動棄牌
	if seatStatus, ok := a.seatStatusGroup.Status[a.uid]; ok &&
		seatStatus.FSM.MustState().(seatstatus.SeatStatusState) != seatstatus.PlayingState {
		a.logger.Info("seat status not in playing state, auto fold", zap.Object("seatStatus", seatStatus))
		actorReqs = append(actorReqs, &core.ActorRequest{
			Delay: a.reqDelay,
			Req: &core.Request{
				Uid: a.uid,
				Msg: &txpokergrpc.FoldRequest{},
			},
		})
	}

	return actorReqs, nil
}
