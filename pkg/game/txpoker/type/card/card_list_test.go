package card

import (
	"errors"
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"card-game-server-prototype/pkg/game/txpoker/type/suit"
	"card-game-server-prototype/pkg/util"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGenerateShuffledDeck(t *testing.T) {
	assert := assert.New(t)
	logger := util.NewTestLogger()

	deck := GenerateShuffledDeck(1)
	assert.Len(deck, 52)
	logger.Info("Result", zap.Array("deck", deck))

	sort.Sort(deck)
	logger.Info("Result", zap.Array("deck", deck))

	sort.Sort(sort.Reverse(deck))
	logger.Info("Result", zap.Array("deck", deck))
}

func TestMatchAllPocketCards(t *testing.T) {
	type args struct {
		pocketCards, cards map[int]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "01",
			args: args{
				pocketCards: map[int]string{
					0: "KD0",
					1: "KC0",
				},
				cards: map[int]string{
					0: "6H0",
					1: "6S0",
					2: "6D0",
					3: "6C0",
					4: "KD0",
				},
			},
			want: false,
		},
		{
			name: "02",
			args: args{
				pocketCards: map[int]string{
					0: "5H0",
					1: "5S0",
				},
				cards: map[int]string{
					0: "5H0",
					1: "5S0",
					2: "5D0",
					3: "5C0",
					4: "9H0",
				},
			},
			want: true,
		},
		{
			name: "03",
			args: args{
				pocketCards: map[int]string{
					0: "5H0",
					1: "5C0",
				},
				cards: map[int]string{
					0: "5H0",
					1: "5C0",
					2: "5D0",
					3: "5S0",
					4: "TS0",
				},
			},
			want: true,
		},
		{
			name: "04",
			args: args{
				pocketCards: map[int]string{
					0: "8D0",
					1: "8S0",
				},
				cards: map[int]string{
					0: "8D0",
					1: "8S0",
					2: "8H0",
					3: "8C0",
					4: "7H0",
				},
			},
			want: true,
		},
		{
			name: "05",
			args: args{
				pocketCards: map[int]string{
					0: "6S0",
					1: "6C0",
				},
				cards: map[int]string{
					0: "6S0",
					1: "6C0",
					2: "6D0",
					3: "6H0",
					4: "AS0",
				},
			},
			want: true,
		},
		{
			name: "06",
			args: args{
				pocketCards: map[int]string{
					0: "QC0",
					1: "QH0",
				},
				cards: map[int]string{
					0: "QC0",
					1: "QH0",
					2: "QD0",
					3: "QS0",
					4: "AS0",
				},
			},
			want: true,
		},
		{
			name: "07",
			args: args{
				pocketCards: map[int]string{
					0: "8C0",
					1: "8H0",
				},
				cards: map[int]string{
					0: "8C0",
					1: "8H0",
					2: "8D0",
					3: "8S0",
					4: "QH0",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cards, err := parseStringsToCards(tt.args.cards)
			assert.NoError(t, err)

			pocketCards, err := parseStringsToCards(tt.args.pocketCards)
			assert.NoError(t, err)

			assert.Equal(t, tt.want, MatchAllPocketCards(cards, pocketCards), "not match all")
		})
	}
}

func parseStringsToCards(s map[int]string) (CardList, error) {
	var cards CardList

	for _, v := range s {
		p := strings.Split(v, "")

		if len(p) != 3 {
			return nil, errors.New("invalid card format")
		}

		deck, err := strconv.Atoi(p[2])
		if err != nil {
			return nil, err
		}

		cards = append(cards, &Card{
			Face: parseStrToFace(p[0]),
			Suit: parseStrToSuit(p[1]),
			Deck: deck,
		})
	}

	return cards, nil
}

var strToFace = map[string]face.Face{
	"A": face.Ace,
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
}

func parseStrToFace(f string) face.Face {
	if v, ok := strToFace[f]; ok {
		return v
	}
	return face.Undefined
}

var strToSuit = map[string]suit.Suit{
	"S": suit.Spades,
	"H": suit.Hearts,
	"D": suit.Diamonds,
	"C": suit.Clubs,
}

func parseStrToSuit(s string) suit.Suit {
	if v, ok := strToSuit[s]; ok {
		return v
	}
	return suit.Undefined
}
