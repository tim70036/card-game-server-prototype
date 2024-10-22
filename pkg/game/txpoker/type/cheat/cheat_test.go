package cheat

import (
	"fmt"
	card2 "card-game-server-prototype/pkg/game/txpoker/type/card"
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"card-game-server-prototype/pkg/game/txpoker/type/suit"
	"strconv"
	"strings"
	"testing"
)

func TestCardToHex(t *testing.T) {
	type args struct {
		cards card2.CardList
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCardToHex",
			args: args{
				cards: card2.CardList{
					&card2.Card{Suit: suit.Spades, Face: face.Ten},
					&card2.Card{Suit: suit.Spades, Face: face.Jack},
					&card2.Card{Suit: suit.Spades, Face: face.Queen},
					&card2.Card{Suit: suit.Diamonds, Face: face.Jack},
					&card2.Card{Suit: suit.Clubs, Face: face.Four},
					&card2.Card{Suit: suit.Spades, Face: face.Ace},
					&card2.Card{Suit: suit.Spades, Face: face.King},
					&card2.Card{Suit: suit.Hearts, Face: face.Jack},
					&card2.Card{Suit: suit.Clubs, Face: face.Jack},
				},
			},
		},
	}
	for _, tt := range tests {
		var arr []string
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.args.cards {
				hex := strings.ToUpper(toHexStr(v))
				t.Logf("%v -> %s", v, hex)

				ca, err := card2.FromHex(hex, 0)
				if err != nil {
					t.Fatalf(err.Error())
				}

				t.Logf(" >> %v", ca)
				arr = append(arr, hex)
			}
			t.Log(arr)
		})
	}
}

func toHexStr(cardVal *card2.Card) string {
	return fmt.Sprintf("%v%v", strconv.FormatInt(int64(cardVal.Suit), 16), strconv.FormatInt(int64(cardVal.Face), 16))
}
