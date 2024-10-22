package model

import (
	gamegrpc "card-game-server-prototype/pkg/grpc/darkchessgrpc"
	"go.uber.org/zap/zapcore"
	"time"
)

type GameSetting struct {
	GameMetaUid     string
	TotalRound      int
	TurnSecond      time.Duration
	ExtraTurnSecond time.Duration
	AnteAmount      int
	EnterLimit      int
	WaterPct        int

	MaxRepeatMoves    int
	MaxChaseSamePiece int

	// Next version: additional rules
	IsCaptureRevealPieces   bool // 連吃
	IsCaptureUnrevealPieces bool // 連吃暗子
	IsCaptureUnrevealPiece  bool // 吃暗子
	HasRookRules            bool // 車走直
	HasBishopRules          bool // 馬走日
}

func (s *GameSetting) ToProto() *gamegrpc.GameSetting {
	return &gamegrpc.GameSetting{
		GameMetaUid:             s.GameMetaUid,
		TotalRoundCount:         int32(s.TotalRound),
		TurnSeconds:             int32(s.TurnSecond.Seconds()),
		ExtraTurnSeconds:        int32(s.ExtraTurnSecond.Seconds()),
		AnteAmount:              int32(s.AnteAmount),
		EnterLimit:              int32(s.EnterLimit),
		WaterPct:                int32(s.WaterPct),
		MaxChaseSamePieceCount:  int32(s.MaxChaseSamePiece),
		MaxRepeatMoves:          int32(s.MaxRepeatMoves),
		IsCaptureTurnedPieces:   s.IsCaptureRevealPieces,
		IsCaptureTurnDownPieces: s.IsCaptureUnrevealPieces,
		IsCaptureTurnDownPiece:  s.IsCaptureUnrevealPiece,
		HasRookRules:            s.HasRookRules,
		HasBishopRules:          s.HasBishopRules,
	}
}

func (s *GameSetting) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("GameMetaUid", s.GameMetaUid)
	enc.AddInt("TotalRound", s.TotalRound)
	enc.AddDuration("TurnSeconds", s.TurnSecond)
	enc.AddDuration("ExtraTurnSeconds", s.ExtraTurnSecond)
	enc.AddInt("AnteAmount", s.AnteAmount)
	enc.AddInt("EnterLimit", s.EnterLimit)
	enc.AddInt("WaterPct", s.WaterPct)
	enc.AddInt("MaxChaseSamePieceCount", s.MaxChaseSamePiece)
	enc.AddInt("MaxRepeatMoves", s.MaxRepeatMoves)
	enc.AddBool("IsCaptureTurnedPieces", s.IsCaptureRevealPieces)
	enc.AddBool("IsCaptureTurnDownPieces", s.IsCaptureUnrevealPieces)
	enc.AddBool("IsCaptureTurnDownPiece", s.IsCaptureUnrevealPiece)
	enc.AddBool("HasRookRules", s.HasRookRules)
	enc.AddBool("HasBishopRules", s.HasBishopRules)
	return nil
}
