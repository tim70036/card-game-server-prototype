package pot

import "go.uber.org/zap/zapcore"

type Winner struct {
	Chip int

	// RawProfit = win chip - bet chip
	RawProfit int

	Water int

	JackpotWater int
}

func (w *Winner) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("Chip", w.Chip)
	enc.AddInt("RawProfit", w.RawProfit)
	enc.AddInt("Water", w.Water)
	if w.JackpotWater != 0 {
		enc.AddInt("JackpotWater", w.JackpotWater)
	}
	return nil
}
