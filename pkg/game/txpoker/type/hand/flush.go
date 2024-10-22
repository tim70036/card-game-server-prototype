package hand

import (
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"github.com/samber/lo"
	"go.uber.org/zap/zapcore"
)

type FlushHand struct {
	baseHand
	Faces []face.Face
}

func (h *FlushHand) Type() HandType {
	return Flush
}

func (h *FlushHand) Less(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return h.Type() < otherHand.Type()
	}

	otherFlushHand := otherHand.(*FlushHand)
	minLen := lo.Min([]int{len(h.Faces), len(otherFlushHand.Faces)})
	for i := 0; i < minLen; i++ {
		if h.Faces[i] != otherFlushHand.Faces[i] {
			return h.Faces[i] < otherFlushHand.Faces[i]
		}
	}

	return false
}

func (h *FlushHand) Equal(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return false
	}

	otherFlushHand := otherHand.(*FlushHand)
	minLen := lo.Min([]int{len(h.Faces), len(otherFlushHand.Faces)})
	for i := 0; i < minLen; i++ {
		if h.Faces[i] != otherFlushHand.Faces[i] {
			return false
		}
	}

	return true
}

func (h *FlushHand) MarshalLogObject(enc zapcore.ObjectEncoder) error {
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
