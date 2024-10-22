package model

import (
	"card-game-server-prototype/pkg/core"
	txpokerconstant "card-game-server-prototype/pkg/game/txpoker/constant"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/constant"
	"card-game-server-prototype/pkg/game/zoomtxpoker/pool/type/participant"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"github.com/qmuntal/stateless"
	"go.uber.org/zap/zapcore"
	"time"
)

// Participant represents a ticket that a player in the matchmaking pool.
type Participant struct {
	Uid core.Uid

	// The state of the participant. It is a state machine.
	FSM *stateless.StateMachine

	// The chip that participant buy in.
	Chip int

	// The chip that participant queued for top up.
	QueuedTopUpChip int

	IdleAt int64

	CountIdleRounds int
}

func (p *Participant) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Uid", p.Uid.String())
	enc.AddInt("Chip", p.Chip)
	enc.AddString("State", p.FSM.MustState().(participant.State).String())
	enc.AddTime("IdleAt", time.Unix(p.IdleAt, 0))
	enc.AddInt("QueuedTopUpChip", p.QueuedTopUpChip)
	enc.AddInt("CountIdleRounds", p.CountIdleRounds)
	return nil
}

func (p *Participant) ToProto() *txpokergrpc.Participant {
	return &txpokergrpc.Participant{
		Uid:      p.Uid.String(),
		Chip:     int32(p.Chip),
		State:    p.FSM.MustState().(participant.State).ToProto(),
		HasTopUp: p.QueuedTopUpChip > 0,
	}
}

func (p *Participant) IsIdlingTimeout(now time.Time) bool {
	if p.IdleAt == 0 {
		return false
	}
	return now.Sub(time.Unix(p.IdleAt, 0)) >= constant.IdlingTimeoutDuration
}

func (p *Participant) IsIdleRoundsReachMax() bool {
	return p.CountIdleRounds >= txpokerconstant.MaxIdleRounds
}

type ParticipantGroup struct {
	Data map[core.Uid]*Participant
}

func ProvideParticipantGroup() *ParticipantGroup {
	return &ParticipantGroup{
		Data: make(map[core.Uid]*Participant),
	}
}

func (g *ParticipantGroup) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for uid, part := range g.Data {
		if err := enc.AddObject(uid.String(), part); err != nil {
			return err
		}
	}
	return nil
}
