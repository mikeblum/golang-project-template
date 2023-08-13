package conf

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mikeblum/golang-project-template/conftest"
	"github.com/stretchr/testify/assert"
)

type ConfTestSuite struct {
	conf *os.File
}

func setupSuite(t *testing.T) (*ConfTestSuite, func(t *testing.T, conf *os.File)) {
	conf, err := conftest.SetupConf()
	assert.Nil(t, err)
	return &ConfTestSuite{
		conf: conf,
	}, teardownSuite
}

func teardownSuite(t *testing.T, _ *os.File) {
	conftest.CleanupConf(t)
}

func TestConf(t *testing.T) {
	// <setup code>
	suite, teardown := setupSuite(t)
	// <teardown code>
	defer teardown(t, suite.conf)
	t.Run("conf=new", NewConfTest)
	t.Run("conf=err", NewConfErrTest)
	t.Run("conf=default", DefaultConfNameTest)
	t.Run("conf=dotenv", DotEnvConfTest)
	t.Run("conf=env-namespace", EnvConfTest)
	t.Run("conf=env-var", GetEnvVarTest)
	t.Run("conf=env-default", GetEnvDefaultTest)
}

func NewConfTest(t *testing.T) {
	conf, err := NewConf(Provider(conftest.TestConfFile))
	assert.Nil(t, err)
	assert.NotNil(t, conf)
}

func NewConfErrTest(t *testing.T) {
	conf, err := NewConf(nil)
	assert.NotNil(t, err)
	assert.Nil(t, conf)
}

func DefaultConfNameTest(t *testing.T) {
	assert.Equal(t, ConfFile, defaultConfName(""))
}

// NOTE: conf file must be populated before calling `NewConf`
func DotEnvConfTest(t *testing.T) {
	expectedKey := "test"
	expectedValue := "test_file_value"
	// !!WARN!! `` injects \tabs
	cfg := fmt.Sprintf("%s=%s", expectedKey, expectedValue)
	err := os.WriteFile(conftest.TestConfFile, []byte(cfg), conftest.TestConfFilePerms)
	assert.Nil(t, err)
	conf, err := NewConf(Provider(conftest.TestConfFile))
	assert.Nil(t, err)
	assert.Equal(t, expectedValue, conf.Get(expectedKey).(string))
}

// NOTE: ENV_VARs must be declared before calling `NewConf`
func EnvConfTest(t *testing.T) {
	envVar := strings.Join([]string{EnvVarNamespace, "TEST"}, "_")
	expectedKey := "test"
	expectedValue := "test_env_value"
	os.Setenv(envVar, expectedValue)
	defer os.Unsetenv(envVar)
	conf, err := NewConf(Provider(conftest.TestConfFile))
	assert.Nil(t, err)
	assert.Equal(t, expectedValue, conf.Get(expectedKey).(string))
}

func GetEnvVarTest(t *testing.T) {
	envShell := "SHELL"
	assert.True(t, len(strings.TrimSpace(GetEnv(envShell, ""))) > 0)
}

func GetEnvDefaultTest(t *testing.T) {
	expected := "default env"
	assert.Equal(t, expected, strings.TrimSpace(GetEnv("", expected)))
}
