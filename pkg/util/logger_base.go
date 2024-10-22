package util

import (
	prettyconsole "github.com/thessem/zap-prettyconsole"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func prettyBaseLogger() *zap.Logger {
	ecfg := prettyconsole.NewEncoderConfig()
	// ecfg.EncodeTime = prettyconsole.DefaultTimeEncoder("2006-01-02T15:04:05.000Z0700") // ISO8601
	ecfg.EncodeTime = prettyconsole.DefaultTimeEncoder("15:04:05.00000") // Who the fuck will want time.Kitchen?

	consoleCore := zapcore.NewCore(
		prettyconsole.NewEncoder(ecfg),
		os.Stdout,
		LoggerLevel,
	)

	// workingDir, _ := os.Getwd()
	// logName := "dump.log"
	// writeSyncer := zapcore.AddSync(&lumberjack.Logger{
	// 	Filename: filepath.Join(workingDir, logName),
	// })

	// fileCore := zapcore.NewCore(
	// 	zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
	// 	writeSyncer,
	// 	LoggerLevel,
	// )

	// combineCores := zapcore.NewTee(consoleCore, fileCore)
	combineCores := zapcore.NewTee(consoleCore)
	options := []zap.Option{
		zap.AddStacktrace(zapcore.WarnLevel), // We want stack trace for warning message.
	}

	logger := zap.New(combineCores, options...)
	return logger
}

func consoleBaseLogger() *zap.Logger {
	// See the documentation for Config and zapcore.EncoderConfig for all the
	// available options.
	cfg := zap.Config{
		Level:            LoggerLevel,
		Development:      false,
		Encoding:         "console",
		Sampling:         nil, // zap.Counter 效能問題，暫時不用。
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			// Keys can be anything except the empty string.
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "name",
			CallerKey:      zapcore.OmitKey,
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger := zap.Must(cfg.Build())
	return logger
}

func jsonBaseLogger() *zap.Logger {
	// See the documentation for Config and zapcore.EncoderConfig for all the
	// available options.
	cfg := zap.Config{
		Level:            LoggerLevel,
		Development:      false,
		Encoding:         "json",
		Sampling:         nil, // zap.Counter 效能問題，暫時不用。
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			// Keys can be anything except the empty string.
			TimeKey:        "timestamp", // todo: 外層有重複記錄，看看是否可以改善。
			LevelKey:       "level",
			NameKey:        "name",
			CallerKey:      zapcore.OmitKey,
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger := zap.Must(cfg.Build())
	return logger
}
