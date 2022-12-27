package conf

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
	assert.True(t, true, "True is true!")
}

func TestGetEnvVar(t *testing.T) {
	envShell := "SHELL"
	assert.True(t, len(strings.TrimSpace(GetEnv(envShell, ""))) > 0)
}

func TestGetEnvDefault(t *testing.T) {
	expected := "default env"
	assert.Equal(t, expected, strings.TrimSpace(GetEnv("", expected)))
}
