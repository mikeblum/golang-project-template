package conftest

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	MockConfFile      = ".test.env"
	MockConfFilePerms = 0600
)

func SetupConf() (*os.File, error) {
	var conf *os.File
	var err error
	_, err = os.Stat(MockConfFile)
	if os.IsNotExist(err) {
		if conf, err = os.Create(MockConfFile); err != nil {
			return nil, err
		}
		defer conf.Close()
	} else {
		if conf, err = os.Open(MockConfFile); err != nil {
			return nil, err
		}
	}
	return conf, err
}

func CleanupConf(t *testing.T) {
	err := os.Remove(MockConfFile)
	assert.Nil(t, err)
}
