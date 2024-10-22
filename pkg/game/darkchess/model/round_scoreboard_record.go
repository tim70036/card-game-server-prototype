package model

import (
	"go.uber.org/zap/zapcore"
)

type RoundScoreboardRecords struct {
	Data []*RoundScoreboard
}

func ProvideRoundScoreboardRecords() *RoundScoreboardRecords {
	return &RoundScoreboardRecords{
		Data: []*RoundScoreboard{},
	}
}

func (s *RoundScoreboardRecords) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, scoreboard := range s.Data {
		_ = enc.AppendObject(scoreboard)
	}
	return nil
}
