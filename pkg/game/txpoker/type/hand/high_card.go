package hand

import (
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"github.com/samber/lo"
	"go.uber.org/zap/zapcore"
)

type HighCardHand struct {
	baseHand
	Faces []face.Face
}

func (h *HighCardHand) Type() HandType {
	return HighCard
}

func (h *HighCardHand) Less(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return h.Type() < otherHand.Type()
	}

	otherHighCardHand := otherHand.(*HighCardHand)
	minLen := lo.Min([]int{len(h.Faces), len(otherHighCardHand.Faces)})
	for i := 0; i < minLen; i++ {
		if h.Faces[i] != otherHighCardHand.Faces[i] {
			return h.Faces[i] < otherHighCardHand.Faces[i]
		}
	}

	return len(h.Faces) < len(otherHighCardHand.Faces)
}

func (h *HighCardHand) Equal(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return false
	}

	otherHighCardHand := otherHand.(*HighCardHand)
	minLen := lo.Min([]int{len(h.Faces), len(otherHighCardHand.Faces)})
	for i := 0; i < minLen; i++ {
		if h.Faces[i] != otherHighCardHand.Faces[i] {
			return false
		}
	}

	return len(h.Faces) == len(otherHighCardHand.Faces)
}

func (h *HighCardHand) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Type", h.Type().String())
	enc.AddArray("Faces", zapcore.ArrayMarshalerFunc(func(enc zapcore.ArrayEncoder) error {
		for _, f := range h.Faces {
			enc.AppendString(f.String())
		}
		return nil
	}))
	enc.AddString("Cards", h.cards.ToString())
	return nil
}
