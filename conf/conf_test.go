package conf

import (
	"os"
	"strings"
	"testing"

	"github.com/mikeblum/golang-project-template/conftest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConf(t *testing.T) {
	// <setup code>
	suite, teardown := conftest.SetupSuite(t)
	// <teardown code>
	defer teardown(t, suite.Conf)
	t.Run("conf=new", NewConfTest)
	t.Run("conf=err", NewConfErrTest)
	t.Run("conf=default", DefaultConfNameTest)
	t.Run("conf=env-namespace", EnvConfTest)
	t.Run("conf=env-var", GetEnvVarTest)
	t.Run("conf=env-default", GetEnvDefaultTest)
}

func NewConfTest(t *testing.T) {
	conf, err := NewConf(Provider(conftest.TestConfFile))
	require.NoError(t, err)
	assert.NotNil(t, conf)
}

func NewConfErrTest(t *testing.T) {
	conf, err := NewConf(nil)
	require.Error(t, err)
	assert.Nil(t, conf)
}

func DefaultConfNameTest(t *testing.T) {
	assert.Equal(t, ConfFile, defaultConfName(""))
}

// NOTE: ENV_VARs must be declared before calling `NewConf`
func EnvConfTest(t *testing.T) {
	envVar := strings.Join([]string{EnvVarNamespace, "TEST"}, "_")
	expectedValue := "test_env_value"
	os.Setenv(envVar, expectedValue)
	defer os.Unsetenv(envVar)
	conf, err := NewConf(Provider(conftest.TestConfFile))
	require.NoError(t, err)
	assert.Equal(t, expectedValue, conf.String(envVar))
}

func GetEnvVarTest(t *testing.T) {
	envShell := "SHELL"
	assert.NotEmpty(t, strings.TrimSpace(GetEnv(envShell, "")))
}

func GetEnvDefaultTest(t *testing.T) {
	expected := "default env"
	assert.Equal(t, expected, strings.TrimSpace(GetEnv("", expected)))
}
