package hand

import (
	card2 "card-game-server-prototype/pkg/game/txpoker/type/card"
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"card-game-server-prototype/pkg/game/txpoker/type/suit"
	"card-game-server-prototype/pkg/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestEvalAllPossibleStraights(t *testing.T) {
	assert := assert.New(t)
	logger := util.NewTestLogger()

	cards := card2.CardList{
		&card2.Card{Suit: suit.Clubs, Face: face.Ace},
		&card2.Card{Suit: suit.Diamonds, Face: face.Ace},
		&card2.Card{Suit: suit.Diamonds, Face: face.Two},
		&card2.Card{Suit: suit.Diamonds, Face: face.Three},
		&card2.Card{Suit: suit.Spades, Face: face.Three},
		&card2.Card{Suit: suit.Diamonds, Face: face.Four},
		&card2.Card{Suit: suit.Diamonds, Face: face.Five},
		&card2.Card{Suit: suit.Hearts, Face: face.Five},
	}
	allStraights := evalAllPossibleStraights(cards)
	assert.Len(allStraights, 8)

	cards = card2.CardList{
		&card2.Card{Suit: suit.Diamonds, Face: face.Ace},
		&card2.Card{Suit: suit.Diamonds, Face: face.Two},
		&card2.Card{Suit: suit.Diamonds, Face: face.Three},
		&card2.Card{Suit: suit.Spades, Face: face.Three},
		&card2.Card{Suit: suit.Diamonds, Face: face.Four},
		&card2.Card{Suit: suit.Diamonds, Face: face.Five},
		&card2.Card{Suit: suit.Diamonds, Face: face.Six},
		&card2.Card{Suit: suit.Diamonds, Face: face.Seven},
	}
	allStraights = evalAllPossibleStraights(cards)
	assert.Len(allStraights, 6)

	cards = card2.CardList{
		&card2.Card{Suit: suit.Spades, Face: face.Ten},
		&card2.Card{Suit: suit.Diamonds, Face: face.Jack},
		&card2.Card{Suit: suit.Diamonds, Face: face.Queen},
		&card2.Card{Suit: suit.Diamonds, Face: face.King},
		&card2.Card{Suit: suit.Diamonds, Face: face.Ace},
	}
	allStraights = evalAllPossibleStraights(cards)
	assert.Len(allStraights, 1)

	cards = card2.CardList{
		&card2.Card{Suit: suit.Spades, Face: face.Eight},
		&card2.Card{Suit: suit.Spades, Face: face.Nine},
		&card2.Card{Suit: suit.Spades, Face: face.Ten},
		&card2.Card{Suit: suit.Diamonds, Face: face.Jack},
		&card2.Card{Suit: suit.Diamonds, Face: face.Queen},
		&card2.Card{Suit: suit.Diamonds, Face: face.King},
		&card2.Card{Suit: suit.Diamonds, Face: face.Ace},
	}
	allStraights = evalAllPossibleStraights(cards)
	assert.Len(allStraights, 3)

	cards = card2.CardList{
		&card2.Card{Suit: suit.Diamonds, Face: face.Jack},
		&card2.Card{Suit: suit.Diamonds, Face: face.Queen},
		&card2.Card{Suit: suit.Diamonds, Face: face.King},
		&card2.Card{Suit: suit.Diamonds, Face: face.Ace},
		&card2.Card{Suit: suit.Diamonds, Face: face.Two},
	}
	allStraights = evalAllPossibleStraights(cards)
	assert.Len(allStraights, 0)

	cards = card2.CardList{
		&card2.Card{Suit: suit.Diamonds, Face: face.Queen},
		&card2.Card{Suit: suit.Diamonds, Face: face.King},
		&card2.Card{Suit: suit.Diamonds, Face: face.Ace},
		&card2.Card{Suit: suit.Diamonds, Face: face.Two},
		&card2.Card{Suit: suit.Diamonds, Face: face.Three},
	}
	allStraights = evalAllPossibleStraights(cards)
	assert.Len(allStraights, 0)

	cards = card2.CardList{
		&card2.Card{Suit: suit.Diamonds, Face: face.Ace},
		&card2.Card{Suit: suit.Diamonds, Face: face.Two},
		&card2.Card{Suit: suit.Diamonds, Face: face.Three},
		&card2.Card{Suit: suit.Diamonds, Face: face.Five},
		&card2.Card{Suit: suit.Diamonds, Face: face.Six},
		&card2.Card{Suit: suit.Diamonds, Face: face.Seven},
		&card2.Card{Suit: suit.Diamonds, Face: face.Eight},
		&card2.Card{Suit: suit.Diamonds, Face: face.Jack},
	}
	allStraights = evalAllPossibleStraights(cards)
	assert.Len(allStraights, 0)

	for _, straight := range allStraights {
		logger.Info("Straight", zap.Array("straight", straight))
	}
}

