package stage

import (
	"fmt"
	role2 "card-game-server-prototype/pkg/game/txpoker/type/role"
	"github.com/qmuntal/stateless"
	"github.com/samber/lo"
)

type Stage int

const (
	UndefinedStage Stage = 0
	AnteStage      Stage = 1
	PreFlopStage   Stage = 2
	FlopStage      Stage = 3
	TurnStage      Stage = 4
	RiverStage     Stage = 5
	ShowdownStage  Stage = 6
)

func NewBetStageFSM() *stateless.StateMachine {
	fsm := stateless.NewStateMachine(AnteStage)

	fsm.Configure(AnteStage).
		Permit(NextStageTrigger, PreFlopStage)

	fsm.Configure(PreFlopStage).
		Permit(NextStageTrigger, FlopStage)

	fsm.Configure(FlopStage).
		Permit(NextStageTrigger, TurnStage)

	fsm.Configure(TurnStage).
		Permit(NextStageTrigger, RiverStage)

	fsm.Configure(RiverStage).
		Permit(NextStageTrigger, ShowdownStage)

	fsm.Configure(ShowdownStage)

	return fsm
}

func (s Stage) String() string {
	name, ok := stateNames[s]
	if !ok {
		return "UndefinedStage"
	}
	return name
}

var stateNames = map[Stage]string{
	UndefinedStage: "UndefinedStage",
	AnteStage:      "AnteStage",
	PreFlopStage:   "PreFlopStage",
	FlopStage:      "FlopStage",
	TurnStage:      "TurnStage",
	RiverStage:     "RiverStage",
	ShowdownStage:  "ShowdownStage",
}

func FirstActionRole(playerCount int, s Stage) (role2.Role, error) {
	roles, err := role2.GetRoles(playerCount)
	if err != nil || len(roles) < 2 {
		return role2.Undefined, fmt.Errorf("cannot found first action role with player count: %d: %w", playerCount, err)
	}

	return lo.Ternary(
		s == PreFlopStage,
		roles[1],
		role2.SB,
	), nil
}
