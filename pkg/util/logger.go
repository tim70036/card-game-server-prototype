package util

import (
	"card-game-server-prototype/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// LoggerLevel Allow changing log level at run time.
	LoggerLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
)

// DebugField Helper function that returns nothing if log level is above debug.
// This can help when you want to ignore some fields in higher level.
// EX: logger.Info("hello", zap.String("1", "1"), util.DebugField(zap.String("2", "2")))
// In this case, "2" will only be printed when log level is debug.
func DebugField(field zap.Field) zap.Field {
	if LoggerLevel.Level() > zapcore.DebugLevel {
		return zap.Skip()
	}
	return field
}

type LoggerFactory struct {
	baseLogger *zap.Logger
}

func (f *LoggerFactory) Create(name string) *zap.Logger {
	return f.baseLogger.Named(name)
}

func ProvideLoggerFactory(logCFG *config.LogConfig, defaultFields []zap.Field) *LoggerFactory {
	baseLogger := jsonBaseLogger()
	if !*(logCFG.Jslog) {
		baseLogger = prettyBaseLogger()
	}

	return &LoggerFactory{
		baseLogger: baseLogger.With(defaultFields...),
	}
}

func NewTestLogger() *zap.Logger {
	return prettyBaseLogger().Named("Test")
}