func TestStraights(t *testing.T) {
	assert := assert.New(t)

	pocketCards := card2.CardList{
		&card2.Card{Suit: suit.Clubs, Face: face.Ace},
		&card2.Card{Suit: suit.Spades, Face: face.Ace},
	}

	communityCards := card2.CardList{
		&card2.Card{Suit: suit.Clubs, Face: face.Ten},
		&card2.Card{Suit: suit.Clubs, Face: face.Three},
		&card2.Card{Suit: suit.Clubs, Face: face.Queen},
		&card2.Card{Suit: suit.Clubs, Face: face.Jack},
		&card2.Card{Suit: suit.Clubs, Face: face.King},
	}

	hand, err := New(pocketCards, communityCards)
	assert.NoError(err)
	assert.IsType(&RoyalFlushHand{}, hand)
	assert.Equal(face.Ace, hand.(*RoyalFlushHand).HighFace)

	pocketCards = card2.CardList{
		&card2.Card{Suit: suit.Spades, Face: face.Ace},
		&card2.Card{Suit: suit.Diamonds, Face: face.Ace},
	}

	communityCards = card2.CardList{
		&card2.Card{Suit: suit.Spades, Face: face.Three},
		&card2.Card{Suit: suit.Clubs, Face: face.Three},
		&card2.Card{Suit: suit.Spades, Face: face.Four},
		&card2.Card{Suit: suit.Spades, Face: face.Five},
		&card2.Card{Suit: suit.Spades, Face: face.Two},
	}

	hand, err = New(pocketCards, communityCards)
	assert.NoError(err)
	assert.IsType(&StraightFlushHand{}, hand)
	assert.Equal(face.Five, hand.(*StraightFlushHand).HighFace)

	pocketCards = card2.CardList{
		&card2.Card{Suit: suit.Spades, Face: face.Ace},
		&card2.Card{Suit: suit.Diamonds, Face: face.Ace},
	}

	communityCards = card2.CardList{
		&card2.Card{Suit: suit.Clubs, Face: face.Two},
		&card2.Card{Suit: suit.Clubs, Face: face.Three},
		&card2.Card{Suit: suit.Spades, Face: face.Three},
		&card2.Card{Suit: suit.Clubs, Face: face.Four},
		&card2.Card{Suit: suit.Clubs, Face: face.Five},
	}

	hand, err = New(pocketCards, communityCards)
	assert.NoError(err)
	assert.IsType(&StraightHand{}, hand)
	assert.Equal(face.Five, hand.(*StraightHand).HighFace)

	pocketCards = card2.CardList{
		&card2.Card{Suit: suit.Clubs, Face: face.Three},
		&card2.Card{Suit: suit.Clubs, Face: face.Five},
	}

	communityCards = card2.CardList{
		&card2.Card{Suit: suit.Clubs, Face: face.Seven},
		&card2.Card{Suit: suit.Diamonds, Face: face.Jack},
		&card2.Card{Suit: suit.Spades, Face: face.Three},
		&card2.Card{Suit: suit.Clubs, Face: face.Four},
		&card2.Card{Suit: suit.Spades, Face: face.Six},
	}

	hand, err = New(pocketCards, communityCards)
	assert.NoError(err)
	assert.IsType(&StraightHand{}, hand)
	assert.Equal(face.Seven, hand.(*StraightHand).HighFace)

	pocketCards = card2.CardList{
		&card2.Card{Suit: suit.Clubs, Face: face.Two},
		&card2.Card{Suit: suit.Diamonds, Face: face.Three},
	}

	communityCards = card2.CardList{
		&card2.Card{Suit: suit.Clubs, Face: face.Three},
		&card2.Card{Suit: suit.Diamonds, Face: face.Four},
		&card2.Card{Suit: suit.Spades, Face: face.Seven},
		&card2.Card{Suit: suit.Hearts, Face: face.Five},
		&card2.Card{Suit: suit.Spades, Face: face.Six},
	}

	hand, err = New(pocketCards, communityCards)
	assert.NoError(err)
	assert.IsType(&StraightHand{}, hand)
	assert.Equal(face.Seven, hand.(*StraightHand).HighFace)

	pocketCards = card2.CardList{
		&card2.Card{Suit: suit.Clubs, Face: face.Ace},
		&card2.Card{Suit: suit.Clubs, Face: face.King},
	}

	communityCards = card2.CardList{
		&card2.Card{Suit: suit.Clubs, Face: face.Queen},
		&card2.Card{Suit: suit.Diamonds, Face: face.Two},
		&card2.Card{Suit: suit.Spades, Face: face.Three},
		&card2.Card{Suit: suit.Clubs, Face: face.Four},
		&card2.Card{Suit: suit.Spades, Face: face.Six},
	}

	hand, err = New(pocketCards, communityCards)
	assert.NoError(err)
	assert.IsType(&HighCardHand{}, hand)
}

