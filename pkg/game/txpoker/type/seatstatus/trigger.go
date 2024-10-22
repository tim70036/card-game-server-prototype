package seatstatus

type SeatStatusStateTrigger int

const (
	UndefinedTrigger    SeatStatusStateTrigger = 0
	StandUpTrigger      SeatStatusStateTrigger = 1
	SitDownTrigger      SeatStatusStateTrigger = 2
	SitOutTrigger       SeatStatusStateTrigger = 3
	RoundStartedTrigger SeatStatusStateTrigger = 4
	RoundEndTrigger     SeatStatusStateTrigger = 5
	BuyInTrigger        SeatStatusStateTrigger = 6
	SuccessTrigger      SeatStatusStateTrigger = 7
	FailedTrigger       SeatStatusStateTrigger = 8
	TimeoutTrigger      SeatStatusStateTrigger = 9
)

func (t SeatStatusStateTrigger) String() string {
	name, ok := tirggerNames[t]
	if !ok {
		return "UndefinedTrigger"
	}
	return name
}

var tirggerNames = map[SeatStatusStateTrigger]string{
	UndefinedTrigger:    "UndefinedTrigger",
	StandUpTrigger:      "StandUpTrigger",
	SitDownTrigger:      "SitDownTrigger",
	SitOutTrigger:       "SitOutTrigger",
	RoundStartedTrigger: "RoundStartedTrigger",
	RoundEndTrigger:     "RoundEndTrigger",
	BuyInTrigger:        "BuyInTrigger",
	SuccessTrigger:      "SuccessTrigger",
	FailedTrigger:       "FailedTrigger",
	TimeoutTrigger:      "TimeoutTrigger",
}
