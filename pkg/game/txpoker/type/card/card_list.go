package card

import (
	"fmt"
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"card-game-server-prototype/pkg/game/txpoker/type/suit"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"card-game-server-prototype/pkg/util"
	"github.com/samber/lo"
	"strconv"
)

type CardList []*Card

func (l CardList) ToProto() []*txpokergrpc.Card {
	return lo.Map(l, func(c *Card, _ int) *txpokergrpc.Card { return c.ToProto() })
}

func (l CardList) ToHexStr() string {
	var buffer []string
	for _, card := range l {
		if card.Face == face.Ace {
			buffer = append(buffer, fmt.Sprintf("%v%v", strconv.FormatInt(int64(card.Suit), 16), strconv.FormatInt(1, 16))) // Ace is 1 in frontend.
		} else {
			buffer = append(buffer, fmt.Sprintf("%v%v", strconv.FormatInt(int64(card.Suit), 16), strconv.FormatInt(int64(card.Face), 16)))
		}
	}
	return util.JoinStrings(buffer)
}

func (l CardList) ToString() string {
	return util.JoinStrings(lo.Map(l, func(c *Card, _ int) string {
		return toLogString(c.Suit, c.Face, c.Deck)
	}), ",")
}

// Shallow clone.
func (l CardList) Clone() CardList {
	var clone = make(CardList, len(l))
	copy(clone, l)
	return clone
}

// Implement sort interface for sorting.
func (l CardList) Len() int           { return len(l) }
func (l CardList) Less(i, j int) bool { return l[i].Face < l[j].Face }
func (l CardList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func GenerateShuffledDeck(deck int) CardList {
	cards := make(CardList, 0)
	for s := suit.Clubs; s <= suit.Spades; s++ {
		for f := face.Two; f <= face.Ace; f++ {
			cards = append(cards, &Card{
				Suit: s,
				Face: f,
				Deck: deck,
			})
		}
	}

	return lo.Shuffle(cards)
}

func MatchAllPocketCards(cards, pocketCards CardList) bool {
	var matchCount int
	for _, pocketCard := range pocketCards {
		for _, c := range cards {
			// 未來規則變更的確認：目前只有一副牌，deck 都是 0。如果未來有多副牌，這裡要確認是不是需要比對 deck。
			if pocketCard.Suit == c.Suit && pocketCard.Face == c.Face && pocketCard.Deck == c.Deck {
				matchCount++
				break
			}
		}
	}

	return matchCount > 0 && matchCount == len(pocketCards)
}
