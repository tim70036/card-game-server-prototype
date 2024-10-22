package state

import (
	"context"
	"fmt"
	commonmodel "card-game-server-prototype/pkg/common/model"
	"card-game-server-prototype/pkg/core"
	actor2 "card-game-server-prototype/pkg/game/txpoker/actor"
	model2 "card-game-server-prototype/pkg/game/txpoker/model"
	service2 "card-game-server-prototype/pkg/game/txpoker/service"
	"card-game-server-prototype/pkg/game/txpoker/type/action"
	"card-game-server-prototype/pkg/game/txpoker/type/card"
	role2 "card-game-server-prototype/pkg/game/txpoker/type/role"
	"card-game-server-prototype/pkg/game/txpoker/type/seatstatus"
	"card-game-server-prototype/pkg/game/txpoker/type/stage"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"go.uber.org/zap/zapcore"
	"reflect"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var GoStartRoundState = &core.StateTrigger{
	Name:      "GoStartRoundState",
	ArgsTypes: []reflect.Type{},
}

type StartRoundState struct {
	core.State

	userGroup       *commonmodel.UserGroup
	seatStatusGroup *model2.SeatStatusGroup
	statsGroup      *model2.StatsGroup
	playerGroup     *model2.PlayerGroup
	actionHintGroup *model2.ActionHintGroup
	statsCacheGroup *model2.StatsCacheGroup
	chipCacheGroup  *model2.ChipCacheGroup

	gameSetting       *model2.GameSetting
	gameInfo          *model2.GameInfo
	table             *model2.Table
	replay            *model2.Replay
	actionHintService *service2.ActionHintService
	gameRepoService   service2.GameRepoService
	seatStatusService service2.SeatStatusService

	actorGroup        *actor2.ActorGroup
	baseActorFactory  *actor2.BaseActorFactory
	dummyActorFactory *actor2.DummyActorFactory
}

func ProvideStartRoundState(
	stateFactory *core.StateFactory,
	userGroup *commonmodel.UserGroup,
	seatStatusGroup *model2.SeatStatusGroup,
	statsGroup *model2.StatsGroup,
	playerGroup *model2.PlayerGroup,
	actionHintGroup *model2.ActionHintGroup,
	statsCacheGroup *model2.StatsCacheGroup,
	chipCacheGroup *model2.ChipCacheGroup,

	gameSetting *model2.GameSetting,
	gameInfo *model2.GameInfo,
	table *model2.Table,
	replay *model2.Replay,
	actionHintService *service2.ActionHintService,
	gameRepoService service2.GameRepoService,
	seatStatusService service2.SeatStatusService,

	actorGroup *actor2.ActorGroup,
	baseActorFactory *actor2.BaseActorFactory,
	dummyActorFactory *actor2.DummyActorFactory,
) *StartRoundState {
	return &StartRoundState{
		State: stateFactory.Create("StartRoundState"),

		userGroup:       userGroup,
		seatStatusGroup: seatStatusGroup,
		statsGroup:      statsGroup,
		playerGroup:     playerGroup,
		actionHintGroup: actionHintGroup,
		statsCacheGroup: statsCacheGroup,
		chipCacheGroup:  chipCacheGroup,

		gameSetting:       gameSetting,
		gameInfo:          gameInfo,
		table:             table,
		replay:            replay,
		actionHintService: actionHintService,
		gameRepoService:   gameRepoService,
		seatStatusService: seatStatusService,

		actorGroup:        actorGroup,
		baseActorFactory:  baseActorFactory,
		dummyActorFactory: dummyActorFactory,
	}
}

func (state *StartRoundState) Run(ctx context.Context, args ...any) error {
	// If player sit out or leave, go to previous state.
	if !state.seatStatusService.IsReadyToStartRound() {
		state.Logger().Debug("not enough joining to start round",
			zap.Object("seatStatus", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
				for _, s := range state.seatStatusGroup.Status {
					enc.AddString(s.Uid.String(), s.FSM.MustState().(seatstatus.SeatStatusState).String())
				}
				return nil
			})),
		)
		state.GameController().GoNextState(GoWaitUserState)
		return nil
	}

	state.gameInfo.RoundId = uuid.NewString()
	state.gameInfo.StreakRoundCount++

	if err := state.seatStatusService.RoundStarted(); err != nil {
		state.Logger().Error(
			"failed to start round for seat status group",
			zap.Error(err),
			zap.Object("seatStatus", state.seatStatusGroup),
			zap.Object("players", state.playerGroup),
		)
		state.GameController().GoErrorState()
		return nil
	}

	if err := state.initModel(); err != nil {
		state.Logger().Error(
			"failed to init model",
			zap.Error(err),
			zap.String("roundId", state.gameInfo.RoundId),
			zap.Object("seatStatus", state.seatStatusGroup),
			zap.Object("players", state.playerGroup),
		)
		state.GameController().GoErrorState()
		return nil
	}

	if err := state.placeStartRoundBet(); err != nil {
		state.Logger().Error(
			"failed to place start round bet",
			zap.Error(err),
			zap.String("roundId", state.gameInfo.RoundId),
			zap.Object("seatStatus", state.seatStatusGroup),
			zap.Object("players", state.playerGroup),
		)
		state.GameController().GoErrorState()
		return nil
	}

	// log merged, no issues reported in a long time.

	state.Logger().Info("starting round, all set",
		zap.Object("gameInfo", state.gameInfo),
		zap.Object("seatStatus", state.seatStatusGroup),
		zap.Object("players", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			for uid, p := range state.playerGroup.Data {
				_ = enc.AddObject(uid.String(), zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
					enc.AddString("Role", p.Role.String())
					if state.playerGroup.LastBBPlayer != nil && uid == state.playerGroup.LastBBPlayer.Uid {
						enc.AddBool("LastBBPlayer", true)
					}
					return nil
				}))
			}
			return nil
		})),
	)

	if err := state.table.BetStageFSM.Fire(stage.NextStageTrigger); err != nil {
		state.Logger().Error(
			"failed to fire next stage trigger",
			zap.Error(err),
			zap.Object("table", state.table),
		)
		state.GameController().GoErrorState()
		return nil
	}

	if err := state.gameRepoService.CreateRound(); err != nil {
		state.Logger().Error(
			"failed to create round",
			zap.Error(err),
			zap.Object("gameInfo", state.gameInfo),
		)
		state.GameController().GoErrorState()
		return nil
	}

	state.GameController().GoNextState(GoDealPocketState)
	return nil
}

