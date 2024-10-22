package hand

import "go.uber.org/zap/zapcore"

type HandList []Hand

func (l HandList) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, card := range l {
		enc.AppendObject(card)
	}
	return nil
}

// Shallow clone.
func (l HandList) Clone() HandList {
	var clone = make(HandList, len(l))
	copy(clone, l)
	return clone
}

// Implement sort interface for sorting.
// Rule reference: https://en.wikipedia.org/wiki/List_of_poker_hand
func (l HandList) Len() int           { return len(l) }
func (l HandList) Less(i, j int) bool { return l[i].Less(l[j]) }
func (l HandList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
