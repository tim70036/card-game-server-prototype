package hand

import (
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"github.com/samber/lo"
	"go.uber.org/zap/zapcore"
)

type ThreeOfAKindHand struct {
	baseHand
	Face        face.Face
	KickerFaces []face.Face
}

func (h *ThreeOfAKindHand) Type() HandType {
	return ThreeOfAKind
}

func (h *ThreeOfAKindHand) Less(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return h.Type() < otherHand.Type()
	}

	otherThreeOfAKindHand := otherHand.(*ThreeOfAKindHand)
	if h.Face != otherThreeOfAKindHand.Face {
		return h.Face < otherThreeOfAKindHand.Face
	}

	minLen := lo.Min([]int{len(h.KickerFaces), len(otherThreeOfAKindHand.KickerFaces)})
	for i := 0; i < minLen; i++ {
		if h.KickerFaces[i] != otherThreeOfAKindHand.KickerFaces[i] {
			return h.KickerFaces[i] < otherThreeOfAKindHand.KickerFaces[i]
		}
	}

	return len(h.KickerFaces) < len(otherThreeOfAKindHand.KickerFaces)
}

func (h *ThreeOfAKindHand) Equal(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return false
	}

	otherThreeOfAKindHand := otherHand.(*ThreeOfAKindHand)
	if h.Face != otherThreeOfAKindHand.Face {
		return false
	}

	minLen := lo.Min([]int{len(h.KickerFaces), len(otherThreeOfAKindHand.KickerFaces)})
	for i := 0; i < minLen; i++ {
		if h.KickerFaces[i] != otherThreeOfAKindHand.KickerFaces[i] {
			return false
		}
	}

	return len(h.KickerFaces) == len(otherThreeOfAKindHand.KickerFaces)
}

func (h *ThreeOfAKindHand) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Type", h.Type().String())
	enc.AddString("Face", h.Face.String())
	enc.AddArray("KickerFaces", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		for _, f := range h.KickerFaces {
			enc.AppendString(f.String())
		}
		return nil
	}))
	enc.AddString("Cards", h.cards.ToString())
	return nil
}
