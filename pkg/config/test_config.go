package config

import (
	"flag"
	"strings"

	"go.uber.org/zap/zapcore"
)

type TestConfig struct {
	LocalMode       *bool
	singleModeGames *string
	AutopilotMode   *bool
	NoAuth          *bool
	cheatModeGame   *string
	CheatData       *string
}

func (c *TestConfig) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddBool("LocalMode", *c.LocalMode)
	enc.AddString("singleModeGames", *c.singleModeGames)
	enc.AddBool("AutopilotMode", *c.AutopilotMode)
	enc.AddBool("NoAuth", *c.NoAuth)
	enc.AddString("cheatModeGame", *c.cheatModeGame)
	enc.AddString("CheatData", *c.CheatData)
	return nil
}

func (c *TestConfig) EnableSingle(gameType string) bool {
	if c.singleModeGames == nil || *c.singleModeGames == "" {
		return false
	}

	for _, singleModeGame := range strings.Split(*c.singleModeGames, ",") {
		if singleModeGame == gameType {
			return true
		}
	}

	return false
}

func (c *TestConfig) EnableCheatMode(gameType string) bool {
	if c.cheatModeGame == nil || *c.cheatModeGame == "" {
		return false
	}

	return gameType == *c.cheatModeGame
}

var TestCFG = &TestConfig{
	LocalMode:       flag.Bool("local-mode", false, "running in local mode, no server integration needed (no agones integration, no calling other server's api)"),
	singleModeGames: flag.String("single-mode-games", "", `run game with only a single user (server will fake other users' behavior). Use game-types, ex. '3' or '1,7'`),
	AutopilotMode:   flag.Bool("autopilot-mode", false, "make all user use auto-mode."),
	NoAuth:          flag.Bool("no-auth", false, "run game without auth (server will not verify jwt)"),
	cheatModeGame:   flag.String("cheat-mode-game", "", "run game in cheat mode, game will probably configure preset data using cheat data"),
	CheatData:       flag.String("cheat-data", "", "cheat data to apply when running in cheat mode"),
}
