package init

import (
	"custom-in-memory-db/internal/server/cmd"
	"log/slog"
	"os"
)

var logLevelMap = map[string]slog.Level{
	"debug": -4,
	"info":  0,
	"warn":  4,
	"error": 8,
}

func Logger(conf cmd.Config) *slog.Logger {
	if validateLoggingConf(conf) {
		return loggerWithConf(conf)
	} else {
		return defaultLogger()
	}
}

func validateLoggingConf(conf cmd.Config) bool {
	if conf.Logging.Level != "debug" &&
		conf.Logging.Level != "info" &&
		conf.Logging.Level != "warn" &&
		conf.Logging.Level != "error" {
		return false
	}
	if conf.Logging.Format != "text" &&
		conf.Logging.Format != "json" {
		return false
	}
	return true
}

func loggerWithConf(conf cmd.Config) *slog.Logger {
	var logLevel = new(slog.LevelVar)
	logLevel.Set(logLevelMap[conf.Logging.Level])

	if conf.Logging.Format == "text" {
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
}

func defaultLogger() *slog.Logger {
	var logLevel = new(slog.LevelVar)
	logLevel.Set(logLevelMap["error"])

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
}