func TestFullHouse(t *testing.T) {
	assert := assert.New(t)

	pocketCards := card2.CardList{
		&card2.Card{Suit: suit.Clubs, Face: face.Two},
		&card2.Card{Suit: suit.Diamonds, Face: face.Three},
	}

	communityCards := card2.CardList{
		&card2.Card{Suit: suit.Clubs, Face: face.Three},
		&card2.Card{Suit: suit.Diamonds, Face: face.Two},
		&card2.Card{Suit: suit.Spades, Face: face.Seven},
		&card2.Card{Suit: suit.Hearts, Face: face.Three},
		&card2.Card{Suit: suit.Spades, Face: face.Two},
	}

	hand, err := New(pocketCards, communityCards)
	assert.NoError(err)
	assert.IsType(&FullHouseHand{}, hand)
	assert.Equal(face.Three, hand.(*FullHouseHand).TripletFace)
	assert.Equal(face.Two, hand.(*FullHouseHand).PairFace)

	pocketCards = card2.CardList{
		&card2.Card{Suit: suit.Clubs, Face: face.Ace},
		&card2.Card{Suit: suit.Diamonds, Face: face.Ace},
	}

	communityCards = card2.CardList{
		&card2.Card{Suit: suit.Clubs, Face: face.Jack},
		&card2.Card{Suit: suit.Diamonds, Face: face.Ace},
		&card2.Card{Suit: suit.Spades, Face: face.Jack},
		&card2.Card{Suit: suit.Hearts, Face: face.Jack},
		&card2.Card{Suit: suit.Spades, Face: face.Three},
	}

	hand, err = New(pocketCards, communityCards)
	assert.NoError(err)
	assert.IsType(&FullHouseHand{}, hand)
	assert.Equal(face.Ace, hand.(*FullHouseHand).TripletFace)
	assert.Equal(face.Jack, hand.(*FullHouseHand).PairFace)
}

