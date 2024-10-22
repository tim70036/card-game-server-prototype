package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/type/action"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/types/known/durationpb"
)

type ActionHint struct {
	Uid core.Uid

	// Chip that player bet (throw on table).
	BetChip int

	// This field only has meaning when action is Raise. It represents
	// the raise portion in BetChip. Used for checking min raise
	// requirement when other player wants to raise further.
	RaiseChip int

	// If player can call, this is the amount of chip that player must call.
	CallingChip int

	// If player can raise, this is the minimum chip that player must
	// raise.  The minimum raise is equal to the size of the previous
	// bet or raise. If someone wishes to re-raise, they must raise at
	// least the amount of the previous raise.
	MinRaisingChip int

	Action action.ActionType

	// Edge case: BB but is all in.
	// Use it in some case:
	//   1. collect pots must handle side pot.
	IsBBAllIn bool

	AvailableActions []action.ActionType

	Duration time.Duration
}

func (a *ActionHint) ToProto() *txpokergrpc.ActionHint {
	return &txpokergrpc.ActionHint{
		Uid:              a.Uid.String(),
		BetChip:          int32(a.BetChip),
		CallingChip:      int32(a.CallingChip),
		MinRaiseChip:     int32(a.MinRaisingChip),
		Action:           a.Action.ToProto(),
		AvailableActions: lo.Map(a.AvailableActions, func(act action.ActionType, _ int) txpokergrpc.BetActionType { return act.ToProto() }),
		Duration:         durationpb.New(a.Duration),
	}
}

func (a *ActionHint) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Uid", a.Uid.String())

	if a.BetChip != 0 {
		enc.AddInt("$Bet", a.BetChip)
	}
	if a.RaiseChip != 0 {
		enc.AddInt("$Raise", a.RaiseChip)
	}
	if a.CallingChip != 0 {
		enc.AddInt("$Calling", a.CallingChip)
	}
	if a.MinRaisingChip != 0 {
		enc.AddInt("$MinRaising", a.MinRaisingChip)
	}
	if len(a.AvailableActions) > 0 {
		enc.AddString("Allow", util.JoinStrings(lo.Map(a.AvailableActions, func(availableAction action.ActionType, _ int) string {
			return availableAction.String()
		}), ","))
	}
	enc.AddString("Act", a.Action.String())
	// enc.AddString("Duration", a.Duration.String())
	return nil
}

type ActionHintGroup struct {
	Hints map[core.Uid]*ActionHint

	// This is the hint of the player who is currently place a
	// bet/raise. Other player that choose to not fold, must follow or
	// raise further. In the case of incomplete raise, this is still
	// the hint of the player who had valid raised.
	RaiserHint *ActionHint
}

func ProvideActionHintGroup() *ActionHintGroup {
	return &ActionHintGroup{
		Hints:      make(map[core.Uid]*ActionHint),
		RaiserHint: nil,
	}
}

func (g *ActionHintGroup) ToProto() *txpokergrpc.ActionHintGroup {
	msg := &txpokergrpc.ActionHintGroup{
		Hints: make(map[string]*txpokergrpc.ActionHint),
	}

	if g.RaiserHint != nil {
		msg.RaiserHint = g.RaiserHint.ToProto()
	}

	for uid, hint := range g.Hints {
		msg.Hints[uid.String()] = hint.ToProto()
	}

	return msg
}

func (g *ActionHintGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	_ = enc.AddArray("Hints", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		for _, hint := range g.Hints {
			_ = enc.AppendObject(hint)
		}
		return nil
	}))

	if g.RaiserHint != nil {
		enc.AddString("RaiserHint", g.RaiserHint.Uid.String())
	}
	return nil
}
