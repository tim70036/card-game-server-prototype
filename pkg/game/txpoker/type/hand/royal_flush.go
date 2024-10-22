package hand

import (
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"go.uber.org/zap/zapcore"
)

type RoyalFlushHand struct {
	baseHand
	HighFace face.Face
}

func (h *RoyalFlushHand) Type() HandType {
	return RoyalFlush
}

func (h *RoyalFlushHand) Less(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return h.Type() < otherHand.Type()
	}

	otherRoyalFlushHand := otherHand.(*RoyalFlushHand)
	return h.HighFace < otherRoyalFlushHand.HighFace
}

func (h *RoyalFlushHand) Equal(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return false
	}

	otherRoyalFlushHand := otherHand.(*RoyalFlushHand)
	return h.HighFace == otherRoyalFlushHand.HighFace
}

func (h *RoyalFlushHand) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Type", h.Type().String())
	enc.AddString("HighFace", h.HighFace.String())
	enc.AddString("Cards", h.cards.ToString())
	return nil
}
