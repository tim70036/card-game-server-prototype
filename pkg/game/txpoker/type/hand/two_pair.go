package hand

import (
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"go.uber.org/zap/zapcore"
)

type TwoPairHand struct {
	baseHand
	HighPairFace face.Face
	LowPairFace  face.Face
	KickerFace   face.Face
}

func (h *TwoPairHand) Type() HandType {
	return TwoPair
}

func (h *TwoPairHand) Less(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return h.Type() < otherHand.Type()
	}

	otherTwoPairHand := otherHand.(*TwoPairHand)
	if h.HighPairFace != otherTwoPairHand.HighPairFace {
		return h.HighPairFace < otherTwoPairHand.HighPairFace
	}

	if h.LowPairFace != otherTwoPairHand.LowPairFace {
		return h.LowPairFace < otherTwoPairHand.LowPairFace
	}

	return h.KickerFace < otherTwoPairHand.KickerFace
}

func (h *TwoPairHand) Equal(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return false
	}

	otherTwoPairHand := otherHand.(*TwoPairHand)
	if h.HighPairFace != otherTwoPairHand.HighPairFace {
		return false
	}

	if h.LowPairFace != otherTwoPairHand.LowPairFace {
		return false
	}

	return h.KickerFace == otherTwoPairHand.KickerFace
}

func (h *TwoPairHand) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Type", h.Type().String())
	enc.AddString("HighPairFace", h.HighPairFace.String())
	enc.AddString("LowPairFace", h.LowPairFace.String())
	enc.AddString("KickerFace", h.KickerFace.String())
	enc.AddString("Cards", h.cards.ToString())
	return nil
}