func TestFourOfAKind(t *testing.T) {
	type args struct {
		pocketCards, communityCards card2.CardList
	}
	type want struct {
		face, kickerFace face.Face
	}
	tests := []struct {
		name string
		args args
		want
	}{
		{
			name: "FourOfAKind 01",
			args: args{
				pocketCards: card2.CardList{
					&card2.Card{Suit: suit.Clubs, Face: face.Two},
					&card2.Card{Suit: suit.Diamonds, Face: face.Three},
				},
				communityCards: card2.CardList{
					&card2.Card{Suit: suit.Clubs, Face: face.Three},
					&card2.Card{Suit: suit.Diamonds, Face: face.Two},
					&card2.Card{Suit: suit.Spades, Face: face.Seven},
					&card2.Card{Suit: suit.Hearts, Face: face.Three},
					&card2.Card{Suit: suit.Spades, Face: face.Three},
				},
			},
			want: want{
				face:       face.Three,
				kickerFace: face.Seven,
			},
		},
		{
			name: "FourOfAKind 02",
			args: args{
				pocketCards: card2.CardList{
					&card2.Card{Suit: suit.Clubs, Face: face.Three},
					&card2.Card{Suit: suit.Diamonds, Face: face.Three},
				},
				communityCards: card2.CardList{
					&card2.Card{Suit: suit.Clubs, Face: face.Ace},
					&card2.Card{Suit: suit.Diamonds, Face: face.Ace},
					&card2.Card{Suit: suit.Spades, Face: face.Ace},
					&card2.Card{Suit: suit.Hearts, Face: face.Ace},
					&card2.Card{Suit: suit.Spades, Face: face.Three},
				},
			},
			want: want{
				face:       face.Ace,
				kickerFace: face.Three,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand, err := New(tt.args.pocketCards, tt.args.communityCards)

			assert.NoError(t, err)
			assert.IsType(t, &FourOfAKindHand{}, hand)
			assert.Equal(t, tt.want.face, hand.(*FourOfAKindHand).Face)
			assert.Equal(t, tt.want.kickerFace, hand.(*FourOfAKindHand).KickerFace)
		})
	}
}

func TestStraightFlush(t *testing.T) {
	type args struct {
		pocketCards, communityCards card2.CardList
	}
	type want struct {
		expectCards card2.CardList
		highFace    face.Face
	}
	tests := []struct {
		name string
		args args
		want
	}{
		{
			name: "RoyalFlush 01",
			args: args{
				pocketCards: card2.CardList{
					&card2.Card{Suit: suit.Spades, Face: face.Ace},
					&card2.Card{Suit: suit.Spades, Face: face.Five},
				},
				communityCards: card2.CardList{
					&card2.Card{Suit: suit.Spades, Face: face.Six},
					&card2.Card{Suit: suit.Spades, Face: face.Four},
					&card2.Card{Suit: suit.Spades, Face: face.Three},
					&card2.Card{Suit: suit.Diamonds, Face: face.Eight},
					&card2.Card{Suit: suit.Spades, Face: face.Two},
				},
			},
			want: want{
				expectCards: card2.CardList{
					&card2.Card{Suit: suit.Spades, Face: face.Six},
					&card2.Card{Suit: suit.Spades, Face: face.Five},
					&card2.Card{Suit: suit.Spades, Face: face.Four},
					&card2.Card{Suit: suit.Spades, Face: face.Three},
					&card2.Card{Suit: suit.Spades, Face: face.Two},
				},
				highFace: face.Six,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand, err := New(tt.args.pocketCards, tt.args.communityCards)

			assert.NoError(t, err)
			assert.IsType(t, &StraightFlushHand{}, hand)
			assert.ElementsMatch(t, tt.want.expectCards, hand.Cards())
			assert.Equal(t, tt.want.highFace, hand.(*StraightFlushHand).HighFace)
		})
	}
}
