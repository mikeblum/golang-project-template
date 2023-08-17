package conftest

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestConfFile = ".test.env"
)

type Suite struct {
	Conf *os.File
}

func SetupSuite(t *testing.T) (*Suite, func(t *testing.T, conf *os.File)) {
	conf, err := SetupConf()
	assert.Nil(t, err)
	assert.NotNil(t, conf)
	return &Suite{
		Conf: conf,
	}, TeardownSuite
}

func TeardownSuite(t *testing.T, _ *os.File) {
	CleanupConf(t)
}

func SetupConf() (*os.File, error) {
	var err error
	_, err = os.Stat(TestConfFile)
	if os.IsNotExist(err) {
		if _, err = os.Create(TestConfFile); err != nil {
			return nil, err
		}
	}
	return os.Open(TestConfFile)
}

func CleanupConf(t *testing.T) {
	err := os.Remove(TestConfFile)
	assert.Nil(t, err)
}
