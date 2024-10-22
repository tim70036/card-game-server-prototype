package model

import (
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/darkchess/constant"
	"card-game-server-prototype/pkg/game/darkchess/type/piece"
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"github.com/samber/lo"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/types/known/durationpb"
)

type ClaimDraw struct {
	ClaimUid   core.Uid
	ClaimTurn  int
	IsAnswered bool
	IsAccepted bool
}

func (c *ClaimDraw) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("ClaimUid", c.ClaimUid.String())
	enc.AddInt("ClaimTurn", c.ClaimTurn)
	enc.AddBool("IsAnswered", c.IsAnswered)
	enc.AddBool("IsAccepted", c.IsAccepted)
	return nil
}

type LastActionCell struct {
	Pos             *Pos
	Piece           piece.Piece
	IsRevealAction  bool
	IsMoveAction    bool
	IsCaptureAction bool
}

func (d *LastActionCell) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	_ = enc.AddObject("Pos", d.Pos)
	enc.AddString("Piece", d.Piece.GetName())
	if d.IsRevealAction {
		enc.AddBool("IsRevealAction", d.IsRevealAction)
	}

	if d.IsMoveAction {
		enc.AddBool("IsMoveAction", d.IsMoveAction)
	}

	if d.IsCaptureAction {
		enc.AddBool("IsCaptureAction", d.IsCaptureAction)
	}

	return nil
}

type ActionHintGroup struct {
	Data             map[core.Uid]*ActionHint
	RepeatMovesCount int
	ClaimDraws       []*ClaimDraw
	LastAction       *LastActionCell
}

func ProvideActionHintGroup() *ActionHintGroup {
	return &ActionHintGroup{
		Data: map[core.Uid]*ActionHint{},
	}
}

func (g *ActionHintGroup) ToProto() *gamegrpc.ActionHintGroup {
	actionHints := make(map[string]*gamegrpc.ActionHint)

	for uid, p := range g.Data {
		claimDraws := lo.Filter(g.ClaimDraws, func(item *ClaimDraw, _ int) bool {
			return item.ClaimUid == uid
		})

		actionHints[uid.String()] = &gamegrpc.ActionHint{
			Uid:                      p.Uid.String(),
			TurnDuration:             durationpb.New(p.TurnDuration),
			TurnCount:                int32(p.TurnCount),
			RemainingTimeExtendCount: int32(constant.TimeExtendCount - len(p.TimeExtendTurns)),
			RemainingDrawOfferCount:  int32(constant.ClaimDrawCount - len(claimDraws)),
		}

		if p.FreezeCell != nil {
			actionHints[uid.String()].FreezeCell = p.FreezeCell.ToProto()
		}
	}

	var lastAction *gamegrpc.LastActionCell
	if g.LastAction != nil {
		lastAction = &gamegrpc.LastActionCell{
			GridPosition: &gamegrpc.GridPosition{
				X: int32(g.LastAction.Pos.X),
				Y: int32(g.LastAction.Pos.Y),
			},
			Piece:           g.LastAction.Piece.ToProto(),
			IsRevealAction:  g.LastAction.IsRevealAction,
			IsMoveAction:    g.LastAction.IsMoveAction,
			IsCaptureAction: g.LastAction.IsCaptureAction,
		}
	}

	return &gamegrpc.ActionHintGroup{
		ActionHints:      actionHints,
		RepeatMovesCount: int32(g.RepeatMovesCount),
		LastAction:       lastAction,
	}
}

func (g *ActionHintGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if len(g.Data) > 0 {
		_ = enc.AddArray("hints", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
			for _, p := range g.Data {
				_ = enc.AppendObject(p)
			}
			return nil
		}))
	}

	enc.AddInt("RepeatMovesCount", g.RepeatMovesCount)

	if len(g.ClaimDraws) > 0 {
		_ = enc.AddArray("ClaimDraws", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
			for _, claimDraw := range g.ClaimDraws {
				_ = enc.AppendObject(claimDraw)
			}
			return nil
		}))
	}

	if g.LastAction != nil {
		_ = enc.AddObject("LastAction", g.LastAction)
	}

	return nil
}
