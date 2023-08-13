package log

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLog(t *testing.T) {
	log := NewLog()
	assert.NotNil(t, log)
}

// TODO: https://pkg.go.dev/golang.org/x/exp/slog/slogtest
func TestWithError(t *testing.T) {
	log := NewLog()
	errLog := log.WithError(fmt.Errorf("test err"))
	assert.NotNil(t, errLog)
	errLog.Error("test")
}

func TestInfof(_ *testing.T) {
	log := NewLog()
	log.Infof("interpolate this: %t", true)
	log.Info("interpolate this: ", true)
}
