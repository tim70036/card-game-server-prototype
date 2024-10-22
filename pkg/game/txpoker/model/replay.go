package model

import (
	"encoding/json"
	"card-game-server-prototype/pkg/core"
	"card-game-server-prototype/pkg/game/txpoker/type/action"
	"card-game-server-prototype/pkg/game/txpoker/type/card"
	"card-game-server-prototype/pkg/game/txpoker/type/role"
	"card-game-server-prototype/pkg/game/txpoker/type/stage"
	"unicode"

	"github.com/samber/lo"
	"go.uber.org/zap/zapcore"
)

type PlayerRecord struct {
	Uid         core.Uid
	Role        role.Role
	SeatId      int
	PocketCards card.CardList

	// Seat status cheap right after round start.
	InitSeatStatusChip int

	// Seat status cheap right before round end (before winner appears).
	BeforeWinSeatStatusChip int

	WinChip int

	IsWinner    bool
	HasShowdown bool
}

func (r *PlayerRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"uid":                     r.Uid,
		"role":                    r.Role,
		"seatId":                  r.SeatId,
		"pocketCards":             r.PocketCards.ToHexStr(),
		"initSeatStatusChip":      r.InitSeatStatusChip,
		"beforeWinSeatStatusChip": r.BeforeWinSeatStatusChip,
		"winChip":                 r.WinChip,
		"isWinner":                r.IsWinner,
		"hasShowdown":             r.HasShowdown,
	})
}

func (r *PlayerRecord) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Uid", r.Uid.String())
	enc.AddString("Role", r.Role.String())
	enc.AddInt("SeatId", r.SeatId)
	enc.AddString("PocketCards", r.PocketCards.ToString())
	enc.AddInt("InitSeatStatusChip", r.InitSeatStatusChip)
	enc.AddInt("BeforeWinSeatStatusChip", r.BeforeWinSeatStatusChip)
	enc.AddInt("WinChip", r.WinChip)
	enc.AddBool("IsWinner", r.IsWinner)
	enc.AddBool("HasShowdown", r.HasShowdown)
	return nil
}

type Replay struct {
	RoundId        string
	ActionLog      map[stage.Stage][]action.ActionRecord
	StagePotChip   map[stage.Stage]int
	PlayerRecords  map[core.Uid]*PlayerRecord
	CommunityCards card.CardList
	PotChipSum     int
}

func ProvideReplay() *Replay {
	return &Replay{
		RoundId:        "Undefined",
		ActionLog:      make(map[stage.Stage][]action.ActionRecord),
		StagePotChip:   make(map[stage.Stage]int),
		PlayerRecords:  make(map[core.Uid]*PlayerRecord),
		CommunityCards: make(card.CardList, 0),
		PotChipSum:     0,
	}
}

func (r *Replay) MarshalJSON() ([]byte, error) {
	rawActionLog := lo.MapKeys(r.ActionLog, func(v []action.ActionRecord, k stage.Stage) string {
		str := k.String()
		for i, v := range str {
			return string(unicode.ToLower(v)) + str[i+1:]
		}
		return ""
	})

	rawStagePotChip := lo.MapKeys(r.StagePotChip, func(v int, k stage.Stage) string {
		str := k.String()
		for i, v := range str {
			return string(unicode.ToLower(v)) + str[i+1:]
		}
		return ""
	})

	return json.Marshal(map[string]any{
		"roundId":        r.RoundId,
		"actionLog":      rawActionLog,
		"stagePotChip":   rawStagePotChip,
		"playerRecords":  r.PlayerRecords,
		"communityCards": r.CommunityCards.ToHexStr(),
		"potChipSum":     r.PotChipSum,
	})
}

func (r *Replay) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("RoundId", r.RoundId)
	enc.AddObject("ActionLog", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for stage, actionLog := range r.ActionLog {
			enc.AddArray(stage.String(), zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
				for _, actionRecord := range actionLog {
					enc.AppendObject(actionRecord)
				}
				return nil
			}))
		}
		return nil
	}))

	enc.AddObject("StagePotChip", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for stage, potChip := range r.StagePotChip {
			enc.AddInt(stage.String(), potChip)
		}
		return nil
	}))

	enc.AddObject("PlayerRecords", zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		for uid, record := range r.PlayerRecords {
			enc.AddObject(uid.String(), record)
		}
		return nil
	}))

	enc.AddString("CommunityCards", r.CommunityCards.ToString())
	enc.AddInt("PotChipSum", r.PotChipSum)
	return nil
}
