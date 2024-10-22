package card

import (
	"fmt"
	"card-game-server-prototype/pkg/game/txpoker/type/face"
	"card-game-server-prototype/pkg/game/txpoker/type/suit"
	"card-game-server-prototype/pkg/grpc/txpokergrpc"
)

type Card struct {
	Suit suit.Suit
	Face face.Face
	Deck int
}

func (c *Card) ToProto() *txpokergrpc.Card {
	return &txpokergrpc.Card{
		Suit: c.Suit.ToProto(),
		Face: c.Face.ToProto(),
		Deck: int32(c.Deck),
	}
}

func toLogString(s suit.Suit, f face.Face, deck int) string {
	// 暫時只有一副 deck，不特別 log
	return fmt.Sprintf("%s%s", s.ShortString(), f.String())
}
