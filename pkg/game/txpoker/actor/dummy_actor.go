package actor

import (
	"fmt"
	"card-game-server-prototype/pkg/core"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/type/action"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"
	"math/rand"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"
)

type DummyActorFactory struct {
	loggerFactory   *util.LoggerFactory
	actionHintGroup *model2.ActionHintGroup
	gameSetting     *model2.GameSetting
}

func ProvideDummyActorFactory(
	loggerFactory *util.LoggerFactory,
	actionHintGroup *model2.ActionHintGroup,
	gameSetting *model2.GameSetting,
) *DummyActorFactory {
	return &DummyActorFactory{
		loggerFactory:   loggerFactory,
		actionHintGroup: actionHintGroup,
		gameSetting:     gameSetting,
	}
}

func (f *DummyActorFactory) Create(uid core.Uid) *DummyActor {
	return &DummyActor{
		uid:      uid,
		reqDelay: 2000 * time.Millisecond,
		logger:   f.loggerFactory.Create(fmt.Sprintf("DummyActor[%s]", uid)),

		actionHintGroup: f.actionHintGroup,
		gameSetting:     f.gameSetting,
	}
}

type DummyActor struct {
	uid      core.Uid
	reqDelay time.Duration
	logger   *zap.Logger

	actionHintGroup *model2.ActionHintGroup
	gameSetting     *model2.GameSetting
}

func (a *DummyActor) Uid() core.Uid { return a.uid }
func (a *DummyActor) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("uid", a.uid.String())
	enc.AddString("actorType", "DummyActor")
	return nil
}

func (a *DummyActor) DecideAction() (core.ActorRequestList, error) {
	// if len(a.actionHintGroup.Hints[a.uid].AvailableActions) > 1 {
	// 	a.actionHintGroup.Hints[a.uid].AvailableActions = lo.Without(a.actionHintGroup.Hints[a.uid].AvailableActions, action.Fold)
	// }

	randAction := lo.Sample(a.actionHintGroup.Hints[a.uid].AvailableActions)
	randNum := rand.Intn(100)
	if randNum > 30 && len(a.actionHintGroup.Hints[a.uid].AvailableActions) > 1 {
		randAction = lo.Sample(lo.Without(a.actionHintGroup.Hints[a.uid].AvailableActions, action.AllIn))
	}

	// randAction := lo.Sample(a.actionHintGroup.Hints[a.uid].AvailableActions)
	// if lo.Contains(a.actionHintGroup.Hints[a.uid].AvailableActions, action.AllIn) {
	// 	randAction = action.AllIn
	// }

	actorReqs := make([]*core.ActorRequest, 0)
	var msg proto.Message = nil
	switch randAction {
	case action.Fold:
		msg = &txpokergrpc.FoldRequest{}
	case action.Check:
		msg = &txpokergrpc.CheckRequest{}
	case action.Bet:
		msg = &txpokergrpc.BetRequest{Chip: int32(a.gameSetting.BigBlind * (1 + rand.Intn(20)))}
	case action.Call:
		msg = &txpokergrpc.CallRequest{}
	case action.Raise:
		raisingChip := a.actionHintGroup.Hints[a.uid].MinRaisingChip + rand.Intn(10)*a.gameSetting.BigBlind
		msg = &txpokergrpc.RaiseRequest{Chip: int32(raisingChip)}
	case action.AllIn:
		msg = &txpokergrpc.AllInRequest{}
	}

	if msg == nil {
		return nil, fmt.Errorf("invalid randAction: %s", randAction.String())
	}

	actorReqs = append(actorReqs, &core.ActorRequest{
		Delay: a.reqDelay,
		Req: &core.Request{
			Uid: a.uid,
			Msg: msg,
		},
	})

	return actorReqs, nil
}
