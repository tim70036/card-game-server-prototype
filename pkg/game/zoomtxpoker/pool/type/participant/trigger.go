package participant

type StateTrigger int

const (
	UndefinedTrigger StateTrigger = 0
	EnterGameTrigger StateTrigger = 1
	ExitGameTrigger  StateTrigger = 2
	BuyInTrigger     StateTrigger = 3
	ExitMatchTrigger StateTrigger = 4
	SuccessTrigger   StateTrigger = 5
	FailedTrigger    StateTrigger = 6
	ExitingTrigger   StateTrigger = 7
	BackTrigger      StateTrigger = 8
)

var triggerNames = map[StateTrigger]string{
	UndefinedTrigger: "UndefinedTrigger",
	EnterGameTrigger: "EnterGameTrigger",
	ExitGameTrigger:  "ExitGameTrigger",
	BuyInTrigger:     "BuyInTrigger",
	ExitMatchTrigger: "ExitMatchTrigger",
	SuccessTrigger:   "SuccessTrigger",
	FailedTrigger:    "FailedTrigger",
	ExitingTrigger:   "ExitingTrigger",
	BackTrigger:      "BackTrigger",
}

func (t StateTrigger) String() string {
	name, ok := triggerNames[t]
	if !ok {
		return "UndefinedTrigger"
	}
	return name
}
