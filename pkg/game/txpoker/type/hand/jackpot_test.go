package hand

import (
	card2 "card-game-server-prototype/pkg/game/txpoker/type/card"
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"card-game-server-prototype/pkg/game/txpoker/type/suit"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTriggerJackpot(t *testing.T) {
	type args struct {
		pocketCards, communityCards card2.CardList
	}
	type want struct {
		handType         Hand
		expectCards      card2.CardList
		isTriggerJackpot bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "01",
			args: args{
				pocketCards: card2.CardList{
					&card2.Card{Suit: suit.Spades, Face: face.Ace},
					&card2.Card{Suit: suit.Clubs, Face: face.Ace},
				},
				communityCards: card2.CardList{
					&card2.Card{Suit: suit.Spades, Face: face.Jack},
					&card2.Card{Suit: suit.Diamonds, Face: face.Jack},
					&card2.Card{Suit: suit.Hearts, Face: face.Jack},
					&card2.Card{Suit: suit.Hearts, Face: face.Six},
					&card2.Card{Suit: suit.Clubs, Face: face.Jack},
				},
			},
			want: want{
				handType: &FourOfAKindHand{},
				expectCards: card2.CardList{
					&card2.Card{Suit: suit.Spades, Face: face.Jack},
					&card2.Card{Suit: suit.Diamonds, Face: face.Jack},
					&card2.Card{Suit: suit.Hearts, Face: face.Jack},
					&card2.Card{Suit: suit.Clubs, Face: face.Jack},
					&card2.Card{Suit: suit.Spades, Face: face.Ace},
				},
				isTriggerJackpot: false, // 沒有用到兩張手牌
			},
		},
		{
			name: "02",
			args: args{
				pocketCards: card2.CardList{
					&card2.Card{Suit: suit.Spades, Face: face.Jack},
					&card2.Card{Suit: suit.Diamonds, Face: face.Jack},
				},
				communityCards: card2.CardList{
					&card2.Card{Suit: suit.Clubs, Face: face.Ace},
					&card2.Card{Suit: suit.Spades, Face: face.Ace},
					&card2.Card{Suit: suit.Hearts, Face: face.Jack},
					&card2.Card{Suit: suit.Hearts, Face: face.Six},
					&card2.Card{Suit: suit.Clubs, Face: face.Jack},
				},
			},
			want: want{
				handType: &FourOfAKindHand{},
				expectCards: card2.CardList{
					&card2.Card{Suit: suit.Spades, Face: face.Jack},
					&card2.Card{Suit: suit.Diamonds, Face: face.Jack},
					&card2.Card{Suit: suit.Hearts, Face: face.Jack},
					&card2.Card{Suit: suit.Clubs, Face: face.Jack},
					&card2.Card{Suit: suit.Clubs, Face: face.Ace},
				},
				isTriggerJackpot: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand, err := New(tt.args.pocketCards, tt.args.communityCards)

			assert.NoError(t, err)
			assert.IsType(t, tt.want.handType, hand)
			assert.ElementsMatch(t, tt.want.expectCards, hand.Cards())

			var isTriggerJackpot bool

			if hand.Type() == RoyalFlush || hand.Type() == StraightFlush {
				// 完成牌型時玩家必須用到兩張手牌
				isTriggerJackpot = card2.MatchAllPocketCards(hand.Cards(), tt.args.pocketCards)
			}

			if hand.Type() == FourOfAKind {
				// 完成牌型時玩家必須用到兩張手牌，且鐵支時必須為手裡對
				isTriggerJackpot = card2.MatchAllPocketCards(hand.Cards(), tt.args.pocketCards) &&
					tt.args.pocketCards[0].Face == tt.args.pocketCards[1].Face
			}

			assert.Equal(t, tt.want.isTriggerJackpot, isTriggerJackpot)
		})
	}
}
