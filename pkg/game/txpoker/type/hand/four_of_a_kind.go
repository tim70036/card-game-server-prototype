package hand

import (
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"go.uber.org/zap/zapcore"
)

type FourOfAKindHand struct {
	baseHand
	Face       face.Face
	KickerFace face.Face
}

func (h *FourOfAKindHand) Type() HandType {
	return FourOfAKind
}

func (h *FourOfAKindHand) Less(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return h.Type() < otherHand.Type()
	}

	otherFourOfAKindHand := otherHand.(*FourOfAKindHand)
	if h.Face != otherFourOfAKindHand.Face {
		return h.Face < otherFourOfAKindHand.Face
	}

	return h.KickerFace < otherFourOfAKindHand.KickerFace
}

func (h *FourOfAKindHand) Equal(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return false
	}

	otherFourOfAKindHand := otherHand.(*FourOfAKindHand)
	return h.Face == otherFourOfAKindHand.Face && h.KickerFace == otherFourOfAKindHand.KickerFace
}

func (h *FourOfAKindHand) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Type", h.Type().String())
	enc.AddString("Face", h.Face.String())
	enc.AddString("KickerFace", h.KickerFace.String())
	enc.AddString("Cards", h.cards.ToString())
	return nil
}
