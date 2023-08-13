package conftest

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestConfFile      = ".test.env"
	TestConfFilePerms = 0600
)

func SetupConf() (*os.File, error) {
	var conf *os.File
	var err error
	_, err = os.Stat(TestConfFile)
	if os.IsNotExist(err) {
		if conf, err = os.Create(TestConfFile); err != nil {
			return nil, err
		}
		defer conf.Close()
	} else {
		if conf, err = os.Open(TestConfFile); err != nil {
			return nil, err
		}
	}
	return conf, err
}

func CleanupConf(t *testing.T) {
	err := os.Remove(TestConfFile)
	assert.Nil(t, err)
}
