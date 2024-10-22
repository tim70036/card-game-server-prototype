package fold

type ShowType int

const (
	ShowNone  ShowType = iota // 0b00
	ShowRight                 // 0b01
	ShowLeft                  // 0b10
	ShowBoth                  // 0b11
)
