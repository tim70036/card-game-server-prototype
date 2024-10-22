package config

import (
	"flag"
	"go.uber.org/zap/zapcore"
	"strings"
)

type LogConfig struct {
	Jslog      *bool
	level      *string
	debugGames *string
	infoGames  *string
	warnGames  *string
}

func (c *LogConfig) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddBool("Jslog", *c.Jslog)
	enc.AddString("level", *c.level)
	enc.AddString("debugGames", *c.debugGames)
	enc.AddString("infoGames", *c.infoGames)
	enc.AddString("warnGames", *c.warnGames)
	return nil
}

func (c *LogConfig) GetLevel(gameType string) zapcore.Level {
	if c.forceDebug(gameType) {
		return zapcore.DebugLevel
	}

	if c.forceInfo(gameType) {
		return zapcore.InfoLevel
	}

	if c.forceWarn(gameType) {
		return zapcore.WarnLevel
	}

	if c.level == nil {
		return zapcore.ErrorLevel
	}

	switch *c.level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	default:
		return zapcore.ErrorLevel
	}
}

func (c *LogConfig) forceDebug(gameType string) bool {
	return c.isInEnableGames(gameType, c.debugGames)
}

func (c *LogConfig) forceInfo(gameType string) bool {
	return c.isInEnableGames(gameType, c.infoGames)
}

func (c *LogConfig) forceWarn(gameType string) bool {
	return c.isInEnableGames(gameType, c.warnGames)
}

func (c *LogConfig) isInEnableGames(gameType string, enableGames *string) bool {
	if enableGames == nil || *enableGames == "" {
		return false
	}

	for _, enableGame := range strings.Split(*enableGames, ",") {
		if enableGame == gameType {
			return true
		}
	}

	return false
}

var LogCFG = &LogConfig{
	Jslog:      flag.Bool("jslog", false, "json logging (faster)"),
	level:      flag.String("log-level", "error", `log level for all games. parameters: 'debug', 'info', 'warn', 'error'`),
	debugGames: flag.String("debug-games", "", `force print debug log for some games. Use game_type, ex. '3' or '1,7'.`),
	infoGames:  flag.String("info-games", "", `force print info log for some games. Use 'game_type', see debug-games.`),
	warnGames:  flag.String("warn-games", "", `force print warn log for some games. Use 'game_type', see debug-games.`),
}
