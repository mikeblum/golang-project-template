package log

import (
	"os"
	"strings"

	"github.com/knadh/koanf"
	"github.com/mikeblum/golang-project-template/conf"
	"github.com/sirupsen/logrus"
)

// logging configuration

const jsonLog = "JSON"
const envLogLevel = "LOG_LEVEL"
const envLogFormat = "LOG_FORMAT"

// NewLog - configure logging
func NewLog() *logrus.Entry {
	var cfg *koanf.Koanf
	var err error
	var cfgLogLevel string
	if cfg, err = conf.NewConf(conf.Provider("")); err != nil {
		// default to INFO
		cfgLogLevel = logrus.InfoLevel.String()
	} else {
		cfgLogLevel = cfg.String(envLogLevel)
	}
	logFormat := conf.GetEnv(envLogFormat, "")
	if strings.EqualFold(logFormat, jsonLog) {
		logrus.SetFormatter(&logrus.JSONFormatter{
			DisableHTMLEscape: true,
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			// timestamp with millisecond precision
			TimestampFormat: "Jan _2 15:04:05.00",
			ForceColors:     true,
		})
	}
	logrus.SetOutput(os.Stdout)
	var logLevel logrus.Level
	logLevel, err = logrus.ParseLevel(conf.GetEnv(envLogLevel, cfgLogLevel))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

	return logrus.WithFields(logrus.Fields{})
}
