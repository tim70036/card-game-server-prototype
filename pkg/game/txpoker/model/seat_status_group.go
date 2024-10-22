package model

import (
	"context"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/type/seatstatus"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"strconv"
	"time"

	"github.com/qmuntal/stateless"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SeatStatus struct {
	Uid core.Uid

	// The amount of chip this user has buy in from wallet when sit
	// down on a seat. These chip are cash out back to wallet when
	// user stand up.
	Chip int

	// The state of the seat. It is a state machine.
	FSM *stateless.StateMachine

	// Record the starting time of sit out.
	SitOutStartTime time.Time

	// Remaining sit out duration that this user has. Used in sit out timer.
	SitOutDuration time.Duration

	// Cancel sit out timer. Typically used when user leaves sit out state.
	CancelSitOutTimer context.CancelFunc

	CancelReservingTimer context.CancelFunc

	ActionExtraDuration time.Duration

	// Used for determine whether this player should place BB for the
	// next game. If the player choose to not place BB (see WaitBB
	// option in PlaySetting), then can't join the next game. Unless
	// the player will get BB role in the next game, or he turn off
	// wait BB option to place a BB in the next game.
	ShouldPlaceBB bool

	// Record how many times player didn't action and wait until
	// time's up. If too many times, the player's will be forced to
	// stand up.
	CountIdleRounds int
}

func (s *SeatStatus) ToProto() *txpokergrpc.SeatStatus {
	return &txpokergrpc.SeatStatus{
		Uid:                  s.Uid.String(),
		Chip:                 int32(s.Chip),
		State:                s.FSM.MustState().(seatstatus.SeatStatusState).ToProto(),
		SitOutStartTimestamp: timestamppb.New(s.SitOutStartTime),
		SitOutDuration:       durationpb.New(s.SitOutDuration),
		ActionExtraDuration:  durationpb.New(s.ActionExtraDuration),
		ShouldPlaceBb:        s.ShouldPlaceBB,
	}
}

func (s *SeatStatus) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if s.Chip != 0 {
		enc.AddInt("Chip", s.Chip)
	}
	enc.AddString("State", s.FSM.MustState().(seatstatus.SeatStatusState).String())
	// enc.AddString("SitOutStart", s.SitOutStartTime.Format("2006-01-02 15:04:05"))
	// enc.AddDuration("SitOutDuration", s.SitOutDuration)
	// enc.AddDuration("ActionExtraDuration", s.ActionExtraDuration)
	enc.AddBool("ShouldPlaceBB", s.ShouldPlaceBB)
	enc.AddInt("CountIdleRounds", s.CountIdleRounds)
	return nil
}

type SeatStatusGroup struct {
	// seat id -> uid
	TableUids map[int]core.Uid

	Status map[core.Uid]*SeatStatus

	// uid -> top up chip
	TopUpQueue map[core.Uid]int

	// Only used for ToProto. For game logic, check
	// GameInfo.Setting.TableSize.
	TableSize int
}

func ProvideSeatStatusGroup() *SeatStatusGroup {
	return &SeatStatusGroup{
		TableUids:  make(map[int]core.Uid),
		Status:     make(map[core.Uid]*SeatStatus),
		TopUpQueue: make(map[core.Uid]int),
		TableSize:  9,
	}
}

func (g *SeatStatusGroup) ToProto() *txpokergrpc.SeatStatusGroup {
	msg := &txpokergrpc.SeatStatusGroup{
		TableUids:  make(map[int32]string),
		Status:     make(map[string]*txpokergrpc.SeatStatus),
		TopUpQueue: make(map[string]int32),
	}

	for seatId := 0; seatId < g.TableSize; seatId++ {
		if uid, ok := g.TableUids[seatId]; ok {
			msg.TableUids[int32(seatId)] = uid.String()
		} else {
			msg.TableUids[int32(seatId)] = ""
		}
	}

	for uid, seatStatus := range g.Status {
		msg.Status[uid.String()] = seatStatus.ToProto()
	}

	for uid, topUpChip := range g.TopUpQueue {
		msg.TopUpQueue[uid.String()] = int32(topUpChip)
	}

	return msg
}

func (g *SeatStatusGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if len(g.Status) > 0 {
		_ = enc.AddObject("Status", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			for uid, seatStatus := range g.Status {
				_ = enc.AddObject(uid.String(), seatStatus)
			}
			return nil
		}))
	}

	if len(g.TableUids) > 0 {
		_ = enc.AddObject("Table", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			for seatId, uid := range g.TableUids {
				enc.AddString(strconv.Itoa(seatId), uid.String())
			}
			return nil
		}))
	}

	if len(g.TopUpQueue) > 0 {
		_ = enc.AddObject("TopUpQueue", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
			for uid, topUpChip := range g.TopUpQueue {
				enc.AddInt(uid.String(), topUpChip)
			}
			return nil
		}))
	}

	enc.AddInt("TableSize", g.TableSize)
	return nil
}
