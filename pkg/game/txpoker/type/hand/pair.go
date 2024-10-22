package hand

import (
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"github.com/samber/lo"
	"go.uber.org/zap/zapcore"
)

type PairHand struct {
	baseHand
	Face        face.Face
	KickerFaces []face.Face
}

func (h *PairHand) Type() HandType {
	return Pair
}

func (h *PairHand) Less(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return h.Type() < otherHand.Type()
	}

	otherPairHand := otherHand.(*PairHand)
	if h.Face != otherPairHand.Face {
		return h.Face < otherPairHand.Face
	}

	minLen := lo.Min([]int{len(h.KickerFaces), len(otherPairHand.KickerFaces)})
	for i := 0; i < minLen; i++ {
		if h.KickerFaces[i] != otherPairHand.KickerFaces[i] {
			return h.KickerFaces[i] < otherPairHand.KickerFaces[i]
		}
	}

	return len(h.KickerFaces) < len(otherPairHand.KickerFaces)
}

func (h *PairHand) Equal(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return false
	}

	otherPairHand := otherHand.(*PairHand)
	if h.Face != otherPairHand.Face {
		return false
	}

	minLen := lo.Min([]int{len(h.KickerFaces), len(otherPairHand.KickerFaces)})
	for i := 0; i < minLen; i++ {
		if h.KickerFaces[i] != otherPairHand.KickerFaces[i] {
			return false
		}
	}

	return len(h.KickerFaces) == len(otherPairHand.KickerFaces)
}

func (h *PairHand) MarshalLogObject(enc zapcore.ObjectEncoder) error {
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
