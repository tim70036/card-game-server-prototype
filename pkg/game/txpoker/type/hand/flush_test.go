package hand

import (
	"errors"
	card2 "card-game-server-prototype/pkg/game/txpoker/type/card"
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"card-game-server-prototype/pkg/game/txpoker/type/suit"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"sort"
	"strings"
	"testing"
)

func TestFlush(t *testing.T) {
	type args struct {
		curCards, otherCards []string
	}
	type want struct {
		isCurLessThanOther, isEqual bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"equal",
			args{
				[]string{
					"JS",
					"9S",
					"8S",
					"7S",
					"4S",
					"3S",
					"4C",
				},
				[]string{
					"JS",
					"9S",
					"8S",
					"7S",
					"4S",
					"AC",
					"JC",
				},
			},
			want{
				isCurLessThanOther: false,
				isEqual:            true,
			},
		},
		{
			"less",
			args{
				[]string{
					"JS",
					"9S",
					"8S",
					"7S",
					"4S",
					"3S",
					"4C",
				},
				[]string{
					"AS",
					"9S",
					"8S",
					"7S",
					"4S",
					"3S",
					"4C",
				},
			},
			want{
				isCurLessThanOther: true,
				isEqual:            false,
			},
		},
		{
			"bigger",
			args{
				[]string{
					"JS",
					"9S",
					"8S",
					"7S",
					"4S",
					"3S",
					"4C",
				},
				[]string{
					"TS",
					"9S",
					"8S",
					"7S",
					"4S",
					"3S",
					"4C",
				},
			},
			want{
				isCurLessThanOther: false,
				isEqual:            false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var curCards, otherCards card2.CardList

			for _, v := range tt.args.curCards {
				curCards = append(curCards, newCardWithoutDeck(v))
			}

			for _, v := range tt.args.otherCards {
				otherCards = append(otherCards, newCardWithoutDeck(v))
			}

			cur, err := newFlush(curCards)
			if err != nil {
				t.Fatal(err)
			}

			other, err := newFlush(otherCards)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.want.isCurLessThanOther, cur.Less(other))
			assert.Equal(t, tt.want.isEqual, cur.Equal(other))
		})
	}
}

// released：不動原程式原則，直接 copy code 在這使用。
// reference: pkg/poker/txpoker/type/hand/hand.go:100
func newFlush(allCards card2.CardList) (Hand, error) {
	suitGroup := lo.GroupBy(allCards, func(c *card2.Card) suit.Suit { return c.Suit })
	flushSuitGroup := lo.PickBy(suitGroup, func(k suit.Suit, v []*card2.Card) bool { return len(v) >= 5 })

	if len(flushSuitGroup) > 0 {
		for _, flush := range flushSuitGroup {
			sort.Sort(sort.Reverse(card2.CardList(flush)))
		}

		var maxFlush card2.CardList = lo.MaxBy(lo.Values(flushSuitGroup), func(flush []*card2.Card, maxFlush []*card2.Card) bool {
			minLen := lo.Min([]int{len(flush), len(maxFlush)})
			for i := 0; i < minLen; i++ {
				if flush[i].Face != maxFlush[i].Face {
					return flush[i].Face > maxFlush[i].Face
				}
			}
			return len(flush) > len(maxFlush)
		})

		maxFlushFaces := lo.Map(maxFlush, func(c *card2.Card, _ int) face.Face { return c.Face })
		return &FlushHand{
			Faces:    maxFlushFaces,
			baseHand: baseHand{cards: maxFlush[:5]},
		}, nil
	}

	return nil, errors.New("no flush")
}

var suits = map[string]suit.Suit{
	"C": suit.Clubs,
	"D": suit.Diamonds,
	"H": suit.Hearts,
	"S": suit.Spades,
}

var faces = map[string]face.Face{
	"2": face.Two,
	"3": face.Three,
	"4": face.Four,
	"5": face.Five,
	"6": face.Six,
	"7": face.Seven,
	"8": face.Eight,
	"9": face.Nine,
	"T": face.Ten,
	"J": face.Jack,
	"Q": face.Queen,
	"K": face.King,
	"A": face.Ace,
}

func newCardWithoutDeck(c string) *card2.Card {
	args := strings.Split(c, "")
	if len(c) < 2 {
		panic("invalid card")
	}

	f, ok := faces[args[0]]
	if !ok {
		panic("invalid face")
	}

	s, ok := suits[args[1]]
	if !ok {
		panic("invalid suit")
	}

	return &card2.Card{
		Suit: s,
		Face: f,
	}
}
