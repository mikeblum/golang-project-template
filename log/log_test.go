package log

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"testing/slogtest"

	"github.com/mikeblum/golang-project-template/conf"
	"github.com/mikeblum/golang-project-template/conftest"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	// <setup code>
	suite, teardown := conftest.SetupSuite(t)
	// <teardown code>
	defer teardown(t, suite.Conf)
	t.Run("log=new", NewLogTest)
	t.Run("log=WithError", WithErrorTest)
	t.Run("log=format;JSON", func(t *testing.T) {
		LogTextTest(t)
		LogFormatConfTest(t, LogFormatJSON)
		LogFormatOptionsTest(t, LogFormatJSON)
	})
	t.Run("log=format;TEXT", func(t *testing.T) {
		LogTextTest(t)
		LogFormatConfTest(t, LogFormatText)
		LogFormatOptionsTest(t, LogFormatText)
	})
	t.Run("log=format;INVALID", LogFormatInvalidTest)
	t.Run("log=replaceAttr;TRACE", ReplaceAttrTraceTest)
	t.Run("log=replaceAttr;FATAL", ReplaceAttrFatalTest)
	t.Run("log=level;TRACE", func(t *testing.T) {
		LogLevelTraceTest(t)
		LogLevelConfTraceTest(t)
	})
	t.Run("log=level;DEBUG", func(t *testing.T) {
		LogLevelDebugTest(t)
		LogLevelConfDebugTest(t)
	})
	t.Run("log=level;INFO", LogLevelInfoTest)
	t.Run("log=level;WARN", LogLevelWarnTest)
	t.Run("log=level;ERROR", LogLevelErrorTest)
	t.Run("log=level;FATAL", LogLevelFatalTest)
	t.Run("log=level;INVALID", LogLevelInvalidTest)
	t.Run("log=level;EMPTY", LogLevelEmptyTest)
	t.Run("log=format;Tracef", LogFormatTracef)
	t.Run("log=format;Debugf", LogFormatDebugf)
	t.Run("log=format;Infof", LogFormatInfof)
	t.Run("log=format;Warnf", LogFormatWarnf)
	t.Run("log=format;Errorf", LogFormatErrorf)
	t.Run("log=format;Fatalf", LogFormatFatalf)
}

func jsonLogOptions() Options {
	return Options{
		Format: LogFormatJSON,
		Out:    os.Stdout,
	}
}

func textLogOptions() Options {
	return Options{
		Format: LogFormatText,
		Out:    os.Stdout,
	}
}

type TestLogEvent struct {
	// error being an interface{} causes encoding/json.Unmarshall to err
	Err string `json:"err"`
	Msg string `json:"msg"`
}

func NewLogTest(t *testing.T) {
	log := NewLog()
	assert.NotNil(t, log)
}

func WithErrorTest(t *testing.T) {
	log := NewLog()
	errLog := log.WithError(errors.New("test"))
	assert.NotNil(t, errLog)
	errLog.Error("test")
}

func LogJSONTest(t *testing.T) {
	log := NewLogWithOptions(jsonLogOptions())
	assert.NotNil(t, log)
}

func LogTextTest(t *testing.T) {
	log := NewLogWithOptions(jsonLogOptions())
	assert.NotNil(t, log)
}

func LogFormatInvalidTest(t *testing.T) {
	assert.Zero(t, LogFormat("ASCII"))
}

func ReplaceAttrTraceTest(t *testing.T) {
	actual := replaceAttr(make([]string, 0), slog.String(slog.LevelKey, LevelTraceName))
	assert.Equal(t, LevelTraceLabel, actual.Value.Resolve().String())
}

func ReplaceAttrFatalTest(t *testing.T) {
	actual := replaceAttr(make([]string, 0), slog.String(slog.LevelKey, LevelFatalName))
	assert.Equal(t, LevelFatalLabel, actual.Value.Resolve().String())
}

func LogLevelTest(t *testing.T, logLevel slog.Level) {
	os.Setenv(envLogLevel, logLevel.String())
	conf, err := conf.NewConf(conf.Provider(conftest.TestConfFile))
	assert.Nil(t, err)
	options := textLogOptions()
	options.Conf = conf
	log := NewLogWithOptions(options)
	assert.True(t, log.Enabled(context.TODO(), logLevel))
	actual, err := LogLevel(logLevel.String())
	assert.Nil(t, err)
	assert.Equal(t, logLevel, actual)
	os.Unsetenv(envLogLevel)
}

