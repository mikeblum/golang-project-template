package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/knadh/koanf"
)

// a logrus-inspired implementation of log/slog
// https://github.com/golang/example/blob/master/slog-handler-guide/README.md

const (
	// default to text loggging format
	LogFormatText = iota
	LogFormatJSON = iota
	textFormat    = "TEXT"
	jsonFormat    = "JSON"
	envLogLevel   = "LOG_LEVEL"
	envLogFormat  = "LOG_FORMAT"
	logAttrError  = "error"

	// default to INFO
	defaultLevel = slog.LevelInfo
	// slog levels are based off OpenTelemetry
	// OpenTelemetry levels: https://opentelemetry.io/docs/specs/otel/logs/data-model/#field-severitynumber
	// log/slog levels: https://cs.opensource.google/go/x/exp/+/d852ddb8:slog/level.go
	LevelTrace      = slog.Level(-8)
	LevelTraceLabel = "TRACE"
	LevelTraceName  = "DEBUG-4"
	LevelFatal      = slog.Level(12)
	LevelFatalLabel = "FATAL"
	LevelFatalName  = "ERROR+4"
)

// drop-in for slog.HandlerOptions
type Options struct {
	Conf        *koanf.Koanf
	AddSource   bool
	Level       slog.Leveler
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
	// LOG_FORMAT supports either text (default) or json
	Format FormatLevel
	// log output defaults to os.Stdout
	Out io.Writer
}

func DefaultOptions() Options {
	return Options{
		// default to INFO
		Level: slog.LevelInfo,
		// default to TEXT
		Format: LogFormatText,
		Out:    os.Stdout,
	}
}

type FormatLevel int

func Levels() []string {
	return []string{
		LevelTraceName,
		slog.LevelDebug.String(),
		slog.LevelInfo.String(),
		slog.LevelWarn.String(),
		slog.LevelError.String(),
		LevelFatalName,
	}
}

func Level(level string) (slog.Level, error) {
	if level == "" {
		return defaultLevel, nil
	}

	logLevels := map[string]slog.Level{
		LevelTraceName:           LevelTrace,
		LevelTraceLabel:          LevelTrace,
		slog.LevelDebug.String(): slog.LevelDebug,
		slog.LevelInfo.String():  slog.LevelInfo,
		slog.LevelWarn.String():  slog.LevelWarn,
		slog.LevelError.String(): slog.LevelError,
		LevelFatalName:           LevelFatal,
		LevelFatalLabel:          LevelFatal,
	}
	var lvl slog.Level
	var ok bool
	if lvl, ok = logLevels[strings.ToUpper(level)]; !ok {
		// warn that this log level isn't supported
		return defaultLevel, fmt.Errorf("Log level: %s not supported. Levels are %s", level, strings.Join(Levels(), ", "))
	}
	return lvl, nil
}

func Formats() []string {
	return []string{
		textFormat, //
		jsonFormat,
	}
}

func Format(format string) FormatLevel {
	logFormatsMap := map[string]FormatLevel{
		// default to text format for the empty zero case
		textFormat: LogFormatText,
		jsonFormat: LogFormatJSON,
	}
	var logFmt FormatLevel
	var ok bool
	if logFmt, ok = logFormatsMap[strings.ToUpper(format)]; !ok {
		return LogFormatText
	}
	return logFmt
}

type Log struct {
	*slog.Logger
	FormatLevel
}

func NewLog() *Log {
	return NewLogWithOptions(DefaultOptions())
}

func NewLogWithOptions(options Options) *Log {
	var log *slog.Logger
	logFormat := options.Format
	logLevel := options.Level

	if options.Conf == nil {
		// fallback to default text log if no conf is provided
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

		// set default logger
		slog.SetDefault(log)
		return &Log{log, LogFormatText}
	}

	// resolve log level
	if options.Conf.Exists(envLogLevel) {
		logLevel, _ = Level(options.Conf.String(envLogLevel))
	}

	// resolve log format
	if options.Conf.Exists(envLogFormat) {
		logFormat = Format(options.Conf.String(envLogFormat))
	}

	handlerOptions := &slog.HandlerOptions{
		Level:       logLevel,
		ReplaceAttr: replaceAttr,
	}

	switch logFormat {
	case LogFormatJSON:
		log = slog.New(slog.NewJSONHandler(options.Out, handlerOptions))
	default:
		log = slog.New(slog.NewTextHandler(options.Out, handlerOptions))
	}

	// set default logger
	slog.SetDefault(log)
	return &Log{log, logFormat}
}

// replaceAttr - resolve log levels for TRACE and FATAL
func replaceAttr(_ []string, attr slog.Attr) slog.Attr {
	switch attr.Key {
	case slog.LevelKey:
		if lvl, err := Level(attr.Value.Resolve().String()); err == nil {
			logLevels := map[slog.Level]string{
				LevelTrace:      LevelTraceLabel,
				slog.LevelDebug: slog.LevelDebug.String(),
				slog.LevelInfo:  slog.LevelInfo.String(),
				slog.LevelWarn:  slog.LevelWarn.String(),
				slog.LevelError: slog.LevelError.String(),
				LevelFatal:      LevelFatalLabel,
			}
			attr.Value = slog.StringValue(logLevels[lvl])
		}
	}
	return attr
}

func exit(code int) {
	os.Exit(code)
}

// Wrappers
// https://pkg.go.dev/golang.org/x/exp/slog#hdr-Wrapping_output_methods

func (log *Log) WithError(err error) *Log {
	return &Log{log.With(logAttrError, err), log.FormatLevel}
}

func (log *Log) Tracef(format string, args ...any) {
	log.Log(context.Background(), LevelTrace, fmt.Sprintf(format, args...))
}

func (log *Log) Debugf(format string, args ...any) {
	log.Log(context.Background(), slog.LevelDebug, fmt.Sprintf(format, args...))
}

func (log *Log) Infof(format string, args ...any) {
	log.Log(context.Background(), slog.LevelInfo, fmt.Sprintf(format, args...))
}

func (log *Log) Warnf(format string, args ...any) {
	log.Log(context.Background(), slog.LevelWarn, fmt.Sprintf(format, args...))
}

func (log *Log) Errorf(format string, args ...any) {
	log.Log(context.Background(), slog.LevelError, fmt.Sprintf(format, args...))
}

func (log *Log) Fatalf(format string, args ...any) {
	log.Log(context.Background(), LevelFatal, fmt.Sprintf(format, args...))
	exit(1)
}