func (state *StartRoundState) Publish(ctx context.Context, args ...any) error {
	playerGroupProto := &txpokergrpc.PlayerGroup{Players: make(map[string]*txpokergrpc.Player)}
	for uid, player := range state.playerGroup.Data {
		playerGroupProto.Players[uid.String()] = player.ToProto()
		_, playerGroupProto.Players[uid.String()].HasShowdown = state.table.ShowdownPocketCards[uid]
	}

	state.MsgBus().Broadcast(core.MessageTopic, &txpokergrpc.Message{
		GameState: state.ToProto(core.Uid("")).(*txpokergrpc.GameState),
		Model: &txpokergrpc.Model{
			SeatStatusGroup: state.seatStatusGroup.ToProto(),
			PlayerGroup:     playerGroupProto,
			ActionHintGroup: state.actionHintGroup.ToProto(),
			GameInfo:        state.gameInfo.ToProto(),
			StatsCacheGroup: state.statsCacheGroup.ToProto(),
			ChipCacheGroup:  state.chipCacheGroup.ToProto(),
		},
	})
	return nil
}

func (state *StartRoundState) ToProto(uid core.Uid) proto.Message {
	return &txpokergrpc.GameState{
		Timestamp: timestamppb.New(state.PublishTime()),
		Name:      state.Name(),
		Context:   &txpokergrpc.GameState_StartRoundStateContext{StartRoundStateContext: &txpokergrpc.StartRoundStateContext{}},
	}
}

