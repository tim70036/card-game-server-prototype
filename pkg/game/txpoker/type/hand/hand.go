package hand

import (
	"fmt"
	card2 "card-game-server-prototype/pkg/game/txpoker/type/card"
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"card-game-server-prototype/pkg/game/txpoker/type/suit"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
	"sort"

	"github.com/samber/lo"
	"go.uber.org/zap/zapcore"
)

type Hand interface {
	Type() HandType
	Cards() card2.CardList
	Less(otherHand Hand) bool
	Equal(otherHand Hand) bool
	MarshalLogObject(enc zapcore.ObjectEncoder) error
	mustEmbedBaseHand()
}

func ToProto(hand Hand) *txpokergrpc.PokerHand {
	return &txpokergrpc.PokerHand{
		Type:  hand.Type().ToProto(),
		Cards: hand.Cards().ToProto(),
	}
}

func New(pocketCards card2.CardList, communityCards card2.CardList) (Hand, error) {
	if pocketCards == nil || communityCards == nil {
		return nil, fmt.Errorf("pocket cards and community cards must not be nil")
	}

	if len(pocketCards) != 2 {
		return nil, fmt.Errorf("pocket cards must be 2, but got %d", len(pocketCards))
	}

	if len(communityCards) != 5 {
		return nil, fmt.Errorf("community cards must be 5, but got %d", len(communityCards))
	}

	allCards := append(pocketCards.Clone(), communityCards...)
	allStraights := evalAllPossibleStraights(allCards)

	// Royal Flush & Flush Straight
	allStraightFlush := lo.Filter(allStraights, func(straight card2.CardList, _ int) bool {
		return lo.EveryBy(straight, func(c *card2.Card) bool { return c.Suit == straight[0].Suit })
	})

	if len(allStraightFlush) > 0 {
		maxStraightFlush := lo.MaxBy(allStraightFlush, func(s card2.CardList, maxS card2.CardList) bool {
			return s[0].Face > maxS[0].Face
		})

		highFace := maxStraightFlush[0].Face
		if highFace == face.Ace {
			return &RoyalFlushHand{HighFace: highFace, baseHand: baseHand{cards: maxStraightFlush[:5]}}, nil
		} else {
			return &StraightFlushHand{HighFace: highFace, baseHand: baseHand{cards: maxStraightFlush[:5]}}, nil
		}
	}

	faceGroup := lo.GroupBy(allCards, func(c *card2.Card) face.Face { return c.Face })

	pairFaceGroup := lo.PickBy(faceGroup, func(k face.Face, v []*card2.Card) bool { return len(v) >= 2 })
	tripletFaceGroup := lo.PickBy(faceGroup, func(k face.Face, v []*card2.Card) bool { return len(v) >= 3 })
	fourFaceGroup := lo.PickBy(faceGroup, func(k face.Face, v []*card2.Card) bool { return len(v) >= 4 })

	// Four of a Kind
	if len(fourFaceGroup) > 0 {
		maxFourFace := lo.MaxBy(lo.Keys(fourFaceGroup), func(f face.Face, maxF face.Face) bool { return f > maxF })
		kickerCandidates := lo.Reject(allCards, func(c *card2.Card, _ int) bool { return lo.Contains(fourFaceGroup[maxFourFace], c) })
		kicker := lo.MaxBy(kickerCandidates, func(c *card2.Card, maxC *card2.Card) bool { return c.Face > maxC.Face })

		cards := append(card2.CardList(fourFaceGroup[maxFourFace]).Clone()[:4], kicker)
		return &FourOfAKindHand{
			Face:       maxFourFace,
			KickerFace: kicker.Face,
			baseHand:   baseHand{cards: cards},
		}, nil
	}

	// Full House
	if len(tripletFaceGroup) > 0 {
		maxTripletFace := lo.MaxBy(lo.Keys(tripletFaceGroup), func(f face.Face, maxF face.Face) bool { return f > maxF })
		if len(lo.Without(lo.Keys(pairFaceGroup), maxTripletFace)) > 0 {
			maxPairFace := lo.MaxBy(lo.Without(lo.Keys(pairFaceGroup), maxTripletFace), func(f face.Face, maxF face.Face) bool { return f > maxF })
			cards := append(card2.CardList(tripletFaceGroup[maxTripletFace]).Clone()[:3], pairFaceGroup[maxPairFace][:2]...)
			return &FullHouseHand{
				TripletFace: maxTripletFace,
				PairFace:    maxPairFace,
				baseHand:    baseHand{cards: cards},
			}, nil
		}

	}

	// Flush
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

	// Straight
	if len(allStraights) > 0 {
		maxStraight := lo.MaxBy(allStraights, func(s card2.CardList, maxS card2.CardList) bool {
			return s[0].Face > maxS[0].Face
		})

		highFace := maxStraight[0].Face
		return &StraightHand{
			HighFace: highFace,
			baseHand: baseHand{cards: maxStraight[:5]},
		}, nil
	}

	// ThreeOfAKind
	if len(tripletFaceGroup) > 0 {
		maxTripletFace := lo.MaxBy(lo.Keys(tripletFaceGroup), func(f face.Face, maxF face.Face) bool { return f > maxF })
		kickerCandidates := lo.Reject(allCards, func(c *card2.Card, _ int) bool { return lo.Contains(tripletFaceGroup[maxTripletFace], c) })
		sort.Sort(sort.Reverse(card2.CardList(kickerCandidates)))
		kickers := kickerCandidates[:2]
		kickerFaces := lo.Map(kickers, func(c *card2.Card, _ int) face.Face { return c.Face })

		cards := append(card2.CardList(tripletFaceGroup[maxTripletFace]).Clone()[:3], kickers...)
		return &ThreeOfAKindHand{
			Face:        maxTripletFace,
			KickerFaces: kickerFaces,
			baseHand:    baseHand{cards: cards},
		}, nil
	}

	// Two Pair
	if len(pairFaceGroup) > 1 {
		highPairFace := lo.MaxBy(lo.Keys(pairFaceGroup), func(f face.Face, maxF face.Face) bool { return f > maxF })
		lowPairFace := lo.MaxBy(lo.Without(lo.Keys(pairFaceGroup), highPairFace), func(f face.Face, maxF face.Face) bool { return f > maxF })

		usedCards := append(card2.CardList(pairFaceGroup[highPairFace]).Clone(), pairFaceGroup[lowPairFace]...)
		kickerCandidates := lo.Reject(allCards, func(c *card2.Card, _ int) bool { return lo.Contains(usedCards, c) })
		kicker := lo.MaxBy(kickerCandidates, func(c *card2.Card, maxC *card2.Card) bool { return c.Face > maxC.Face })

		cards := append(card2.CardList(pairFaceGroup[highPairFace]).Clone()[:2], pairFaceGroup[lowPairFace][:2]...)
		cards = append(cards, kicker)
		return &TwoPairHand{
			HighPairFace: highPairFace,
			LowPairFace:  lowPairFace,
			KickerFace:   kicker.Face,
			baseHand:     baseHand{cards: cards},
		}, nil
	}

	// Pair
	if len(pairFaceGroup) > 0 {
		maxPairFace := lo.MaxBy(lo.Keys(pairFaceGroup), func(f face.Face, maxF face.Face) bool { return f > maxF })
		kickerCandidates := lo.Reject(allCards, func(c *card2.Card, _ int) bool { return lo.Contains(pairFaceGroup[maxPairFace], c) })
		sort.Sort(sort.Reverse(card2.CardList(kickerCandidates)))
		kickers := kickerCandidates[:3]
		kickerFaces := lo.Map(kickers, func(c *card2.Card, _ int) face.Face { return c.Face })

		cards := append(card2.CardList(pairFaceGroup[maxPairFace]).Clone()[:2], kickers...)
		return &PairHand{
			Face:        maxPairFace,
			KickerFaces: kickerFaces,
			baseHand:    baseHand{cards: cards},
		}, nil
	}

	// High Card
	sort.Sort(sort.Reverse(allCards))
	highCards := allCards[:5]
	highCardFaces := lo.Map(highCards, func(c *card2.Card, _ int) face.Face { return c.Face })
	return &HighCardHand{
		Faces:    highCardFaces,
		baseHand: baseHand{cards: highCards},
	}, nil
}

