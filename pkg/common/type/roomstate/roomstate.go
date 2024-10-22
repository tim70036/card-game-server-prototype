package roomstate

type RoomState int

const (
	Running RoomState = 0
	Closing RoomState = 1
	Finish  RoomState = 2
	Invalid RoomState = 3
)

func (s RoomState) Int() int {
	return int(s)
}
