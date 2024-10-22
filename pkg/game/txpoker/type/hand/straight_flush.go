package hand

import (
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"go.uber.org/zap/zapcore"
)

type StraightFlushHand struct {
	baseHand
	HighFace face.Face
}

func (h *StraightFlushHand) Type() HandType {
	return StraightFlush
}

func (h *StraightFlushHand) Less(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return h.Type() < otherHand.Type()
	}

	otherStraightFlushHand := otherHand.(*StraightFlushHand)
	return h.HighFace < otherStraightFlushHand.HighFace
}

func (h *StraightFlushHand) Equal(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return false
	}

	otherStraightFlushHand := otherHand.(*StraightFlushHand)
	return h.HighFace == otherStraightFlushHand.HighFace
}

func (h *StraightFlushHand) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Type", h.Type().String())
	enc.AddString("HighFace", h.HighFace.String())
	enc.AddString("Cards", h.cards.ToString())
	return nil
}