func evalAllPossibleStraights(cards card2.CardList) []card2.CardList {
	allStraights := []card2.CardList{}
	faceGroup := lo.GroupBy(cards, func(c *card2.Card) face.Face { return c.Face })

	for face := range faceGroup {
		dfsStraight(faceGroup, face, card2.CardList{}, &allStraights)
	}
	return allStraights
}

func dfsStraight(faceGroup map[face.Face][]*card2.Card, curFace face.Face, curStraight card2.CardList, allStraights *[]card2.CardList) {
	if len(curStraight) >= 5 {
		*allStraights = append(*allStraights, curStraight)
		return
	}

	if _, ok := faceGroup[curFace]; !ok {
		return
	}

	for _, card := range faceGroup[curFace] {
		nextStraight := append(curStraight.Clone(), card)

		// Under high rules, an ace can rank either high (as in A♦ K♣
		// Q♣ J♦ 10♠, an ace-high straight) or low (as in 5♣ 4♦ 3♥ 2♥
		// A♠, a five-high straight), but cannot simultaneously rank
		// both high and low
		nextFace := lo.Ternary(
			curFace == face.Two && len(nextStraight) == 4,
			face.Ace,
			curFace-1,
		)

		dfsStraight(faceGroup, nextFace, nextStraight, allStraights)
	}
}

type baseHand struct {
	cards card2.CardList
}

func (h *baseHand) Cards() card2.CardList { return h.cards }
func (h *baseHand) mustEmbedBaseHand()    {}
