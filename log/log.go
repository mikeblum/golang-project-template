package log

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/knadh/koanf"
	"github.com/mikeblum/golang-project-template/conf"
)

// logging configuration

const (
	jsonLog      = "JSON"
	envLogLevel  = "LOG_LEVEL"
	envLogFormat = "LOG_FORMAT"
	logAttrError = "error"
)

type Log struct {
	*slog.Logger
}

// NewLog - configure logging
func NewLog() *Log {
	var cfg *koanf.Koanf
	var err error
	var log *slog.Logger
	var cfgLogLevel slog.Level
	logLevels := map[string]slog.Level{
		slog.LevelDebug.String(): slog.LevelDebug,
	}
	if cfg, err = conf.NewConf(conf.Provider("")); err != nil {
		// default to INFO
		cfgLogLevel = slog.LevelInfo
	} else {
		cfgLogLevel = logLevels[cfg.String(envLogLevel)]
	}
	handler := &slog.HandlerOptions{
		Level: cfgLogLevel,
	}
	logFormat := conf.GetEnv(envLogFormat, "")
	if strings.EqualFold(logFormat, jsonLog) {
		// review html escape
		log = slog.New(slog.NewJSONHandler(os.Stdout, handler))
	} else {
		// implement in slog
		// logrus.SetFormatter(&logrus.TextFormatter{
		// 	FullTimestamp: true,
		// 	// timestamp with millisecond precision
		// 	TimestampFormat: "Jan _2 15:04:05.00",
		// 	ForceColors:     true,
		// })
		log = slog.New(slog.NewTextHandler(os.Stdout, handler))
	}

	// set default logger
	slog.SetDefault(log)
	return &Log{log}
}

// Wrappers
// https://pkg.go.dev/golang.org/x/exp/slog#hdr-Wrapping_output_methods

func (log *Log) WithError(err error) *Log {
	ctx := log.With(logAttrError, err)
	return &Log{ctx}
}

func (log *Log) Infof(format string, args ...any) {
	log.Info(fmt.Sprintf(format, args...))
}
