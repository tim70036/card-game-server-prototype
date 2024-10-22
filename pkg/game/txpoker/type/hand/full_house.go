package hand

import (
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"go.uber.org/zap/zapcore"
)

type FullHouseHand struct {
	baseHand
	TripletFace face.Face
	PairFace    face.Face
}

func (h *FullHouseHand) Type() HandType {
	return FullHouse
}

func (h *FullHouseHand) Less(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return h.Type() < otherHand.Type()
	}

	otherFullHouseHand := otherHand.(*FullHouseHand)
	if h.TripletFace != otherFullHouseHand.TripletFace {
		return h.TripletFace < otherFullHouseHand.TripletFace
	}

	return h.PairFace < otherFullHouseHand.PairFace
}

func (h *FullHouseHand) Equal(otherHand Hand) bool {
	if h.Type() != otherHand.Type() {
		return false
	}

	otherFullHouseHand := otherHand.(*FullHouseHand)
	return h.TripletFace == otherFullHouseHand.TripletFace && h.PairFace == otherFullHouseHand.PairFace
}

func (h *FullHouseHand) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Type", h.Type().String())
	enc.AddString("TripletFace", h.TripletFace.String())
	enc.AddString("PairFace", h.PairFace.String())
	enc.AddString("Cards", h.cards.ToString())
	return nil
}
