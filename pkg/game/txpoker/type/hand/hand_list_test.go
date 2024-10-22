package hand

import (
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSameTypeSort(t *testing.T) {
	assert := assert.New(t)

	sortedHands := HandList{
		&StraightHand{HighFace: face.King},
		&StraightHand{HighFace: face.Three},
		&StraightHand{HighFace: face.Jack},
		&StraightHand{HighFace: face.Seven},
	}

	sort.Sort(sortedHands)

	for idx, face := range []face.Face{face.Three, face.Seven, face.Jack, face.King} {
		assert.Equal(face, sortedHands[idx].(*StraightHand).HighFace)
	}

	sortedHands = HandList{
		&PairHand{Face: face.Three, KickerFaces: []face.Face{face.Five, face.Four, face.Two}},
		&PairHand{Face: face.Three, KickerFaces: []face.Face{face.Five, face.Four, face.Three}},
	}

	sort.Sort(sortedHands)
	for idx, face := range []face.Face{face.Two, face.Three} {
		assert.Equal(face, sortedHands[idx].(*PairHand).KickerFaces[2])
	}
}

func TestSort(t *testing.T) {
	assert := assert.New(t)

	hands := HandList{
		&StraightHand{HighFace: face.King},
		&PairHand{Face: face.Three, KickerFaces: []face.Face{face.Two, face.Four, face.Five}},
		&RoyalFlushHand{HighFace: face.Ace},
		&TwoPairHand{HighPairFace: face.Jack, LowPairFace: face.Seven, KickerFace: face.Two},
		&HighCardHand{Faces: []face.Face{face.Two, face.Four, face.Five, face.Six, face.Seven}},
		&ThreeOfAKindHand{Face: face.Seven, KickerFaces: []face.Face{face.Two, face.Four}},
		&FullHouseHand{TripletFace: face.Jack, PairFace: face.Seven},
		&FlushHand{Faces: []face.Face{face.Two, face.Four, face.Five, face.Six, face.Seven}},
		&StraightFlushHand{HighFace: face.Jack},
		&FourOfAKindHand{Face: face.Seven, KickerFace: face.Two},
	}

	sort.Sort(hands)

	idx := 0
	for i := HighCard; i <= RoyalFlush; i++ {
		assert.Equal(i, hands[idx].Type())
		idx++
	}
}
