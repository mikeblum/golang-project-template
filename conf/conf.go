package conf

import (
	"os"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

const (
	EnvConfigPath   = "CONF_PATH"
	EnvVarNamespace = ""
	EnvDelimiter    = "_"
	PropDelimiter   = "."
	ConfFile        = ".env"
	cwd             = "."
)

type Conf struct{}

// NewConf instantiates a new dotenv config with environment variables for context
// !!Note!! environment variables must be configured BEFORE calling NewConf
func NewConf(provider koanf.Provider) (*koanf.Koanf, error) {
	koan := koanf.New(cwd)
	if err := koan.Load(provider, dotenv.Parser()); err != nil {
		return nil, err
	}
	// load env variables under EnvVarNamespace namespace`
	err := koan.Load(env.Provider(EnvVarNamespace, EnvDelimiter, processEnvVar), nil)
	return koan, err
}

func Provider(confName string) *file.File {
	return file.Provider(defaultConfName(confName))
}

func defaultConfName(confName string) string {
	if len(strings.TrimSpace(confName)) == 0 {
		return ConfFile
	}
	return confName
}

func processEnvVar(s string) string {
	return strings.TrimPrefix(strings.Replace(strings.ToLower(
		strings.TrimPrefix(s, EnvVarNamespace)), EnvDelimiter, PropDelimiter, -1), PropDelimiter)
}

// GetEnv lookup an environment variable or fallback
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