func LogLevelConfTest(t *testing.T, logLevel slog.Level) {
	os.Unsetenv(envLogLevel)
	confFile, err := os.Create(conftest.TestConfFile)
	assert.Nil(t, err)
	assert.NotNil(t, confFile)
	confFile.WriteString(fmt.Sprintf("%s=%s", envLogLevel, logLevel.String()))
	confFile.Sync()
	confFile.Close()
	conf, err := conf.NewConf(conf.Provider(conftest.TestConfFile))
	assert.Nil(t, err)
	options := textLogOptions()
	options.Conf = conf
	log := NewLogWithOptions(options)
	assert.True(t, log.Enabled(context.TODO(), logLevel))
}

func LogLevelOptionsTest(t *testing.T, logLevel slog.Level) {
	err := os.Unsetenv(envLogLevel)
	assert.Nil(t, err)
	err = os.Remove(conftest.TestConfFile)
	assert.Nil(t, err)
	conf, err := conf.NewConf(conf.Provider(conftest.TestConfFile))
	assert.Nil(t, err)
	options := textLogOptions()
	options.Conf = conf
	options.Level = logLevel
	log := NewLogWithOptions(options)
	assert.True(t, log.Enabled(context.TODO(), logLevel))
}

func LogFormatConfTest(t *testing.T, logFormat Format) {
	os.Unsetenv(envLogFormat)
	os.Unsetenv(envLogLevel)
	confFile, err := os.Create(conftest.TestConfFile)
	assert.Nil(t, err)
	assert.NotNil(t, confFile)
	confFile.WriteString(fmt.Sprintf("%s=%s", envLogFormat, LogFormats()[logFormat]))
	confFile.Sync()
	confFile.Close()
	conf, err := conf.NewConf(conf.Provider(conftest.TestConfFile))
	assert.Nil(t, err)
	log := NewLogWithOptions(Options{
		Conf:   conf,
		Format: logFormat,
		Level:  slog.LevelInfo,
	})
	assert.Equal(t, logFormat, log.Format)
}

func LogFormatOptionsTest(t *testing.T, logLevel slog.Level) {
	err := os.Unsetenv(envLogFormat)
	assert.Nil(t, err)
	// clear the .env file
	err = os.Truncate(conftest.TestConfFile, 0)
	assert.Nil(t, err)
	conf, err := conf.NewConf(conf.Provider(conftest.TestConfFile))
	assert.Nil(t, err)
	options := textLogOptions()
	options.Conf = conf
	options.Level = logLevel
	log := NewLogWithOptions(options)
	assert.True(t, log.Enabled(context.TODO(), logLevel))
}

func LogLevelTraceTest(t *testing.T) {
	lvl, err := LogLevel(LevelTraceName)
	assert.Nil(t, err)
	LogLevelTest(t, lvl)
}

func LogLevelConfTraceTest(t *testing.T) {
	lvl, err := LogLevel(LevelTraceName)
	assert.Nil(t, err)
	LogLevelConfTest(t, lvl)
}

func LogLevelOptionsTraceTest(t *testing.T) {
	lvl, err := LogLevel(LevelTraceName)
	assert.Nil(t, err)
	LogLevelOptionsTest(t, lvl)
}

func LogLevelDebugTest(t *testing.T) {
	lvl, err := LogLevel(slog.LevelDebug.String())
	assert.Nil(t, err)
	LogLevelTest(t, lvl)
}

func LogLevelConfDebugTest(t *testing.T) {
	lvl, err := LogLevel(slog.LevelDebug.String())
	assert.Nil(t, err)
	LogLevelConfTest(t, lvl)
}

func LogLevelOptionsDebugTest(t *testing.T) {
	lvl, err := LogLevel(slog.LevelDebug.String())
	assert.Nil(t, err)
	LogLevelOptionsTest(t, lvl)
}

func LogLevelInfoTest(t *testing.T) {
	lvl, err := LogLevel(slog.LevelInfo.String())
	assert.Nil(t, err)
	LogLevelTest(t, lvl)
}

