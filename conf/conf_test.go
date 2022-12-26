package conf

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	config, err := NewConfig()
	assert.Nil(t, err)
	assert.NotNil(t, config)
}

func TestGetEnvVar(t *testing.T) {
	envShell := "SHELL"
	assert.True(t, len(strings.TrimSpace(GetEnv(envShell, ""))) > 0)
}

func TestGetEnvDefault(t *testing.T) {
	expected := "default env"
	assert.Equal(t, expected, strings.TrimSpace(GetEnv("", expected)))
}
