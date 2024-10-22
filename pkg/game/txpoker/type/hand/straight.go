package hand

import (
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"go.uber.org/zap/zapcore"
)

type StraightHand struct {
	baseHand
	HighFace face.Face
}

func (h *StraightHand) Type() HandType {
	return Straight
}

func (h *StraightHand) Less(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return h.Type() < otherHand.Type()
	}

	otherStraightHand := otherHand.(*StraightHand)
	return h.HighFace < otherStraightHand.HighFace
}

func (h *StraightHand) Equal(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return false
	}

	otherStraightHand := otherHand.(*StraightHand)
	return h.HighFace == otherStraightHand.HighFace
}

func (h *StraightHand) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Type", h.Type().String())
	enc.AddString("HighFace", h.HighFace.String())
	enc.AddString("Cards", h.cards.ToString())
	return nil
}