func LogLevelWarnTest(t *testing.T) {
	lvl, err := LogLevel(slog.LevelWarn.String())
	assert.Nil(t, err)
	LogLevelTest(t, lvl)
}

func LogLevelErrorTest(t *testing.T) {
	lvl, err := LogLevel(slog.LevelError.String())
	assert.Nil(t, err)
	LogLevelTest(t, lvl)
}

func LogLevelFatalTest(t *testing.T) {
	lvl, err := LogLevel(LevelFatalName)
	assert.Nil(t, err)
	LogLevelTest(t, lvl)
}

func LogLevelInvalidTest(t *testing.T) {
	logFmt, err := LogLevel("NONE")
	assert.NotNil(t, err)
	assert.Equal(t, slog.LevelInfo, logFmt)
}

func LogLevelEmptyTest(t *testing.T) {
	logFmt, err := LogLevel("")
	assert.Nil(t, err)
	assert.Equal(t, slog.LevelInfo, logFmt)
}

func ParseAttrs(t *testing.T, buffer bytes.Buffer) []map[string]any {
	var attrs []map[string]any
	for _, line := range bytes.Split(buffer.Bytes(), []byte{'\n'}) {
		if len(line) == 0 {
			continue
		}
		var props map[string]any
		if err := json.Unmarshal(line, &props); err != nil {
			t.Fatal(err)
		}
		attrs = append(attrs, props)
	}
	return attrs
}

func LogFormatTest(t *testing.T, logLevelLabel string, logFormat string, args ...any) {
	// enable log level
	os.Setenv(envLogLevel, logLevelLabel)
	var buf bytes.Buffer
	options := jsonLogOptions()
	options.Out = &buf
	// write to buffer for test hook
	log := NewLogWithOptions(options)
	switch logLevelLabel {
	case LevelTraceLabel:
		log.Tracef(logFormat, args...)
	case slog.LevelDebug.String():
		log.Debugf(logFormat, args...)
	case slog.LevelInfo.String():
		log.Infof(logFormat, args...)
	case slog.LevelWarn.String():
		log.Warnf(logFormat, args...)
	case slog.LevelError.String():
		log.Errorf(logFormat, args...)
	case LevelFatalLabel:
		log.Fatalf(logFormat, args...)
	default:
		t.Fatalf("%s not a supported LOG_LEVEL", logLevelLabel)
	}
	// run LogHandlerTest harness for JSON logs
	if LogFormat(logFormat) == LogFormatJSON {
		err := slogtest.TestHandler(log.Handler(), func() []map[string]any {
			var params []map[string]any
			fmt.Println(buf.String())
			err := json.Unmarshal(buf.Bytes(), &params)
			assert.Nil(t, err)
			return params
		})
		assert.Nil(t, err)
	}
	os.Unsetenv(envLogLevel)
}

func LogFormatTracef(t *testing.T) {
	LogFormatTest(t, LevelTraceLabel, "trace log: %t", true)
}

func LogFormatDebugf(t *testing.T) {
	LogFormatTest(t, slog.LevelDebug.String(), "debug log: %t", true)
}

func LogFormatInfof(t *testing.T) {
	LogFormatTest(t, slog.LevelInfo.String(), "info log: %t", true)
}

func LogFormatWarnf(t *testing.T) {
	LogFormatTest(t, slog.LevelWarn.String(), "warn log: %t", true)
}

func LogFormatErrorf(t *testing.T) {
	LogFormatTest(t, slog.LevelError.String(), "error log: %t", true)
}

// NOTE: this test throws off the test coverage since its spawning a new `go test` proc
func LogFormatFatalf(t *testing.T) {
	envExit := "EXIT"
	exitCode := 1
	if os.Getenv(envExit) != strconv.Itoa(exitCode) {
		return
	}
	LogFormatTest(t, LevelFatalLabel, "fatal log: %t", true)
	// capture os.Exit escape
	cmd := exec.Command(os.Args[0], os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], "-test.run=LogFormatFatalf")
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%d", envExit, exitCode))
	err := cmd.Run()
	assert.NotNil(t, err)
	e, ok := err.(*exec.ExitError)
	assert.True(t, ok)
	assert.True(t, !e.Success())
}
