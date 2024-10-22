package card

import (
	"fmt"
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"card-game-server-prototype/pkg/game/txpoker/type/suit"
	"strconv"
)

func FromHex(rawCard string, deck int) (*Card, error) {
	if len(rawCard) != 2 {
		return nil, fmt.Errorf("decode rawCard %+v, each card length must be 2, but is: %v", rawCard, len(rawCard))
	}

	rawSuit, err := strconv.ParseInt(string(rawCard[0]), 16, 0)
	if err != nil {
		return nil, fmt.Errorf("decode rawCard %+v, parse raw suit %+v error: %w", rawCard, rawSuit, err)
	}

	if rawSuit < int64(suit.Clubs) || rawSuit > int64(suit.Spades) {
		return nil, fmt.Errorf("decode rawCard %+v, parse raw suit %+v out of range", rawCard, rawSuit)
	}

	rawFace, err := strconv.ParseInt(string(rawCard[1]), 16, 0)
	if err != nil {
		return nil, fmt.Errorf("decode rawCard %+v, parse raw face %+v error: %w", rawCard, rawFace, err)
	}

	if rawFace < int64(face.Two) || rawFace > int64(face.Ace) {
		return nil, fmt.Errorf("decode rawCard %+v, parse raw face %+v out of range", rawCard, rawFace)
	}

	return &Card{
		Suit: suit.Suit(rawSuit),
		Face: face.Face(rawFace),
		Deck: deck,
	}, nil
}
