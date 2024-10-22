package state

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/model"
	"card-game-server-prototype/pkg/game/txpoker/service"
	"card-game-server-prototype/pkg/game/txpoker/type/hand"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"reflect"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoJackpotState = &core.StateTrigger{
	Name:      "GoJackpotState",
	ArgsTypes: []reflect.Type{},
}

type JackpotState struct {
	core.State

	playerGroup *model.PlayerGroup
	gameRepo    service.GameRepoService

	jackpotResult map[core.Uid]int
}

func ProvideJackpotState(
	stateFactory *core.StateFactory,
	playerGroup *model.PlayerGroup,
	gameRepo service.GameRepoService,
) *JackpotState {
	return &JackpotState{
		State: stateFactory.Create("JackpotState"),

		playerGroup: playerGroup,
		gameRepo:    gameRepo,
	}
}

func (state *JackpotState) Run(ctx context.Context, args ...any) error {
	var err error
	state.jackpotResult, err = state.gameRepo.TriggerJackpot()
	if err != nil {
		state.Logger().Error("trigger jackpot failed", zap.Error(err))
		state.GameController().GoErrorState()
		return nil
	}

	if len(state.jackpotResult) <= 0 {
		state.GameController().GoNextState(GoEndRoundState)
		return nil
	}

	duration := time.Duration(len(state.jackpotResult)*3500) * time.Millisecond
	state.GameController().RunTimer(duration, func() {
		state.GameController().GoNextState(GoEndRoundState)
	})
	return nil
}

func (state *JackpotState) Publish(ctx context.Context, args ...any) error {
	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
	})
	return nil
}

func (state *JackpotState) ToProto(uid core.Uid) proto.Message {
	jackpots := lo.MapToSlice(state.jackpotResult, func(uid core.Uid, amount int) *txpokergrpc.Jackpot {
		return &txpokergrpc.Jackpot{
			Uid:    uid.String(),
			Hand:   hand.ToProto(state.playerGroup.Data[uid].Hand),
			Amount: int32(amount),
		}
	})

	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context: &txpokergrpc.GameState_JackpotStateContext{JackpotStateContext: &txpokergrpc.JackpotStateContext{
			Jackpots: jackpots,
		}},
	}
}
