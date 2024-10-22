package stage

type StageTrigger int

const (
	UndefinedTrigger StageTrigger = 0
	NextStageTrigger StageTrigger = 1
)

func (t StageTrigger) String() string {
	name, ok := triggerNames[t]
	if !ok {
		return "UndefinedTrigger"
	}
	return name
}

var triggerNames = map[StageTrigger]string{
	UndefinedTrigger: "UndefinedTrigger",
	NextStageTrigger: "NextStageTrigger",
}