func (state *StartRoundState) initModel() error {
	var playingUids core.UidList = lo.FilterMap(lo.Values(state.seatStatusGroup.Status), func(s *model2.SeatStatus, _ int) (core.Uid, bool) {
		return s.Uid, s.FSM.MustState().(seatstatus.SeatStatusState) == seatstatus.PlayingState
	})

	seatIdMap := lo.Invert(state.seatStatusGroup.TableUids)
	for _, uid := range playingUids {
		state.playerGroup.Data[uid] = &model2.Player{
			Uid:         uid,
			Role:        role2.Undefined,
			SeatId:      seatIdMap[uid],
			PocketCards: card.CardList{},
			Hand:        nil,
		}

		state.actionHintGroup.Hints[uid] = &model2.ActionHint{
			Uid:              uid,
			BetChip:          0,
			RaiseChip:        0,
			CallingChip:      0,
			MinRaisingChip:   0,
			Action:           action.Undefined,
			AvailableActions: []action.ActionType{},
			Duration:         state.gameSetting.TurnDuration,
		}

		state.statsCacheGroup.Data[uid] = state.statsGroup.Data[uid]
		state.chipCacheGroup.SeatStatusChips[uid] = state.seatStatusGroup.Status[uid].Chip

		if state.userGroup.Data[uid].IsAI {
			state.actorGroup.Data[uid] = state.dummyActorFactory.Create(uid)
		} else {
			state.actorGroup.Data[uid] = state.baseActorFactory.Create(uid)
		}
	}

	// Role assignment
	lastBBSeatId := -1
	if state.playerGroup.LastBBPlayer != nil {
		lastBBSeatId = state.playerGroup.LastBBPlayer.SeatId
	}

	roleAssignment, err := role2.EvalRoleAssignment(playingUids, state.seatStatusGroup.TableUids, lastBBSeatId)
	if err != nil {
		return fmt.Errorf("failed to eval role assignment: %w", err)
	}

	for uid, role := range roleAssignment {
		state.playerGroup.Data[uid].Role = role
	}

	bbUid := lo.Invert(roleAssignment)[role2.BB]
	state.playerGroup.LastBBPlayer = state.playerGroup.Data[bbUid]
	// log merged
	state.Logger().Debug("roles assigned",
		zap.Object("players", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			for uid, p := range state.playerGroup.Data {
				_ = enc.AddObject(uid.String(), zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
					enc.AddString("Role", p.Role.String())
					// enc.AddInt("SeatId", p.SeatId)
					if state.playerGroup.LastBBPlayer != nil && uid == state.playerGroup.LastBBPlayer.Uid {
						enc.AddBool("LastBBPlayer", true)
					}
					return nil
				}))
			}
			return nil
		})),
	)

	state.replay.RoundId = state.gameInfo.RoundId
	for uid, player := range state.playerGroup.Data {
		state.replay.PlayerRecords[uid] = &model2.PlayerRecord{
			Uid:                     uid,
			Role:                    player.Role,
			SeatId:                  player.SeatId,
			PocketCards:             card.CardList{},
			InitSeatStatusChip:      state.seatStatusGroup.Status[uid].Chip,
			BeforeWinSeatStatusChip: 0,
			IsWinner:                false,
			HasShowdown:             false,
		}
	}
	return nil
}

// Place sb, bb bets
// Do it in this state to prevent race condition such as:
// BB player stand up before placing BB bet.
func (state *StartRoundState) placeStartRoundBet() error {
	bbPlayer, ok := lo.Find(lo.Values(state.playerGroup.Data), func(p *model2.Player) bool { return p.Role == role2.BB })
	if !ok {
		return fmt.Errorf("bb player not found")
	}

	if err := state.actionHintService.OpenBet(bbPlayer.Uid, action.BB, state.gameSetting.BigBlind); err != nil {
		return fmt.Errorf("failed to open bet for bb: %w", err)
	}

	sbPlayer, ok := lo.Find(lo.Values(state.playerGroup.Data), func(p *model2.Player) bool { return p.Role == role2.SB })
	if !ok {
		return fmt.Errorf("sb player not found")
	}

	if err := state.actionHintService.FollowBet(sbPlayer.Uid, action.SB); err != nil {
		return fmt.Errorf("failed to follow bet for sb: %w", err)
	}

	// Place BB bet for wait BB players to. (First round is skipped,
	// doesn't make sense to place BB bet for them)
	if state.gameInfo.StreakRoundCount > 1 {
		placingBBUids := lo.Filter(lo.Keys(state.playerGroup.Data), func(uid core.Uid, _ int) bool {
			seatStatus := state.seatStatusGroup.Status[uid]
			player := state.playerGroup.Data[uid]
			return seatStatus.ShouldPlaceBB && player.Role != role2.BB && player.Role != role2.SB
		})

		for _, uid := range placingBBUids {
			if err := state.actionHintService.FollowBet(uid, action.BB); err != nil {
				return fmt.Errorf("failed to follow bet for place bb playyer %v: %w", uid.String(), err)
			}
		}
	}

	for uid := range state.playerGroup.Data {
		state.seatStatusGroup.Status[uid].ShouldPlaceBB = false
	}

	// Mark those who should place BB in the next game. Sitting out
	// players should place BB next game, only when the BB role going
	// pass them.
	for i := 0; i < state.gameSetting.TableSize; i++ {
		seatId := (bbPlayer.SeatId - i + state.gameSetting.TableSize) % state.gameSetting.TableSize
		if uid, ok := state.seatStatusGroup.TableUids[seatId]; ok {
			if seatStatus, ok := state.seatStatusGroup.Status[uid]; ok {
				if seatStatus.FSM.MustState().(seatstatus.SeatStatusState) == seatstatus.SittingOutState {
					seatStatus.ShouldPlaceBB = true
					state.Logger().Debug("marked should place BB", zap.Object("seatStatus", seatStatus))
				} else {
					break
				}
			}

		}
	}

	// log merged
	state.Logger().Debug("bets placed", zap.Object("actionHints", state.actionHintGroup))
	return nil
}
